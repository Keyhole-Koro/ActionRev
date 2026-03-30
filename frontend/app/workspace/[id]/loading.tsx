export default function WorkspaceLoading() {
  return (
    <div className="relative h-screen bg-white">
      <header className="absolute inset-x-0 top-0 z-20">
        <div className="mx-auto flex h-12 max-w-none items-center justify-between gap-4 px-5">
          <div className="h-4 w-48 animate-pulse rounded bg-slate-200" />
          <div className="h-8 w-20 animate-pulse rounded-lg bg-slate-100" />
        </div>
      </header>

      <main className="absolute inset-0 flex items-center justify-center">
        <div className="text-sm text-slate-300">Loading…</div>
      </main>
    </div>
  )
}
