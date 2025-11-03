package parser

// NOTE: This file implements the XML-based (.gxl) node parsing using encoding/xml.
// It was previously named mdvue_parser.go; renamed to reflect the GXL format explicitly.

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/ryo-arima/goxcel/pkg/config"
	"github.com/ryo-arima/goxcel/pkg/model"
	"github.com/ryo-arima/goxcel/pkg/util"
)

type GxlRepository interface {
	ReadGxl() (model.GXL, error)
	// FormatGxl pretty-prints the GXL template referenced by Conf.FilePath and returns the formatted bytes.
	// Note: This does not write to disk; the caller decides where to write.
	FormatGxl() ([]byte, error)
}

type gxlRepository struct {
	Conf config.BaseConfig
}

// NewGxlRepository creates a repository from config.
func NewGxlRepository(conf config.BaseConfig) GxlRepository {
	return &gxlRepository{Conf: conf}
}

// ReadGxl reads and parses the .gxl file from config.
func (r *gxlRepository) ReadGxl() (model.GXL, error) {
	if strings.TrimSpace(r.Conf.FilePath) == "" {
		r.Conf.Logger.ERROR(util.RP2, "File path is not set in config", nil)
		return model.GXL{}, fmt.Errorf("file path is not set in config")
	}
	r.Conf.Logger.DEBUG(util.RP1, "Reading GXL file", map[string]interface{}{"file": r.Conf.FilePath})
	gxl, err := ReadGxlFromFile(r.Conf.FilePath, r.Conf.Logger)
	if err != nil {
		r.Conf.Logger.ERROR(util.RP2, "Failed to read GXL file", map[string]interface{}{"file": r.Conf.FilePath, "error": err.Error()})
		return model.GXL{}, err
	}
	r.Conf.Logger.INFO(util.RP1, "GXL file parsed successfully", map[string]interface{}{"file": r.Conf.FilePath, "sheets": len(gxl.Sheets)})
	return gxl, nil
}

// FormatGxl pretty-prints the .gxl file specified in the repository config and returns the formatted bytes.
func (r *gxlRepository) FormatGxl() ([]byte, error) {
	if strings.TrimSpace(r.Conf.FilePath) == "" {
		r.Conf.Logger.ERROR(util.RP2, "File path is not set in config", nil)
		return nil, fmt.Errorf("file path is not set in config")
	}
	r.Conf.Logger.DEBUG(util.RP1, "Formatting GXL file", map[string]interface{}{"file": r.Conf.FilePath})
	f, err := os.Open(r.Conf.FilePath)
	if err != nil {
		r.Conf.Logger.ERROR(util.FSR2, "Failed to open file for formatting", map[string]interface{}{"file": r.Conf.FilePath, "error": err.Error()})
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()
	out, err := prettyFormatGXL(f)
	if err != nil {
		r.Conf.Logger.ERROR(util.XMLU2, "Failed to format GXL XML", map[string]interface{}{"file": r.Conf.FilePath, "error": err.Error()})
		return nil, err
	}
	r.Conf.Logger.INFO(util.GXLP1, "GXL file formatted successfully", map[string]interface{}{"file": r.Conf.FilePath})
	return out, nil
}

// ReadGxlFromFile reads and parses a .gxl file from the given path.
func ReadGxlFromFile(filePath string, logger util.LoggerInterface) (model.GXL, error) {
	if strings.TrimSpace(filePath) == "" {
		logger.ERROR(util.FSR2, "File path is empty", nil)
		return model.GXL{}, fmt.Errorf("file path is empty")
	}
	logger.DEBUG(util.FSO1, "Opening GXL file", map[string]interface{}{"file": filePath})
	file, err := os.Open(filePath)
	if err != nil {
		logger.ERROR(util.FSR2, "Failed to open file", map[string]interface{}{"file": filePath, "error": err.Error()})
		return model.GXL{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	logger.DEBUG(util.XMLU1, "Parsing GXL XML content", nil)
	gxl, err := parseGXL(file)
	if err != nil {
		logger.ERROR(util.XMLU2, "Failed to parse GXL XML", map[string]interface{}{"file": filePath, "error": err.Error()})
		return model.GXL{}, err
	}
	logger.INFO(util.GXLP1, "GXL file parsed successfully", map[string]interface{}{"sheets": len(gxl.Sheets)})
	return gxl, nil
}

// parseGXL parses XML content into GXL structure.
func parseGXL(r io.Reader) (model.GXL, error) {
	decoder := xml.NewDecoder(r)
	decoder.Strict = false
	var gxl model.GXL
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return model.GXL{}, fmt.Errorf("failed to decode XML: %w", err)
		}
		switch se := token.(type) {
		case xml.StartElement:
			switch se.Name.Local {
			case "Book":
				if name := getAttr(se, "name"); name != "" {
					gxl.BookTag.Name = name
				}
			case "Sheet":
				sheet, err := parseSheetTag(decoder, se)
				if err != nil {
					return model.GXL{}, err
				}
				gxl.Sheets = append(gxl.Sheets, sheet)
			}
		}
	}
	return gxl, nil
}

// writeAlignedGrid formats pipe-delimited rows so that '|' columns align.
// It indents each produced line by indentLevel (2 spaces per level).
func writeAlignedGrid(buf *bytes.Buffer, content string, indentLevel int) {
	// Split into lines and collect grid rows (lines starting with '|', allowing leading spaces)
	lines := strings.Split(content, "\n")
	var rows [][]string
	var isGridLine []bool

	// Preprocess: trim right whitespace of each line; detect grid lines
	for _, ln := range lines {
		// normalize CRLF/CR handled by split; trim right spaces
		lnr := strings.TrimRight(ln, " \t\r")
		ltrim := strings.TrimLeft(lnr, " \t")
		if strings.HasPrefix(ltrim, "|") {
			// parse cells
			parts := strings.Split(ltrim, "|")
			if len(parts) > 0 && parts[0] == "" {
				parts = parts[1:]
			}
			if len(parts) > 0 && parts[len(parts)-1] == "" {
				parts = parts[:len(parts)-1]
			}
			var cells []string
			for _, p := range parts {
				cells = append(cells, strings.TrimSpace(p))
			}
			rows = append(rows, cells)
			isGridLine = append(isGridLine, true)
		} else {
			rows = append(rows, []string{lnr})
			isGridLine = append(isGridLine, false)
		}
	}

	// Compute max width per column across grid lines
	var colWidths []int
	for i, cells := range rows {
		if !isGridLine[i] {
			continue
		}
		for c, cell := range cells {
			l := runeLen(cell)
			if c >= len(colWidths) {
				colWidths = append(colWidths, l)
			} else if l > colWidths[c] {
				colWidths[c] = l
			}
		}
	}

	// Writer helpers
	writeIndent := func(level int) {
		for i := 0; i < level; i++ {
			buf.WriteString("  ")
		}
	}

	// Emit lines; collapse multiple blank lines
	blankPending := false
	wroteAny := false
	for i, cells := range rows {
		if !isGridLine[i] {
			// Non-grid line
			if strings.TrimSpace(cells[0]) == "" {
				// blank line: defer to avoid double newlines
				if !blankPending {
					blankPending = true
				}
				continue
			}
			// flush pending blank
			if blankPending {
				if wroteAny { // only emit if something has been written already
					buf.WriteByte('\n')
				}
				buf.WriteByte('\n')
				blankPending = false
			}
			buf.WriteByte('\n')
			writeIndent(indentLevel)
			buf.WriteString(cells[0])
			wroteAny = true
			continue
		}

		// grid line
		// flush pending blank before a grid line
		if blankPending {
			if wroteAny { // suppress leading blank before first grid row
				buf.WriteByte('\n')
			}
			blankPending = false
		}
		buf.WriteByte('\n')
		writeIndent(indentLevel)
		buf.WriteByte('|')
		for c, cell := range cells {
			// one space after pipe
			buf.WriteByte(' ')
			buf.WriteString(cell)
			// pad to column width
			pad := 0
			if c < len(colWidths) {
				pad = colWidths[c] - runeLen(cell)
			}
			for k := 0; k < pad; k++ {
				buf.WriteByte(' ')
			}
			// one space then closing pipe
			buf.WriteByte(' ')
			buf.WriteByte('|')
		}
		// ensure trailing spaces trimmed (no-op here) and keep newline added at the start of next iteration/after loop
		wroteAny = true
	}
	// End with a newline so caller can place closing tag on next line
	buf.WriteByte('\n')
}

// runeLen returns the number of runes (code points) in s.
// Note: This does not account for East Asian full-width display; it's a simple approximation.
func runeLen(s string) int {
	return len([]rune(s))
}

// parseSheetTag parses a <Sheet> element.
func parseSheetTag(decoder *xml.Decoder, start xml.StartElement) (model.SheetTag, error) {
	sheet := model.SheetTag{
		Name: getAttr(start, "name"),
	}

	// Optional defaults
	if cw := getAttr(start, "col_width"); cw != "" {
		if w, _ := util.ParseColWidth(cw); w > 0 {
			if sheet.Config == nil {
				sheet.Config = &model.SheetConfigTag{}
			}
			sheet.Config.DefaultColumnWidth = w
		}
	}
	// Support both row_height (canonical) and row_heigh (typo-compatible)
	if rh := getAttr(start, "row_height"); rh != "" {
		if h, _ := util.ParseRowHeight(rh); h > 0 {
			if sheet.Config == nil {
				sheet.Config = &model.SheetConfigTag{}
			}
			sheet.Config.DefaultRowHeight = h
		}
	} else if rh2 := getAttr(start, "row_heigh"); rh2 != "" {
		if h, _ := util.ParseRowHeight(rh2); h > 0 {
			if sheet.Config == nil {
				sheet.Config = &model.SheetConfigTag{}
			}
			sheet.Config.DefaultRowHeight = h
		}
	}

	for {
		token, err := decoder.Token()
		if err != nil {
			return sheet, err
		}
		switch se := token.(type) {
		case xml.StartElement:
			node, err := parseNodeTag(decoder, se)
			if err != nil {
				return sheet, err
			}
			if node != nil {
				sheet.Nodes = append(sheet.Nodes, node)
			}
		case xml.EndElement:
			if se.Name.Local == "Sheet" {
				return sheet, nil
			}
		}
	}
}

// parseNodeTag parses individual node elements.
func parseNodeTag(decoder *xml.Decoder, start xml.StartElement) (any, error) {
	switch start.Name.Local {
	case "Anchor":
		node := model.AnchorTag{Ref: getAttr(start, "ref")}
		if err := skipToEnd(decoder, "Anchor"); err != nil {
			return nil, err
		}
		return node, nil

	case "Merge":
		node := model.MergeTag{Range: getAttr(start, "range")}
		if err := skipToEnd(decoder, "Merge"); err != nil {
			return nil, err
		}
		return node, nil

	case "Grid":
		return parseGridTag(decoder, start)

	case "Image":
		node := model.ImageTag{
			Ref: getAttr(start, "ref"),
			Src: getAttr(start, "src"),
		}
		if w := getAttr(start, "width"); w != "" {
			node.Width, _ = strconv.Atoi(w)
		}
		if h := getAttr(start, "height"); h != "" {
			node.Height, _ = strconv.Atoi(h)
		}
		if err := skipToEnd(decoder, "Image"); err != nil {
			return nil, err
		}
		return node, nil

	case "Shape":
		node := model.ShapeTag{
			Ref:   getAttr(start, "ref"),
			Kind:  getAttr(start, "kind"),
			Text:  getAttr(start, "text"),
			Style: getAttr(start, "style"),
		}
		if w := getAttr(start, "width"); w != "" {
			node.Width, _ = strconv.Atoi(w)
		}
		if h := getAttr(start, "height"); h != "" {
			node.Height, _ = strconv.Atoi(h)
		}
		if err := skipToEnd(decoder, "Shape"); err != nil {
			return nil, err
		}
		return node, nil

	case "Chart":
		node := model.ChartTag{
			Ref:       getAttr(start, "ref"),
			Type:      getAttr(start, "type"),
			DataRange: getAttr(start, "dataRange"),
			Title:     getAttr(start, "title"),
		}
		if w := getAttr(start, "width"); w != "" {
			node.Width, _ = strconv.Atoi(w)
		}
		if h := getAttr(start, "height"); h != "" {
			node.Height, _ = strconv.Atoi(h)
		}
		if err := skipToEnd(decoder, "Chart"); err != nil {
			return nil, err
		}
		return node, nil

	case "Pivot":
		node := model.PivotTag{
			Ref:         getAttr(start, "ref"),
			SourceRange: getAttr(start, "sourceRange"),
			Rows:        getAttr(start, "rows"),
			Columns:     getAttr(start, "columns"),
			Values:      getAttr(start, "values"),
			Filters:     getAttr(start, "filters"),
			Options:     getAttr(start, "options"),
		}
		if err := skipToEnd(decoder, "Pivot"); err != nil {
			return nil, err
		}
		return node, nil

	case "For":
		return parseForTag(decoder, start)

	case "If":
		return parseIfTag(decoder, start)

	case "Style":
		node := model.StyleTag{
			Selector: getAttr(start, "selector"),
			Name:     getAttr(start, "name"),
			ID:       getAttr(start, "id"),
			Class:    getAttr(start, "class"),
		}
		if err := skipToEnd(decoder, "Style"); err != nil {
			return nil, err
		}
		return node, nil

	default:
		if err := skipToEnd(decoder, start.Name.Local); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

// parseGridTag parses <Grid> content.
func parseGridTag(decoder *xml.Decoder, start xml.StartElement) (model.GridTag, error) {
	var content strings.Builder

	// Extract ref attribute if present
	var ref string
	var fontName string
	var fontSize int
	var fontColor string
	var fillColor string
	var borderStyle string
	var borderColor string
	var borderSides string
	for _, attr := range start.Attr {
		if attr.Name.Local == "ref" {
			ref = attr.Value
		} else if attr.Name.Local == "font" || attr.Name.Local == "font_name" || attr.Name.Local == "fontName" {
			fontName = attr.Value
		} else if attr.Name.Local == "font_size" || attr.Name.Local == "fontSize" || attr.Name.Local == "text_size" {
			if v, err := strconv.Atoi(attr.Value); err == nil {
				fontSize = v
			}
		} else if attr.Name.Local == "font_color" || attr.Name.Local == "fontColor" || attr.Name.Local == "text_color" {
			fontColor = sanitizeColor(attr.Value)
		} else if attr.Name.Local == "fill_color" || attr.Name.Local == "fillColor" || attr.Name.Local == "color" {
			fillColor = sanitizeColor(attr.Value)
		} else if attr.Name.Local == "border" || attr.Name.Local == "border_style" || attr.Name.Local == "borderStyle" {
			borderStyle = strings.ToLower(strings.TrimSpace(attr.Value))
		} else if attr.Name.Local == "border_color" || attr.Name.Local == "borderColor" {
			borderColor = sanitizeColor(attr.Value)
		} else if attr.Name.Local == "border_sides" || attr.Name.Local == "borderSides" {
			borderSides = strings.ToLower(strings.TrimSpace(attr.Value))
		}
	}

	for {
		token, err := decoder.Token()
		if err != nil {
			return model.GridTag{}, err
		}

		switch t := token.(type) {
		case xml.CharData:
			content.Write(t)
		case xml.EndElement:
			if t.Name.Local == "Grid" {
				gridContent := content.String()
				rows := parseGridContent(gridContent)
				return model.GridTag{
					Content:     gridContent,
					Rows:        rows,
					Ref:         ref,
					FontName:    fontName,
					FontSize:    fontSize,
					FontColor:   fontColor,
					FillColor:   fillColor,
					BorderStyle: borderStyle,
					BorderColor: borderColor,
					BorderSides: borderSides,
				}, nil
			}
		}
	}
}

// parseGridContent parses pipe-delimited grid content.
func parseGridContent(content string) []model.GridRowTag {
	var rows []model.GridRowTag
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "|") {
			continue
		}

		// Split by pipe and process
		parts := strings.Split(line, "|")

		// Remove first and last element if they are empty (from leading/trailing pipes)
		if len(parts) > 0 && parts[0] == "" {
			parts = parts[1:]
		}
		if len(parts) > 0 && parts[len(parts)-1] == "" {
			parts = parts[:len(parts)-1]
		}

		// Build cells array, preserving empty cells
		var cells []string
		for _, part := range parts {
			// Trim spaces but keep the cell even if empty
			cells = append(cells, strings.TrimSpace(part))
		}

		if len(cells) > 0 {
			rows = append(rows, model.GridRowTag{Cells: cells})
		}
	}

	return rows
}

// parseForTag parses <For> loop.
func parseForTag(decoder *xml.Decoder, start xml.StartElement) (model.ForTag, error) {
	forTag := model.ForTag{
		Each: getAttr(start, "each"),
	}

	for {
		token, err := decoder.Token()
		if err != nil {
			return forTag, err
		}

		switch se := token.(type) {
		case xml.StartElement:
			node, err := parseNodeTag(decoder, se)
			if err != nil {
				return forTag, err
			}
			if node != nil {
				forTag.Body = append(forTag.Body, node)
			}
		case xml.EndElement:
			if se.Name.Local == "For" {
				return forTag, nil
			}
		}
	}
}

// parseIfTag parses <If> conditional.
func parseIfTag(decoder *xml.Decoder, start xml.StartElement) (model.IfTag, error) {
	ifTag := model.IfTag{
		Cond: getAttr(start, "cond"),
	}

	inElse := false

	for {
		token, err := decoder.Token()
		if err != nil {
			return ifTag, err
		}

		switch se := token.(type) {
		case xml.StartElement:
			if se.Name.Local == "Else" {
				inElse = true
				continue
			}

			node, err := parseNodeTag(decoder, se)
			if err != nil {
				return ifTag, err
			}
			if node != nil {
				if inElse {
					ifTag.Else = append(ifTag.Else, node)
				} else {
					ifTag.Then = append(ifTag.Then, node)
				}
			}
		case xml.EndElement:
			if se.Name.Local == "If" {
				return ifTag, nil
			}
		}
	}
}

// getAttr retrieves attribute value.
func getAttr(start xml.StartElement, name string) string {
	for _, attr := range start.Attr {
		if attr.Name.Local == name {
			return attr.Value
		}
	}
	return ""
}

// skipToEnd skips to matching end element.
func skipToEnd(decoder *xml.Decoder, elementName string) error {
	depth := 1
	for depth > 0 {
		token, err := decoder.Token()
		if err != nil {
			return err
		}

		switch t := token.(type) {
		case xml.StartElement:
			if t.Name.Local == elementName {
				depth++
			}
		case xml.EndElement:
			if t.Name.Local == elementName {
				depth--
			}
		}
	}
	return nil
}

// sanitizeColor removes leading '#' and uppercases hex
func sanitizeColor(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "#")
	return strings.ToUpper(s)
}

// prettyFormatGXL pretty-prints XML tags with indentation while preserving character data and comments.
func prettyFormatGXL(r io.Reader) ([]byte, error) {
	dec := xml.NewDecoder(r)
	dec.Strict = false

	var buf bytes.Buffer

	type frame struct {
		name       string
		hasContent bool
		hasChild   bool
		isGrid     bool
		gridBuf    bytes.Buffer
	}
	var stack []frame
	depth := 0
	lastWasNewline := true // xml.Header ends with a newline

	writeIndent := func() {
		for i := 0; i < depth; i++ {
			buf.WriteString("  ")
		}
	}

	// Always write XML header
	buf.WriteString(xml.Header)

	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("decode xml: %w", err)
		}

		switch t := tok.(type) {
		case xml.StartElement:
			// starting a child element -> ensure newline and indent (except root)
			if depth > 0 {
				if !lastWasNewline {
					buf.WriteByte('\n')
					lastWasNewline = true
				}
				writeIndent()
			}
			// write start tag
			buf.WriteByte('<')
			buf.WriteString(t.Name.Local)
			for _, a := range t.Attr {
				buf.WriteByte(' ')
				buf.WriteString(a.Name.Local)
				buf.WriteString("=\"")
				var esc bytes.Buffer
				if err := xml.EscapeText(&esc, []byte(a.Value)); err != nil {
					return nil, fmt.Errorf("escape attr: %w", err)
				}
				buf.Write(esc.Bytes())
				buf.WriteByte('"')
			}
			buf.WriteByte('>')
			lastWasNewline = false

			// mark parent hasChild
			if n := len(stack); n > 0 {
				stack[n-1].hasChild = true
			}
			// push frame
			stack = append(stack, frame{name: t.Name.Local, isGrid: t.Name.Local == "Grid"})
			depth++

		case xml.EndElement:
			// close current element
			if depth > 0 {
				depth--
			}
			// pop frame
			var fr frame
			if n := len(stack); n > 0 {
				fr = stack[n-1]
				stack = stack[:n-1]
			}

			if fr.isGrid {
				// Handle Grid content alignment
				content := fr.gridBuf.String()
				trimmed := strings.TrimSpace(content)
				if trimmed == "" {
					// treat as empty element -> inline close
					buf.WriteString(" </")
					buf.WriteString(t.Name.Local)
					buf.WriteByte('>')
					lastWasNewline = false
				} else {
					// write aligned grid lines at indent level depth+1
					indentLevel := depth + 1
					writeAlignedGrid(&buf, content, indentLevel)
					// after content, write closing tag on its own line
					lastWasNewline = true
					writeIndent()
					buf.WriteString("</")
					buf.WriteString(t.Name.Local)
					buf.WriteByte('>')
					lastWasNewline = false
				}

			} else if !fr.hasContent && !fr.hasChild {
				// empty element -> inline close with a single space
				buf.WriteString(" </")
				buf.WriteString(t.Name.Local)
				buf.WriteByte('>')
				lastWasNewline = false
			} else if fr.hasChild && !fr.hasContent {
				// children present, no direct text -> close on its own line
				if !lastWasNewline {
					buf.WriteByte('\n')
					lastWasNewline = true
				}
				writeIndent()
				buf.WriteString("</")
				buf.WriteString(t.Name.Local)
				buf.WriteByte('>')
				lastWasNewline = false
			} else {
				// had text content -> close inline with text
				buf.WriteString("</")
				buf.WriteString(t.Name.Local)
				buf.WriteByte('>')
				lastWasNewline = false
			}

		case xml.CharData:
			data := []byte(t)
			if n := len(stack); n > 0 && stack[n-1].isGrid {
				// Buffer inside Grid for alignment later
				stack[n-1].gridBuf.Write(data)
				stack[n-1].hasContent = stack[n-1].hasContent || len(bytes.TrimSpace(data)) > 0
			} else {
				if len(bytes.TrimSpace(data)) == 0 {
					// ignore pure whitespace to avoid double newlines/extra spaces
					continue
				}
				// write text as-is
				buf.Write(data)
				if n := len(stack); n > 0 {
					stack[n-1].hasContent = true
				}
				if len(data) > 0 && data[len(data)-1] == '\n' {
					lastWasNewline = true
				} else {
					lastWasNewline = false
				}
			}

		case xml.Comment:
			// treat comments like children on their own line
			if !lastWasNewline {
				buf.WriteByte('\n')
				lastWasNewline = true
			}
			writeIndent()
			buf.WriteString("<!--")
			buf.Write([]byte(t))
			buf.WriteString("-->")
			lastWasNewline = false
			if n := len(stack); n > 0 {
				stack[n-1].hasChild = true
			}

		case xml.Directive:
			// Skip other directives (we already wrote header)
		default:
			// ignore
		}
	}

	if !lastWasNewline {
		buf.WriteByte('\n')
	}
	return buf.Bytes(), nil
}
