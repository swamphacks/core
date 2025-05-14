import { Button } from "@/components/ui/Button/Button";
import { createFileRoute } from "@tanstack/react-router";
import auth from "@/lib/auth";

export const Route = createFileRoute("/")({
  component: Index,
});

function Index() {
  return (
    <div className="p-2">
      <h3>SwampHacks Portal</h3>
      <Button
        color="primary"
        onPress={() => auth.signInWithOauth({ provider: "discord" })}
      >
        Hello everyone!
      </Button>
    </div>
  );
}
