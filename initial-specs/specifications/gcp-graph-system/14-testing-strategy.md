# 14. Testing Strategy

## Overview

本章ではコードレベルの品質保証方針を定義する。LLM 出力の品質評価は [13-evaluation-data.md](13-evaluation-data.md) を参照。

---

## テストの分類

### ユニットテスト

各コンポーネントを独立して検証する。

| 対象 | テスト内容 |
| --- | --- |
| テキスト抽出 | PDF / Markdown / TXT / CSV からテキストが正しく取り出せるか |
| セマンティックチャンク分割 | Gemini モックを使い、チャンク構造が正しく保存されるか |
| JSON パース・正規化 | 不正 JSON・欠損フィールド・型ミスへの耐性 |
| ノード統合ロジック | 重複判定・canonical 化のロジック |
| zip 展開 | 対応ファイルのみ処理されるか、ネスト zip は除外されるか |
| source_filename 付与 | zip 内ファイルの chunk に正しく source_filename が付くか |

### 結合テスト

複数コンポーネントを組み合わせて検証する。

| 対象 | テスト内容 |
| --- | --- |
| アップロード → GCS 保存 | ファイルが正しいパスに保存され、document レコードが作成されるか |
| 抽出 → BigQuery 書き込み | node / edge / chunk が正しいスキーマで保存されるか |
| BigQuery → Spanner Graph 同期 | canonical node / edge が探索用 graph に正しく反映されるか |
| GetGraph RPC | source_filename_filter / node_category_filter が正しく動くか |
| ExpandNeighbors / FindPaths RPC | 近傍展開と経路検索が期待する subgraph と path evidence ref を返すか |
| GetGraphEntityDetail RPC | 参照対象に応じた詳細・隣接エッジ・出典チャンク・supporting edge が返るか |
| ToolService | 正規化ツールの dry-run / 本実行・差分保存の一連フロー |

### E2E テスト

実際のファイルを投入してパイプライン全体を検証する。Gemini 呼び出しはモックを使う。

| シナリオ | 検証内容 |
| --- | --- |
| PDF 単体アップロード | アップロードから graph 取得まで一気通貫で動くか |
| zip アップロード | zip 内ファイルが source_filename 付きで処理されるか |
| 対話的探索 | graph 表示後に近傍展開と経路検索が行えるか |
| 処理失敗シナリオ | Gemini 失敗時に status が `failed` になり再処理可能か |
| 再処理（force_reprocess） | 既存ノード・エッジが上書きされるか |

---

## LLM モック戦略

Gemini 呼び出しはテストで実行しない。以下のいずれかでモックする。

**固定レスポンス（ユニット・結合テスト）**
- 事前に用意した JSON を返すスタブを使う
- ゴールドドキュメント（[13-evaluation-data.md](13-evaluation-data.md)）の期待値と同じ形式

**録画・再生（E2E テスト）**
- 実際の Gemini レスポンスを一度録画して保存する
- テスト時は録画済みレスポンスを再生する
- プロンプト変更時のみ録画を更新する

---

## テスト実行方針

### 通常実行

```
日常開発時:
  - ユニットテスト（全件）
  - 結合テスト（モック使用）

リリース前:
  - E2E テスト（モック使用）
```

### 定期実行

```
週次:
  - 評価データセットを使った抽出品質の自動計測（13-evaluation-data.md 参照）
  - 結果を BigQuery に記録してトレンドを追う

変更時:
  - プロンプト、JSON Schema、抽出戦略、canonical 化閾値を変更した場合は手動で評価データセットを実行する
```

### BigQuery 結合テスト方針

- BigQuery エミュレータ前提では設計しない
- PR では BigQuery をモックし、I/O 契約とクエリ組み立てを中心に検証する
- `main` マージ後と nightly では専用 GCP テストプロジェクトに対して実 BigQuery 結合テストを実行する
- テストプロジェクトの dataset は実行ごとにプレフィックス付きで作成し、終了時に削除する
- 認証情報は CI の secret と Workload Identity で供給し、ローカル開発者に本番権限を要求しない

### Spanner Graph 結合テスト方針

- `Spanner Graph` はエミュレータまたは専用テストインスタンスのいずれかで再現可能にする
- PR では探索 RPC のロジックをモックまたはローカル graph fixture で検証する
- `main` マージ後と定期 CI では専用 GCP テストプロジェクト上の `Spanner Graph` に対して同期・探索の結合テストを実行する
- `BigQuery` から `Spanner Graph` への同期整合性は E2E で最低 1 シナリオ確認する

### E2E 環境方針

- 初期は常設 staging を持たず、CI から起動する一時環境または専用 test project を使う
- E2E は main / release 前に限定し、PR ごとには実行しない
- Cloud Run staging を常設するのは、運用コストに見合う不具合が観測された時点で再検討する
- E2E では `BigQuery` と `Spanner Graph` の両方を含む二層構成の整合性を検証する

---

## Proto 契約テスト

Connect RPC の proto 変更がフロントエンドを壊さないことを確認する。

- `buf breaking` で後方互換性を検証する
- フィールド追加は許容、フィールド削除・型変更は契約テストで検出してブロックする

---

## テストデータの管理

- fixtures/ にサンプルファイル（PDF / MD / CSV / zip）を配置する
- Gemini モックレスポンスは `testdata/gemini_responses/` に保存する
- 評価用 gold データは `eval/` に保存する（[13-evaluation-data.md](13-evaluation-data.md) 参照）

---

## Open Issues

- Go のテストフレームワーク選定（標準 `testing` パッケージのみか、testify を使うか）
- buf の導入タイミング（proto 管理パイプラインの自動化）
