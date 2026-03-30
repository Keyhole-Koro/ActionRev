'use client'

import Link from 'next/link'
import { useEffect, useRef, useState } from 'react'
import { GraphCanvasPanel } from '@/features/graph/components/graph-canvas-panel'
import { useGetGraph } from '@/features/graph/hooks/use-get-graph'
import { getWorkspaceCard, updateWorkspaceName } from '../data/get-workspaces'
import type { WorkspaceCard } from '../types/workspace-card'

type WorkspaceGraphPageProps = {
  workspaceId: string
}

export function WorkspaceGraphPage({ workspaceId }: WorkspaceGraphPageProps) {
  const [workspace, setWorkspace] = useState<WorkspaceCard | null>(null)
  const [pageError, setPageError] = useState<string | null>(null)
  const [isWorkspaceLoading, setIsWorkspaceLoading] = useState(true)
  const [draftName, setDraftName] = useState('')
  const [isEditingName, setIsEditingName] = useState(false)
  const [isRenaming, setIsRenaming] = useState(false)
  const nameInputRef = useRef<HTMLInputElement | null>(null)

  useEffect(() => {
    let cancelled = false

    async function load() {
      setIsWorkspaceLoading(true)
      const result = await getWorkspaceCard(workspaceId)

      if (cancelled) {
        return
      }

      setWorkspace(result.workspace)
      setDraftName(result.workspace?.name ?? '')
      setPageError(result.error ?? (result.workspace ? null : 'Workspace not found.'))
      setIsWorkspaceLoading(false)
    }

    void load()

    return () => {
      cancelled = true
    }
  }, [workspaceId])

  useEffect(() => {
    if (!isEditingName) {
      return
    }

    nameInputRef.current?.focus()
    nameInputRef.current?.select()
  }, [isEditingName])

  const { graph, isLoading, error } = useGetGraph({
    workspaceId: workspace?.id ?? workspaceId,
    documentId: workspace?.documentId ?? null,
  })
  const combinedError = pageError ?? error
  const emptyMessage = workspace?.documentId
    ? 'No graph data.'
    : workspace
      ? 'この workspace にはまだ document がありません。'
      : 'Workspace not found.'

  async function commitWorkspaceName() {
    if (!workspace) {
      return
    }

    const name = draftName.trim()
    if (!name) {
      setPageError('Workspace name is required.')
      return
    }

    if (name === workspace.name) {
      setIsEditingName(false)
      return
    }

    try {
      setIsRenaming(true)
      setPageError(null)
      const updatedWorkspace = await updateWorkspaceName(workspace.id, name)
      setWorkspace(updatedWorkspace)
      setDraftName(updatedWorkspace.name)
      setIsEditingName(false)
    } catch (renameError) {
      setPageError(renameError instanceof Error ? renameError.message : String(renameError))
    } finally {
      setIsRenaming(false)
    }
  }

  function handleStartEditing() {
    if (!workspace || isWorkspaceLoading || isRenaming) {
      return
    }

    setDraftName(workspace.name)
    setIsEditingName(true)
  }

  function handleCancelEditing() {
    setDraftName(workspace?.name ?? '')
    setIsEditingName(false)
  }

  return (
    <div className="relative h-screen bg-white">
      {/* Floating nav */}
      <header className="absolute inset-x-0 top-0 z-20">
        <div className="mx-auto flex min-h-12 max-w-none items-start justify-between gap-4 px-5 py-3">
          <div className="flex min-w-0 flex-col gap-2">
            <nav className="flex items-center gap-2 text-sm">
              <span className="h-1.5 w-1.5 rounded-full bg-teal-500" />
              <Link href="/dashboard" className="text-slate-400 transition-colors hover:text-slate-600">
                Workspaces
              </Link>
              <span className="text-slate-300">/</span>
              {isEditingName ? (
                <input
                  ref={nameInputRef}
                  value={draftName}
                  onChange={(event) => setDraftName(event.target.value)}
                  onBlur={() => void commitWorkspaceName()}
                  onKeyDown={(event) => {
                    if (event.key === 'Enter') {
                      event.preventDefault()
                      void commitWorkspaceName()
                    }

                    if (event.key === 'Escape') {
                      event.preventDefault()
                      handleCancelEditing()
                    }
                  }}
                  disabled={isRenaming}
                  className="min-w-0 rounded-md border border-slate-200 bg-white/80 px-2 py-1 text-slate-600 outline-none transition-colors focus:border-slate-400"
                />
              ) : (
                <button
                  type="button"
                  onClick={handleStartEditing}
                  disabled={!workspace || isWorkspaceLoading || isRenaming}
                  className="truncate rounded-md px-2 py-1 text-slate-600 transition-colors hover:bg-white/70 disabled:cursor-default"
                >
                  {workspace?.name ?? 'Workspace'}
                </button>
              )}
            </nav>
          </div>

          <Link
            href="/dashboard"
            className="rounded-lg border border-slate-200/70 bg-white/70 px-3 py-1.5 text-xs text-slate-400 shadow-sm backdrop-blur-sm transition-colors hover:text-slate-600"
          >
            ← Back
          </Link>
        </div>
      </header>

      {/* Full-screen canvas */}
      <main className="absolute inset-0">
        <GraphCanvasPanel
          canvas={graph?.canvas ?? null}
          isLoading={isWorkspaceLoading || isLoading}
          error={combinedError}
          emptyMessage={emptyMessage}
        />
      </main>
    </div>
  )
}
