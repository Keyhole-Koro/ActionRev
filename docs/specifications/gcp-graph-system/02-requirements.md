# 02. Requirements

## Functional Requirements

### File Upload

- ユーザーは Web UI からファイルをアップロードできること
- アップロードされたファイルは `Cloud Storage` に保存されること
- 保存時に一意な `document_id` を採番すること
- ファイル名、ファイルサイズ、MIME type、保存先 URI、作成日時を記録すること

### Text Extraction

- PDF、Markdown、TXT、CSV からテキストを抽出できること
- 抽出結果を意味的またはサイズ上の単位で chunk に分割できること
- 各 chunk に `chunk_id` と `chunk_index` を付与すること
- 可能な範囲でページ番号やオフセットなどの出典位置情報を保持すること

### Node and Edge Extraction by Gemini

- chunk または chunk 群を Gemini に入力できること
- Gemini の返却形式は JSON であること
- 最低限以下の概念を抽出対象とすること
- `abstract` ノード
- `concrete` ノード
- ノード説明
- ノード間エッジ
- 出典 chunk
- 同一 document 内で重複する概念は統合可能であること

### Graph Persistence

- 抽出した document、chunk、node、edge を `BigQuery` に保存できること
- エッジには少なくとも `source_node_id`、`target_node_id`、`edge_type` を保持すること
- ノードには少なくとも `label`、`type`、`description` を保持すること
- ノードおよびエッジは出典 chunk を追跡できること

### Graph Retrieval

- document 単位でノード一覧を取得できること
- document 単位でエッジ一覧を取得できること
- node type によるフィルタができること
- 将来的に件数制限、深さ制限、relation type 指定ができること

### Visualization

- ノードエッジ形式でグラフを表示できること
- 抽象ノードと具体ノードを見た目で区別できること
- ノードクリックで詳細情報を表示できること
- 出典 chunk を表示できること

## Non-Functional Requirements

### Scalability

- 初期は MB 級から数百 MB 級を対象とする
- 将来的に GB 級ファイルや大量投入に対応できる構成へ拡張可能であること

### Availability and Operability

- 初期構成はサーバレス中心とし、常時運用管理を極小化すること
- 処理失敗時に再処理可能であること
- ログとステータスにより処理進行が把握できること

### Data Quality

- LLM 出力揺れに対応し、再抽出および再統合が可能であること
- 出典追跡により生成結果の監査性を確保すること

### Security

- ストレージバケットは非公開とすること
- API からのみファイルにアクセスできること
- 認証認可は将来的に導入可能な構造とすること
- 個人情報や機密情報を含む場合の保存ポリシーを定義可能であること

## Input and Output Specification

### Supported Input Formats

- PDF
- Markdown
- TXT
- CSV

### Future Input Formats

- DOCX
- HTML
- JSON

### Output Format

- フロントとバックエンド間の主たる契約は `Protocol Buffers` とする
- ブラウザとの同期通信には `Connect RPC` を使用する
- グラフ表示用レスポンスには `nodes` と `edges` を含めること

```json
{
  "document_id": "doc_001",
  "nodes": [
    {
      "id": "n1",
      "label": "販売戦略",
      "type": "abstract",
      "description": "販売拡大のための上位方針"
    }
  ],
  "edges": [
    {
      "id": "e1",
      "source": "n1",
      "target": "n2",
      "type": "abstract_to_concrete"
    }
  ]
}
```

## Acceptance Criteria

- ユーザーがファイルをアップロードできること
- upload 後に document が記録されること
- テキスト抽出と chunk 保存が行われること
- Gemini からノード・エッジ JSON を取得できること
- BigQuery に document/chunk/node/edge が保存されること
- フロントで抽象/具体ノードを含むグラフが表示されること
- ノード詳細から出典 chunk を参照できること
