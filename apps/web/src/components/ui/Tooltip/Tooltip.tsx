import type { TooltipTriggerProps } from "react-aria";
import {
  Tooltip as RAC_Tooltip,
  TooltipTrigger,
  type TooltipProps as RAC_ToolTipProps,
} from "react-aria-components";

// Use this to extend or customize trigger props in the future
// interface TriggerProps extends TooltipTriggerProps {}

interface TooltipPanelProps extends Omit<RAC_ToolTipProps, "children"> {
  label: string;
}

/**
 * Usage:
 * <Tooltip
 *  tooltipProps={{ placement: "top", label: "This is a tooltip" }}
 *  triggerProps={{ delay: 500 }}
 * >
 *  <button>Hover me</button>
 * </Tooltip>
 */
export interface TooltipProps {
  triggerProps?: TooltipTriggerProps;
  tooltipProps: TooltipPanelProps;
  children: React.ReactNode;
}

const Tooltip = ({ children, triggerProps, tooltipProps }: TooltipProps) => {
  return (
    <TooltipTrigger {...triggerProps}>
      {children}
      <RAC_Tooltip
        className="bg-surface border-[1px] border-border/40 px-2 py-1 text-sm rounded-xs"
        {...tooltipProps}
      >
        {tooltipProps.label}
      </RAC_Tooltip>
    </TooltipTrigger>
  );
};

Tooltip.displayName = "Tooltip";

export { Tooltip };
