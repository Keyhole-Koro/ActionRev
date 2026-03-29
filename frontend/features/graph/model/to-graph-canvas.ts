import { MarkerType, type Edge, type Node } from '@xyflow/react'
import { EdgeType, NodeCategory, type Graph } from '@/src/generated/synthify/graph/v1/graph_types_pb'
import type { GraphCanvas, GraphCanvasNodeData } from '../types/graph-canvas'

const categoryLabels: Record<number, string> = {
  [NodeCategory.CONCEPT]: 'Concept',
  [NodeCategory.CLAIM]: 'Claim',
  [NodeCategory.EVIDENCE]: 'Evidence',
  [NodeCategory.ENTITY]: 'Entity',
  [NodeCategory.METRIC]: 'Metric',
  [NodeCategory.ACTION]: 'Action',
}

const edgeLabels: Record<number, string> = {
  [EdgeType.HIERARCHICAL]: 'hierarchical',
  [EdgeType.RELATED_TO]: 'related_to',
  [EdgeType.SUPPORTS]: 'supports',
  [EdgeType.CONTRADICTS]: 'contradicts',
  [EdgeType.CAUSES]: 'causes',
  [EdgeType.MEASURED_BY]: 'measured_by',
  [EdgeType.MENTIONS]: 'mentions',
}

const categoryColors: Record<string, string> = {
  Concept: '#0f172a',
  Claim: '#2563eb',
  Evidence: '#059669',
  Entity: '#7c3aed',
  Metric: '#d97706',
  Action: '#dc2626',
  Unknown: '#475569',
}

export function toGraphCanvas(graph: Graph | undefined): GraphCanvas {
  if (!graph) {
    return { nodes: [], edges: [] }
  }

  const levelGroups = new Map<number, typeof graph.nodes>()
  for (const node of graph.nodes) {
    const level = node.level || 0
    const group = levelGroups.get(level) ?? []
    group.push(node)
    levelGroups.set(level, group)
  }

  const sortedLevels = [...levelGroups.keys()].sort((a, b) => a - b)
  const nodes: Node<GraphCanvasNodeData>[] = []

  for (const level of sortedLevels) {
    const group = levelGroups.get(level) ?? []
    group.forEach((node, index) => {
      const category = categoryLabels[node.category] ?? 'Unknown'
      nodes.push({
        id: node.id,
        position: {
          x: level * 280,
          y: index * 160,
        },
        data: {
          label: node.label,
          category,
          level,
          description: node.description,
        },
        draggable: false,
        selectable: true,
        style: {
          width: 220,
          borderRadius: 24,
          border: `1px solid ${categoryColors[category] ?? categoryColors.Unknown}`,
          background: 'rgba(255,255,255,0.92)',
          boxShadow: '0 20px 45px rgba(15, 23, 42, 0.08)',
          padding: 16,
        },
      })
    })
  }

  const edges: Edge[] = graph.edges.map((edge) => ({
    id: edge.id,
    source: edge.source,
    target: edge.target,
    type: 'smoothstep',
    animated: edge.type !== EdgeType.HIERARCHICAL,
    label: edgeLabels[edge.type] ?? 'edge',
    markerEnd: {
      type: MarkerType.ArrowClosed,
      width: 18,
      height: 18,
      color: '#94a3b8',
    },
    style: {
      stroke: '#94a3b8',
      strokeWidth: edge.type === EdgeType.HIERARCHICAL ? 2 : 2.5,
    },
    labelStyle: {
      fill: '#475569',
      fontSize: 12,
      fontWeight: 600,
    },
    labelBgStyle: {
      fill: '#ffffff',
      fillOpacity: 0.9,
    },
  }))

  return { nodes, edges }
}
