import { WorkspaceDashboard } from '@/features/workspaces/components/workspace-dashboard'
import { mockWorkspaces } from '@/features/workspaces/data/mock-workspaces'

export default function DashboardPage() {
  return <WorkspaceDashboard workspaces={mockWorkspaces} />
}
