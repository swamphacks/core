import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_protected/_admin/test2')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/_protected/_admin/test2"!</div>
}
