export function HeroPanel() {
  return (
    <section className="overflow-hidden rounded-[2rem] bg-slate-950 px-8 py-10 text-slate-50 shadow-panel">
      <p className="text-sm uppercase tracking-[0.35em] text-teal-200">Synthify</p>
      <h1 className="mt-4 max-w-2xl text-4xl font-semibold tracking-tight md:text-6xl">
        Structured documents, rendered as a graph you can actually inspect.
      </h1>
      <p className="mt-6 max-w-xl text-base leading-7 text-slate-300 md:text-lg">
        The current shell fetches the live GetGraph stub, maps it into a document graph view model, and renders it on a React Flow canvas.
      </p>
      <div className="mt-8 flex flex-wrap gap-3 text-sm text-slate-200">
        <span className="rounded-full border border-white/15 bg-white/5 px-4 py-2">Next.js App Router</span>
        <span className="rounded-full border border-white/15 bg-white/5 px-4 py-2">Tailwind CSS</span>
        <span className="rounded-full border border-white/15 bg-white/5 px-4 py-2">Connect RPC</span>
        <span className="rounded-full border border-white/15 bg-white/5 px-4 py-2">React Flow</span>
      </div>
    </section>
  )
}
