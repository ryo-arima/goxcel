# GXL Specification (Markdown + Vue-like)

Version: 0.1 (Draft)

## Goals
- 表形式（シート/行/列）の構造をそのまま表現できるテンプレート
- Markdown互換の読みやすさを維持しつつ、構造化された宣言（カスタムタグ）を使用
- 値の埋め込みや制御構文（for/if）、画像/図形/グラフ/ピボット等のコンポーネントをサポート

## File
- Extension: `.gxl`
- Encoding: UTF-8 (BOM なし)

## High-level Structure
- Markdown文書内に、以下のカスタムタグを使用:
  - `<Sheet name="..."> ... </Sheet>`: シート定義（複数可）
  - `<Anchor ref="A1" />`: アンカー（起点セルの絶対参照）
  - `<Merge range="A1:C1" />`: セル結合
  - `<Style selector="A1" name="title" />`: スタイル付与（実際のスタイルはレンダラ/Writer側が解決）
  - `<Grid> ... </Grid>`: グリッド行をパイプ区切り（`|`）で記述
  - `<For each="item in items"> ... </For>`: 繰り返し（下方向に行を増やす）
  - `<If cond="..."> ... <Else> ... </Else> </If>`: 条件分岐
  - `<Image ... />`, `<Shape ... />`, `<Chart ... />`, `<Pivot ... />`: コンポーネント
- 通常のMarkdown記法はそのまま使用可能（ただし `<Grid>` 内の表記のみセルとして解釈）

## Expressions
- 値の埋め込み: `{{ expr }}`
  - 例: `Hello {{ user.name }}`
- 属性値内でも使用可能: `<Chart dataRange="A1:C{{ endRow }}" />`
- 型の扱い: レンダラで推定（数値/日付/真偽/文字列）。ヒントが必要な場合はセルに `type` を指定（将来拡張）。

## Tags and Attributes
### Sheet
- `<Sheet name="Invoice"> ... </Sheet>`
- name: 文字列（重複不可）。重複があればパーサ/バリデータでエラー。

### Anchor
- `<Anchor ref="A1" />`
- ref: 絶対セル参照（A1記法）。以後の `<Grid>` はこのアンカーからの相対位置にマップ。

### Grid
- `<Grid>` ブロック内は、行単位で `|` 区切りのセル列を記述。
- 例:
  ```
  <Grid>
    | Name | Qty | Price |
    | {{ item.name }} | {{ item.qty }} | {{ item.price }} |
  </Grid>
  ```
- それぞれの行はアンカーの行からの相対行に、セルは左からの相対列に配置。
- `=` で始まるセルは Excel 式として扱う（`=SUM(A1:A10)`）。式内に `{{ }}` を含めてよい。

### Merge
- `<Merge range="A1:C1" />` … 絶対参照の結合範囲。

### Style
- `<Style selector="A1" name="title" />`
- `<Style selector="A1:C3" id="title" class="heading,emphasis" fontFamily="Inter" fontSize="14" color="#333333" fillColor="#FFF8E1" bold italic />`
- selector: 絶対参照や将来的に相対指定を想定（`[rowOffset,colOffset]` など）。
- name: 任意のスタイルキー。実体はWriter側で解決。
- id / class を指定可能。複数 class はカンマ区切り。
- フォント/色などのプロパティをインライン指定可能：
  - `fontFamily`, `fontSize`, `bold`, `italic`, `underline`
  - `color`（文字色, `#RRGGBB`）, `fillColor`（背景色, `#RRGGBB`）
  - `hAlign`（left|center|right）, `vAlign`（top|middle|bottom）

### For
- `<For each="item in items"> ... </For>`
- 構文: `each="<var> in <path>"`
- 各反復で `<Grid>` の行数分だけカーソルは下方向に進む。
- 組み込み `loop` 変数（レンダリング時に提供）:
  - `loop.index` (0-based), `loop.number` (1-based)
  - `loop.startRow`, `loop.endRow`（絶対行番号; 最終的に分かるため、後続参照で使用可能）

### If
- `<If cond="..."> ... <Else> ... </Else> </If>`
- cond は真偽に評価される式（空/ゼロ/falseは偽と判定）。

### Image
- `<Image ref="B3" src="assets/logo.png" width="120" height="60" />`
- ref: アンカーセル
- src: 画像パスまたはリソースキー
- width/height: px 単位（省略可）。
- その他オプションは `options`（key=valueのCSV などを検討; 将来拡張）。

### Shape
- `<Shape ref="D3" kind="rectangle" text="Hello" width="100" height="40" style="banner" />`
- kind: rectangle|rounded|ellipse|arrow|...（Writer依存）

### Chart
- `<Chart ref="F3" type="column" dataRange="A9:C20" title="Sales" width="420" height="240" />`
- type: column|bar|line|pie|scatter|...
- dataRange: A1記法の範囲。`{{ }}` を含めてよい。

### Pivot
- `<Pivot ref="F15" sourceRange="A9:C200" rows="Name" columns="Category" values="SUM:Price" filters="Year" />`
- rows/columns/values/filters: いずれもカンマ区切り複数可。
- valuesは集計関数:フィールド（例: `SUM:Price`, `COUNT:ID`）。

## AST Mapping
- `<Sheet>` → `model.SheetTemplate`
- `<Anchor>` → `model.AnchorNode`
- `<Grid>` の各行 → `model.GridRowNode`（行ごとに `[]CellTemplate`）
- `<Merge>` → `model.MergeNode`
- `<Style>` → `model.StyleNode`（`ID`, `Class`, `Props` を含む）
- `<For>` → `model.ForNode`（Body は入れ子の Node 列）
- `<If>` → `model.IfNode`
- `<Image>` → `model.ImageNode`
- `<Shape>` → `model.ShapeNode`
- `<Chart>` → `model.ChartNode`
- `<Pivot>` → `model.PivotNode`

## Validation Rules (抜粋)
- `<Sheet>` は少なくとも1つ必要。
- `<Grid>` は直前に有効な `<Anchor>` が必要（未指定なら暗黙に `A1` を初期値としてもよい）。
- セル参照/範囲の妥当性チェック（A1記法）。
- 未知のタグ・不足属性はエラー。
- シート名重複はエラー。

## Rendering Semantics (抜粋)
- アンカーは現在位置を絶対リセット。`<Grid>` はアンカーからの相対セルにマップ。
- `<For>` は Body 内の `<Grid>` の総行数分、反復毎に下へ進む。
- `{{ }}` の評価には提供されたデータコンテキストを使用。属性内の `{{ }}` も同様。

## Examples
- `.examples/markdown_vue.gxl` を参照。
- 旧DSL（行頭 `%%` ディレクティブ/パイプ行）は後方互換の検討対象（現状はREADMEに記載）。
