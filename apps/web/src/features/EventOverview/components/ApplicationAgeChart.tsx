import ECharts from "@/components/ECharts";
import { useTheme } from "@/components/ThemeProvider";
import { Card } from "@/components/ui/Card";
import type { components } from "@/lib/openapi/schema";

interface Props {
  data?: components["schemas"]["sqlc.GetApplicationAgeSplitRow"];
}

export default function ApplicationAgeChart({ data }: Props) {
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
    <Card className="w-full h-72 md:h-80 bg-input-bg">
      <ECharts
        className=" w-full h-full"
        option={{
          color: [
            "#00B894", // mint
            "#55EFC4", // light teal
            "#0984E3", // vivid blue
            "#6C5CE7", // violet
            "#E84393", // pink
            "#D63031", // red
            "#FDCB6E", // warm yellow
          ],
          title: {
            text: "Application Ages",
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

              return `${p.marker} ${p.value} (${percent}%)`;
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
              data: [
                {
                  value: 1,
                  name: "<18 yrs",
                },
                {
                  value: 30,
                  name: "18 yrs",
                },
                {
                  value: 20,
                  name: "19 yrs",
                },
                {
                  value: 15,
                  name: "20 yrs",
                },
                {
                  value: 13,
                  name: "21 yrs",
                },
                {
                  value: 5,
                  name: "22 yrs",
                },
                {
                  value: 2,
                  name: "23+ yrs",
                },
              ].filter((item) => item.value > 0),
            },
          ],
        }}
      />
    </Card>
  );
}
