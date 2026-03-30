import type { WorkspaceCard } from '../types/workspace-card'

export const mockWorkspaces: WorkspaceCard[] = [
  {
    id: 'ws_demo',
    name: 'Growth Strategy Review',
    summary: '営業戦略資料から抽出した document graph を確認するワークスペース。',
    ownerLabel: 'Owner: unix',
    graphNodeCount: 3,
    graphEdgeCount: 2,
    documentId: 'doc_demo',
    badge: 'Live stub',
  },
  {
    id: 'ws_metrics',
    name: 'Metrics Exploration',
    summary: 'KPI と evidence を横断的に追うための検証用ワークスペース。',
    ownerLabel: 'Owner: unix',
    graphNodeCount: 3,
    graphEdgeCount: 2,
    documentId: 'doc_demo',
    badge: 'Next target',
  },
]

export function getWorkspaceCard(workspaceId: string) {
  return mockWorkspaces.find((workspace) => workspace.id === workspaceId) ?? null
}
