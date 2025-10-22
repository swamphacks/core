import { Heading } from "react-aria-components";
import { useApplicationStatistics } from "../hooks/useApplicationStatistics";
import { Card } from "@/components/ui/Card";
import * as echarts from "echarts/core";

import type { ComposeOption } from "echarts/core";

import { PieChart } from "echarts/charts";

import {
  TitleComponent,
  TooltipComponent,
  GridComponent,
  DatasetComponent,
  TransformComponent,
} from "echarts/components";

import type {
  // The component option types are defined with the ComponentOption suffix
  TitleComponentOption,
  TooltipComponentOption,
  GridComponentOption,
  DatasetComponentOption,
} from "echarts/components";

import type {
  // The series option types are defined with the SeriesOption suffix
  PieSeriesOption,
} from "echarts/charts";

import { LabelLayout, UniversalTransition } from "echarts/features";
import { SVGRenderer } from "echarts/renderers";

echarts.use([
  PieChart,
  TitleComponent,
  TooltipComponent,
  GridComponent,
  DatasetComponent,
  TransformComponent,
  LabelLayout,
  UniversalTransition,
  SVGRenderer,
]);

type ECOption = ComposeOption<
  | PieSeriesOption
  | TitleComponentOption
  | TooltipComponentOption
  | GridComponentOption
  | DatasetComponentOption
>;

interface Props {
  eventId: string;
}

export default function StaffOverview({ eventId }: Props) {
  const { data, isLoading, isError, error } = useApplicationStatistics(eventId);

  let myChart = echarts.init(document.getElementById("chart"), {
    width: 600,
    height: 400,
  });

  const option: ECOption = {
    title: {
      text: "Application Status Overview",
      left: "center",
    },
    tooltip: {
      trigger: "item",
    },
    legend: {
      orient: "vertical",
      left: "left",
    },
    series: [
      {
        name: "Applications",
        type: "pie",
        radius: "50%",
        data: [
          { value: 100, name: "Total Applications" },
          { value: 30, name: "Accepted Applications" },
          { value: 50, name: "Rejected Applications" },
          { value: 20, name: "Pending Applications" },
        ],
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: "rgba(0, 0, 0, 0.5)",
          },
        },
      },
    ],
  };

  myChart.setOption(option);

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
                className="w-96 h-64 bg-neutral-200 dark:bg-neutral-800 animate-pulse rounded"
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
      <section>
        <Card className="w-96 h-80" id="chart"></Card>
      </section>
    </main>
  );
}
