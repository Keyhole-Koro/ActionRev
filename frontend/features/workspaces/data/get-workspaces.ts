import type { Workspace } from '@/src/generated/synthify/graph/v1/workspace_pb'
import { workspaceClient } from '@/lib/rpc-client'
import type { WorkspaceCard } from '../types/workspace-card'

type WorkspaceCardsResult = {
  workspaces: WorkspaceCard[]
  error: string | null
}

type WorkspaceCardResult = {
  workspace: WorkspaceCard | null
  error: string | null
}

async function toWorkspaceCard(workspace: Workspace): Promise<WorkspaceCard> {
  return {
    id: workspace.workspaceId,
    name: workspace.name,
    summary: 'まだドキュメントがありません。',
    ownerLabel: 'Workspace',
    graphNodeCount: 0,
    graphEdgeCount: 0,
    documentId: null,
    badge: 'Empty',
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
