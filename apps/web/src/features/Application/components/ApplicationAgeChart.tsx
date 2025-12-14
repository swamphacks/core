import ECharts, { pieChartColors } from "@/components/ECharts";
import { useTheme } from "@/components/ThemeProvider";
import { Card } from "@/components/ui/Card";
import type { components } from "@/lib/openapi/schema";

interface Props {
  data: components["schemas"]["sqlc.GetApplicationAgeSplitRow"];
}

export default function ApplicationAgeChart({ data }: Props) {
  const { theme } = useTheme();

  const isDark =
    theme === "dark" ||
    (theme === "system" &&
      window.matchMedia("(prefers-color-scheme: dark)").matches);

  return (
    <Card className="w-full h-72 md:h-80 bg-input-bg">
      <ECharts
        className=" w-full h-full"
        option={{
          color: pieChartColors,
          title: {
            text: "Ages",
            textStyle: {
              color: isDark ? "#e4e4e7" : "#18181b",
              fontFamily: "Figtree",
              fontSize: 18,
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
              data: [
                {
                  value: data.underage,
                  name: "<18 yrs",
                },
                {
                  value: data.age_18,
                  name: "18 yrs",
                },
                {
                  value: data.age_19,
                  name: "19 yrs",
                },
                {
                  value: data.age_20,
                  name: "20 yrs",
                },
                {
                  value: data.age_21,
                  name: "21 yrs",
                },
                {
                  value: data.age_22,
                  name: "22 yrs",
                },
                {
                  value: data.age_23_plus,
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
