import { createFileRoute } from "@tanstack/react-router";
import { ComingSoonPage } from "@/modules/ComingSoon";

export const Route = createFileRoute("/coming-soon")({
  component: ComingSoonPage,
});
