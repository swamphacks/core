import { Button } from "@/components/ui/Button";
import { useMyApplication } from "@/modules/Application/hooks/useMyApplication";
import { createFileRoute, Link, Outlet } from "@tanstack/react-router";

export const Route = createFileRoute("/_protected/team/join")({
  component: RouteComponent,
});

function RouteComponent() {
  const { data, isPending } = useMyApplication();

  if (isPending) {
    return (
      <div className="w-full h-full sm:max-w-180 mx-auto font-figtree p-2 relative flex justify-center items-center pb-50">
        <p>Checking application status...</p>
      </div>
    );
  }

  if (!data || data.status === "started") {
    return (
      <div className="w-full h-full sm:max-w-180 mx-auto font-figtree p-2 relative flex justify-center items-center pb-50">
        <div>
          <p>You must submit an application before joining a team.</p>
          <Link to="/application">
            <Button className="mt-2">Apply Now!</Button>
          </Link>
        </div>
      </div>
    );
  }

  if (data.status === "submitted") {
    return <Outlet />;
  }
}
