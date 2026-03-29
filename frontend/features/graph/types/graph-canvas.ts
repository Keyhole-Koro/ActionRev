import type { Edge, Node } from '@xyflow/react'
import type { Graph } from '@/src/generated/synthify/graph/v1/graph_types_pb'

export type GraphCanvasNodeData = {
  label: string
  category: string
  level: number
  description: string
}

export type GraphCanvas = {
  nodes: Node<GraphCanvasNodeData>[]
  edges: Edge[]
}

export type LoadedGraph = {
  graph: Graph
  canvas: GraphCanvas
}
