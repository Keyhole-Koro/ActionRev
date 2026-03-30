import { NodeCategory, EdgeType } from '@/src/generated/synthify/graph/v1/graph_types_pb'
import type { GraphFilters } from '../hooks/use-get-graph'

const CATEGORY_OPTIONS = [
  { value: NodeCategory.CONCEPT, label: 'Concept', color: '#0f172a' },
  { value: NodeCategory.CLAIM, label: 'Claim', color: '#2563eb' },
  { value: NodeCategory.EVIDENCE, label: 'Evidence', color: '#059669' },
  { value: NodeCategory.ENTITY, label: 'Entity', color: '#7c3aed' },
  { value: NodeCategory.METRIC, label: 'Metric', color: '#d97706' },
  { value: NodeCategory.ACTION, label: 'Action', color: '#dc2626' },
]

const EDGE_TYPE_OPTIONS = [
  { value: EdgeType.HIERARCHICAL, label: 'hierarchical' },
  { value: EdgeType.RELATED_TO, label: 'related_to' },
  { value: EdgeType.SUPPORTS, label: 'supports' },
  { value: EdgeType.CONTRADICTS, label: 'contradicts' },
  { value: EdgeType.CAUSES, label: 'causes' },
  { value: EdgeType.MEASURED_BY, label: 'measured_by' },
  { value: EdgeType.MENTIONS, label: 'mentions' },
]

const LIMIT_OPTIONS = [25, 50, 100, 200]

type GraphFilterPanelProps = {
  filters: GraphFilters
  onChange: (filters: GraphFilters) => void
  onClose: () => void
}

function toggle<T>(arr: T[], value: T): T[] {
  return arr.includes(value) ? arr.filter((v) => v !== value) : [...arr, value]
}

export function GraphFilterPanel({ filters, onChange, onClose }: GraphFilterPanelProps) {
  const categoryFilters = filters.categoryFilters ?? []
  const edgeTypeFilters = filters.edgeTypeFilters ?? []
  const limit = filters.limit ?? 50

  const activeFilterCount =
    categoryFilters.length + edgeTypeFilters.length + (limit !== 50 ? 1 : 0)

  function handleCategoryToggle(value: number) {
    onChange({ ...filters, categoryFilters: toggle(categoryFilters, value) })
  }

  function handleEdgeTypeToggle(value: number) {
    onChange({ ...filters, edgeTypeFilters: toggle(edgeTypeFilters, value) })
  }

  function handleLimitChange(value: number) {
    onChange({ ...filters, limit: value })
  }

  function handleReset() {
    onChange({ categoryFilters: [], edgeTypeFilters: [], limit: 50 })
  }

  return (
    <div className="absolute left-5 top-20 z-30 w-[300px] rounded-2xl border border-slate-200 bg-white/95 p-4 shadow-xl backdrop-blur-sm">
      <div className="flex items-center justify-between gap-3">
        <div className="flex items-center gap-2">
          <p className="text-sm font-semibold text-slate-900">Filter graph</p>
          {activeFilterCount > 0 && (
            <span className="rounded-full bg-slate-900 px-1.5 py-0.5 text-[10px] font-bold text-white">
              {activeFilterCount}
            </span>
          )}
        </div>
        <div className="flex items-center gap-2">
          {activeFilterCount > 0 && (
            <button
              type="button"
              onClick={handleReset}
              className="text-xs text-slate-400 transition-colors hover:text-slate-600"
            >
              Reset
            </button>
          )}
          <button
            type="button"
            onClick={onClose}
            className="text-xs text-slate-400 transition-colors hover:text-slate-600"
          >
            Close
          </button>
        </div>
      </div>

      <div className="mt-4 space-y-4">
        {/* Node categories */}
        <div>
          <p className="mb-2 text-xs font-semibold uppercase tracking-[0.18em] text-slate-400">
            Node type
            {categoryFilters.length > 0 && (
              <span className="ml-1.5 normal-case font-normal text-slate-400">
                ({categoryFilters.length} selected)
              </span>
            )}
          </p>
          <div className="flex flex-wrap gap-1.5">
            {CATEGORY_OPTIONS.map(({ value, label, color }) => {
              const active = categoryFilters.includes(value)
              return (
                <button
                  key={value}
                  type="button"
                  onClick={() => handleCategoryToggle(value)}
                  style={active ? { borderColor: color, color } : undefined}
                  className={`rounded-lg border px-2.5 py-1 text-xs font-medium transition-colors ${
                    active
                      ? 'bg-white'
                      : 'border-slate-200 bg-slate-50 text-slate-500 hover:border-slate-300 hover:text-slate-700'
                  }`}
                >
                  {label}
                </button>
              )
            })}
          </div>
          {categoryFilters.length === 0 && (
            <p className="mt-1.5 text-[11px] text-slate-400">All types shown</p>
          )}
        </div>

        {/* Edge types */}
        <div>
          <p className="mb-2 text-xs font-semibold uppercase tracking-[0.18em] text-slate-400">
            Edge type
            {edgeTypeFilters.length > 0 && (
              <span className="ml-1.5 normal-case font-normal text-slate-400">
                ({edgeTypeFilters.length} selected)
              </span>
            )}
          </p>
          <div className="flex flex-wrap gap-1.5">
            {EDGE_TYPE_OPTIONS.map(({ value, label }) => {
              const active = edgeTypeFilters.includes(value)
              return (
                <button
                  key={value}
                  type="button"
                  onClick={() => handleEdgeTypeToggle(value)}
                  className={`rounded-lg border px-2.5 py-1 font-mono text-xs transition-colors ${
                    active
                      ? 'border-slate-700 bg-slate-900 text-white'
                      : 'border-slate-200 bg-slate-50 text-slate-500 hover:border-slate-300 hover:text-slate-700'
                  }`}
                >
                  {label}
                </button>
              )
            })}
          </div>
          {edgeTypeFilters.length === 0 && (
            <p className="mt-1.5 text-[11px] text-slate-400">All edge types shown</p>
          )}
        </div>

        {/* Limit */}
        <div>
          <p className="mb-2 text-xs font-semibold uppercase tracking-[0.18em] text-slate-400">
            Node limit
          </p>
          <div className="flex gap-1.5">
            {LIMIT_OPTIONS.map((option) => (
              <button
                key={option}
                type="button"
                onClick={() => handleLimitChange(option)}
                className={`rounded-lg border px-2.5 py-1 text-xs font-medium transition-colors ${
                  limit === option
                    ? 'border-slate-700 bg-slate-900 text-white'
                    : 'border-slate-200 bg-slate-50 text-slate-500 hover:border-slate-300 hover:text-slate-700'
                }`}
              >
                {option}
              </button>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}
