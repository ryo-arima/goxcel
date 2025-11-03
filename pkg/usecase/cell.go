package usecase

import (
	"fmt"
	"regexp"
	"strings"
)

// CellUsecase handles cell-level operations (value expansion, data resolution)
type CellUsecase interface {
	ExpandMustache(ctxStack []map[string]any, template string) string
	ResolvePath(ctxStack []map[string]any, path string) any
}

// DefaultCellUsecase is the default implementation of CellUsecase
type DefaultCellUsecase struct {
	mustacheRe *regexp.Regexp
}

// NewDefaultCellUsecase creates a new DefaultCellUsecase
func NewDefaultCellUsecase() *DefaultCellUsecase {
	return &DefaultCellUsecase{
		mustacheRe: regexp.MustCompile(`\{\{\s*([^}]+?)\s*\}\}`),
	}
}

// ExpandMustache replaces {{ varPath }} expressions with values from context stack
func (u *DefaultCellUsecase) ExpandMustache(ctxStack []map[string]any, template string) string {
	return u.mustacheRe.ReplaceAllStringFunc(template, func(match string) string {
		// Extract the variable path from {{ varPath }}
		submatch := u.mustacheRe.FindStringSubmatch(match)
		if len(submatch) < 2 {
			return match
		}

		path := strings.TrimSpace(submatch[1])
		value := u.ResolvePath(ctxStack, path)

		// Convert value to string
		return u.valueToString(value)
	})
}

// ResolvePath resolves a dot-separated path from the context stack
// Searches from the innermost (most recent) context to the outermost
func (u *DefaultCellUsecase) ResolvePath(ctxStack []map[string]any, path string) any {
	parts := strings.Split(path, ".")

	// Search from innermost to outermost context
	for i := len(ctxStack) - 1; i >= 0; i-- {
		value := u.resolveInContext(ctxStack[i], parts)
		if value != nil {
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
