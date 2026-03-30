export type WorkspaceCard = {
  id: string
  name: string
  summary: string
  ownerLabel: string
  graphNodeCount: number
  graphEdgeCount: number
  documentId: string | null
  badge: string
}

export type WorkspaceDocument = {
  id: string
  filename: string
  mimeType: string
  fileSize: number
  status: string
}
