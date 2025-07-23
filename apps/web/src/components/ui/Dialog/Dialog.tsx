import { cn } from "@/utils/cn";
import { type DialogProps, Dialog as RACDialog } from "react-aria-components";

export function Dialog(props: DialogProps) {
  return (
    <RACDialog
      {...props}
      className={cn(
        "outline outline-0 p-6 [[data-placement]>&]:p-4 max-h-[inherit] overflow-auto relative",
        props.className,
      )}
    />
  );
}
