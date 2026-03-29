'use client'

import { useEffect, useState } from 'react'
import { graphClient } from '@/lib/graph-client'
import { toGraphCanvas } from '../model/to-graph-canvas'
import type { LoadedGraph } from '../types/graph-canvas'
import type { GraphSummary } from '../types/graph-summary'

type UseGetGraphResult = {
  graph: LoadedGraph | null
  summary: GraphSummary | null
  isLoading: boolean
  error: string | null
}

export function useGetGraph(): UseGetGraphResult {
  const [graph, setGraph] = useState<LoadedGraph | null>(null)
  const [summary, setSummary] = useState<GraphSummary | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let cancelled = false

    async function load() {
      try {
        setIsLoading(true)
        const response = await graphClient.getGraph({
          workspaceId: 'ws_demo',
          documentId: 'doc_demo',
          categoryFilters: [],
          levelFilters: [],
          edgeTypeFilters: [],
          limit: 50,
          sourceFilename: '',
          resolveAliases: false,
        })

        if (cancelled) {
          return
        }

        if (response.graph) {
          setGraph({
            graph: response.graph,
            canvas: toGraphCanvas(response.graph),
          })

          setSummary({
            documentId: response.documentId,
            nodeCount: response.graph.nodes.length,
            edgeCount: response.graph.edges.length,
            nodeLabels: response.graph.nodes.map((node) => node.label),
          })
        } else {
          setGraph(null)
          setSummary(null)
        }

        setError(null)
      } catch (err) {
        if (cancelled) {
          return
        }
        setError(err instanceof Error ? err.message : 'failed to load graph')
      } finally {
        if (!cancelled) {
          setIsLoading(false)
        }
      }
    }

    void load()

    return () => {
      cancelled = true
    }
  }, [])

  return { graph, summary, isLoading, error }
}
