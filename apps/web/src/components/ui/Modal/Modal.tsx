import {
  Modal as RAC_Modal,
  ModalOverlay,
  Dialog,
  Heading,
} from "react-aria-components";
import { tv, type VariantProps } from "tailwind-variants";
import { type PropsWithChildren } from "react";

const modal = tv({
  base: "w-full max-h-screen sm:max-h-[90vh] overflow-hidden flex flex-col bg-surface text-left align-middle shadow-xs",
  variants: {
    size: {
      sm: "max-w-sm",
      md: "max-w-md",
      lg: "max-w-lg",
      xl: "max-w-xl",
      "2xl": "max-w-2xl",
      "3xl": "max-w-3xl",
    },
    padding: {
      none: "p-0",
      sm: "p-4",
      md: "p-6",
      lg: "p-8",
    },
    responsive: {
      // Standard modal behavior - rounded on all screens
      standard: "rounded-md",
      // Mobile-first - full screen on mobile, modal on desktop
      adaptive:
        "rounded-none sm:rounded-md h-full sm:h-auto max-h-screen sm:max-h-[90vh]",
      // Bottom sheet style - slides up from bottom on mobile
      sheet:
        "rounded-t-xl sm:rounded-md fixed bottom-0 sm:relative sm:bottom-auto max-h-[85vh]",
      // Full screen on mobile, centered modal on tablet+
      fullMobile:
        "rounded-none md:rounded-md h-screen md:h-auto w-screen md:w-full",
    },
  },
  defaultVariants: {
    size: "md",
    padding: "md",
    responsive: "standard",
  },
});

const overlay = tv({
  base: "fixed inset-0 z-10 bg-black/15 backdrop-blur-sm",
  variants: {
    responsive: {
      standard:
        "overflow-y-auto flex min-h-full items-center justify-center p-4 text-center",
      adaptive:
        "overflow-y-auto flex min-h-full items-end sm:items-center justify-center p-0 sm:p-4 text-center",
      sheet:
        "overflow-y-auto flex min-h-full items-end sm:items-center justify-center p-0 sm:p-4 text-center",
      fullMobile:
        "overflow-y-auto flex min-h-full items-center justify-center p-0 md:p-4 text-center",
    },
  },
  defaultVariants: {
    responsive: "standard",
  },
});

export type ModalProps = {
  title?: string;
  withCloseButton?: boolean;
  className?: string;
  minimumSheetHeight?: string;
  isDismissible?: boolean;
} & VariantProps<typeof modal>;

export function Modal({
  title,
  children,
  isDismissible = true,
  minimumSheetHeight = "65vh",
  size,
  padding,
  responsive,
  className,
}: PropsWithChildren<ModalProps>) {
  const getSheetClassnames = () => {
    if (responsive !== "sheet") return "";

    return `min-h-[${minimumSheetHeight}]`;
  };
  return (
    <ModalOverlay
      isDismissable={isDismissible}
      className={({ isEntering, isExiting }) => `
        ${overlay({ responsive })}
        ${isEntering ? "animate-in fade-in duration-200 ease-out" : ""}
        ${isExiting ? "animate-out fade-out duration-100 ease-in" : ""}
      `}
    >
      <RAC_Modal
        className={({ isEntering, isExiting }) => {
          const baseClasses = modal({ size, padding, responsive });
          const animationClasses =
            responsive === "sheet"
              ? (isEntering
                  ? "animate-in slide-in-from-bottom md:slide-in-from-bottom-0 duration-200 ease-out"
                  : "") +
                (isExiting
                  ? "animate-out slide-out-to-bottom md:slide-out-to-bottom-0 duration-100 ease-in"
                  : "")
              : (isEntering
                  ? "animate-in zoom-in-95 ease-out duration-200"
                  : "") +
                (isExiting
                  ? "animate-out zoom-out-95 ease-in duration-100"
                  : "");

          return `${baseClasses} ${className ?? ""} ${animationClasses} ${getSheetClassnames}`;
        }}
      >
        {title && (
          <Heading className=" text-text-primary text-lg">{title}</Heading>
        )}
        <Dialog className="outline-none relative overflow-y-auto">
          {children}
        </Dialog>
      </RAC_Modal>
    </ModalOverlay>
  );
}
