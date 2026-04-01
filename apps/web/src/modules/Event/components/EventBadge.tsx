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

export const eventBadge = tv({
  extend: badge,
  variants: {
    type: {}, // override the `type` variant in the Badge component
    status: applicationStatusVariants,
  },
});

interface EventBadgeProps extends Omit<BadgeProps, "type"> {
  status: ApplicationStatusTypes;
}

const EventBadge = ({
  status: statusProp,
  size,
  border,
  className,
}: EventBadgeProps) => {
  const eventBadgeClassname = eventBadge({
    status: statusProp,
    size,
    border,
    className,
  });
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
      {BadgeIcon && <BadgeIcon width="1.2em" height="1.2em" />}
      <span
        style={{
          marginTop: "0.5px",
        }}
      >
        {status.text}
      </span>
    </Badge>
  );
};

export { EventBadge };
