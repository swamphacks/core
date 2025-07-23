import { cn } from "@/utils/cn";
import { type DialogProps, Dialog as RACDialog } from "react-aria-components";

export function Dialog(props: DialogProps) {
  return (
    <RACDialog
      {...props}
      className={cn(
        "outline-0 [[data-placement]>&]:p-2 max-h-[inherit] overflow-auto relative",
        props.className,
      )}
    />
  );
}
