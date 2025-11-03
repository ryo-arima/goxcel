# Render Phase

## Overview

Evaluates AST with data context to produce populated sheet models.

## Process

1. **Context Initialization**: Load JSON data into context stack
2. **AST Traversal**: Walk tree depth-first
3. **Expression Evaluation**: Resolve `{{ varPath }}` expressions
4. **Loop Expansion**: Iterate For loops with data
5. **Grid Positioning**: Calculate cell positions (row/col)
6. **Type Inference**: Determine cell types (number, bool, formula, etc.)
7. **Style Parsing**: Extract markdown styles (**bold**, _italic_)
8. **Cell Creation**: Build Cell models with values and metadata

## Input/Output

**Input**: `model.GXL` (AST) + JSON data  
**Output**: `model.Sheet` with populated cells

## Key Components

- `SheetUsecase.RenderSheet()`: Main renderer
- `CellUsecase.ExpandMustacheWithType()`: Expression evaluator
- `CellUsecase.InferCellType()`: Type inference
- `CellUsecase.ParseMarkdownStyle()`: Style parser

## Example

```
Template: | {{ name }} |
Data: {name: "Alice"}
Result: Cell{Ref: "A1", Value: "Alice", Type: String}
```
