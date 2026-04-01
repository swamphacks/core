import ECharts from "@/components/ECharts";
import { useTheme } from "@/components/ThemeProvider";
import type { components } from "@/lib/openapi/schema";

interface Props {
  data: components["schemas"]["sqlc.GetApplicationSchoolSplitRow"][];
}

export default function ApplicationSchoolsChart({ data }: Props) {
  const { theme } = useTheme();

  const isDark =
    theme === "dark" ||
    (theme === "system" &&
      window.matchMedia("(prefers-color-scheme: dark)").matches);

  const chartData = data
    .map((d) => ({ name: d.school, value: d.count ?? 0 }))
    .sort((a, b) => a.value - b.value);

  const names = chartData.map((d) => d.name);
  const values = chartData.map((d) => d.value);

  return (
    <div className="inline-block bg-input-bg rounded-md w-full h-full shadow-xs border-neutral-200 dark:border-neutral-800">
      <ECharts
        className="w-full h-full"
        option={{
          title: {
            text: "Schools",
            textStyle: {
              color: isDark ? "#e4e4e7" : "#18181b",
              fontFamily: "Figtree",
              fontSize: 18,
            },
            padding: 15,
          },
          tooltip: {
            trigger: "axis",
            axisPointer: {
              type: "shadow",
            },
            backgroundColor: isDark ? "#333333" : "#FFFFFF",
            textStyle: {
              color: isDark ? "#FFFFFF" : "#000000",
            },
            formatter: (params) => {
              const p = Array.isArray(params) ? params[0] : params;
              const value = p.value ?? 0;
              const total = values.reduce((s, v) => s + v, 0) || 1;
              const percent = (((value as number) / total) * 100).toFixed(1);
              return `${p.marker} ${p.name}: ${value} (${percent}%)`;
            },
          },
          grid: {
            left: 16,
            right: 24,
            top: 60,
            bottom: 24,
          },
          xAxis: {
            type: "value",
            axisLine: { lineStyle: { color: isDark ? "#888888" : "#666666" } },
            axisLabel: { color: isDark ? "#FFFFFF" : "#000000" },
            splitLine: {
              lineStyle: {
                color: isDark ? "#2b2b2b" : "#f1f1f1",
              },
            },
          },
          yAxis: {
            type: "category",
            data: names,
            axisLine: { lineStyle: { color: isDark ? "#888888" : "#666666" } },
            axisLabel: {
              color: isDark ? "#FFFFFF" : "#000000",
              formatter: (v: string) => v,
            },
            inverse: false,
          },
          series: [
            {
              type: "bar",
              data: values,
              barCategoryGap: "40%",
              itemStyle: {
                color: isDark ? "#60a5fa" : "#2563eb",
              },
              label: {
                show: true,
                position: "right",
                color: isDark ? "#FFFFFF" : "#000000",
                formatter: "{c}",
              },
            },
          ],
        }}
      />
    </div>
  );
}
