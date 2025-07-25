import { createFileRoute } from "@tanstack/react-router";
import { build } from "@/features/FormBuilder/build";
import data from "@/features/FormBuilder/stories/applicationFormExample.json";

export const Route = createFileRoute("/application")({
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <div className="w-full h-screen bg-surface">
      <div className="w-full bg-surface transition-[background]">
        <ApplicationForm />
      </div>
    </div>
  );
}

export const ApplicationForm = () => {
  const Form = build(data);

  return <Form />;
};
