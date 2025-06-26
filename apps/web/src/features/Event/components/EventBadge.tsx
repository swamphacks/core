import { Badge, badge, type BadgeProps } from "@/components/ui/Badge";
import { tv } from "tailwind-variants";
import applicationStatus from "../applicationStatus";

type ApplicationStatusTypes = keyof typeof applicationStatus;

/*
  This is how we go from this object

  {
    rejected: {
      className: "bg-badge-bg-rejected text-badge-text-rejected",
      text: "Rejected",
      icon: null,
    },
  }
  
  to this

  {
    rejected: "bg-badge-bg-rejected text-badge-text-rejected",
  }

  which we can use as tailwind variants.
*/
const applicationStatusVariants = Object.fromEntries(
  Object.entries(applicationStatus).map(([key, value]) => [
    key,
    value.className,
  ]),
) as {
  [K in ApplicationStatusTypes]: string;
};

const eventBadge = tv({
  extend: badge,
  variants: {
    type: {}, // override the `type` variant in the Badge component
    status: applicationStatusVariants,
  },
});

interface EventBadgeProps extends Omit<BadgeProps, "type"> {
  status: ApplicationStatusTypes;
}

const EventBadge = ({ status: statusProp, size }: EventBadgeProps) => {
  const eventBadgeClassname = eventBadge({ status: statusProp, size });
  const status = applicationStatus[statusProp];

  if (!status) {
    console.error(
      `Incorrect status prop passed to EventBadge component: ${statusProp}`,
    );
    return null;
  }

  return (
    <Badge className={eventBadgeClassname}>
      {statusProp && status.icon?.()}
      {status.text}
    </Badge>
  );
};

export { EventBadge };
