export default function DashboardLoading() {
  return (
    <div className="min-h-screen bg-slate-50">
      <header className="border-b border-slate-200 bg-white">
        <div className="mx-auto flex h-14 max-w-5xl items-center gap-3 px-6 md:px-10">
          <span className="h-2 w-2 rounded-full bg-teal-500" />
          <span className="text-sm font-semibold text-slate-900">Synthify</span>
        </div>
      </header>

      <main className="mx-auto max-w-5xl px-6 py-12 md:px-10">
        <div className="mb-10">
          <div className="h-8 w-40 animate-pulse rounded bg-slate-200" />
          <div className="mt-3 h-4 w-72 animate-pulse rounded bg-slate-100" />
        </div>

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
      </main>
    </div>
  )
}
