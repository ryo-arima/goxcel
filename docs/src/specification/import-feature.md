# Import Feature Specification

## Overview

The Import feature allows `.gxl` template files to import sheets from other `.gxl` files. Each import creates a new sheet in the output workbook, enabling modular template composition and reusability.

**Key Characteristics:**
- Import tags are placed at the **Book level** only (not inside Sheet tags)
- Each `<Import>` creates a **new sheet** in the workbook
- Imports are resolved at render time with circular dependency detection
- Relative paths are resolved from the importing file's directory

## Syntax

```xml
<Book name="BookName">
  <Import src="path/to/external.gxl" sheet="SheetName" />
  
  <Sheet name="MainSheet">
    <!-- sheet content -->
  </Sheet>
</Book>
```

### Attributes

- `src` (required): Path to the external `.gxl` file
  - Relative paths are resolved from the importing file's directory
  - Absolute paths are supported
  - Example: `./common.gxl`, `../shared/headers.gxl`, `/path/to/template.gxl`
  
- `sheet` (required): Name of the sheet to import
  - Must exactly match a sheet name in the imported file
  - Case-sensitive
  - Example: `"Headers"`, `"CompanyHeader"`, `"PageFooter"`

### Placement Rules

✅ **Valid** - Import at Book level:
```xml
<Book name="Report">
  <Import sProcessing Flow

1. **Parse Time**: Parser records Import tags in `GXL.Imports[]` (book level only)
2. **Render Time**: Book usecase processes imports before rendering main sheets:
   - Resolves file path relative to template directory
   - Loads and parses the imported `.gxl` file
   - Finds the specified sheet by name
   - Checks for circular imports
   - Renders the imported sheet
   - Adds it as a new sheet to the workbook
3. **Output**: Each import creates a separate sheet in the final Excel file

### Content Expansion

**Sheet Creation Model:**
- Each `<Import>` tag creates a **new sheet** in the output workbook
- The sheet retains its original name from the imported file
- Imported sheets appear in the order they are declared
- Main file sheets are added after all imports

**Example:**
```xml
<Book name="Report">
  <Import src="./headers.gxl" sheet="CompanyHeader" />
  <Import src="./footers.gxl" sheet="PageFooter" />
  <Sheet name="Data">...</Sheet>
</Book>
```

**Result:** 3 sheets in output Excel:
1. `CompanyHeader` (from headers.gxl)
2. `PageFooter` (from footers.gxl)
3. `Data` (from main file)rser records the import reference
2. During rendering (usecase layer), the system:
   - Resolves the file path relative to the current template
   - Loads and parses the imported `.gxl` file
   - Finds the specified sheet by name
   - Validates for circular imports and invalid nesting
   -parser enforces structural constraints:

1. **Book-level only**: Import tags must be direct children of `<Book>` tag
2. **No Sheet nesting**: `<Sheet>` tags cannot appear inside other `<Sheet>` tags
3. **No Book nesting**: `<Book>` tags cannot appear inside `<Sheet>` tags
4. **Required attributes**: Both `src` and `sheet` attributes are mandatory

**Parse-time validation:**
- Import inside Sheet → Parse error
- Missing src/sheet attribute → Parse error
- Sheet inside Sheet → Parse error
- Book inside Sheet → Parse error

**Render-time validation:**
Circular imports are detected and prevented:

```
A.gxl → imports B.gxl → imports A.gxl  ❌ Circular import error
```

**Implementation:**
- Import resolution maintains a `visitedFiles` map
- Maximum import depth is limited to 10 levels
- Circular reference triggers immediate error with file path chain
- **Parse time**: Parser rejects `<Sheet>` or `<Book>` tags inside `<Sheet>` elements
- **Import time**: Imported content is validated before expansion

### Circular Import Detection

The system must detect and prevent circular imports:

```
A.gxl imports B.gxl
B.gxl imports A.gxl  ❌ Error: circular import detected
```

Implementation maintains an import stack during resolution. If a file already exists in the stack, a circular import error is raised.

### Import Scope

- Imports **must** occur at the **Book level** (directly under `<Book>` tags)
- Imports **cannot** occur inside `<Sheet>` tags
- Multiple `<Import>` tags are allowed at book level
- Each import creates a new sheet in the final workbook

## Examples

### Basic Import

**common.gxl** (reusable component - separate file):
```xml
<Book name="Common">
  <Sheet name="Headers">
    <Grid>
    | Header | Value |
    | ------| ------|
    </Grid>
  </Sheet>
</Book>
```

**main.gxl** (imports common.gxl - separate file):
```xml
<Book name="Main">in the workbook):
```
Sheet 1: "Headers" (imported from common.gxl)
  | Header | Value |
  | ------| ------|

Sheet 2: "Report" (from main.gxl)
  | Custom | Data |
  | ------ | ---- |
  | Value1 | 100  |
```

Note: Imported sheets are created first, then main file sheets./Sheet>
</Book>
```

**Result after import** (2 sheets created in the workbook):
```
Final Workbook:
  
  Sheet 1: "Report" (from main.gxl)
    Row 1, Col A: (Anchor position)
    Row 5, Col A: Custom | Data
    Row 6, Col A: Value1 | 100
  
  Sheet 2: "Headers" (imported from common.gxl)
    Row 1, Col A: Header | Value
```

### Multiple Imports from Different Sheets

**components.gxl** (library file with multiple sheets):
```xml
<Book name="Components">
  <Sheet name="Header">
    <Grid>
    | Title | Date |
    </Grid>
  </Sheet>
  
  <Sheet name="Footer">
    <Grid>
    | Copyright | Page |
    </Grid>
  </Sheet>
</Book>
```

**Example showing multiple imports at book level**:
```xml
<Book name="Report">
  <Import src="./components.gxl" sheet="Header" />
  <Import src="./metrics.gxl" sheet="Sales" />
  <Import src="./components.gxl" sheet="Footer" />
  
  <Sheet name="Dashboard">
    <Grid>
    | Main Content |
    </Grid>
  </Sheet>
</Book>
```in final workbook (in this order):
1. `Header` (imported from components.gxl)
2. `Sales` (imported from metrics.gxl)
3. `Footer` (imported from components.gxl)
4. `Dashboard` (main sheet)
- "Sales" (imported from metrics.gxl)
- "FLI Usage

### Basic Import Examplebe at book level, not inside Sheet tags
2. **Sheet-level granularity**: Import entire sheets only (not individual nodes)
3. **Required attributes**: Both `src` and `sheet` are mandatory
4. **Creates new sheets**: Each import adds a new sheet to the workbook
5. **No nested structures**: Cannot import Sheet/Book tags inside sheets
6. **File system only**: Only local file paths supported (no URLs)
7. **No circular imports**: A → B → A circular references are detected and rejected
8. **Import depth limit**: Maximum 10 levels of nested imports
9. **Render-time resolution**: Imports resolved during rendering, not parsing
10. **Path resolution**: Relative paths resolved from importing file's directory

**Directory structure:**
```
project/
├── templates/
│   ├── main.gxl          # Main template with imports
│   ├── components/
│   │   ├── headers.gxl   # Reusable header sheets
│   │   └── footers.gxl   # Reusable footer sheets
│   └── shared/
│       └── common.gxl    # Common components
└── output/
    └── report.xlsx       # Generated file
```

**main.gxl:**
```xml
<Book name="Report">
  <Import src="./components/headers.gxl" sheet="CompanyHeader" />
  <Import src="./shared/common.gxl" sheet="DataFormat" />
  <SArchitecture

**Model Layer** (`pkg/model/gxl.go`):
```go
type ImportTag struct {
    Src   string  // File path
    Sheet string  // Sheet name to import
}

type GXL struct {
    BookTag BookTag
    Imports []ImportTag  // Book-level imports
    Sheets  []SheetTag
}
```

**Repository Layer** (`pkg/repository/gxl.go`):
- Parse `<Import>` tags at book level
- Record in `GXL.Imports[]` array
- Validate: reject Import inside Sheet tags

**Usecase Layer** (`pkg/usecase/book.go`):
- Process `gxl.Imports` before rendering main sheets
- For each import:
  1. Resolve file path (relative to template dir)
  2. Load and parse imported .gxl file
  3. Find specified sheet by name
  4. Detect circular imports (visitedFiles map)
  5. Render sheet
  6. Add to workbook as new sheet
- Circular detection with max depth of 10
Potential improvements for future versions:

1. **Import aliasing**: Rename imported sheets
   ```xml
   <Import src="file.gxl" sheet="Data" name="ImportedData" />
   ```

2. **Default sheet**: Import first sheet without specifying name
   ```xml
   <Import src="file.gxl" />
   ```

3. **Import caching**: Cache parsed files to avoid redundant parsing

4. **Wildcard imports**: Import all sheets from a file
   ```xml
   <Import src="file.gxl" sheet="*" />
   ```

5. **Conditional imports**: Import based on conditions
   ```xml
   <If condition="{{ .includeHeaders }}">
     <Import src="headers.gxl" sheet="Header" />
   </If>
   ```

6. **Import from URLs**: Remote template loading with security controls

7. **Partial imports**: Import specific node ranges within sheets inside Sheet tags (enforced by parser)
6. **File must exist**: Import fails if the referenced file does not exist
7. **Sheet must exist**: Import fails if the specified sheet is not found in the imported file
8. **No circular imports**: A → B → A is not allowed (detected at import resolution time)
9. **Path resolution**: Only file system paths supported (no URLs or network resources)
10. **Resolution timing**: Imports are resolved at render time, not parse time
11. **Import depth limit**: Maximum import depth of 10 levels to prevent deeply nested structures

## Error Handling

| Error Condition | Behavior |
|----------------|----------|
| File not found | Return error with clear message indicating missing file |
| Circular import | Return error with import chain showing the cycle |
| Invalid .gxl syntax | Return parsing error from imported file |
| Import depth exceeded | Return error indicating maximum depth reached |
| Permission denied | Return error indicating file access issue |
| Sheet nested in Sheet | Parse error: invalid nesting detected |
| Book nested in Sheet | Parse error: invalid nesting detected |
| Import inside Sheet | Parse error: Import tags must be at book level |
| Missing src attribute | Parse/render error: src attribute is required |
| Missing sheet attribute | Render error: sheet attribute is required |
| Sheet not found | Render error: specified sheet does not exist in imported file |

## Implementation Notes

### Model Layer
- Add `ImportTag` struct to `pkg/model/gxl.go`
- ImportTag contains: `Src string` (file path) and `Sheet string` (sheet name)

### Repository Layer  
- Add `parseImportTag` function to `pkg/repository/gxl.go`
- Parse `<Import>` tags into `ImportTag` structs with both `src` and `sheet` attributes
- No file loading at parse time

### Usecase Layer
- Add import resolution logic to sheet rendering
- Implement circular import detection with visited file set
- **Sheet lookup**: Find and extract nodes from the specified sheet only
- Validate that the specified sheet exists in the imported file
- Recursively load and expand imported files
- Merge imported nodes into parent sheet's node list

### Testing
- Unit tests for parser recognizing `<Import>` tag
- Integration tests for import resolution
- Error case tests: circular imports, missing files, invalid paths
- Performance tests with multiple nested imports

## Future Enhancements

- Import aliasing: `<Import src="file.gxl" sheet="Data" as="namespace" />`
- Default sheet: Allow omitting `sheet` attribute to import the first/default sheet
- Import caching to avoid reparsing same files
- Import from remote URLs with security controls
- Partial sheet import: Import specific node ranges within a sheet
