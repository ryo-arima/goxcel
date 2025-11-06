package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/ryo-arima/goxcel/pkg/config"
	"github.com/ryo-arima/goxcel/pkg/model"
	"github.com/ryo-arima/goxcel/pkg/util"
)

// sheetRenderer is an internal renderer for sheet-level operations
type sheetRenderer struct {
	conf   config.BaseConfig
	logger util.Logger
	cell   *cellHelper
}

// newSheetRenderer creates a new internal sheet renderer
func newSheetRenderer(conf config.BaseConfig) *sheetRenderer {
	return &sheetRenderer{
		conf:   conf,
		logger: conf.Logger,
		cell:   newCellHelper(conf),
	}
}

// RenderSheet renders a SheetTag with data context into a Sheet
func (rcv *sheetRenderer) RenderSheet(ctx context.Context, sheetTag *model.SheetTag, data map[string]any) (*model.Sheet, error) {
	if sheetTag == nil {
		return nil, fmt.Errorf("sheet tag is nil")
	}

	sheet := model.NewSheet(sheetTag.Name)

	// Apply sheet-level defaults from tag config (if provided)
	if sheetTag.Config != nil {
		if sheetTag.Config.DefaultRowHeight > 0 {
			sheet.Config.DefaultRowHeight = sheetTag.Config.DefaultRowHeight
		}
		if sheetTag.Config.DefaultColumnWidth > 0 {
			sheet.Config.DefaultColumnWidth = sheetTag.Config.DefaultColumnWidth
		}
	}

	// Initialize render state
	state := &renderState{
		sheet:     sheet,
		anchorRow: 1,
		anchorCol: 1,
		rowOffset: 0,
	}

	// Create context stack from data
	ctxStack := []map[string]any{data}

	// Render all nodes in the sheet
	if err := rcv.renderNodes(state, ctxStack, sheetTag.Nodes); err != nil {
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
func (rcv *sheetRenderer) renderNodes(state *renderState, ctxStack []map[string]any, nodes []any) error {
	for _, node := range nodes {
		if err := rcv.renderNode(state, ctxStack, node); err != nil {
			return err
		}
	}
	return nil
}

// renderNode processes a single node and renders it to the sheet
func (rcv *sheetRenderer) renderNode(state *renderState, ctxStack []map[string]any, node any) error {
	switch v := node.(type) {
	case model.AnchorTag:
		return rcv.handleAnchor(state, v)
	case model.GridTag:
		return rcv.handleGrid(state, ctxStack, v)
	case model.GridRowTag:
		return rcv.handleGridRow(state, ctxStack, v)
	case model.MergeTag:
		return rcv.handleMerge(state, v)
	case model.ForTag:
		return rcv.handleFor(state, ctxStack, v)
	case model.ImageTag:
		return rcv.handleImage(state, v)
	case model.ShapeTag:
		return rcv.handleShape(state, v)
	case model.ChartTag:
		return rcv.handleChart(state, v)
	case model.PivotTag:
		return rcv.handlePivot(state, v)
	default:
		// Unknown node type, skip
		return nil
	}
}

// handleAnchor sets the anchor position
func (rcv *sheetRenderer) handleAnchor(state *renderState, tag model.AnchorTag) error {
	rcv.logger.DEBUG(util.USA1, fmt.Sprintf("Setting anchor position: %s", tag.Ref), nil)

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
func (rcv *sheetRenderer) handleGrid(state *renderState, ctxStack []map[string]any, tag model.GridTag) error {
	rcv.logger.DEBUG(util.USG1, fmt.Sprintf("Rendering grid with %d rows", len(tag.Rows)), nil)
	if tag.Ref != "" {
		return rcv.handleGridWithRef(state, ctxStack, tag)
	}
	return rcv.handleGridSequential(state, ctxStack, tag)
}

// handleGridWithRef renders a grid at an absolute position
func (rcv *sheetRenderer) handleGridWithRef(state *renderState, ctxStack []map[string]any, tag model.GridTag) error {
	row, col, err := parseA1Ref(tag.Ref)
	if err != nil {
		return fmt.Errorf("invalid grid ref %q: %w", tag.Ref, err)
	}

	// Save and restore state for absolute positioning
	return rcv.withSavedState(state, func() error {
		state.anchorRow = row
		state.anchorCol = col
		state.rowOffset = 0
		base := rcv.gridTagToStyle(ctxStack, tag)
		return rcv.renderGridRowsWithStyle(state, ctxStack, tag.Rows, base)
	})
}

// handleGridSequential renders a grid at the current position
func (rcv *sheetRenderer) handleGridSequential(state *renderState, ctxStack []map[string]any, tag model.GridTag) error {
	base := rcv.gridTagToStyle(ctxStack, tag)
	if base == nil {
		// No style attributes - use legacy path for compatibility
		return rcv.renderGridRows(state, ctxStack, tag.Rows)
	}
	return rcv.renderGridRowsWithStyle(state, ctxStack, tag.Rows, base)
}

// renderGridRows renders all rows in a grid (legacy wrapper for compatibility)
func (rcv *sheetRenderer) renderGridRows(state *renderState, ctxStack []map[string]any, rows []model.GridRowTag) error {
	// legacy: no base style
	return rcv.renderGridRowsWithStyle(state, ctxStack, rows, nil)
}

// renderGridRowsWithStyle renders all rows with a provided base style
func (rcv *sheetRenderer) renderGridRowsWithStyle(state *renderState, ctxStack []map[string]any, rows []model.GridRowTag, baseStyle *model.CellStyle) error {
	for _, row := range rows {
		if err := rcv.handleGridRowWithStyle(state, ctxStack, row, baseStyle); err != nil {
			return err
		}
	}
	return nil
}

// withSavedState executes a function while preserving the render state
func (rcv *sheetRenderer) withSavedState(state *renderState, fn func() error) error {
	savedAnchorRow := state.anchorRow
	savedAnchorCol := state.anchorCol
	savedRowOffset := state.rowOffset

	err := fn()

	state.anchorRow = savedAnchorRow
	state.anchorCol = savedAnchorCol
	state.rowOffset = savedRowOffset

	return err
}

// handleGridRow renders a single row of cells
func (rcv *sheetRenderer) handleGridRow(state *renderState, ctxStack []map[string]any, row model.GridRowTag) error {
	return rcv.handleGridRowWithStyle(state, ctxStack, row, nil)
}

func (rcv *sheetRenderer) handleGridRowWithStyle(state *renderState, ctxStack []map[string]any, row model.GridRowTag, baseStyle *model.CellStyle) error {
	currentRow := state.anchorRow + state.rowOffset

	for colIndex, cellValue := range row.Cells {
		col := state.anchorCol + colIndex
		cell := rcv.createCell(currentRow, col, cellValue, ctxStack, baseStyle)
		state.sheet.AddCell(cell)
	}

	state.rowOffset++
	return nil
}

// createCell creates a cell with proper type and style
func (rcv *sheetRenderer) createCell(row, col int, cellValue string, ctxStack []map[string]any, baseStyle *model.CellStyle) *model.Cell {
	ref := toA1Ref(row, col)

	// Expand mustache templates and infer cell type
	expandedValue, cellType := rcv.cell.ExpandMustacheWithType(ctxStack, cellValue)

	// Parse markdown style formatting
	cleanValue, style := rcv.cell.ParseMarkdownStyle(expandedValue)
	// Merge grid-level base style
	eff := mergeStyles(baseStyle, style)

	return &model.Cell{
		Ref:   ref,
		Value: cleanValue,
		Type:  cellType,
		Style: eff,
	}
}

// gridTagToStyle converts GridTag style hints into a CellStyle pointer with mustache expansion
func (rcv *sheetRenderer) gridTagToStyle(ctxStack []map[string]any, tag model.GridTag) *model.CellStyle {
	has := false
	st := &model.CellStyle{}
	if tag.FontName != "" {
		// Expand mustache in font name
		st.FontName = rcv.cell.ExpandMustache(ctxStack, tag.FontName)
		has = true
	}
	if tag.FontSize > 0 {
		st.FontSize = tag.FontSize
		has = true
	}
	if tag.FontColor != "" {
		// Expand mustache in color
		st.FontColor = rcv.cell.ExpandMustache(ctxStack, tag.FontColor)
		has = true
	}
	if tag.FillColor != "" {
		// Expand mustache in fill color
		st.FillColor = rcv.cell.ExpandMustache(ctxStack, tag.FillColor)
		has = true
	}
	if tag.BorderStyle != "" {
		b := &model.CellBorder{Style: tag.BorderStyle, Color: tag.BorderColor}
		// sides
		sides := strings.Split(tag.BorderSides, ",")
		if len(sides) == 0 || tag.BorderSides == "" || tag.BorderSides == "all" {
			b.Top, b.Right, b.Bottom, b.Left = true, true, true, true
		} else {
			for _, s := range sides {
				switch strings.TrimSpace(s) {
				case "all":
					b.Top, b.Right, b.Bottom, b.Left = true, true, true, true
				case "top":
					b.Top = true
				case "right":
					b.Right = true
				case "bottom":
					b.Bottom = true
				case "left":
					b.Left = true
				}
			}
		}
		st.Border = b
		has = true
	}
	if !has {
		return nil
	}
	return st
}

// mergeStyles overlays b over a (grid base a, cell-specific b). Returns nil if both are nil.
func mergeStyles(a, b *model.CellStyle) *model.CellStyle {
	if a == nil && b == nil {
		return nil
	}
	if a == nil {
		// Return a copy of b
		c := *b
		return &c
	}
	if b == nil {
		c := *a
		return &c
	}
	c := *a // start from base (grid)
	// Overlay booleans: if true in b, set true
	if b.Bold {
		c.Bold = true
	}
	if b.Italic {
		c.Italic = true
	}
	if b.Underline {
		c.Underline = true
	}
	// Overlay font attributes if specified in b
	if b.FontName != "" {
		c.FontName = b.FontName
	}
	if b.FontSize > 0 {
		c.FontSize = b.FontSize
	}
	if b.FontColor != "" {
		c.FontColor = b.FontColor
	}
	if b.FillColor != "" {
		c.FillColor = b.FillColor
	}
	// Border overlay: if b has border, override entirely
	if b.Border != nil {
		// Copy
		cb := *b.Border
		c.Border = &cb
	}
	// Alignment (future)
	if b.HAlign != "" {
		c.HAlign = b.HAlign
	}
	if b.VAlign != "" {
		c.VAlign = b.VAlign
	}
	return &c
}

// handleMerge adds a cell merge to the sheet
func (rcv *sheetRenderer) handleMerge(state *renderState, tag model.MergeTag) error {
	state.sheet.AddMerge(model.Merge{Range: tag.Range})
	return nil
}

// handleFor processes a for loop and renders its body multiple times
func (rcv *sheetRenderer) handleFor(state *renderState, ctxStack []map[string]any, tag model.ForTag) error {
	rcv.logger.DEBUG(util.USF1, fmt.Sprintf("Processing for loop: %s", tag.Each), nil)

	varName, dataPath, err := rcv.parseForSyntax(tag.Each)
	if err != nil {
		return err
	}

	items := rcv.cell.ResolvePath(ctxStack, dataPath)
	return rcv.iterateAndRender(state, ctxStack, varName, items, tag.Body)
}

// parseForSyntax parses "varName in dataPath" syntax
func (rcv *sheetRenderer) parseForSyntax(each string) (varName, dataPath string, err error) {
	parts := strings.Fields(each)
	if len(parts) != 3 || parts[1] != "in" {
		return "", "", fmt.Errorf("invalid For syntax: %q (expected: 'varName in dataPath')", each)
	}
	return parts[0], parts[2], nil
}

// iterateAndRender iterates over items and renders the body for each
func (rcv *sheetRenderer) iterateAndRender(state *renderState, ctxStack []map[string]any, varName string, items any, body []any) error {
	switch arr := items.(type) {
	case []any:
		return rcv.renderLoop(state, ctxStack, varName, arr, body)
	case []map[string]any:
		return rcv.renderMapLoop(state, ctxStack, varName, arr, body)
	default:
		// Not an iterable type, skip
		return nil
	}
}

// renderLoop renders loop body for []any array
func (rcv *sheetRenderer) renderLoop(state *renderState, ctxStack []map[string]any, varName string, items []any, body []any) error {
	for i, item := range items {
		scope := rcv.createLoopScope(varName, item, i)
		newStack := append(ctxStack, scope)
		if err := rcv.renderNodes(state, newStack, body); err != nil {
			return err
		}
	}
	return nil
}

// renderMapLoop renders loop body for []map[string]any array
func (rcv *sheetRenderer) renderMapLoop(state *renderState, ctxStack []map[string]any, varName string, items []map[string]any, body []any) error {
	for i, item := range items {
		scope := rcv.createLoopScope(varName, item, i)
		newStack := append(ctxStack, scope)
		if err := rcv.renderNodes(state, newStack, body); err != nil {
			return err
		}
	}
	return nil
}

// createLoopScope creates a loop variable scope
func (rcv *sheetRenderer) createLoopScope(varName string, item any, index int) map[string]any {
	return map[string]any{
		varName: item,
		"loop": map[string]any{
			"index":  index,
			"number": index + 1,
		},
	}
}

// handleImage adds an image to the sheet
func (rcv *sheetRenderer) handleImage(state *renderState, tag model.ImageTag) error {
	state.sheet.AddImage(model.Image{
		Ref:      tag.Ref,
		Source:   tag.Src,
		WidthPx:  tag.Width,
		HeightPx: tag.Height,
	})
	return nil
}

// handleShape adds a shape to the sheet
func (rcv *sheetRenderer) handleShape(state *renderState, tag model.ShapeTag) error {
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
func (rcv *sheetRenderer) handleChart(state *renderState, tag model.ChartTag) error {
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
func (rcv *sheetRenderer) handlePivot(state *renderState, tag model.PivotTag) error {
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
