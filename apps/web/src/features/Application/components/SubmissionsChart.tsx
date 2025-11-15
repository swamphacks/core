import ECharts from "@/components/ECharts";
import { useTheme } from "@/components/ThemeProvider";
import { useMemo } from "react";

interface SubmissionsChartProps {
  submission_stats?: {
    count: number;
    day: string;
  }[];
}

export default function SubmissionsChart({
  submission_stats,
}: SubmissionsChartProps) {
  // TODO: Make this more elegant, throws error when submission stats is null/undefined
  if (!submission_stats) {
    return null;
  }

  const { theme } = useTheme();

  const isDark = theme === "dark";

  const chartData = useMemo(() => {
    const labels: string[] = [];
    const values: number[] = [];

    if (submission_stats) {
      for (const submission of submission_stats) {
        labels.push(
          new Date(submission.day).toLocaleDateString(undefined, {
            weekday: "narrow",
            month: "short",
            day: "numeric",
            timeZone: "UTC",
          }),
        );
        values.push(submission.count);
      }
    }

    return { labels, values };
  }, [submission_stats]);

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
          barWidth: "80%",
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
    <div className="w-full max-w-full h-80 mt-3 overflow-hidden">
      <ECharts className="w-full h-full" option={chartOption as any} />
    </div>
  );
}
