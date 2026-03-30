import { MarkerType, type Edge, type Node } from '@xyflow/react'
import { EdgeType, GraphProjectionScope, NodeCategory, type Graph } from '@/src/generated/synthify/graph/v1/graph_types_pb'
import type { GraphCanvas, GraphCanvasNodeData, GraphCanvasSourceDocument } from '../types/graph-canvas'

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

const scopeLabels: Record<number, string> = {
  [GraphProjectionScope.DOCUMENT]: 'Document',
  [GraphProjectionScope.CANONICAL]: 'Canonical',
}

const COLLAPSED_NODE_WIDTH = 220
const EXPANDED_NODE_WIDTH = 360
const COLLAPSED_NODE_HEIGHT = 92
const COLUMN_GAP = 96
const ROW_GAP = 28
const SOURCE_ROW_HEIGHT = 52
const METRICS_SECTION_HEIGHT = 88
const EXPANDED_PADDING_HEIGHT = 48

type ToGraphCanvasOptions = {
  expandedNodeIds?: Iterable<string>
  sourceDocumentsByNodeId?: Record<string, GraphCanvasSourceDocument[]>
}

function estimateExpandedHeight(description: string, sourceDocuments: GraphCanvasSourceDocument[], chunkCount: number) {
  const normalizedDescription = description.trim()
  const estimatedLines = Math.max(1, Math.ceil(normalizedDescription.length / 42))
  const descriptionHeight = estimatedLines * 20 + 30
  const sourcesHeight = Math.max(40, sourceDocuments.length * SOURCE_ROW_HEIGHT)
  const chunksHeight = chunkCount > 0 ? 0 : 0
  return COLLAPSED_NODE_HEIGHT + EXPANDED_PADDING_HEIGHT + descriptionHeight + METRICS_SECTION_HEIGHT + sourcesHeight + chunksHeight
}

export function toGraphCanvas(graph: Graph | undefined, options: ToGraphCanvasOptions = {}): GraphCanvas {
  if (!graph) {
    return { nodes: [], edges: [] }
  }

  const expandedNodeIds = new Set(options.expandedNodeIds ?? [])
  const sourceDocumentsByNodeId = options.sourceDocumentsByNodeId ?? {}
  const levelGroups = new Map<number, typeof graph.nodes>()
  for (const node of graph.nodes) {
    const level = node.level || 0
    const group = levelGroups.get(level) ?? []
    group.push(node)
    levelGroups.set(level, group)
  }

  const sortedLevels = [...levelGroups.keys()].sort((a, b) => a - b)
  const nodes: Node<GraphCanvasNodeData>[] = []
  let currentX = 0

  for (const level of sortedLevels) {
    const group = levelGroups.get(level) ?? []
    let currentY = 0
    let columnWidth = COLLAPSED_NODE_WIDTH

    group.forEach((node) => {
      const category = categoryLabels[node.category] ?? 'Unknown'
      const expanded = expandedNodeIds.has(node.id)
      const sourceDocuments = sourceDocumentsByNodeId[node.id] ?? []
      const width = expanded ? EXPANDED_NODE_WIDTH : COLLAPSED_NODE_WIDTH
      const height = expanded
        ? estimateExpandedHeight(node.description, sourceDocuments, node.sourceChunkIds.length)
        : COLLAPSED_NODE_HEIGHT

      columnWidth = Math.max(columnWidth, width)
      nodes.push({
        id: node.id,
        type: 'graphNode',
        position: {
          x: currentX,
          y: currentY,
        },
        data: {
          label: node.label,
          category,
          level,
          description: node.description,
          documentId: node.documentId || null,
          scope: scopeLabels[node.scope] ?? 'Unknown',
          sourceChunkIds: node.sourceChunkIds ?? [],
          expanded,
          sourceDocuments,
        },
        draggable: false,
        selectable: true,
        style: {
          width,
          minHeight: height,
          borderRadius: 24,
          border: `1px solid ${categoryColors[category] ?? categoryColors.Unknown}`,
          background: 'rgba(255,255,255,0.92)',
          boxShadow: expanded ? '0 28px 60px rgba(15, 23, 42, 0.16)' : '0 20px 45px rgba(15, 23, 42, 0.08)',
          zIndex: expanded ? 30 : 10,
        },
      })

      currentY += height + ROW_GAP
    })

    currentX += columnWidth + COLUMN_GAP
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
