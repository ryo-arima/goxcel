# Write Phase

## Overview

Converts sheet models to OOXML format and writes .xlsx files.

## Process

1. **OOXML Generation**: Create XML structures
   - `workbook.xml`: Workbook metadata
   - `sheet1.xml`: Worksheet with cells
   - `styles.xml`: Font and cell styles
   - `sharedStrings.xml`: Shared string table
2. **Cell Encoding**: Convert cells by type
   - Numbers: `<c t="n"><v>123</v></c>`
   - Strings: `<c t="s"><v>0</v></c>` (index to sharedStrings)
   - Formulas: `<c><f>=SUM(A1:A10)</f></c>`
   - Booleans: `<c t="b"><v>1</v></c>`
3. **Style Application**: Apply bold/italic via styleId
4. **ZIP Packaging**: Bundle XML files into .xlsx
5. **File Writing**: Save to disk

## Input/Output

**Input**: `model.Sheet` with cells  
**Output**: `.xlsx` file (ZIP archive with OOXML)

## Key Functions

- `repository.WriteXlsx()`: Main writer
- `writeSheet()`: Generate sheet XML
- `writeStyles()`: Generate styles XML
- `createXMLCell()`: Cell-specific XML generation

## OOXML Structure

```
output.xlsx (ZIP)
├── [Content_Types].xml
├── _rels/.rels
├── xl/
│   ├── workbook.xml
│   ├── styles.xml
│   ├── sharedStrings.xml
│   ├── worksheets/sheet1.xml
│   └── _rels/workbook.xml.rels
```
