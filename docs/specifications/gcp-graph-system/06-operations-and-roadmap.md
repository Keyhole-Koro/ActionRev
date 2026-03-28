# 06. Operations and Roadmap

## Error Handling

### Failure Types

- アップロード失敗
- テキスト抽出失敗
- Gemini 呼び出し失敗
- 正規化ツール生成失敗
- サンドボックス実行失敗
- JSON parse 失敗
- BigQuery 書き込み失敗
- RPC ハンドラ失敗
- 非同期ジョブ起動失敗

### Failure Policy

- `documents.status` を `failed` に更新する
- 失敗理由をログに記録する
- 再処理可能な設計とする
- `StartProcessing` は upload 完了済みかつ `documents.status=uploaded` の document のみ受け付ける
- upload 未完了や不正な状態遷移はジョブを起動せず、同期エラーとして返す

## Monitoring and Logging

### Initial Monitoring

- `Cloud Logging` に API 実行ログを出力する
- document 単位で処理開始、成功、失敗を記録する
- RPC メソッド単位でレイテンシと失敗率を記録する
- job 単位で開始、完了、失敗を記録する
- ツール生成、承認、dry-run、本実行のイベントを記録する

### Future Monitoring

- 処理時間監視
- エラー率監視
- Gemini 失敗率監視
- コスト監視

## Security and Access Control

### Initial Policy

- Cloud Storage バケットは非公開
- API 経由でのみファイルアクセス
- サービス間アクセスは GCP IAM により制御

### Future Policy

- ユーザー認証導入
- ドキュメント単位のアクセス制御
- 監査ログの強化

## Extension Roadmap

### Near-Term Extensions

- `Cloud Tasks` による非同期ジョブ化
- 大きなファイルの chunk 並列処理
- ノード重複統合ルールの強化
- フロントでの検索・フィルタ・折りたたみ
- `buf` を使った proto 管理とコード生成パイプライン
- 正規化ツールのレビュー、承認、再利用フロー

### Mid-Term Extensions

- `Pub/Sub` による大量投入
- `Cloud Run Jobs` による再処理バッチ
- `Memorystore` によるレスポンスキャッシュ
- `BigQuery` の scheduled query による夜間再集計

### Advanced Extensions

- `Vertex AI Embeddings` または類似技術による類似ノード探索
- 複数 document 横断の概念統合
- `Spanner Graph` への移行または併用
- 高度なグラフ探索 API の追加

## Open Issues

- PDF の抽出品質をどのライブラリで担保するか
- 正規化ツールの approval workflow をどこまで厳密にするか
- Gemini の出力スキーマとリトライ戦略の詳細
- ノード統合ルールの厳密度
- フロントの可視化ライブラリ選定
- 認証導入のタイミング
- proto の package 分割方針
