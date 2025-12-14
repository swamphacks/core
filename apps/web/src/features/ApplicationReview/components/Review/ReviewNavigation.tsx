import TablerCheck from "~icons/tabler/check";
import TablerArrowLeft from "~icons/tabler/arrow-left";
import TablerArrowRight from "~icons/tabler/arrow-right";
import TablerRefresh from "~icons/tabler/refresh";
import { Button } from "@/components/ui/Button";

export type ButtonMode = "submit" | "completed-dirty" | "completed-clean";

interface ReviewNavigationProps {
  isLast: boolean;
  mode: ButtonMode;
  allowSubmit: boolean;
  currentIndex: number;
  back: () => void;
  onSubmit: () => void;
  next: () => void;
  reset: () => void;
}

export function ReviewNavigation({
  isLast,
  mode,
  allowSubmit,
  currentIndex,
  back,
  onSubmit,
  next,
  reset,
}: ReviewNavigationProps) {
  //   if (currentIndex > 0) {
  //     let backButton = (
  //       <Button
  //         variant="secondary"
  //         className="flex gap-1 items-center text-lg bg-button-secondary rounded-md py-2 px-3"
  //         onClick={back}
  //       >
  //         <TablerArrowLeft />
  //         Back
  //       </Button>
  //     );
  //   } else <div />;

  const backButton =
    currentIndex > 0 ? (
      <Button
        variant="secondary"
        className="flex gap-1 items-center text-lg rounded-md py-2 px-3"
        onClick={back}
      >
        <TablerArrowLeft />
        Back
      </Button>
    ) : (
      <div />
    );

  const buttons = (() => {
    switch (mode) {
      case "submit":
        return (
          <Button
            isDisabled={!allowSubmit}
            className={`flex gap-1 items-center text-lg rounded-md py-2 px-3 ${
              !allowSubmit ? "opacity-30" : ""
            }`}
            variant={isLast ? "success" : "primary"}
            onClick={onSubmit}
          >
            {isLast ? "Finish" : "Submit and Proceed"}
            {isLast ? <TablerCheck /> : <TablerArrowRight />}
          </Button>
        );
      case "completed-clean":
        return (
          <Button
            variant={isLast ? "success" : "secondary"}
            className="flex gap-1 items-center text-lg rounded-md py-2 px-3"
            onClick={next}
          >
            {isLast ? "Finish" : "Next"}
            {isLast ? <TablerCheck /> : <TablerArrowRight />}
          </Button>
        );
      case "completed-dirty":
        return (
          <div className="flex flex-row gap-4">
            <Button
              variant="icon"
              className="flex gap-1 items-center text-lg rounded-md py-2 px-3"
              onClick={reset}
            >
              <TablerRefresh />
            </Button>
            <Button
              variant={isLast ? "success" : "primary"}
              isDisabled={!allowSubmit}
              className="flex gap-1 items-center text-lg rounded-md py-2 px-3"
              onClick={onSubmit}
            >
              {isLast ? "Save and Finish" : "Save and Proceed"}
              {isLast ? <TablerCheck /> : <TablerArrowRight />}
            </Button>
          </div>
        );
    }
  })();

  return (
    <div className="w-full flex flex-row justify-between mt-4">
      {backButton}
      {buttons}
    </div>
  );
}
