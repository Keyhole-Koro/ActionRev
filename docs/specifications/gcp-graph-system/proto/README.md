# Proto Draft

## Purpose

このディレクトリには `Go + Connect RPC` 前提の `Protocol Buffers` 叩き台を配置する。

## Layout

- [actionrev/graph/v1/common.proto](/home/unix/ActionRev/docs/specifications/gcp-graph-system/proto/actionrev/graph/v1/common.proto)
- [actionrev/graph/v1/document.proto](/home/unix/ActionRev/docs/specifications/gcp-graph-system/proto/actionrev/graph/v1/document.proto)
- [actionrev/graph/v1/graph.proto](/home/unix/ActionRev/docs/specifications/gcp-graph-system/proto/actionrev/graph/v1/graph.proto)
- [actionrev/graph/v1/node.proto](/home/unix/ActionRev/docs/specifications/gcp-graph-system/proto/actionrev/graph/v1/node.proto)
- [actionrev/graph/v1/job.proto](/home/unix/ActionRev/docs/specifications/gcp-graph-system/proto/actionrev/graph/v1/job.proto)
- [actionrev/graph/v1/tool.proto](/home/unix/ActionRev/docs/specifications/gcp-graph-system/proto/actionrev/graph/v1/tool.proto)

## Design Policy

- package は `actionrev.graph.v1` とする
- service は用途ごとに分離する
- frontend が `React Flow` に直接マップしやすい message 形状を優先する
- 長時間処理は unary RPC で閉じず、job 起動と status 参照に分割する
- 初期段階ではシンプルさを優先し、将来の field 追加を見込んで optional 拡張しやすい構造にする

## Notes

- 実ファイル upload は RPC 本体に載せず、`CreateDocument` で発行した署名付き URL 経由で行う
- `buf` の導入は後続タスクとし、ここでは `.proto` の契約を先に固定する
