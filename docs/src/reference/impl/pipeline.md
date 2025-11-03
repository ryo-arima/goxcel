# Rendering Pipeline

## Three-Phase Process

### 1. Parse Phase
**Input**: .gxl file  
**Output**: GXL AST  
**Operations**: XML parsing, tag recognition, syntax validation

### 2. Render Phase
**Input**: GXL AST + JSON data  
**Output**: Sheet model with cells  
**Operations**: Expression evaluation, loop expansion, grid positioning

### 3. Write Phase
**Input**: Sheet model  
**Output**: .xlsx file  
**Operations**: OOXML generation, zip packaging, file writing

## Data Flow

```
.gxl → Parser → AST → Renderer → Sheet → Writer → .xlsx
              ↑                    ↑
         Repository          Usecase
```

See individual pipeline stages:
- [Parse](./pipeline/parse.md)
- [Render](./pipeline/render.md)
- [Write](./pipeline/write.md)
