# Tree-Native Proto 草案

日付: 2026-04-01

## 目的

このドキュメントは、Synthify を tree-native backend に刷新するための proto / API shape の草案をまとめる。

前提:

- 主表示は graph ではなく tree
- backend は tree に近い shape を返す
- frontend は大きな `Graph -> Tree` 変換責務を持たない
- AI note / action request は first-class に扱う

## 基本方針

新しい API は、旧 `GraphService` の延長として考えない。

代わりに、次の 3 レイヤーで分ける。

1. workspace tree
2. paper node
3. action / note

旧 graph API は移行期間だけ残し、最終的には主役から外す想定。

## 中心 message

### WorkspaceTree

workspace 全体の tree 状態を返すトップレベル message。

```proto
message WorkspaceTree {
  string workspace_id = 1;
  string root_node_id = 2;
  repeated PaperNode nodes = 3;
  repeated SourceDocument documents = 4;
}
```

意図:

- tree 全体を一括取得できる
- frontend はそのまま renderer に流しやすい
- document provenance も同時に返せる

### PaperNode

tree 上に表示される最小単位。

```proto
message PaperNode {
  string paper_node_id = 1;
  string workspace_id = 2;
  string parent_id = 3;
  repeated string child_ids = 4;

  string title = 5;
  string description = 6;
  string content = 7;

  PaperNodeCategory category = 8;
  PaperNodeScope scope = 9;

  uint32 display_order = 10;
  PaperNodeStatus status = 11;

  repeated string source_document_ids = 12;
  repeated string source_chunk_ids = 13;

  PaperMeta meta = 14;
  repeated PaperBlock blocks = 15;
}
```

### PaperMeta

renderer や app が参照する補助情報。

```proto
message PaperMeta {
  repeated string badges = 1;
  bool is_new = 2;
  bool is_loading = 3;
  repeated NodeRelationSummary relations = 4;
}
```

意図:

- `blocks` ほど重くない補助情報を持つ
- renderer 側の装飾や小さな state を受け渡す

### PaperBlock

node 内に出す view model。

```proto
message PaperBlock {
  oneof kind {
    NoteBlock note = 1;
    MetricBlock metric = 2;
    RelationsBlock relations = 3;
    DocumentsBlock documents = 4;
    WarningBlock warning = 5;
    MiniGraphBlock mini_graph = 6;
  }
}
```

意図:

- backend は「何を見せるか」を返す
- frontend は block type ごとに描画する

## Note / Action 系 message

### PaperNote

AI から人間への付箋。

```proto
message PaperNote {
  string note_id = 1;
  string workspace_id = 2;
  string paper_node_id = 3;

  PaperNoteKind kind = 4;
  string title = 5;
  string body = 6;
  NotePriority priority = 7;

  google.protobuf.Timestamp created_at = 8;
  google.protobuf.Timestamp updated_at = 9;
}
```

### ActionRequest

状態遷移を持つ human action request。

```proto
message ActionRequest {
  string action_request_id = 1;
  string workspace_id = 2;
  string paper_node_id = 3;
  string note_id = 4;

  ActionRequestType type = 5;
  string title = 6;
  string body = 7;
  ActionRequestPriority priority = 8;
  ActionRequestStatus status = 9;

  string requested_by = 10;
  string assigned_to = 11;

  google.protobuf.Timestamp created_at = 12;
  google.protobuf.Timestamp resolved_at = 13;
}
```

## Relation 系 message

tree が主表示でも、relation は残す。

ただし edge を主役にはしない。

```proto
message NodeRelationSummary {
  string target_node_id = 1;
  RelationType relation_type = 2;
  string label = 3;
}
```

必要なら詳細取得用に別 API を切る。

## Source Document message

document provenance は引き続き重要。

```proto
message SourceDocument {
  string document_id = 1;
  string workspace_id = 2;
  string filename = 3;
  string mime_type = 4;
  int64 file_size = 5;
  DocumentLifecycleState status = 6;
}
```

## enum 草案

### PaperNodeCategory

```proto
enum PaperNodeCategory {
  PAPER_NODE_CATEGORY_UNSPECIFIED = 0;
  PAPER_NODE_CATEGORY_CONCEPT = 1;
  PAPER_NODE_CATEGORY_CLAIM = 2;
  PAPER_NODE_CATEGORY_EVIDENCE = 3;
  PAPER_NODE_CATEGORY_ENTITY = 4;
  PAPER_NODE_CATEGORY_METRIC = 5;
  PAPER_NODE_CATEGORY_ACTION = 6;
}
```

### PaperNodeScope

```proto
enum PaperNodeScope {
  PAPER_NODE_SCOPE_UNSPECIFIED = 0;
  PAPER_NODE_SCOPE_WORKSPACE = 1;
  PAPER_NODE_SCOPE_DOCUMENT = 2;
  PAPER_NODE_SCOPE_CANONICAL = 3;
}
```

### PaperNodeStatus

```proto
enum PaperNodeStatus {
  PAPER_NODE_STATUS_UNSPECIFIED = 0;
  PAPER_NODE_STATUS_READY = 1;
  PAPER_NODE_STATUS_DRAFT = 2;
  PAPER_NODE_STATUS_PENDING = 3;
  PAPER_NODE_STATUS_ARCHIVED = 4;
}
```

### PaperNoteKind

```proto
enum PaperNoteKind {
  PAPER_NOTE_KIND_UNSPECIFIED = 0;
  PAPER_NOTE_KIND_CONFIRMATION = 1;
  PAPER_NOTE_KIND_DECISION = 2;
  PAPER_NOTE_KIND_MISSING_EVIDENCE = 3;
  PAPER_NOTE_KIND_REVIEW = 4;
  PAPER_NOTE_KIND_NEXT_ACTION = 5;
}
```

### ActionRequestStatus

```proto
enum ActionRequestStatus {
  ACTION_REQUEST_STATUS_UNSPECIFIED = 0;
  ACTION_REQUEST_STATUS_OPEN = 1;
  ACTION_REQUEST_STATUS_ACKNOWLEDGED = 2;
  ACTION_REQUEST_STATUS_RESOLVED = 3;
  ACTION_REQUEST_STATUS_DISMISSED = 4;
}
```

## RPC 草案

### WorkspaceTreeService

```proto
service WorkspaceTreeService {
  rpc GetWorkspaceTree(GetWorkspaceTreeRequest) returns (GetWorkspaceTreeResponse);
  rpc ListPaperNodeChildren(ListPaperNodeChildrenRequest) returns (ListPaperNodeChildrenResponse);
  rpc GetPaperNode(GetPaperNodeRequest) returns (PaperNode);
  rpc CreatePaperNode(CreatePaperNodeRequest) returns (PaperNode);
  rpc UpdatePaperNode(UpdatePaperNodeRequest) returns (PaperNode);
  rpc ReorderPaperNode(ReorderPaperNodeRequest) returns (ReorderPaperNodeResponse);
}
```

### NotesService

```proto
service PaperNoteService {
  rpc ListNodeNotes(ListNodeNotesRequest) returns (ListNodeNotesResponse);
  rpc CreateNodeNote(CreateNodeNoteRequest) returns (PaperNote);
}
```

### ActionRequestService

```proto
service ActionRequestService {
  rpc ListNodeActionRequests(ListNodeActionRequestsRequest) returns (ListNodeActionRequestsResponse);
  rpc ResolveActionRequest(ResolveActionRequestRequest) returns (ActionRequest);
  rpc DismissActionRequest(DismissActionRequestRequest) returns (ActionRequest);
}
```

## `GetWorkspaceTree` の shape

最初の MVP ではこれで十分。

```proto
message GetWorkspaceTreeRequest {
  string workspace_id = 1;
}

message GetWorkspaceTreeResponse {
  Workspace workspace = 1;
  WorkspaceTree tree = 2;
}
```

## `ReorderPaperNode` の shape

drag-and-drop に対応するため、並び替えは first-class にする。

```proto
message ReorderPaperNodeRequest {
  string workspace_id = 1;
  string paper_node_id = 2;
  string new_parent_id = 3;
  string insert_before_id = 4;
}

message ReorderPaperNodeResponse {
  PaperNode node = 1;
  repeated string sibling_ids = 2;
}
```

## `ResolveActionRequest` の shape

```proto
message ResolveActionRequestRequest {
  string workspace_id = 1;
  string action_request_id = 2;
  string resolution_note = 3;
}
```

## frontend から見た利点

この shape にすると frontend は次をしなくてよい。

- graph edge を見て hierarchy を推論する
- document/canonical/workspace の混在を苦労して tree に直す
- sticky note を別の補助 API から無理に合成する

つまり renderer のための大きな adapter が不要になる。

## 実装上の次の判断

proto を本当に切り始める前に、次の 3 点を決める必要がある。

1. `PaperNode` に `meta` と `blocks` を両方持たせるか
2. `PaperNote` と `ActionRequest` を分けるか、統合するか
3. relation を summary だけ返すか、詳細 graph 用 API を残すか

## 現時点のおすすめ

現時点では次がおすすめ。

1. `meta` と `blocks` は分ける
2. `PaperNote` と `ActionRequest` は分ける
3. relation は summary を主にし、詳細 graph API は後方互換的に残してもよい
