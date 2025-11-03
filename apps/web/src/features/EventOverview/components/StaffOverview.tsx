import { Heading } from "react-aria-components";
import { Card } from "@/components/ui/Card";
import ApplicationStats from "@/features/EventOverview/components/ApplicationStats";
import EventDetails from "@/features/EventOverview/components/EventDetails";
import { useEventOverview } from "@/features/EventOverview/hooks/useEventOverview";

interface Props {
  eventId: string;
}

export default function StaffOverview({ eventId }: Props) {
  const { data, isLoading, isError, error } = useEventOverview(eventId);

  if (isError && error) {
    return (
      <main>
        <Heading className="text-2xl lg:text-3xl font-semibold mb-6">
          Overview
        </Heading>
        <Card className="p-6 bg-red-100 text-red-800">
          <p>Error loading application statistics: {error.message}</p>
        </Card>
      </main>
    );
  }

  if (isLoading) {
    return (
      <main>
        <Heading className="text-2xl lg:text-3xl font-semibold mb-6">
          Overview
        </Heading>

        {/* Skeleton Cards */}
        <section>
          <div className="flex flex-row gap-6">
            {[...Array(3)].map((_, i) => (
              <div
                key={i}
                className="md:w-96 h-64 bg-neutral-200 dark:bg-neutral-800 animate-pulse rounded"
              />
            ))}
          </div>
        </section>
      </main>
    );
  }

  if (!data) {
    return <p>Something went wrong :(</p>;
  }

  return (
    <main>
      <Heading className="text-2xl lg:text-3xl font-semibold mb-6 flex flex-col gap-2">
        Overview
        <p className="text-lg text-text-secondary">
          Pssst, we look better on Laptops/Desktops...
        </p>
      </Heading>

      <div className="flex flex-col space-y-3 max-w-150">
        <EventDetails data={data} />
        <ApplicationStats data={data} />
      </div>
    </main>
  );
}
