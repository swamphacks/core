import ECharts from "@/components/ECharts";
import { useTheme } from "@/components/ThemeProvider";
import type { components } from "@/lib/openapi/schema";

interface Props {
  data?: components["schemas"]["sqlc.GetApplicationMajorSplitRow"][];
}

export default function ApplicationMajorsChart({ data }: Props) {
  const { theme } = useTheme();

  const isDark =
    theme === "dark" ||
    (theme === "system" &&
      window.matchMedia("(prefers-color-scheme: dark)").matches);

  if (!data) {
    return (
      <div>
        <p>Loading...</p>
      </div>
    );
  }

  return (
    <div className="inline-block bg-input-bg rounded-md w-full h-full shadow-xs border-neutral-200 dark:border-neutral-800">
      <ECharts
        className=" w-full h-full"
        option={{
          title: {
            text: "Application Majors",
            textStyle: {
              color: isDark ? "#FFFFFF" : "#000000",
              fontFamily: "Figtree",
              fontSize: 22,
            },
            padding: 15,
          },
          tooltip: {
            show: true,
            trigger: "item",
            backgroundColor: isDark ? "#333333" : "#FFFFFF",
            textStyle: {
              color: isDark ? "#FFFFFF" : "#000000",
            },
            formatter: (params) => {
              const p = Array.isArray(params) ? params[0] : params;
              const percent = p.percent?.toFixed(1);

              return `${p.marker} ${p.name}: ${p.value} (${percent}%)`;
            },
          },
          series: [
            {
              type: "pie",
              center: ["50%", "55%"],
              label: {
                show: true,
                color: isDark ? "#FFFFFF" : "#000000",
              },
              data: data.map((row) => ({
                value: row.count,
                name: row.major,
              })),
            },
          ],
        }}
      />
    </div>
  );
}
