package usecase

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/ryo-arima/goxcel/pkg/model"
)

// Renderer converts a GXL and input data into a model.Book workbook.
type Renderer interface {
	Render(ctx context.Context, t *model.GXL, data any) (*model.Book, error)
}

// DefaultRenderer is a placeholder implementation to be extended.
type DefaultRenderer struct{}

// Render renders the template into a minimal empty workbook for now.
func (DefaultRenderer) Render(ctx context.Context, t *model.GXL, data any) (*model.Book, error) {
	_ = ctx
	if t == nil {
		return nil, errors.New("renderer: template is nil")
	}
	b := model.NewBook()
	for _, sh := range t.Sheets {
		s := model.NewSheet(sh.Name)
		st := renderState{sheet: s, anchorRow: 1, anchorCol: 1, rowOffset: 0}
		// root context
		ctxStack := []map[string]any{normalizeData(data)}
		if err := renderNodes(&st, ctxStack, sh.Nodes); err != nil {
			return nil, err
		}
		b.AddSheet(s)
	}
	return b, nil
}

// parseA1Ref parses an A1-style cell reference (e.g., "B3") into 1-based row and column.
func parseA1Ref(ref string) (int, int, error) {
	if ref == "" {
		return 0, 0, fmt.Errorf("empty ref")
	}
	// split letters and digits
	i := 0
	for i < len(ref) && ((ref[i] >= 'A' && ref[i] <= 'Z') || (ref[i] >= 'a' && ref[i] <= 'z')) {
		i++
	}
	if i == 0 || i == len(ref) {
		return 0, 0, fmt.Errorf("invalid ref: %s", ref)
	}
	colLetters := ref[:i]
	rowDigits := ref[i:]
	// column letters to number
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

// toA1Ref converts 1-based row and column to an A1 cell reference.
func toA1Ref(row, col int) string {
	// column
	c := col
	letters := make([]byte, 0, 4)
	for c > 0 {
		c--
		letters = append([]byte{byte('A' + (c % 26))}, letters...)
		c /= 26
	}
	return fmt.Sprintf("%s%d", string(letters), row)
}

// --- minimal rendering engine helpers ---

type renderState struct {
	sheet     *model.Sheet
	anchorRow int
	anchorCol int
	rowOffset int
}

func renderNodes(st *renderState, ctxStack []map[string]any, nodes []any) error {
	for _, n := range nodes {
		switch v := n.(type) {
		case model.AnchorTag:
			r, c, err := parseA1Ref(v.Ref)
			if err != nil {
				return fmt.Errorf("renderer: invalid anchor ref %q: %w", v.Ref, err)
			}
			st.anchorRow, st.anchorCol = r, c
			st.rowOffset = 0
		case model.GridTag:
			// GridTag contains parsed GridRowTags, render them
			for _, row := range v.Rows {
				r := st.anchorRow + st.rowOffset
				for j, cellVal := range row.Cells {
					ref := toA1Ref(r, st.anchorCol+j)
					val := expandMustache(ctxStack, cellVal)
					st.sheet.AddCell(&model.Cell{Ref: ref, Value: val, Type: model.CellTypeString})
				}
				st.rowOffset++
			}
		case model.GridRowTag:
			r := st.anchorRow + st.rowOffset
			for j, cellVal := range v.Cells {
				ref := toA1Ref(r, st.anchorCol+j)
				val := expandMustache(ctxStack, cellVal)
				st.sheet.AddCell(&model.Cell{Ref: ref, Value: val, Type: model.CellTypeString})
			}
			st.rowOffset++
		case model.MergeTag:
			st.sheet.AddMerge(model.Merge{Range: v.Range})
		case model.ForTag:
			// Parse "Each" field like "item in items"
			parts := strings.Fields(v.Each)
			if len(parts) != 3 || parts[1] != "in" {
				return fmt.Errorf("invalid For syntax: %q", v.Each)
			}
			varName := parts[0]
			dataKey := parts[2]
			items := resolvePath(ctxStack, dataKey)
			switch arr := items.(type) {
			case []any:
				for i, it := range arr {
					scope := map[string]any{varName: it, "loop": map[string]any{"index": i, "number": i + 1}}
					if err := renderNodes(st, append(ctxStack, scope), v.Body); err != nil {
						return err
					}
				}
			case []map[string]any:
				for i, it := range arr {
					scope := map[string]any{varName: it, "loop": map[string]any{"index": i, "number": i + 1}}
					if err := renderNodes(st, append(ctxStack, scope), v.Body); err != nil {
						return err
					}
				}
			default:
				// ignore if not iterable
			}
		case model.ImageTag:
			st.sheet.AddImage(model.Image{Ref: v.Ref, Source: v.Src, WidthPx: v.Width, HeightPx: v.Height})
		case model.ShapeTag:
			st.sheet.AddShape(model.Shape{Ref: v.Ref, Kind: v.Kind, Text: v.Text, WidthPx: v.Width, HeightPx: v.Height, Style: v.Style})
		case model.ChartTag:
			st.sheet.AddChart(model.Chart{Ref: v.Ref, Type: v.Type, DataRange: v.DataRange, Title: v.Title, WidthPx: v.Width, HeightPx: v.Height})
		case model.PivotTag:
			// Parse comma-separated strings into slices
			rows := parseCommaSeparated(v.Rows)
			cols := parseCommaSeparated(v.Columns)
			vals := parseCommaSeparated(v.Values)
			filt := parseCommaSeparated(v.Filters)
			st.sheet.AddPivot(model.PivotTable{Ref: v.Ref, SourceRange: v.SourceRange, Rows: rows, Columns: cols, Values: vals, Filters: filt})
		}
	}
	return nil
}

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

func normalizeData(data any) map[string]any {
	if m, ok := data.(map[string]any); ok {
		return m
	}
	return map[string]any{"data": data}
}

var mustacheRe = regexp.MustCompile(`\{\{\s*([^}]+?)\s*\}\}`)

func expandMustache(ctxStack []map[string]any, s string) string {
	return mustacheRe.ReplaceAllStringFunc(s, func(m string) string {
		sub := mustacheRe.FindStringSubmatch(m)
		if len(sub) < 2 {
			return m
		}
		v := resolvePath(ctxStack, strings.TrimSpace(sub[1]))
		switch vv := v.(type) {
		case nil:
			return ""
		case string:
			return vv
		default:
			return fmt.Sprint(vv)
		}
	})
}

func resolvePath(ctxStack []map[string]any, path string) any {
	parts := strings.Split(path, ".")
	for i := len(ctxStack) - 1; i >= 0; i-- {
		var cur any = ctxStack[i]
		ok := true
		for _, p := range parts {
			if m, mok := cur.(map[string]any); mok {
				cur, ok = m[p]
				if !ok {
					cur = nil
					break
				}
				continue
			}
			ok = false
			break
		}
		if ok {
			return cur
		}
	}
	return nil
}
