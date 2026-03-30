import { WorkspaceGraphPage } from '@/features/workspaces/components/workspace-graph-page'

type WorkspacePageProps = {
  params: Promise<{
    id: string
  }>
}

export default async function WorkspacePage(props: WorkspacePageProps) {
  const { id } = await props.params

  return <WorkspaceGraphPage workspaceId={id} />
}
