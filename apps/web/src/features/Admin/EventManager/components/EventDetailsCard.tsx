import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import type { Event } from "@/lib/openapi/types";
import {
  DialogTrigger,
  Group,
  Heading,
  Text,
  Tooltip,
  TooltipTrigger,
} from "react-aria-components";

import TablerTrash from "~icons/tabler/trash";
import TablerLocation from "~icons/tabler/location";
import TablerCalendarEventFilled from "~icons/tabler/calendar-event-filled";
import TablerFileExport from "~icons/tabler/file-export";
import TablerUsers from "~icons/tabler/users";
import TablerUserCog from "~icons/tabler/user-cog";
import { DeleteEventDialog } from "./DeleteEventDialog";
import { ManageEventStaffDialog } from "./ManageEventStaffDialog";

const prettyPrintDate = (input: string): string => {
  const date = new Date(input);
  return new Intl.DateTimeFormat("en-US", {
    month: "short",
    day: "numeric",
    hour: "numeric",
    minute: "2-digit",
    hour12: true,
  }).format(date);
};

interface EventDetailsCardProps {
  event: Event;
}

function EventDetailsCard({ event }: EventDetailsCardProps) {
  const {
    id,
    name,
    description,
    location,
    start_time,
    end_time,
    application_open,
    application_close,
    max_attendees,
  } = event;

  return (
    <Card className="p-4 flex flex-col gap-6 justify-between">
      <div className="flex flex-col gap-1">
        <Heading className="text-xl font-semibold">{name}</Heading>
        <Text className="text-text-secondary">
          {description || "No description provided"}
        </Text>

        <span className="flex flex-row items-center gap-1.5 text-text-secondary">
          <TablerLocation className="h-4 w-4" />
          {location ? location : "Unknown"}
        </span>
      </div>

      <div className="flex flex-col gap-1 text-sm text-text-secondary">
        <Text>Details:</Text>

        <span className="flex flex-row items-center gap-1.5">
          <TablerCalendarEventFilled className="h-4 w-4" /> Event Dates:{" "}
          {prettyPrintDate(start_time)} – {prettyPrintDate(end_time)}
        </span>

        <span className="flex flex-row items-center gap-1.5">
          <TablerFileExport className="h-4 w-4" /> Apply Dates:{" "}
          {prettyPrintDate(application_open)} –{" "}
          {prettyPrintDate(application_close)}
        </span>

        {max_attendees && (
          <span className="flex flex-row items-center gap-1.5">
            <TablerUsers className="h-4 w-4" /> Max Attendees: {max_attendees}
          </span>
        )}
      </div>

      <Group className="flex flex-row gap-3 mt-4">
        <DialogTrigger>
          <TooltipTrigger delay={250} closeDelay={250}>
            <Tooltip
              offset={5}
              className="bg-surface border-input-border border-2 flex justify-center items-center py-1 px-2 rounded-md"
            >
              <Text>Manage Users</Text>
            </Tooltip>

            <Button variant="secondary">
              <TablerUserCog />
            </Button>
          </TooltipTrigger>

          <ManageEventStaffDialog eventId={id} />
        </DialogTrigger>

        <DialogTrigger>
          <TooltipTrigger delay={250} closeDelay={250}>
            <Tooltip
              offset={5}
              className="bg-surface border-input-border border-2 flex justify-center items-center py-1 px-2 rounded-md"
            >
              <Text>Delete</Text>
            </Tooltip>

            <Button variant="danger">
              <TablerTrash />
            </Button>
          </TooltipTrigger>

          <DeleteEventDialog event={event} />
        </DialogTrigger>
      </Group>
    </Card>
  );
}

export { EventDetailsCard };
