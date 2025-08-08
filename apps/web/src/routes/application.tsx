import { createFileRoute } from "@tanstack/react-router";
import { build } from "@/features/FormBuilder/build";
import data from "@/features/FormBuilder/stories/applicationFormExample.json";

export const Route = createFileRoute("/application")({
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <div className="w-full h-screen bg-surface">
      {/* this padding left prevent the page being shifted when the form fully loads because of the scrollbar */}
      <div className="w-full bg-surface transition-[background] sm:pl-[calc(100vw-100%)]">
        <ApplicationForm />
      </div>
    </div>
  );
}

export const ApplicationForm = () => {
  const Form = build(data);

  return <Form />;
};
