import TablerCalendarDue from "~icons/tabler/calendar-due";
import TablerLocation from "~icons/tabler/location";
import TablerInfoCircle from "~icons/tabler/info-circle";
import { EventBadge } from "./EventBadge";
import { EventButton } from "./EventButton";

import imageFile from "./placeholder.jpg";
import type applicationStatus from "../applicationStatus";
import { Separator } from "@/components/ui/Seperator";
import { Card } from "@/components/ui/Card";

interface EventCardProps {
  status: keyof typeof applicationStatus;
  title: string;
  description: string;
  date: string;
  location: string;
}

const EventCard = ({
  status,
  title,
  description,
  date,
  location,
}: EventCardProps) => {
  return (
    <Card>
      <div className="w-full">
        <img
          className="w-full h-40 object-cover rounded-t-md"
          src={imageFile}
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
            {/* TODO: what does this Event Details button do? Display a model or go to new page? */}
            <p>Event Details</p>
          </div>
        </div>

        <EventButton className="w-full mt-4" status={status} />
      </div>
    </Card>
  );
};

export { EventCard };
