// TODO: Remove the eslint-disable line and fix the fast refresh issue
import { tv } from "tailwind-variants";
import applicationStatus from "../applicationStatus";
import { Button, button, type ButtonProps } from "@/components/ui/Button";
import { cn } from "@/utils/cn";

type ApplicationStatusTypes = keyof typeof applicationStatus;

const applicationStatusVariants = Object.fromEntries(
  Object.entries(applicationStatus).map(([key, value]) => [
    key,
    value.button.className,
  ]),
) as {
  [K in ApplicationStatusTypes]: (typeof applicationStatus)[K]["button"]["className"];
};

export const eventButton = tv({
  extend: button,
  variants: {
    color: {},
    status: applicationStatusVariants,
  },
});

interface EventButtonProps extends ButtonProps {
  status: ApplicationStatusTypes;
  text?: string;
}

const EventButton = ({
  status: statusProp,
  className,
  text,
}: EventButtonProps) => {
  const eventButtonClassName = eventButton({
    status: statusProp,
    className,
  });

  return (
    <Button className={cn(eventButtonClassName, "font-semibold")}>
      {text || applicationStatus[statusProp].button.text}
    </Button>
  );
};

export { EventButton };
