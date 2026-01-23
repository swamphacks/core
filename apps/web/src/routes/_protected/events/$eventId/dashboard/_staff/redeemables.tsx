import { createFileRoute } from "@tanstack/react-router";
import { RedeemableCard } from "@/features/Redeemables/components/RedeemableCard";
import { CreateRedeemableModal } from "@/features/Redeemables/components/CreateRedeemableModal";
import { useRedeemables } from "@/features/Redeemables/hooks/useRedeemables";
import { Button } from "@/components/ui/Button";
import { DialogTrigger, Heading, Text } from "react-aria-components";

export const Route = createFileRoute(
  "/_protected/events/$eventId/dashboard/_staff/redeemables",
)({
  component: RouteComponent,
});

function RouteComponent() {
  const { eventId } = Route.useParams();
  const { eventRole } = Route.useRouteContext();
  const { data: redeemables, isLoading, isError } = useRedeemables(eventId);

  if (isLoading) {
    return (
      <div className="flex flex-col gap-8">
        <div className="flex justify-between items-center">
          <Heading className="text-3xl text-text-main">Redeemables</Heading>
        </div>
        <div className="flex flex-row gap-6">
          {[...Array(3)].map((_, i) => (
            <div
              key={i}
              className="w-96 h-64 bg-neutral-200 dark:bg-neutral-800 animate-pulse rounded"
            />
          ))}
        </div>
      </div>
    );
  }

  if (isError) {
    return (
      <div className="flex flex-col gap-4">
        <Heading className="text-3xl text-text-main">Redeemables</Heading>
        <Text className="text-red-500">
          Error loading redeemables. Please refresh and try again.
        </Text>
      </div>
    );
  }

  return (
    <div className="flex flex-col gap-8">
      <div className="flex justify-between items-center">
        <Heading className="text-3xl text-text-main">Redeemables</Heading>
        <DialogTrigger>
          <Button variant="primary">Create New Redeemable</Button>
          <CreateRedeemableModal eventId={eventId} />
        </DialogTrigger>
      </div>

      {redeemables && redeemables.length > 0 ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {redeemables.map((redeemable) => (
            <RedeemableCard
              key={redeemable.id}
              id={redeemable.id}
              name={redeemable.name}
              totalStock={redeemable.total_stock}
              maxUserAmount={redeemable.max_user_amount}
              totalRedeemed={
                typeof redeemable.total_redeemed === "string"
                  ? parseInt(redeemable.total_redeemed) || 0
                  : redeemable.total_redeemed
              }
              eventId={eventId}
              eventRole={eventRole}
            />
          ))}
        </div>
      ) : (
        <div className="flex flex-col items-center justify-center py-12 text-center">
          <Text className="text-text-secondary text-lg mb-4">
            No redeemables found
          </Text>
          <DialogTrigger>
            <Button variant="primary">Create Your First Redeemable</Button>
            <CreateRedeemableModal eventId={eventId} />
          </DialogTrigger>
        </div>
      )}
    </div>
  );
}
