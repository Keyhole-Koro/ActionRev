import { notFound } from 'next/navigation'
import { WorkspaceGraphPage } from '@/features/workspaces/components/workspace-graph-page'
import { getWorkspaceCard } from '@/features/workspaces/data/mock-workspaces'

type WorkspacePageProps = {
  params: Promise<{
    id: string
  }>
}

export default async function WorkspacePage(props: WorkspacePageProps) {
  const { id } = await props.params
  const workspace = getWorkspaceCard(id)

  if (!workspace) {
    notFound()
  }

  return <WorkspaceGraphPage workspace={workspace} />
}
