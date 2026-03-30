'use client'

import Link from 'next/link'
import { useEffect, useRef, useState } from 'react'
import type { NodeMouseHandler } from '@xyflow/react'
import { GraphCanvasPanel } from '@/features/graph/components/graph-canvas-panel'
import { useGetGraph } from '@/features/graph/hooks/use-get-graph'
import { toGraphCanvas } from '@/features/graph/model/to-graph-canvas'
import {
  getWorkspaceCard,
  getWorkspaceDocuments,
  updateWorkspaceName,
  uploadWorkspaceDocument,
} from '../data/get-workspaces'
import type { WorkspaceCard, WorkspaceDocument } from '../types/workspace-card'

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
  const [documents, setDocuments] = useState<WorkspaceDocument[]>([])
  const [isUploadOpen, setIsUploadOpen] = useState(false)
  const [selectedFile, setSelectedFile] = useState<File | null>(null)
  const [isUploading, setIsUploading] = useState(false)
  const [graphRefreshKey, setGraphRefreshKey] = useState(0)
  const [expandedNodeIds, setExpandedNodeIds] = useState<string[]>([])
  const nameInputRef = useRef<HTMLInputElement | null>(null)
  const renameInFlightRef = useRef(false)
  const uploadInFlightRef = useRef(false)
  const refreshInFlightRef = useRef(false)

  useEffect(() => {
    let cancelled = false

    async function load() {
      setIsWorkspaceLoading(true)
      const [result, workspaceDocuments] = await Promise.all([getWorkspaceCard(workspaceId), getWorkspaceDocuments(workspaceId)])

      if (cancelled) {
        return
      }

      setWorkspace(result.workspace)
      setDraftName(result.workspace?.name ?? '')
      setDocuments(workspaceDocuments)
      setPageError(result.error ?? (result.workspace ? null : 'Workspace not found.'))
      setIsWorkspaceLoading(false)
    }

    void load()

    return () => {
      cancelled = true
    }
  }, [workspaceId])

  useEffect(() => {
    if (!documents.some((document) => ['UPLOADED', 'PROCESSING', 'PENDING_NORMALIZATION'].includes(document.statusCode))) {
      return
    }

    const intervalId = window.setInterval(() => {
      void refreshWorkspaceState()
    }, 1500)

    return () => {
      window.clearInterval(intervalId)
    }
  }, [documents, workspaceId])

  useEffect(() => {
    if (!isEditingName) {
      return
    }

    nameInputRef.current?.focus()
    nameInputRef.current?.select()
  }, [isEditingName])

  const { graph, isLoading, error } = useGetGraph({
    workspaceId: workspace?.id ?? workspaceId,
    documentId: null,
    refreshKey: graphRefreshKey,
  })
  const combinedError = pageError ?? error
  const emptyMessage = workspace ? 'この workspace にはまだ document がありません。' : 'Workspace not found.'
  const expandedNodeIdSet = new Set(expandedNodeIds)
  const sourceDocumentsByNodeId = graph
    ? Object.fromEntries(
        graph.graph.nodes.map((node) => {
          const sourceDocuments = node.documentId
            ? documents.filter((document) => document.id === node.documentId)
            : documents

          return [
            node.id,
            sourceDocuments.map((document) => ({
              id: document.id,
              filename: document.filename,
              status: document.status,
            })),
          ]
        }),
      )
    : {}
  const expandedCanvas = graph
    ? toGraphCanvas(graph.graph, {
        expandedNodeIds: expandedNodeIdSet,
        sourceDocumentsByNodeId,
      })
    : null

  useEffect(() => {
    if (!graph?.graph.nodes.length) {
      setExpandedNodeIds([])
      return
    }

    const graphNodeIds = new Set(graph.graph.nodes.map((node) => node.id))
    setExpandedNodeIds((current) => current.filter((nodeId) => graphNodeIds.has(nodeId)))
  }, [graph])

  async function refreshWorkspaceState() {
    if (refreshInFlightRef.current) {
      return
    }

    try {
      refreshInFlightRef.current = true
      const [result, workspaceDocuments] = await Promise.all([getWorkspaceCard(workspaceId), getWorkspaceDocuments(workspaceId)])
      setWorkspace(result.workspace)
      setDraftName(result.workspace?.name ?? '')
      setDocuments(workspaceDocuments)
      setPageError(result.error ?? (result.workspace ? null : 'Workspace not found.'))
      setGraphRefreshKey((value) => value + 1)
    } finally {
      refreshInFlightRef.current = false
    }
  }

  async function commitWorkspaceName() {
    if (!workspace || renameInFlightRef.current) {
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
      renameInFlightRef.current = true
      setIsRenaming(true)
      setPageError(null)
      const updatedWorkspace = await updateWorkspaceName(workspace.id, name)
      setWorkspace(updatedWorkspace)
      setDraftName(updatedWorkspace.name)
      setIsEditingName(false)
    } catch (renameError) {
      setPageError(renameError instanceof Error ? renameError.message : String(renameError))
    } finally {
      renameInFlightRef.current = false
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

  async function handleUploadDocument(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()

    if (!workspace || !selectedFile || uploadInFlightRef.current) {
      return
    }

    try {
      uploadInFlightRef.current = true
      setIsUploading(true)
      setPageError(null)
      await uploadWorkspaceDocument(workspace.id, selectedFile)
      await refreshWorkspaceState()
      setSelectedFile(null)
      setIsUploadOpen(false)
    } catch (uploadError) {
      setPageError(uploadError instanceof Error ? uploadError.message : String(uploadError))
    } finally {
      uploadInFlightRef.current = false
      setIsUploading(false)
    }
  }

  const handleNodeSelect: NodeMouseHandler = (_, node) => {
    setExpandedNodeIds((current) =>
      current.includes(node.id) ? current.filter((nodeId) => nodeId !== node.id) : [...current, node.id],
    )
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

          <div className="flex items-center gap-2">
            {expandedNodeIds.length > 0 && (
              <button
                type="button"
                onClick={() => setExpandedNodeIds([])}
                className="rounded-lg border border-slate-200/70 bg-white/80 px-3 py-1.5 text-xs text-slate-500 shadow-sm backdrop-blur-sm transition-colors hover:text-slate-700"
              >
                Collapse all
              </button>
            )}
            <button
              type="button"
              onClick={() => setIsUploadOpen((open) => !open)}
              className="rounded-lg border border-slate-200/70 bg-white/80 px-3 py-1.5 text-xs text-slate-500 shadow-sm backdrop-blur-sm transition-colors hover:text-slate-700"
            >
              Upload file
            </button>
            <Link
              href="/dashboard"
              className="rounded-lg border border-slate-200/70 bg-white/70 px-3 py-1.5 text-xs text-slate-400 shadow-sm backdrop-blur-sm transition-colors hover:text-slate-600"
            >
              ← Back
            </Link>
          </div>
        </div>
      </header>

      {/* Full-screen canvas */}
      <main className="absolute inset-0">
        <GraphCanvasPanel
          canvas={expandedCanvas}
          isLoading={isWorkspaceLoading || isLoading}
          error={combinedError}
          emptyMessage={emptyMessage}
          onNodeSelect={handleNodeSelect}
        />
      </main>

      {isUploadOpen && (
        <div className="absolute right-5 top-20 z-30 w-[360px] rounded-2xl border border-slate-200 bg-white/95 p-4 shadow-xl backdrop-blur-sm">
          <div className="flex items-start justify-between gap-3">
            <div>
              <p className="text-sm font-semibold text-slate-900">Upload file</p>
              <p className="mt-1 text-xs text-slate-500">
                ファイルを追加すると mock processing 後に graph の対象 document が切り替わります。
              </p>
            </div>
            <button
              type="button"
              onClick={() => setIsUploadOpen(false)}
              className="text-xs text-slate-400 transition-colors hover:text-slate-600"
            >
              Close
            </button>
          </div>

          <form onSubmit={handleUploadDocument} className="mt-4 space-y-4">
            <label className="flex cursor-pointer flex-col items-center justify-center rounded-xl border border-dashed border-slate-300 bg-slate-50 px-4 py-6 text-center transition-colors hover:border-slate-400 hover:bg-white">
              <span className="text-sm font-medium text-slate-700">
                {selectedFile ? selectedFile.name : 'Select a file'}
              </span>
              <span className="mt-1 text-xs text-slate-400">
                PDF, text, or any mock file metadata
              </span>
              <input
                type="file"
                className="hidden"
                onChange={(event) => setSelectedFile(event.target.files?.[0] ?? null)}
              />
            </label>

            <button
              type="submit"
              disabled={!selectedFile || isUploading}
              className="inline-flex w-full items-center justify-center rounded-xl bg-slate-900 px-4 py-3 text-sm font-medium text-white transition-colors hover:bg-slate-700 disabled:cursor-not-allowed disabled:opacity-60"
            >
              {isUploading ? 'Uploading…' : 'Upload and process'}
            </button>
          </form>

          <div className="mt-4 border-t border-slate-100 pt-4">
            <p className="text-xs font-semibold uppercase tracking-[0.18em] text-slate-400">Documents</p>
            <div className="mt-3 space-y-2">
              {documents.length === 0 ? (
                <p className="text-sm text-slate-400">まだ document はありません。</p>
              ) : (
                documents.slice(0, 4).map((document) => (
                  <div key={document.id} className="rounded-xl border border-slate-100 bg-slate-50 px-3 py-2">
                    <div className="flex items-center justify-between gap-3">
                      <p className="truncate text-sm font-medium text-slate-700">{document.filename}</p>
                      <span className="text-xs text-slate-400">{document.status}</span>
                    </div>
                    <p className="mt-1 text-xs text-slate-400">{document.id}</p>
                  </div>
                ))
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
