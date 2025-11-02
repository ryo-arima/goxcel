package parser

// NOTE: This file implements the XML-based (.gxl) node parsing using encoding/xml.
// It was previously named mdvue_parser.go; renamed to reflect the GXL format explicitly.

import (
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
}

type gxlRepository struct {
	Conf config.BaseConfig
}

// NewGxlRepository creates a repository from config.
func NewGxlRepository(conf config.BaseConfig) GxlRepository {
	return &gxlRepository{
		Conf: conf,
	}
}

// ReadGxl reads and parses the .gxl file.
func (r *gxlRepository) ReadGxl() (model.GXL, error) {
	if r.Conf.FilePath == "" {
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

// ReadGxlFromFile reads and parses a .gxl file from the given path.
func ReadGxlFromFile(filePath string, logger util.LoggerInterface) (model.GXL, error) {
	if filePath == "" {
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

// parseSheetTag parses a <Sheet> element.
func parseSheetTag(decoder *xml.Decoder, start xml.StartElement) (model.SheetTag, error) {
	sheet := model.SheetTag{
		Name: getAttr(start, "name"),
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
					Content: gridContent,
					Rows:    rows,
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

		parts := strings.Split(line, "|")
		var cells []string
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part != "" {
				cells = append(cells, part)
			}
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
