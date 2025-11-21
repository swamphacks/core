import { Heading, Text } from "react-aria-components";
import { toast, type ToastIcon, type TypeOptions } from "react-toastify";

export interface ToastProps {
  title?: string;
  message: string;
  icon?: ToastIcon;
  duration?: number; // ms
  position?: "top-right" | "top-left" | "bottom-right" | "bottom-left";
  type?: TypeOptions;
}

export function showToast({
  title,
  message,
  duration,
  icon,
  position = "bottom-right",
  type,
}: ToastProps) {
  toast(
    <div className="flex flex-col px-3 text-text-main">
      {title && <Heading className="text-md font-semibold">{title}</Heading>}
      <Text>{message}</Text>
    </div>,
    {
      autoClose: duration ?? 5000, // Default 5 seconds,
      position: position,
      type: type,
      icon: icon,
    },
  );
}
