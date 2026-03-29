import { createClient } from '@connectrpc/connect'
import { createConnectTransport } from '@connectrpc/connect-web'
import { GraphService } from '../src/generated/synthify/graph/v1/graph_pb'

const baseUrl = process.env.NEXT_PUBLIC_API_BASE_URL ?? 'http://localhost:8080'

const transport = createConnectTransport({
  baseUrl,
})

export const graphClient = createClient(GraphService, transport)
