import { Badge, badge, type BadgeProps } from "@/components/ui/Badge";
import { tv } from "tailwind-variants";
import applicationStatus from "../applicationStatus";

type ApplicationStatusTypes = keyof typeof applicationStatus;

/*
 * Transforms the `applicationStatus` object into a flattened variant mapping
 * where each status key maps to its corresponding className.
 *
 * For example, it converts:
 *
 * {
 *    rejected: { className: "bg-badge-bg-rejected text-badge-text-rejected", ... },
 * }
 *
 * to:
 *
 * {
 *    rejected: "bg-badge-bg-rejected text-badge-text-rejected",
 * }
 */
const applicationStatusVariants = Object.fromEntries(
  Object.entries(applicationStatus).map(([key, value]) => [
    key,
    value.className,
  ]),
) as {
  [K in ApplicationStatusTypes]: (typeof applicationStatus)[K]["className"];
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
  const BadgeIcon = status?.icon;

  if (!status) {
    console.error(
      `Incorrect status prop passed to EventBadge component: ${statusProp}`,
    );
    return null;
  }

  return (
    <Badge className={eventBadgeClassname}>
      {BadgeIcon && <BadgeIcon />}
      {status.text}
    </Badge>
  );
};

export { EventBadge };
