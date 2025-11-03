package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/ryo-arima/goxcel/pkg/model"
)

// SheetUsecase handles sheet-level rendering operations
type SheetUsecase interface {
	RenderSheet(ctx context.Context, sheetTag *model.SheetTag, data map[string]any) (*model.Sheet, error)
}

// DefaultSheetUsecase is the default implementation of SheetUsecase
type DefaultSheetUsecase struct {
	cellUsecase CellUsecase
}

// NewDefaultSheetUsecase creates a new DefaultSheetUsecase
func NewDefaultSheetUsecase() *DefaultSheetUsecase {
	return &DefaultSheetUsecase{
		cellUsecase: NewDefaultCellUsecase(),
	}
}

// RenderSheet renders a single sheet from a SheetTag
func (u *DefaultSheetUsecase) RenderSheet(ctx context.Context, sheetTag *model.SheetTag, data map[string]any) (*model.Sheet, error) {
	sheet := model.NewSheet(sheetTag.Name)

	// Initialize rendering state
	state := &renderState{
		sheet:     sheet,
		anchorRow: 1,
		anchorCol: 1,
		rowOffset: 0,
	}

	// Create initial context stack
	ctxStack := []map[string]any{data}

	// Render all nodes in the sheet
	if err := u.renderNodes(state, ctxStack, sheetTag.Nodes); err != nil {
		return nil, fmt.Errorf("render sheet %q: %w", sheetTag.Name, err)
	}

	return sheet, nil
}

// renderState holds the current rendering position and context
type renderState struct {
	sheet     *model.Sheet
	anchorRow int // Current anchor row (1-based)
	anchorCol int // Current anchor column (1-based)
	rowOffset int // Offset from anchor for sequential content
}

// renderNodes processes a list of nodes (tags) and renders them to the sheet
func (u *DefaultSheetUsecase) renderNodes(state *renderState, ctxStack []map[string]any, nodes []any) error {
	for _, node := range nodes {
		if err := u.renderNode(state, ctxStack, node); err != nil {
			return err
		}
	}
	return nil
}

// renderNode processes a single node and renders it to the sheet
func (u *DefaultSheetUsecase) renderNode(state *renderState, ctxStack []map[string]any, node any) error {
	switch v := node.(type) {
	case model.AnchorTag:
		return u.handleAnchor(state, v)
	case model.GridTag:
		return u.handleGrid(state, ctxStack, v)
	case model.GridRowTag:
		return u.handleGridRow(state, ctxStack, v)
	case model.MergeTag:
		return u.handleMerge(state, v)
	case model.ForTag:
		return u.handleFor(state, ctxStack, v)
	case model.ImageTag:
		return u.handleImage(state, v)
	case model.ShapeTag:
		return u.handleShape(state, v)
	case model.ChartTag:
		return u.handleChart(state, v)
	case model.PivotTag:
		return u.handlePivot(state, v)
	default:
		// Unknown node type, skip
		return nil
	}
}

// handleAnchor sets the anchor position
func (u *DefaultSheetUsecase) handleAnchor(state *renderState, tag model.AnchorTag) error {
	row, col, err := parseA1Ref(tag.Ref)
	if err != nil {
		return fmt.Errorf("invalid anchor ref %q: %w", tag.Ref, err)
	}
	state.anchorRow = row
	state.anchorCol = col
	state.rowOffset = 0
	return nil
}

// handleGrid renders a grid (table) to the sheet
func (u *DefaultSheetUsecase) handleGrid(state *renderState, ctxStack []map[string]any, tag model.GridTag) error {
	for _, row := range tag.Rows {
		if err := u.handleGridRow(state, ctxStack, row); err != nil {
			return err
		}
	}
	return nil
}

// handleGridRow renders a single row of cells
func (u *DefaultSheetUsecase) handleGridRow(state *renderState, ctxStack []map[string]any, row model.GridRowTag) error {
	currentRow := state.anchorRow + state.rowOffset

	for colIndex, cellValue := range row.Cells {
		col := state.anchorCol + colIndex
		ref := toA1Ref(currentRow, col)

		// Expand mustache templates and add cell
		expandedValue := u.cellUsecase.ExpandMustache(ctxStack, cellValue)
		cell := &model.Cell{
			Ref:   ref,
			Value: expandedValue,
			Type:  model.CellTypeString,
		}
		state.sheet.AddCell(cell)
	}

	state.rowOffset++
	return nil
}

// handleMerge adds a cell merge to the sheet
func (u *DefaultSheetUsecase) handleMerge(state *renderState, tag model.MergeTag) error {
	state.sheet.AddMerge(model.Merge{Range: tag.Range})
	return nil
}

// handleFor processes a for loop and renders its body multiple times
func (u *DefaultSheetUsecase) handleFor(state *renderState, ctxStack []map[string]any, tag model.ForTag) error {
	// Parse "Each" field like "item in items"
	parts := strings.Fields(tag.Each)
	if len(parts) != 3 || parts[1] != "in" {
		return fmt.Errorf("invalid For syntax: %q (expected: 'varName in dataPath')", tag.Each)
	}

	varName := parts[0]
	dataPath := parts[2]

	// Resolve the array from data context
	items := u.cellUsecase.ResolvePath(ctxStack, dataPath)

	// Handle different array types
	switch arr := items.(type) {
	case []any:
		for i, item := range arr {
			scope := map[string]any{
				varName: item,
				"loop": map[string]any{
					"index":  i,
					"number": i + 1,
				},
			}
			newStack := append(ctxStack, scope)
			if err := u.renderNodes(state, newStack, tag.Body); err != nil {
				return err
			}
		}
	case []map[string]any:
		for i, item := range arr {
			scope := map[string]any{
				varName: item,
				"loop": map[string]any{
					"index":  i,
					"number": i + 1,
				},
			}
			newStack := append(ctxStack, scope)
			if err := u.renderNodes(state, newStack, tag.Body); err != nil {
				return err
			}
		}
	default:
		// Not an iterable type, skip
	}

	return nil
}

// handleImage adds an image to the sheet
func (u *DefaultSheetUsecase) handleImage(state *renderState, tag model.ImageTag) error {
	state.sheet.AddImage(model.Image{
		Ref:      tag.Ref,
		Source:   tag.Src,
		WidthPx:  tag.Width,
		HeightPx: tag.Height,
	})
	return nil
}

// handleShape adds a shape to the sheet
func (u *DefaultSheetUsecase) handleShape(state *renderState, tag model.ShapeTag) error {
	state.sheet.AddShape(model.Shape{
		Ref:      tag.Ref,
		Kind:     tag.Kind,
		Text:     tag.Text,
		WidthPx:  tag.Width,
		HeightPx: tag.Height,
		Style:    tag.Style,
	})
	return nil
}

// handleChart adds a chart to the sheet
func (u *DefaultSheetUsecase) handleChart(state *renderState, tag model.ChartTag) error {
	state.sheet.AddChart(model.Chart{
		Ref:       tag.Ref,
		Type:      tag.Type,
		DataRange: tag.DataRange,
		Title:     tag.Title,
		WidthPx:   tag.Width,
		HeightPx:  tag.Height,
	})
	return nil
}

// handlePivot adds a pivot table to the sheet
func (u *DefaultSheetUsecase) handlePivot(state *renderState, tag model.PivotTag) error {
	rows := parseCommaSeparated(tag.Rows)
	cols := parseCommaSeparated(tag.Columns)
	vals := parseCommaSeparated(tag.Values)
	filt := parseCommaSeparated(tag.Filters)

	state.sheet.AddPivot(model.PivotTable{
		Ref:         tag.Ref,
		SourceRange: tag.SourceRange,
		Rows:        rows,
		Columns:     cols,
		Values:      vals,
		Filters:     filt,
	})
	return nil
}

// parseCommaSeparated splits a comma-separated string into a slice
func parseCommaSeparated(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// parseA1Ref parses an A1-style cell reference (e.g., "B3") into 1-based row and column
func parseA1Ref(ref string) (int, int, error) {
	if ref == "" {
		return 0, 0, fmt.Errorf("empty ref")
	}

	// Split letters and digits
	i := 0
	for i < len(ref) && ((ref[i] >= 'A' && ref[i] <= 'Z') || (ref[i] >= 'a' && ref[i] <= 'z')) {
		i++
	}
	if i == 0 || i == len(ref) {
		return 0, 0, fmt.Errorf("invalid ref: %s", ref)
	}

	colLetters := ref[:i]
	rowDigits := ref[i:]

	// Convert column letters to number
	col := 0
	for _, ch := range colLetters {
		uc := ch
		if uc >= 'a' && uc <= 'z' {
			uc = uc - 'a' + 'A'
		}
		if uc < 'A' || uc > 'Z' {
			return 0, 0, fmt.Errorf("invalid column in ref: %s", ref)
		}
		col = col*26 + int(uc-'A'+1)
	}

	// Convert row digits to number
	var row int
	for _, ch := range rowDigits {
		if ch < '0' || ch > '9' {
			return 0, 0, fmt.Errorf("invalid row in ref: %s", ref)
		}
		row = row*10 + int(ch-'0')
	}

	if row <= 0 || col <= 0 {
		return 0, 0, fmt.Errorf("non-positive ref: %s", ref)
	}

	return row, col, nil
}

// toA1Ref converts 1-based row and column to an A1 cell reference
func toA1Ref(row, col int) string {
	c := col
	letters := make([]byte, 0, 4)
	for c > 0 {
		c--
		letters = append([]byte{byte('A' + (c % 26))}, letters...)
		c /= 26
	}
	return fmt.Sprintf("%s%d", string(letters), row)
}
