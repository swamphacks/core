import TablerCalendarDue from "~icons/tabler/calendar-due";
import TablerLocation from "~icons/tabler/location";
import TablerInfoCircle from "~icons/tabler/info-circle";
import { EventBadge } from "./EventBadge";
import { EventButton } from "./EventButton";
import { Link } from "react-aria-components";

import imageFile from "./placeholder.jpg";
import type applicationStatus from "../applicationStatus";
import { Separator } from "@/components/ui/Seperator";
import { Card } from "@/components/ui/Card";

export interface EventCardProps {
  eventId: string;
  status: keyof typeof applicationStatus;
  title: string;
  description: string;
  date: string;
  location: string;
  banner: string | null;
}

const EventCard = ({
  eventId,
  status,
  title,
  description,
  date,
  location,
  banner,
}: EventCardProps) => {
  return (
    <Card className="border">
      <div className="w-full">
        <img
          className="w-full h-40 object-cover rounded-t-md"
          src={banner ?? imageFile}
          alt={`${title} Image`}
        />
      </div>
      <div className="p-3 space-y-1">
        <div className="flex justify-between">
          <p className="text-text-main font-semibold text-xl font-eventCard-title">
            {title}
          </p>
          <EventBadge status={status} />
        </div>
        <div className="mt-2 w-96">
          <p className="text-text-secondary text-sm text-wrap">{description}</p>
        </div>
        <div className="flex text-sm justify-around items-center text-text-main pt-4">
          <div className="flex flex-col items-center gap-1">
            <TablerCalendarDue />
            <p>{date}</p>
          </div>
          <Separator />
          <div className="flex flex-col items-center gap-1">
            <TablerLocation />
            <p>{location}</p>
          </div>
          <Separator />
          <div className="flex flex-col items-center gap-1">
            <TablerInfoCircle />
            <Link
              href={`/events/${eventId}`}
              className="underline text-text-link"
            >
              Event Details
            </Link>
          </div>
        </div>

        {status === "accepted" ? (
          <div className="flex gap-2">
            <EventButton
              className="w-full mt-4"
              status={status}
              eventId={eventId}
            />
            <EventButton
              className="w-full mt-4"
              status="notGoing"
              text="Not Going"
              eventId={eventId}
            />
          </div>
        ) : (
          <EventButton
            className="w-full mt-4"
            status={status}
            eventId={eventId}
          />
        )}
      </div>
    </Card>
  );
};

export { EventCard };
