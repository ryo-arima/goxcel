# Changelog

All notable changes to goxcel will be documented in this file.

## [0.1.1] - 2025-11-04

### Added
- Config-based dependency injection for usecase layer
- Logger integration throughout codebase with message codes
- Cell styling support: FontName, FontSize, FontColor, FillColor
- Sheet configuration: Row heights, column widths, freeze panes
- Style collector for dynamic OOXML style generation
- GXL extensions: SheetConfigTag, ColumnTag, RowHeightTag, StyleTag

### Changed
- Refactored usecase layer to use config.BaseConfig
- Removed nil checks for logger (always instantiated)
- Updated documentation structure (removed guide and development sections)
- Consolidated non-specification documentation

### Fixed
- XML structure conflict in XMLFillColor (separated XMLBgColor)
- Interface implementation: Added RenderSheet method to DefaultSheetUsecase
- Nil pointer dereference in logger usage

## [0.1.0] - 2025-11-03

### Added
- Core GXL template parsing
- Grid layout with pipe-delimited syntax
- Mustache-style expression interpolation `{{ .path }}`
- For loop iteration `<For each="item in items">`
- Anchor positioning `<Anchor ref="A1" />`
- Cell type inference (string, number, boolean, date, formula)
- Type hints `{{ .value:number }}`
- Markdown styling `**bold**`, `_italic_`
- Cell merging `<Merge range="A1:B2" />`
- Formula support `=SUM(A1:A10)`
- JSON and YAML data support
- CLI with generate command
- Dry-run mode for template validation
- Comprehensive logging system with message codes

### Implemented
- Clean architecture (config, controller, usecase, repository, model layers)
- GXL XML parser
- XLSX file writer using OOXML format
- Context stack for nested data access
- Dynamic style management
- Complete test coverage for core functionality

## [Planned for v1.1]

### To Be Added
- If/Else conditional structures
- Enhanced style tag with direct attributes
- Number formatting patterns
- Date formatting
- Currency formatting
- Sheet-level configuration via GXL
- Column width and row height specification
- Enhanced color support (themes, indexed colors)

### Under Consideration
- Image embedding `<Image>`
- Chart generation `<Chart>`
- Pivot tables `<Pivot>`
- Data validation
- Conditional formatting
- Named ranges
- Multiple data source support
- Template inheritance
- Macros/VBA support

## Version History

- **0.1.1** (2025-11-04): Config refactoring, styling features, documentation updates
- **0.1.0** (2025-11-03): Initial stable release with v1.0 features
