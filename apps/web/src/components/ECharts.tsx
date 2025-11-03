import * as echarts from "echarts/core";

import { CanvasRenderer } from "echarts/renderers";

import {
  TitleComponent,
  TooltipComponent,
  GridComponent,
  // Dataset
  DatasetComponent,
  // Built-in transform (filter, sort)
  TransformComponent,
} from "echarts/components";

import type { ComposeOption } from "echarts/core";

import { PieChart, LineChart, BarChart } from "echarts/charts";

import { LabelLayout, UniversalTransition } from "echarts/features";

import type {
  PieSeriesOption,
  LineSeriesOption,
  BarSeriesOption,
} from "echarts/charts";

import type {
  // The component option types are defined with the ComponentOption suffix
  TitleComponentOption,
  TooltipComponentOption,
  GridComponentOption,
  DatasetComponentOption,
} from "echarts/components";
import { useEffect, useRef } from "react";

type ECOption = ComposeOption<
  | PieSeriesOption
  | TitleComponentOption
  | TooltipComponentOption
  | GridComponentOption
  | DatasetComponentOption
  | LineSeriesOption
  | BarSeriesOption
>;

echarts.use([
  TitleComponent,
  TooltipComponent,
  GridComponent,
  DatasetComponent,
  TransformComponent,
  PieChart,
  LineChart,
  BarChart,
  LabelLayout,
  UniversalTransition,
  CanvasRenderer,
]);

interface Props {
  option: ECOption;
  className?: string;
}

export default function ECharts({ option, className }: Props) {
  const chartRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const chartInstance = echarts.init(chartRef.current);

    chartInstance.setOption(option);

    const handleResize = () => {
      chartInstance.resize();
    };
    window.addEventListener("resize", handleResize);

    return () => {
      window.removeEventListener("resize", handleResize);
      chartInstance.dispose();
    };
  }, [option]);

  return <div ref={chartRef} className={className}></div>;
}

export const pieChartColors = [
  "#00B894", // mint
  "#55EFC4", // light teal
  "#0984E3", // vivid blue
  "#6C5CE7", // violet
  "#E84393", // pink
  "#D63031", // red
  "#FDCB6E", // warm yellow
];
