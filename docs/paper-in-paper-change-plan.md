# Paper-in-Paper 変更計画

日付: 2026-04-01

## 目的

このドキュメントは、移行作業を次の 2 つに分けて整理するためのものです。

- `vendor/paper-in-paper` 側で変更すること
- Synthify 側で変更すること

後方互換は不要です。

また、この計画は `paper-in-paper` を既存ライブラリとして無理に適応させるのではなく、Synthify 向けの tree renderer として直接再設計していく前提で書いています。

## 現在の方針

`paper-in-paper` は、固定された外部ライブラリとして扱わない。

代わりに、次のようなものとして扱う。

- internal submodule
- 再利用可能な土台
- Synthify 向けに直接設計変更してよい実装

つまり、優先すべきなのは

- `adapter を頑張ること`

ではなく

- `renderer とデータモデルを Synthify 向けに再定義すること`

です。

## Paper-in-Paper が持つべき責務

`paper-in-paper` は、最終的に次の責務を持つべきです。

- 階層 tree の描画
- node の展開 / 折りたたみ
- primary child / focused branch の制御
- drag / reorder の UI と状態遷移
- tree の motion と layout transition
- node-level presentation primitive
- tree-level の selection / highlight / dim 制御

つまり、Synthify の workspace tree を描くための専用 visual engine になるべきです。

## Synthify が持つべき責務

Synthify 側は、引き続き次を持ちます。

- graph の取得
- workspace / document の取得
- upload
- rename
- search query の状態
- filter の状態
- backend の `ExpandNeighbors` 呼び出し
- graph → tree への変換ルール
- domain-specific metadata
- side panel / top bar / workspace page shell

## Paper-in-Paper 側で必要な変更

### 1. core node type を作り直す

現在の `Paper` 型は小さすぎます。

- `id`
- `title`
- `description`
- `content`
- `parentId`
- `childIds`

Synthify では、これに加えて次のような情報が必要です。

- `category`
- `scope`
- `documentId`
- `documentIds`
- `sourceDocuments`
- `sourceChunkIds`
- `relatedNodes`
- `badges`
- `isNew`
- `isLoading`
- `actions`

やり方は 2 通りあります。

- `Paper` を直接拡張する
- `meta` のような汎用フィールドを追加する

最初の段階では、直接拡張した方が単純です。

ただし、Synthify 側の要求を考えると `meta` フィールドを持たせる案も強いです。

たとえば次のような block list を `meta` に入れられるようにすると、backend は「何を見せるか」だけ返し、frontend は「どう描画するか」に集中できます。

- note
- metric
- relations list
- source documents
- warning
- action suggestion
- mini graph

補足:

- backend が node の open / close を直接制御する必要はない
- open / close は frontend / renderer 側で持つ
- backend は `meta.blocks` のような view model を返すだけでよい
- mini graph は可能だが、最初からやると複雑なので後回しでよい

ここでいう `note` は、単なる説明文ではなく、主に `AI が人間にアクションを求めるための付箋` を想定する。

例:

- 確認してほしい論点
- 判断してほしい分岐
- 不足している根拠の追加依頼
- 次に取ってほしいアクション
- レビューしてほしい矛盾やリスク

つまり付箋は、

- backend / AI 側が生成する `human action request`

として扱い、

- frontend はそれを node 内で見やすく出す

という分担が自然。

### 2. 外部 callback を追加する

`PaperCanvas` は外から操作を受け取れる必要があります。

必要な callback:

- `onNodeClick`
- `onNodeOpen`
- `onNodeClose`
- `onNodeReorder`
- `onNodeAction`
- `onRequestExpandNeighbors`

これがないと、Synthify 側の backend-driven な挙動と綺麗に接続できません。

### 3. custom node rendering を追加する

Synthify は node の見た目を package 外から制御したいです。

必要な extension point:

- `renderNode`
- `renderNodeHeader`
- `renderNodeBody`
- `renderNodeFooter`

最低限必要なのは、

- node data と interaction state を受け取る `renderNode`

です。

理由:

- source documents を独自カードで見せたい
- category / scope を独自表示したい
- `Expand neighbors` ボタンを条件付きで出したい
- `New` / `Loading` バッジを app 側の状態で出したい
- `meta.blocks` に入った note / metric / relations / mini graph を描き分けたい

### 4. 外部制御の visual state を受け取れるようにする

Synthify 側ではすでに search / filter / highlight の状態を持っています。

そのため `PaperCanvas` は次のような状態を外から受け取れるべきです。

- `expandedNodeIds`
- `selectedNodeId`
- `highlightedNodeIds`
- `dimmedNodeIds`
- `hiddenNodeIds`

こうしておくと、app 側のロジックを renderer 内部に押し込まずに済みます。

### 5. node-level action を扱えるようにする

tree node の中に action を差し込める必要があります。

最低限必要なもの:

- node ごとの action button
- action の loading state
- action の disabled state

これは次のために必要です。

- `Expand neighbors`
- 将来的な node action

### 6. tree 以外の relation を出せるようにする

Synthify の graph には、tree 以外の relation もあります。

renderer 側で full edge 描画までやらなくてよいですが、少なくとも次のどれかは出せるべきです。

- related node の要約
- relation badge
- expanded node 内の relation section

これがないと、graph の意味がかなり失われます。

### 7. demo 前提の構造を薄くする

今の package には demo 的な前提や内部構造の癖がまだ残っています。

整理したいもの:

- stable な tree primitive
- stable な public type
- stable な render hook
- demo と本体の責務分離

### 8. drag-and-drop の意味を見直す

`REORDER` が入ったのは良い進展です。

ただし、まだ検討すべき点があります。

- drag を prop で無効化できるか
- reorder は opt-in にすべきか
- app 側が reorder を拒否できるか
- parent scope によって reorder を制限できるか

Synthify では、どこでも自由に tree を変えてよいとは限りません。

### 9. styling contract を明確にする

package 側でスタイルの境界をはっきりさせる必要があります。

欲しいもの:

- 安定した class 名または slot prop
- theme variable
- spacing / width / typography の上書き手段

現状のデフォルト見た目は、Synthify の最終形とはみなさない方がよいです。

## Synthify 側で必要な変更

補足:

- 現時点では `frontend に graph -> tree の変換レイヤーを厚く置く` よりも、
- `backend / DB 自体を tree-native に刷新する` 方が本筋になりつつある

そのため、このドキュメントの frontend 変換レイヤー案は暫定であり、backend 側の再設計メモも併読すること:

- `docs/tree-native-backend-redesign.md`

### 1. React Flow ベースの型と renderer を捨てる

削除対象:

- `frontend/features/graph/components/graph-canvas-panel.tsx`
- `frontend/features/graph/types/graph-canvas.ts`
- `frontend/features/graph/model/to-graph-canvas.ts`
- `@xyflow/react`
- `frontend/app/globals.css` の React Flow CSS import

### 2. tree model 変換レイヤーを追加する

新しく必要な変換レイヤー:

- `frontend/features/graph/model/to-paper-tree-model.ts`

この層がやること:

- backend graph を `paper-in-paper` 用 tree node に変換する
- renderer に必要な domain metadata を保持する
- root / parent / child を決める
- tree に載らない relation をどう持つか決める

### 3. graph → tree の変換ルールを決める

Synthify 側で決めるべきルール:

- どの edge type を hierarchy とみなすか
- hierarchy edge が無い場合に parent をどう補完するか
- root を何にするか
- document-scope node をどこにぶら下げるか
- `ExpandNeighbors` で増えた node をどこに入れるか

初期案:

- `HIERARCHICAL` を最優先
- 無ければ `level` から補完
- canonical を document より上に置く

### 4. workspace page を新 renderer 契約に合わせて簡素化する

`frontend/features/workspaces/components/workspace-graph-page.tsx` はかなり整理できます。

今後持つべきもの:

- search state
- filter state
- selected / expanded state
- source document data
- expand-neighbors の wiring

逆に、React Flow node style の直接注入はやめる。

### 5. backend graph API は当面そのままでよい

proto を今すぐ変える必要はありません。

backend は引き続き次を返してよいです。

- `Graph`
- `Node`
- `Edge`

tree への変換は、最初の段階では frontend 側で行う。

### 6. search / filter を tree UI 向けに再設計する

今の検索・フィルタは graph 前提です。

tree UI では次を決める必要があります。

- unmatched node を隠すか、薄くするか
- descendant が match したとき parent chain を残すか
- filter は subtree を消すのか、注記だけにするのか

### 7. tree に乗らない relation を UI 上どこに出すか決める

必要な判断:

- node の中に inline 表示
- 右側 panel
- expanded body 内の relation section
- 将来の overlay / secondary view

これは renderer だけの話ではなく、プロダクト上の見せ方の判断です。

### 8. テストを置き換える

React Flow 前提のテストは捨てる。

必要なもの:

- 新 tree UI に対する renderer-level test
- workspace の E2E test

## UI の見た目方針

目指す見た目は、`graph tool` より `research desk` に近い。

キーワード:

- editorial
- paper workspace
- sticky annotation
- calm intelligence

### node 本体

- 白〜生成りの paper card
- 角丸は控えめ
- 影は浅く、紙が重なっている程度
- 情報はタイポ中心で静かに配置する
- category / scope は細いラベルで補助的に見せる

### AI 付箋

- paper node とは色を分ける
- 薄い黄、淡い橙、薄い青などの柔らかい色を使う
- 少しだけ「貼ってある」感じを出す
- 本体より感情があるが、騒がしくしない
- 一文で `AI が人間に求めるアクション` を示す

### 動き

- node 展開は静かに開く
- 付箋は少し遅れて自然に出る
- hover の反応は小さく抑える

### 避けるもの

- チャット UI っぽい見た目
- kanban / task board っぽい見た目
- 派手な neon / gradient
- いかにも AI 的な紫発光表現

## 実装順

おすすめの順序:

1. `paper-in-paper` の型を作り直す
2. `PaperCanvas` の extension point を追加する
3. custom node rendering を入れる
4. 外部制御の visual state を受けられるようにする
5. Synthify 側に graph → tree 変換を作る
6. `GraphCanvasPanel` を置き換える
7. `@xyflow/react` を削除する
8. テストを更新する

## 今すぐ着手するべき具体タスク

最初の 4 つはこれです。

1. `vendor/paper-in-paper/src/lib/core/types.ts` を変更する
   - Synthify 用 metadata を node 型に追加する
2. `vendor/paper-in-paper/src/lib/react/PaperCanvas.tsx` を変更する
   - custom renderer と callback を受けられるようにする
3. `vendor/paper-in-paper/src/lib/react/internal` 配下を確認・変更する
   - richer な node data を visible UI に流す
4. `frontend/features/graph/model/to-paper-tree-model.ts` を作る
   - backend graph を新しい tree 入力モデルに変換する

## 最初の段階でやらなくてよいこと

初回実装では次を優先しない。

- npm 公開 API の互換維持
- React Flow 的な挙動の再現
- 汎用ライブラリとしての美しさ
- backend schema の再設計

優先すべきなのは次です。

- Synthify を React Flow から外すこと
- tree renderer を十分 expressive にすること
- data flow を理解しやすく保つこと
