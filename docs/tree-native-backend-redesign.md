# Tree-Native Backend 再設計メモ

日付: 2026-04-01

## 目的

Synthify の次の UI は graph viewer ではなく、`paper-in-paper` をベースにした tree / paper workspace になる想定である。

そのため、frontend 側で

- `Graph -> Tree`

の変換レイヤーを持つのではなく、backend / DB / API の時点で tree-native な設計へ刷新する。

つまり、これからの主語は

- `graph`
- `edge`

ではなく

- `workspace tree`
- `paper node`
- `AI note`
- `action request`

である。

## この方針で捨てるもの

次の発想は捨てる。

- 旧 graph schema を frontend で tree に変換する
- graph node / edge を中心に UI を設計する
- backend が relation graph を返し、frontend が苦労して tree に落とす

これは移行コストが高いわりに、最終形のドメインに合っていない。

## 新しい中心ドメイン

これからの中心エンティティは次の 4 つ。

1. `workspace`
2. `paper_node`
3. `paper_note`
4. `action_request`

必要に応じてこれに

- `source_document`
- `source_chunk`
- `related_reference`

をぶら下げる。

## Paper Node とは何か

`paper_node` は、tree 上に表示される最小単位。

例:

- 概念
- 主張
- 根拠
- 指標
- アクション候補

これは旧 graph の `Node` に近いが、重要なのは

- 最初から親子関係を持つ
- UI に出すことが前提
- 補助情報を `meta` / `blocks` として持てる

という点。

## AI Note とは何か

`paper_note` は、単なる説明ではない。

これは主に

- AI が人間に対して何か確認・判断・行動を求めるための付箋

として扱う。

例:

- この論点を確認してください
- この分岐を選択してください
- この根拠を追加してください
- この矛盾をレビューしてください
- 次のアクションを選んでください

したがって note は、本文の補足ではなく

- `human action request`

に近い意味を持つ。

## Action Request とは何か

`action_request` は、AI から人間への要求を状態付きで持つエンティティ。

最低限ほしい属性:

- `id`
- `workspace_id`
- `paper_node_id`
- `type`
- `title`
- `body`
- `priority`
- `status`
- `requested_by`
- `created_at`
- `resolved_at`

`paper_note` と違って、こちらは workflow を持つ。

例:

- `open`
- `acknowledged`
- `resolved`
- `dismissed`

UI では sticky note として見せるが、backend / DB 上ではちゃんと状態遷移できる形にする。

## 想定する API モデル

旧:

- `GetGraph`
- `ExpandNeighbors`

新:

- `GetWorkspaceTree`
- `GetPaperNode`
- `ListPaperNodeChildren`
- `ListNodeNotes`
- `ListNodeActionRequests`
- `ResolveActionRequest`
- `CreatePaperNode`
- `ReorderPaperNode`

必要なら後で relation 系 API を追加する。

## `GetWorkspaceTree` が返すもの

最初の段階では、少なくとも次を返せる必要がある。

- `workspace`
- `root_node_id`
- `paper_nodes`
- 各 node の `meta`
- 各 node の `blocks`
- note / action request の概要

ここでいう `blocks` は、node 内に出す view model。

例:

- note
- metric
- relation summary
- source documents
- warning
- mini graph

重要なのは、

- backend は「何を見せるか」を返す
- frontend は「どう見せるか」を決める

という分担にすること。

## DB の基本設計案

### `workspaces`

- `id`
- `name`
- `created_at`
- `updated_at`

### `paper_nodes`

- `id`
- `workspace_id`
- `parent_id`
- `title`
- `description`
- `content`
- `category`
- `scope`
- `display_order`
- `status`
- `meta_json`
- `created_at`
- `updated_at`

ここでは `parent_id` を first-class に持つ。

つまり tree は DB レベルで表現される。

### `paper_notes`

- `id`
- `workspace_id`
- `paper_node_id`
- `kind`
- `title`
- `body`
- `priority`
- `meta_json`
- `created_at`
- `updated_at`

### `action_requests`

- `id`
- `workspace_id`
- `paper_node_id`
- `note_id` nullable
- `type`
- `title`
- `body`
- `priority`
- `status`
- `requested_by`
- `assigned_to` nullable
- `created_at`
- `resolved_at` nullable

### `source_documents`

- `id`
- `workspace_id`
- `filename`
- `mime_type`
- `file_size`
- `status`
- `created_at`
- `updated_at`

### `paper_node_document_refs`

- `paper_node_id`
- `document_id`
- `chunk_id` nullable
- `role`

これで node と document の provenance を持つ。

### `paper_node_relations`

tree の親子ではない relation は、ここで持つ。

- `id`
- `workspace_id`
- `source_node_id`
- `target_node_id`
- `relation_type`
- `weight` nullable
- `meta_json`

これは UI 主表示の主役ではない。
ただし relation summary を出す元データとして必要。

## UI との対応

frontend 側では `paper-in-paper` が tree を描く。

その時点で backend はすでに

- root
- parent/child
- note
- action request
- source document summary

を返しているので、frontend には大きな変換責務を持たせない。

frontend がやるのは:

- 開閉状態
- 選択状態
- hover / highlight
- 表示のスタイリング
- button click からの API 呼び出し

## 開閉状態は backend に持たせない

重要:

- backend は `node を開け / 閉じろ` を返さない

理由:

- それは表示状態であり UI ローカルの責務だから
- 勝手に開閉が起こると UX が悪いから

backend が返すべきなのは、

- 開いた時に見せるための材料

だけ。

## `meta` と `blocks`

node に直接たくさんの列を増やしすぎると硬くなるので、柔らかい情報は `meta` と `blocks` に寄せる。

おすすめの考え方:

- `paper_node` の主要属性
  - 構造・検索・並び替えに必要なもの
- `meta`
  - renderer が見る追加情報
- `blocks`
  - node 内で表示する小さい view model

たとえば `blocks` はこういう union を想定できる。

- `note`
- `metric`
- `relations`
- `documents`
- `warning`
- `mini_graph`

## mini graph について

mini graph は可能だが、最初から主役にしない方がよい。

理由:

- nested interaction が重くなる
- tree の読みやすさを壊しやすい
- まず note / metric / relations の方が価値が高い

よって優先度は低い。

## 移行の意味

この刷新をすると、実質的に次のようになる。

旧システム:

- relation graph が中心
- frontend で整形して見せる

新システム:

- tree-native workspace が中心
- backend が表示に近い構造を返す
- relation は補助情報として扱う

## 最初に決めるべきこと

設計着手前に固定したい点:

1. `paper_node` の最小必須属性
2. `paper_note` と `action_request` の境界
3. `blocks` の union 種別
4. `paper_node_relations` をどこまで保存するか
5. `display_order` を sibling 単位でどう扱うか

## 実装の最初の一歩

次にやるべきなのはコードではなく、まずこれ。

1. proto / API の新 shape を文章で切る
2. DB table / entity の草案を切る
3. `paper_note` と `action_request` の状態遷移を決める

proto / API shape の草案は次を参照:

- `docs/tree-native-proto-draft.md`

その後で:

4. backend mock を新ドメインに寄せる
5. frontend を新 tree API 前提に寄せる
