import ApplicationGenderChart from "./ApplicationGenderChart";
import ApplicationAgeChart from "./ApplicationAgeChart";
import ApplicationRaceChart from "./ApplicationRaceChart";
import ApplicationMajorsChart from "./ApplicationMajorsChart";
import ApplicationSchoolsChart from "./ApplicationSchoolsChart";
import { Heading } from "react-aria-components";
import { useApplicationStatistics } from "@/features/Application/hooks/useApplicationStatistics";
import { Card } from "@/components/ui/Card";
import { useEventOverview } from "@/features/EventOverview/hooks/useEventOverview";
import ApplicationStats from "./ApplicationStats";

interface ApplicationStatisticsProps {
  eventId: string;
}

export default function ApplicationStatistics({
  eventId,
}: ApplicationStatisticsProps) {
  const { data, isLoading, isError, error } = useApplicationStatistics(eventId);
  const {
    data: eventData,
    isLoading: isEventLoading,
    isError: isEventError,
    error: eventError,
  } = useEventOverview(eventId);

  if ((isError && error) || (isEventError && eventError)) {
    return (
      <main>
        <Heading className="text-2xl lg:text-3xl font-semibold mb-6">
          Application Statistics
        </Heading>
        <Card className="p-6 bg-red-100 text-red-800">
          <p>
            Error loading application statistics:{" "}
            {(error || (eventError as Error)).message}
          </p>
        </Card>
      </main>
    );
  }

  if (isLoading || isEventLoading) {
    return (
      <main>
        <Heading className="text-2xl lg:text-3xl font-semibold mb-6">
          Application Statistics
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

  if (!data || !eventData) {
    return <p>Something went wrong :(</p>;
  }

  return (
    <main>
      <Heading className="text-2xl lg:text-3xl font-semibold mb-6 flex flex-col gap-2">
        Application Statistics
        <p className="text-lg text-text-secondary">
          Pssst, we look better on Laptops/Desktops...
        </p>
      </Heading>

      <section className="grid gap-3 grid-cols-1 lg:grid-cols-4 auto-rows-[minmax(300px,auto)]">
        {/* Row 1: small charts */}
        <div className="md:col-span-1 h-full">
          <ApplicationStats data={data} eventData={eventData} />
        </div>
        <div className="md:col-span-1 h-full">
          <ApplicationGenderChart data={data.gender_stats} />
        </div>
        <div className="md:col-span-1">
          <ApplicationAgeChart data={data.age_stats} />
        </div>
        <div className="md:col-span-1">
          <ApplicationRaceChart data={data.race_stats} />
        </div>

        {/* Row 2 & 3: Majors chart spans 2 columns and 2 rows on md+ */}
        <div className="w-full h-full md:col-span-2 md:row-span-2">
          <ApplicationMajorsChart data={data.major_stats} />
        </div>

        {/* Row 2 & 3: Schools chart spans 2 columns and 2 rows on md+ */}
        <div className="w-full h-full md:col-span-2 md:row-span-2">
          <ApplicationSchoolsChart data={data.school_stats} />
        </div>
      </section>
    </main>
  );
}
