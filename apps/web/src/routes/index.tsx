import { createFileRoute } from "@tanstack/react-router";
import { Button } from "@/components/button";

export const Route = createFileRoute("/")({
  component: Index,
});

function Index() {
  return (
    <div className="p-2">
      <h3>SwampHacks Portal</h3>
      <Button
        onClick={() => console.log("Button clicked!")}
        className="bg-blue-300 rounded-md p-2 cursor-pointer hover:bg-blue-400"
      >
        Click me!
      </Button>
    </div>
  );
}
