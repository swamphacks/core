import TablerClipboardData from "~icons/tabler/clipboard-data";
import ECharts from "@/components/ECharts";
import { useTheme } from "@/components/ThemeProvider";
import { useMemo } from "react";
import type { EventOverview } from "@/features/EventOverview/hooks/useEventOverview";

interface ApplicationStatsProps {
  data: EventOverview;
}

export default function ApplicationStats({ data }: ApplicationStatsProps) {
  const { theme } = useTheme();

  const isDark =
    theme === "dark" ||
    (theme === "system" &&
      window.matchMedia("(prefers-color-scheme: dark)").matches);

  const chartData = useMemo(() => {
    const labels: string[] = [];
    const values: number[] = [];

    for (const submission of data.application_submission_stats) {
      labels.push(
        new Date(submission.day).toLocaleDateString(undefined, {
          month: "short",
          day: "numeric",
          timeZone: "UTC",
        }),
      );
      values.push(submission.count);
    }
    return { labels, values };
  }, [data]);

  const chartOption = useMemo(
    () => ({
      color: isDark ? ["#0984E3", "#0984E3"] : ["#0984E3", "#00B894"],
      title: {
        text: "Daily Submissions (EST Timezone)",
        textStyle: {
          color: isDark ? "#e4e4e7" : "#000000",
          fontFamily: "Figtree",
          fontSize: 16,
          fontWeight: "normal",
        },
        left: "auto",
        padding: [5, 5, 40, 5],
      },
      tooltip: {
        trigger: "axis",
        axisPointer: {
          type: "shadow",
        },
      },
      grid: {
        left: 12,
        right: 12,
        top: 60,
        bottom: 24,
      },
      xAxis: {
        type: "category",
        data: chartData.labels,
        axisLine: {
          lineStyle: {
            color: isDark ? "#444" : "#ddd",
          },
        },
        axisLabel: {
          rotate: 45,
          interval: Math.ceil(chartData.labels.length / 10),
        },
      },
      yAxis: {
        type: "value",
        axisLine: {
          lineStyle: {
            color: isDark ? "#444" : "#ddd",
          },
        },
        splitLine: {
          lineStyle: {
            color: isDark ? "#2b2b2b" : "#f1f1f1",
          },
        },
      },
      series: [
        {
          name: "Submissions",
          type: "bar",
          data: chartData.values,
          barWidth: "60%",
          itemStyle: {
            borderRadius: [4, 4, 0, 0],
          },
          emphasis: {
            itemStyle: {
              opacity: 0.9,
            },
          },
        },
      ],
    }),
    [chartData, isDark],
  );

  return (
    <div className="border border-input-border rounded-md px-4 py-3 w-full">
      <div className="flex justify-between items-start relative">
        <p className="flex items-center gap-1 bg-badge-bg-attending rounded-md px-2 py-1 w-fit">
          <TablerClipboardData /> Applications
        </p>
        <div className="flex gap-2 items-center text-text-secondary text-sm">
          <span>Deadline:</span>
          {new Intl.DateTimeFormat("en-US", {
            year: "numeric",
            month: "short",
            day: "numeric",
            hour: "numeric",
            minute: "2-digit",
            hour12: true,
          }).format(new Date(data.event_details.application_close))}
        </div>
      </div>

      <div className="flex gap-5 mt-3 items-center min-w-0">
        <div className="bg-surface py-2 px-4 rounded-md">
          <p className="text-text-secondary">Total</p>
          <p className="text-2xl">
            {data.application_status_stats.started +
              data.application_status_stats.submitted}
          </p>
        </div>

        <div>
          <p className="text-xl">=</p>
        </div>

        <div className="flex gap-3 items-center">
          <div className="bg-surface py-2 px-4 rounded-md">
            <p className="text-text-secondary">Started</p>
            <p className="text-2xl">{data.application_status_stats.started}</p>
          </div>

          <div>
            <p className="text-xl">+</p>
          </div>

          <div className="bg-surface py-2 px-4 rounded-md">
            <p className="text-text-secondary">Submitted</p>
            <p className="text-2xl">
              {data.application_status_stats.submitted}
            </p>
          </div>
        </div>
      </div>

      <div className="w-full max-w-full h-80 mt-3 overflow-hidden">
        <ECharts className="w-full h-full" option={chartOption as any} />
      </div>
    </div>
  );
}
