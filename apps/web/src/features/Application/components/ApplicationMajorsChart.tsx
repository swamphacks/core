import ECharts from "@/components/ECharts";
import { useTheme } from "@/components/ThemeProvider";
import type { components } from "@/lib/openapi/schema";

interface Props {
  data: components["schemas"]["sqlc.GetApplicationMajorSplitRow"][];
}

export default function ApplicationMajorsChart({ data }: Props) {
  const { theme } = useTheme();

  const isDark =
    theme === "dark" ||
    (theme === "system" &&
      window.matchMedia("(prefers-color-scheme: dark)").matches);

  const sortedData = data
    .slice()
    .sort((a, b) => (a.count ?? 0) - (b.count ?? 0));

  const names = sortedData.map((row) => row.major);
  const counts = sortedData.map((row) => row.count);
  const total = counts.reduce((s, v) => s + (v ?? 0), 0);

  return (
    <div className="inline-block bg-input-bg rounded-md w-full h-full shadow-xs border-neutral-200 dark:border-neutral-800">
      <ECharts
        className="w-full h-full"
        option={{
          title: {
            text: "Majors",
            textStyle: {
              color: isDark ? "#e4e4e7" : "#18181b",
              fontFamily: "Figtree",
              fontSize: 18,
            },
            padding: 15,
          },
          tooltip: {
            trigger: "axis",
            backgroundColor: isDark ? "#333333" : "#FFFFFF",
            textStyle: {
              color: isDark ? "#FFFFFF" : "#000000",
            },
            axisPointer: {
              type: "shadow",
            },
            formatter: (params) => {
              const p = Array.isArray(params) ? params[0] : params;
              const value = p.value ?? 0;
              const percent = total
                ? (((value as number) / total) * 100).toFixed(1)
                : "0.0";
              return `${p.marker} ${p.name}: ${value} (${percent}%)`;
            },
          },
          grid: {
            left: 16,
            right: 24,
            top: 60,
            bottom: 24,
            containLabel: true,
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
            inverse: false,
            axisLine: { lineStyle: { color: isDark ? "#888888" : "#666666" } },
            axisLabel: {
              color: isDark ? "#FFFFFF" : "#000000",
              formatter: (v: string) => v,
            },
          },
          series: [
            {
              type: "bar",
              data: counts,
              barCategoryGap: "50%",
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
