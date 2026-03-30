'use client'

import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { useEffect, useState } from 'react'
import { createWorkspace, getWorkspaceCards } from '../data/get-workspaces'
import type { WorkspaceCard } from '../types/workspace-card'

export function WorkspaceDashboard() {
  const router = useRouter()
  const [workspaces, setWorkspaces] = useState<WorkspaceCard[]>([])
  const [error, setError] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [isCreating, setIsCreating] = useState(false)

  useEffect(() => {
    let cancelled = false

    async function load() {
      setIsLoading(true)
      const result = await getWorkspaceCards()

      if (cancelled) {
        return
      }

      setWorkspaces(result.workspaces)
      setError(result.error)
      setIsLoading(false)
    }

    void load()

    return () => {
      cancelled = true
    }
  }, [])

  async function handleCreateWorkspace() {
    try {
      setIsCreating(true)
      setError(null)
      const workspace = await createWorkspace('Untitled workspace')
      router.push(`/workspace/${workspace.id}`)
    } catch (createError) {
      setError(createError instanceof Error ? createError.message : String(createError))
      setIsCreating(false)
    }
  }

  return (
    <div className="min-h-screen bg-slate-50">
      {/* Nav */}
      <header className="border-b border-slate-200 bg-white">
        <div className="mx-auto flex h-14 max-w-5xl items-center gap-3 px-6 md:px-10">
          <span className="h-2 w-2 rounded-full bg-teal-500" />
          <span className="text-sm font-semibold text-slate-900">Synthify</span>
        </div>
      </header>

      <main className="mx-auto max-w-5xl px-6 py-12 md:px-10">
        {/* Page title */}
        <div className="mb-10">
          <div className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
            <div>
              <h1 className="text-2xl font-bold text-slate-900">Workspaces</h1>
              <p className="mt-1 text-sm text-slate-500">
                ドキュメントから抽出したグラフを確認・探索する
              </p>
            </div>
            <button
              type="button"
              onClick={handleCreateWorkspace}
              disabled={isCreating}
              className="inline-flex items-center justify-center rounded-xl bg-slate-900 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-slate-700"
            >
              {isCreating ? 'Creating…' : 'Create workspace'}
            </button>
          </div>
        </div>

        {/* Workspace list */}
        {isLoading ? (
          <div className="space-y-2">
            {Array.from({ length: 3 }).map((_, index) => (
              <div
                key={index}
                className="rounded-xl border border-slate-200 bg-white px-5 py-4 shadow-sm"
              >
                <div className="h-4 w-48 animate-pulse rounded bg-slate-200" />
                <div className="mt-3 h-4 w-full animate-pulse rounded bg-slate-100" />
              </div>
            ))}
          </div>
        ) : error ? (
          <div className="rounded-2xl border border-red-200 bg-white p-6 shadow-sm">
            <p className="text-sm font-semibold text-red-600">Workspace list unavailable</p>
            <p className="mt-1 text-sm text-slate-500">{error}</p>
          </div>
        ) : workspaces.length === 0 ? (
          <div className="rounded-2xl border border-slate-200 bg-white p-6 text-sm text-slate-500 shadow-sm">
            表示できる workspace がまだありません。
          </div>
        ) : (
          <div className="space-y-2">
            {workspaces.map((ws) => (
              <Link
                key={ws.id}
                href={`/workspace/${ws.id}`}
                className="group flex items-center justify-between gap-6 rounded-xl border border-slate-200 bg-white px-5 py-4 shadow-sm transition-colors hover:border-slate-300 hover:bg-slate-50"
              >
                <div className="min-w-0">
                  <div className="flex items-center gap-2.5">
                    <h2 className="truncate text-sm font-semibold text-slate-900">{ws.name}</h2>
                    <span
                      className={`shrink-0 rounded-md px-2 py-0.5 text-xs font-medium ${
                        ws.badge === 'Ready'
                          ? 'bg-teal-50 text-teal-700'
                          : ws.badge === 'Failed'
                            ? 'bg-red-50 text-red-600'
                            : ws.badge === 'Empty'
                              ? 'bg-slate-100 text-slate-500'
                              : 'bg-amber-50 text-amber-700'
                      }`}
                    >
                      {ws.badge}
                    </span>
                  </div>
                  <p className="mt-1 truncate text-sm text-slate-500">{ws.summary}</p>
                </div>

                <div className="flex shrink-0 items-center gap-6 text-sm">
                  <div className="hidden items-center gap-4 text-slate-400 sm:flex">
                    <span>
                      <span className="font-medium text-slate-700">{ws.graphNodeCount}</span> nodes
                    </span>
                    <span>
                      <span className="font-medium text-slate-700">{ws.graphEdgeCount}</span> edges
                    </span>
                    <span className="font-mono text-xs text-slate-300">{ws.documentId ?? 'no-doc'}</span>
                  </div>
                  <span className="text-slate-300 transition-colors group-hover:text-teal-500">→</span>
                </div>
              </Link>
            ))}
          </div>
        )}
      </main>
    </div>
  )
}
