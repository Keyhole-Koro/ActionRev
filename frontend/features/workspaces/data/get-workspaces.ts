import { DocumentLifecycleState, type Document } from '@/src/generated/synthify/graph/v1/document_pb'
import type { Workspace } from '@/src/generated/synthify/graph/v1/workspace_pb'
import { documentClient, workspaceClient } from '@/lib/rpc-client'
import type { WorkspaceCard, WorkspaceDocument } from '../types/workspace-card'

type WorkspaceCardsResult = {
  workspaces: WorkspaceCard[]
  error: string | null
}

type WorkspaceCardResult = {
  workspace: WorkspaceCard | null
  error: string | null
}

async function toWorkspaceCard(workspace: Workspace): Promise<WorkspaceCard> {
  const documents = await getWorkspaceDocuments(workspace.workspaceId)
  const latestDocument = documents[0] ?? null

  return {
    id: workspace.workspaceId,
    name: workspace.name,
    summary: latestDocument ? `${latestDocument.filename} · ${latestDocument.status}` : 'まだドキュメントがありません。',
    ownerLabel: 'Workspace',
    graphNodeCount: 0,
    graphEdgeCount: 0,
    documentId: latestDocument?.id ?? null,
    badge: latestDocument?.status ?? 'Empty',
  }
}

export async function getWorkspaceCards(): Promise<WorkspaceCardsResult> {
  try {
    const response = await workspaceClient.listWorkspaces({})
    const workspaces = await Promise.all(response.workspaces.map(toWorkspaceCard))

    return {
      workspaces,
      error: null,
    }
  } catch (error) {
    return {
      workspaces: [],
      error: error instanceof Error ? error.message : 'Failed to load workspaces.',
    }
  }
}

export async function getWorkspaceCard(workspaceId: string): Promise<WorkspaceCardResult> {
  try {
    const workspace = await workspaceClient.getWorkspace({ workspaceId })

    return {
      workspace: await toWorkspaceCard(workspace),
      error: null,
    }
  } catch (error) {
    if (error instanceof Error && /not found/i.test(error.message)) {
      return {
        workspace: null,
        error: 'Workspace not found.',
      }
    }

    return {
      workspace: null,
      error: error instanceof Error ? error.message : 'Failed to load workspace.',
    }
  }
}

export async function createWorkspace(name: string) {
  const workspace = await workspaceClient.createWorkspace({ name })
  return toWorkspaceCard(workspace)
}

export async function updateWorkspaceName(workspaceId: string, name: string) {
  const workspace = await workspaceClient.updateWorkspace({ workspaceId, name })
  return toWorkspaceCard(workspace)
}

function getDocumentStatusLabel(status: DocumentLifecycleState) {
  switch (status) {
    case DocumentLifecycleState.UPLOADED:
      return 'Uploaded'
    case DocumentLifecycleState.PENDING_NORMALIZATION:
      return 'Queued'
    case DocumentLifecycleState.PROCESSING:
      return 'Processing'
    case DocumentLifecycleState.COMPLETED:
      return 'Ready'
    case DocumentLifecycleState.FAILED:
      return 'Failed'
    default:
      return 'Idle'
  }
}

function mapDocument(document: Document): WorkspaceDocument {
  return {
    id: document.documentId,
    filename: document.filename,
    mimeType: document.mimeType,
    fileSize: Number(document.fileSize),
    status: getDocumentStatusLabel(document.status),
  }
}

export async function getWorkspaceDocuments(workspaceId: string): Promise<WorkspaceDocument[]> {
  const response = await documentClient.listDocuments({ workspaceId })
  return response.documents.map(mapDocument)
}

export async function uploadWorkspaceDocument(workspaceId: string, file: File) {
  const response = await documentClient.createDocument({
    workspaceId,
    filename: file.name,
    mimeType: file.type || 'application/octet-stream',
    fileSize: BigInt(file.size),
  })

  const document = response.document
  if (!document) {
    throw new Error('Failed to create document.')
  }

  await documentClient.startProcessing({
    workspaceId,
    documentId: document.documentId,
    forceReprocess: false,
    extractionDepth: 1,
  })

  return mapDocument({
    ...document,
    status: DocumentLifecycleState.COMPLETED,
  })
}
