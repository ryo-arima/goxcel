package usecase

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/ryo-arima/goxcel/pkg/model"
	"github.com/ryo-arima/goxcel/pkg/util"
)

// CellUsecase handles cell-level operations (value expansion, data resolution)
type CellUsecase interface {
	ExpandMustache(ctxStack []map[string]any, template string) string
	ResolvePath(ctxStack []map[string]any, path string) any
	InferCellType(value string) model.CellType
	ParseTypeHint(template string) (string, model.CellType)
	ExpandMustacheWithType(ctxStack []map[string]any, template string) (string, model.CellType)
	ParseMarkdownStyle(text string) (string, *model.CellStyle)
}

// DefaultCellUsecase is the default implementation of CellUsecase
type DefaultCellUsecase struct {
	mustacheRe *regexp.Regexp
	typeHintRe *regexp.Regexp
	numberRe   *regexp.Regexp
	dateRe     *regexp.Regexp
	boldRe     *regexp.Regexp
	italicRe   *regexp.Regexp
	logger     util.LoggerInterface
}

// NewDefaultCellUsecase creates a new DefaultCellUsecase
func NewDefaultCellUsecase() *DefaultCellUsecase {
	return NewDefaultCellUsecaseWithLogger(nil)
}

// NewDefaultCellUsecaseWithLogger creates a new DefaultCellUsecase with logger
func NewDefaultCellUsecaseWithLogger(logger util.LoggerInterface) *DefaultCellUsecase {
	return &DefaultCellUsecase{
		mustacheRe: regexp.MustCompile(`\{\{\s*([^}]+?)\s*\}\}`),
		// Type hint syntax: {{ .value:type }}
		typeHintRe: regexp.MustCompile(`:\s*(int|float|number|bool|boolean|date|string)\s*$`),
		// Number pattern: optional sign, digits, optional decimal part
		numberRe: regexp.MustCompile(`^-?\d+(\.\d+)?$`),
		// Date pattern: ISO format YYYY-MM-DD or YYYY-MM-DDTHH:MM:SS
		dateRe: regexp.MustCompile(`^\d{4}-\d{2}-\d{2}(T\d{2}:\d{2}:\d{2})?`),
		// Markdown style patterns
		boldRe:   regexp.MustCompile(`\*\*(.+?)\*\*`),
		italicRe: regexp.MustCompile(`_(.+?)_`),
		logger:   logger,
	}
}

// ExpandMustache replaces {{ varPath }} expressions with values from context stack
func (u *DefaultCellUsecase) ExpandMustache(ctxStack []map[string]any, template string) string {
	result, _ := u.ExpandMustacheWithType(ctxStack, template)
	return result
}

// ExpandMustacheWithType replaces {{ varPath }} expressions and returns the value with its type
func (u *DefaultCellUsecase) ExpandMustacheWithType(ctxStack []map[string]any, template string) (string, model.CellType) {
	u.logger.DEBUG(util.UCE1, fmt.Sprintf("Expanding mustache template: %s", template), nil)

	var detectedType model.CellType = model.CellTypeAuto
	expansionCount := 0

	result := u.mustacheRe.ReplaceAllStringFunc(template, func(match string) string {
		expansionCount++

		// Extract expression from {{ }}
		expr := u.extractExpression(match)
		if expr == "" {
			return match
		}

		// Parse type hint if present
		cleanPath, typeHint := u.ParseTypeHint(expr)
		if typeHint != model.CellTypeAuto {
			detectedType = typeHint
		}

		// Resolve value from context
		value := u.ResolvePath(ctxStack, cleanPath)

		// Convert to string
		return u.valueToString(value)
	})

	// Determine final cell type
	finalType := u.determineFinalType(result, detectedType, expansionCount)
	u.logger.DEBUG(util.UCE2, fmt.Sprintf("Expanded result: %s (type: %s)", result, finalType), nil)
	return result, finalType
}

// extractExpression extracts the expression from {{ expr }}
func (u *DefaultCellUsecase) extractExpression(match string) string {
	submatch := u.mustacheRe.FindStringSubmatch(match)
	if len(submatch) < 2 {
		return ""
	}
	return strings.TrimSpace(submatch[1])
}

// determineFinalType determines the final cell type based on context
func (u *DefaultCellUsecase) determineFinalType(result string, detectedType model.CellType, expansionCount int) model.CellType {
	// Multiple expansions → always string
	if expansionCount > 1 {
		return model.CellTypeString
	}

	// Explicit type hint → use it
	if detectedType != model.CellTypeAuto {
		return detectedType
	}

	// Otherwise → infer from result
	return u.InferCellType(result)
}

// ParseTypeHint extracts type hint from a mustache expression
// Input: ".value:int" -> Output: (".value", CellTypeNumber)
func (u *DefaultCellUsecase) ParseTypeHint(expr string) (string, model.CellType) {
	matches := u.typeHintRe.FindStringSubmatch(expr)
	if len(matches) == 0 {
		return expr, model.CellTypeAuto
	}

	// Remove type hint from expression
	cleanExpr := u.typeHintRe.ReplaceAllString(expr, "")
	cleanExpr = strings.TrimSpace(cleanExpr)

	// Map type hint to CellType
	typeHint := strings.ToLower(matches[1])
	switch typeHint {
	case "int", "float", "number":
		return cleanExpr, model.CellTypeNumber
	case "bool", "boolean":
		return cleanExpr, model.CellTypeBoolean
	case "date":
		return cleanExpr, model.CellTypeDate
	case "string":
		return cleanExpr, model.CellTypeString
	default:
		return cleanExpr, model.CellTypeAuto
	}
}

// InferCellType infers the cell type from a string value
func (u *DefaultCellUsecase) InferCellType(value string) model.CellType {
	u.logger.DEBUG(util.UCT1, fmt.Sprintf("Inferring cell type for value: %s", value), nil)

	if value == "" {
		return model.CellTypeString
	}

	if u.isFormula(value) {
		return model.CellTypeFormula
	}

	if u.isBoolean(value) {
		return model.CellTypeBoolean
	}

	if u.isNumber(value) {
		return model.CellTypeNumber
	}

	if u.isDate(value) {
		return model.CellTypeDate
	}

	return model.CellTypeString
}

// isFormula checks if value starts with =
func (u *DefaultCellUsecase) isFormula(value string) bool {
	return strings.HasPrefix(value, "=")
}

// isBoolean checks if value is true or false
func (u *DefaultCellUsecase) isBoolean(value string) bool {
	lowerValue := strings.ToLower(strings.TrimSpace(value))
	return lowerValue == "true" || lowerValue == "false"
}

// isNumber checks if value matches number pattern
func (u *DefaultCellUsecase) isNumber(value string) bool {
	return u.numberRe.MatchString(value)
}

// isDate checks if value matches ISO date pattern and is parseable
func (u *DefaultCellUsecase) isDate(value string) bool {
	if !u.dateRe.MatchString(value) {
		return false
	}

	// Validate it's actually parseable as a date
	_, err := time.Parse("2006-01-02", value[:10])
	return err == nil
}

// ResolvePath resolves a dot-separated path from the context stack
// Searches from the innermost (most recent) context to the outermost
// Also handles string literals like "text" and numeric literals
func (u *DefaultCellUsecase) ResolvePath(ctxStack []map[string]any, path string) any {
	u.logger.DEBUG(util.UCR1, fmt.Sprintf("Resolving path: %s", path), nil)

	// Try to resolve as literal first
	if literal := u.tryResolveLiteral(path); literal != nil {
		return literal
	}

	// Resolve from context stack
	return u.resolveFromContext(ctxStack, path)
}

// tryResolveLiteral attempts to resolve the path as a literal value
func (u *DefaultCellUsecase) tryResolveLiteral(path string) any {
	// String literal: "text" or 'text'
	if u.isStringLiteral(path) {
		return path[1 : len(path)-1]
	}

	// Numeric literal
	if u.numberRe.MatchString(path) {
		return path
	}

	// Boolean literal
	lowerPath := strings.ToLower(path)
	if lowerPath == "true" || lowerPath == "false" {
		return lowerPath
	}

	return nil
}

// isStringLiteral checks if path is a quoted string
func (u *DefaultCellUsecase) isStringLiteral(path string) bool {
	return (strings.HasPrefix(path, `"`) && strings.HasSuffix(path, `"`)) ||
		(strings.HasPrefix(path, `'`) && strings.HasSuffix(path, `'`))
}

// resolveFromContext resolves a path from the context stack
func (u *DefaultCellUsecase) resolveFromContext(ctxStack []map[string]any, path string) any {
	// Remove leading dot (e.g., ".quantity" -> "quantity")
	cleanPath := strings.TrimPrefix(path, ".")
	parts := strings.Split(cleanPath, ".")

	// Search from innermost to outermost context
	for i := len(ctxStack) - 1; i >= 0; i-- {
		if value := u.resolveInContext(ctxStack[i], parts); value != nil {
			return value
		}
	}

	return nil
}

// resolveInContext resolves a path within a single context map
func (u *DefaultCellUsecase) resolveInContext(context map[string]any, parts []string) any {
	var current any = context

	for _, part := range parts {
		contextMap, ok := current.(map[string]any)
		if !ok {
			return nil
		}

		value, exists := contextMap[part]
		if !exists {
			return nil
		}

		current = value
	}

	return current
}

// valueToString converts a value to its string representation
func (u *DefaultCellUsecase) valueToString(value any) string {
	switch v := value.(type) {
	case nil:
		return ""
	case string:
		return v
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", v)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%v", v)
	case bool:
		if v {
			return "true"
		}
		return "false"
	default:
		// For other types, use Sprint
		return fmt.Sprint(v)
	}
}

// ParseMarkdownStyle parses markdown-style formatting and returns clean text with style
// Supports: **bold**, _italic_
func (u *DefaultCellUsecase) ParseMarkdownStyle(text string) (string, *model.CellStyle) {
	u.logger.DEBUG(util.UCS1, fmt.Sprintf("Parsing markdown style: %s", text), nil)

	style := &model.CellStyle{}
	cleanText := text

	// Check for bold: **text**
	if u.boldRe.MatchString(text) {
		style.Bold = true
		cleanText = u.boldRe.ReplaceAllString(cleanText, "$1")
	}

	// Check for italic: _text_
	if u.italicRe.MatchString(cleanText) {
		style.Italic = true
		cleanText = u.italicRe.ReplaceAllString(cleanText, "$1")
	}

	// Return nil style if no formatting was found
	if !style.Bold && !style.Italic && !style.Underline {
		return text, nil
	}

	return cleanText, style
}
