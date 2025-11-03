package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/ryo-arima/goxcel/pkg/config"
	"github.com/ryo-arima/goxcel/pkg/model"
	"github.com/ryo-arima/goxcel/pkg/util"
)

// SheetUsecase handles sheet-level rendering operations
type SheetUsecase interface {
	RenderSheet(ctx context.Context, sheetTag *model.SheetTag, data map[string]any) (*model.Sheet, error)
}

// DefaultSheetUsecase is the default implementation of SheetUsecase
type DefaultSheetUsecase struct {
	conf        config.BaseConfig
	logger      util.LoggerInterface
	cellUsecase CellUsecase
}

// NewSheetUsecase creates a new SheetUsecase with config
func NewSheetUsecase(conf config.BaseConfig) SheetUsecase {
	return &DefaultSheetUsecase{
		conf:        conf,
		logger:      conf.Logger,
		cellUsecase: NewCellUsecase(conf),
	}
}

// NewDefaultSheetUsecase creates a new DefaultSheetUsecase (deprecated: use NewSheetUsecase)
func NewDefaultSheetUsecase() *DefaultSheetUsecase {
	conf := config.NewBaseConfig()
	return &DefaultSheetUsecase{
		conf:        conf,
		logger:      conf.Logger,
		cellUsecase: NewCellUsecase(conf),
	}
}

// RenderSheet renders a SheetTag with data context into a Sheet
func (u *DefaultSheetUsecase) RenderSheet(ctx context.Context, sheetTag *model.SheetTag, data map[string]any) (*model.Sheet, error) {
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
	gridStyle *model.CellStyle // Default style for current Grid
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
	u.logger.DEBUG(util.USA1, fmt.Sprintf("Setting anchor position: %s", tag.Ref), nil)

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
	u.logger.DEBUG(util.USG1, fmt.Sprintf("Rendering grid with %d rows", len(tag.Rows)), nil)
	if tag.Ref != "" {
		return u.handleGridWithRef(state, ctxStack, tag)
	}
	return u.handleGridSequential(state, ctxStack, tag)
}

// handleGridWithRef renders a grid at an absolute position
func (u *DefaultSheetUsecase) handleGridWithRef(state *renderState, ctxStack []map[string]any, tag model.GridTag) error {
	row, col, err := parseA1Ref(tag.Ref)
	if err != nil {
		return fmt.Errorf("invalid grid ref %q: %w", tag.Ref, err)
	}

	// Save and restore state for absolute positioning
	return u.withSavedState(state, func() error {
		state.anchorRow = row
		state.anchorCol = col
		state.rowOffset = 0
		state.gridStyle = gridTagToStyle(tag)
		return u.renderGridRows(state, ctxStack, tag.Rows)
	})
}

// handleGridSequential renders a grid at the current position
func (u *DefaultSheetUsecase) handleGridSequential(state *renderState, ctxStack []map[string]any, tag model.GridTag) error {
	// Save/restore grid style around this grid
	saved := state.gridStyle
	state.gridStyle = gridTagToStyle(tag)
	err := u.renderGridRows(state, ctxStack, tag.Rows)
	state.gridStyle = saved
	return err
}

// renderGridRows renders all rows in a grid
func (u *DefaultSheetUsecase) renderGridRows(state *renderState, ctxStack []map[string]any, rows []model.GridRowTag) error {
	for _, row := range rows {
		if err := u.handleGridRow(state, ctxStack, row); err != nil {
			return err
		}
	}
	return nil
}

// withSavedState executes a function while preserving the render state
func (u *DefaultSheetUsecase) withSavedState(state *renderState, fn func() error) error {
	savedAnchorRow := state.anchorRow
	savedAnchorCol := state.anchorCol
	savedRowOffset := state.rowOffset
	savedGridStyle := state.gridStyle

	err := fn()

	state.anchorRow = savedAnchorRow
	state.anchorCol = savedAnchorCol
	state.rowOffset = savedRowOffset
	state.gridStyle = savedGridStyle

	return err
}

// handleGridRow renders a single row of cells
func (u *DefaultSheetUsecase) handleGridRow(state *renderState, ctxStack []map[string]any, row model.GridRowTag) error {
	currentRow := state.anchorRow + state.rowOffset

	for colIndex, cellValue := range row.Cells {
		col := state.anchorCol + colIndex
		cell := u.createCell(currentRow, col, cellValue, ctxStack, state.gridStyle)
		state.sheet.AddCell(cell)
	}

	state.rowOffset++
	return nil
}

// createCell creates a cell with proper type and style
func (u *DefaultSheetUsecase) createCell(row, col int, cellValue string, ctxStack []map[string]any, baseStyle *model.CellStyle) *model.Cell {
	ref := toA1Ref(row, col)

	// Expand mustache templates and infer cell type
	expandedValue, cellType := u.cellUsecase.ExpandMustacheWithType(ctxStack, cellValue)

	// Parse markdown style formatting
	cleanValue, style := u.cellUsecase.ParseMarkdownStyle(expandedValue)
	// Merge grid-level base style
	eff := mergeStyles(baseStyle, style)

	return &model.Cell{
		Ref:   ref,
		Value: cleanValue,
		Type:  cellType,
		Style: eff,
	}
}

// gridTagToStyle converts GridTag style hints into a CellStyle pointer (or nil if none)
func gridTagToStyle(tag model.GridTag) *model.CellStyle {
	has := false
	st := &model.CellStyle{}
	if tag.FontName != "" {
		st.FontName = tag.FontName
		has = true
	}
	if tag.FontSize > 0 {
		st.FontSize = tag.FontSize
		has = true
	}
	if tag.FontColor != "" {
		st.FontColor = tag.FontColor
		has = true
	}
	if tag.FillColor != "" {
		st.FillColor = tag.FillColor
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
func (u *DefaultSheetUsecase) handleMerge(state *renderState, tag model.MergeTag) error {
	state.sheet.AddMerge(model.Merge{Range: tag.Range})
	return nil
}

// handleFor processes a for loop and renders its body multiple times
func (u *DefaultSheetUsecase) handleFor(state *renderState, ctxStack []map[string]any, tag model.ForTag) error {
	u.logger.DEBUG(util.USF1, fmt.Sprintf("Processing for loop: %s", tag.Each), nil)

	varName, dataPath, err := u.parseForSyntax(tag.Each)
	if err != nil {
		return err
	}

	items := u.cellUsecase.ResolvePath(ctxStack, dataPath)
	return u.iterateAndRender(state, ctxStack, varName, items, tag.Body)
}

// parseForSyntax parses "varName in dataPath" syntax
func (u *DefaultSheetUsecase) parseForSyntax(each string) (varName, dataPath string, err error) {
	parts := strings.Fields(each)
	if len(parts) != 3 || parts[1] != "in" {
		return "", "", fmt.Errorf("invalid For syntax: %q (expected: 'varName in dataPath')", each)
	}
	return parts[0], parts[2], nil
}

// iterateAndRender iterates over items and renders the body for each
func (u *DefaultSheetUsecase) iterateAndRender(state *renderState, ctxStack []map[string]any, varName string, items any, body []any) error {
	switch arr := items.(type) {
	case []any:
		return u.renderLoop(state, ctxStack, varName, arr, body)
	case []map[string]any:
		return u.renderMapLoop(state, ctxStack, varName, arr, body)
	default:
		// Not an iterable type, skip
		return nil
	}
}

// renderLoop renders loop body for []any array
func (u *DefaultSheetUsecase) renderLoop(state *renderState, ctxStack []map[string]any, varName string, items []any, body []any) error {
	for i, item := range items {
		scope := u.createLoopScope(varName, item, i)
		newStack := append(ctxStack, scope)
		if err := u.renderNodes(state, newStack, body); err != nil {
			return err
		}
	}
	return nil
}

// renderMapLoop renders loop body for []map[string]any array
func (u *DefaultSheetUsecase) renderMapLoop(state *renderState, ctxStack []map[string]any, varName string, items []map[string]any, body []any) error {
	for i, item := range items {
		scope := u.createLoopScope(varName, item, i)
		newStack := append(ctxStack, scope)
		if err := u.renderNodes(state, newStack, body); err != nil {
			return err
		}
	}
	return nil
}

// createLoopScope creates a loop variable scope
func (u *DefaultSheetUsecase) createLoopScope(varName string, item any, index int) map[string]any {
	return map[string]any{
		varName: item,
		"loop": map[string]any{
			"index":  index,
			"number": index + 1,
		},
	}
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
