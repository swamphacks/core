import TablerCalendarDue from "~icons/tabler/calendar-due";
import TablerLocation from "~icons/tabler/location";
import TablerInfoCircle from "~icons/tabler/info-circle";
import { EventBadge } from "./EventBadge";
import { EventButton } from "./EventButton";
import { DialogTrigger, Link } from "react-aria-components";

import imageFile from "./placeholder.jpg";
import type applicationStatus from "../applicationStatus";
import { Separator } from "@/components/ui/Seperator";
import { Card } from "@/components/ui/Card";
import EventDetailsModal from "./EventDetailsModal";

export interface EventCardProps {
  eventId: string;
  status: keyof typeof applicationStatus;
  title: string;
  description: string;
  date: string;
  location: string;
  banner: string | null;

  // Used for EventDetailsModal
  application_close: Date;
  application_open: Date;
  end_time: Date;
  start_time: Date;
  location_url?: string;
  decision_release?: Date;
  rsvp_deadline?: Date;
  website_url?: string;
}

const EventCard = ({
  eventId,
  status,
  title,
  description,
  date,
  location,
  banner,
  website_url,
  application_close,
  application_open,
  end_time,
  start_time,
  location_url,
  decision_release,
  rsvp_deadline,
}: EventCardProps) => {
  //FIX: Sort of a crude calculation to determine is description will be clamped
  const descriptionLength = description.length;
  const isClamped = descriptionLength > 105;
  const currentTime: Date = new Date();

  //TODO: Optimize so that multiple modals aren't created for each card

  return (
    <Card className="border">
      <div className="w-full">
        <img
          className="w-full h-40 object-cover rounded-t-md"
          src={banner ?? imageFile}
          alt={`${title} Image`}
        />
      </div>
      <div className="p-3">
        <div className="space-y-1">
          <div className="flex justify-between">
            <p className="text-text-main font-semibold text-xl font-eventCard-title">
              {title}
            </p>
            <EventBadge status={status} />
          </div>
          <div className="mt-2 w-full">
            <p className="text-text-secondary text-sm break-words overflow-hidden text-ellipsis line-clamp-2">
              {description}
            </p>{" "}
            {isClamped && (
              <DialogTrigger>
                <Link className="underline text-text-link hover:cursor-pointer text-sm">
                  Read More
                </Link>

                <EventDetailsModal
                  title={title}
                  description={description}
                  location={location}
                  application_close={application_close}
                  application_open={application_open}
                  end_time={end_time}
                  start_time={start_time}
                  location_url={location_url}
                  decision_release={decision_release}
                  rsvp_deadline={rsvp_deadline}
                  website_url={website_url}
                  id={eventId}
                />
              </DialogTrigger>
            )}
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
              <DialogTrigger>
                <Link className="underline text-text-link hover:cursor-pointer">
                  Event Details
                </Link>

                <EventDetailsModal
                  title={title}
                  description={description}
                  location={location}
                  application_close={application_close}
                  application_open={application_open}
                  end_time={end_time}
                  start_time={start_time}
                  location_url={location_url}
                  decision_release={decision_release}
                  rsvp_deadline={rsvp_deadline}
                  website_url={website_url}
                  id={eventId}
                />
              </DialogTrigger>
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
          ) : (status === "notApplied" &&
              currentTime.getTime() < application_open.getTime()) ||
            currentTime.getTime() > application_close.getTime() ? (
            <EventButton
              className="w-full mt-4"
              status={"upcoming"}
              eventId={eventId}
            />
          ) : (
            <EventButton
              className="w-full mt-4"
              status={status}
              eventId={eventId}
            />
          )}
        </div>
      </div>
    </Card>
  );
};

export { EventCard };
