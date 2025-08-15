import { Heading, Text } from "react-aria-components";
import { toast, type TypeOptions } from "react-toastify";

export interface ToastProps {
  title?: string;
  message: string;
  icon?: React.ReactNode;
  duration?: number; // ms
  position?: "top-right" | "top-left" | "bottom-right" | "bottom-left";
  type?: TypeOptions;
}

export function showToast({
  title,
  message,
  duration,
  position,
  type,
}: ToastProps) {
  toast(
    <div className="p-4 bg-surface border-input-border border-2 flex flex-col gap-2">
      {title && <Heading>{title}</Heading>}
      <Text>{message}</Text>
    </div>,
    {
      autoClose: duration ?? 5000, // Default 5 seconds,
      position: position,
      type: type,
    },
  );
}
