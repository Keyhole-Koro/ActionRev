import type { Edge, Node } from '@xyflow/react'
import type { Graph } from '@/src/generated/synthify/graph/v1/graph_types_pb'

export type GraphCanvasNodeData = {
  label: string
  category: string
  level: number
  description: string
  documentId: string | null
  scope: string
  sourceChunkIds: string[]
  expanded?: boolean
  isExpanding?: boolean
  isNew?: boolean
  sourceDocuments?: Array<{
    id: string
    filename: string
    status: string
  }>
  onExpandNeighbors?: () => void
}

export type GraphCanvas = {
  nodes: Node<GraphCanvasNodeData>[]
  edges: Edge[]
}

export type GraphCanvasSourceDocument = NonNullable<GraphCanvasNodeData['sourceDocuments']>[number]

export type LoadedGraph = {
  graph: Graph
  canvas: GraphCanvas
}
