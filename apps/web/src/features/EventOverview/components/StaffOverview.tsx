import { Heading } from "react-aria-components";
import { useApplicationStatistics } from "../hooks/useApplicationStatistics";
import { Card } from "@/components/ui/Card";
import ApplicationGenderChart from "./ApplicationGenderChart";
import ApplicationAgeChart from "./ApplicationAgeChart";
import ApplicationRaceChart from "./ApplicationRaceChart";
import ApplicationMajorsChart from "./ApplicationMajorsChart";

interface Props {
  eventId: string;
}

export default function StaffOverview({ eventId }: Props) {
  const { data, isLoading, isError, error } = useApplicationStatistics(eventId);

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

  return (
    <main>
      <Heading className="text-2xl lg:text-3xl font-semibold mb-6">
        Overview
      </Heading>

      {/* Stat cards for applications */}
      <section className="grid gap-6 grid-cols-1 md:grid-cols-2 xl:grid-cols-4">
        <ApplicationGenderChart data={data?.gender_stats} />
        <ApplicationAgeChart data={data?.age_stats} />
        <ApplicationRaceChart data={data?.race_stats} />
        <ApplicationMajorsChart data={data?.major_stats} />
      </section>
    </main>
  );
}
