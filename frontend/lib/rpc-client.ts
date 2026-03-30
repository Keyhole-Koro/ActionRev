import { createClient } from '@connectrpc/connect'
import { createConnectTransport } from '@connectrpc/connect-web'
import { DocumentService } from '@/src/generated/synthify/graph/v1/document_pb'
import { GraphService } from '@/src/generated/synthify/graph/v1/graph_pb'
import { WorkspaceService } from '@/src/generated/synthify/graph/v1/workspace_pb'

export function resolveApiBaseUrl() {
  const publicBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL ?? 'http://localhost:8080'
  const internalBaseUrl = process.env.INTERNAL_API_BASE_URL ?? 'http://backend:8080'

  return typeof window === 'undefined' ? internalBaseUrl : publicBaseUrl
}

const transport = createConnectTransport({
  baseUrl: resolveApiBaseUrl(),
})

export const graphClient = createClient(GraphService, transport)
export const documentClient = createClient(DocumentService, transport)
export const workspaceClient = createClient(WorkspaceService, transport)
