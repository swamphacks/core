import { Modal } from "@/components/ui/Modal";
import { Heading } from "react-aria-components";
import TablerCalendarEventFilled from "~icons/tabler/calendar-event-filled";
import TablerFileSpark from "~icons/tabler/file-spark";
import TablerMailOpened from "~icons/tabler/mail-opened";
import TablerTicket from "~icons/tabler/ticket";
import TablerLocation from "~icons/tabler/location";
import TablerWorld from "~icons/tabler/world";

interface Props {
  title: string;
  description: string;
  application_close: Date;
  application_open: Date;
  end_time: Date;
  start_time: Date;
  location?: string;
  location_url?: string;
  decision_release?: Date;
  rsvp_deadline?: Date;
  website_url?: string;
  id: string;
}

const formatDate = (date: Date) =>
  date.toLocaleString(undefined, {
    year: "numeric",
    month: "short",
    day: "numeric",
    hour: "numeric",
    minute: "2-digit",
  });

const EventDetailsModal = ({
  title,
  description,
  application_close,
  application_open,
  end_time,
  start_time,
  location,
  location_url,
  website_url,
  decision_release,
  rsvp_deadline,
}: Props) => {
  return (
    <Modal size="xl" isDismissible>
      <div className="h-full w-full">
        <Heading className="text-text-main text-xl">{title}</Heading>

        <p className="text-text-secondary mt-4 whitespace-pre-wrap">
          {description}
        </p>

        <div className="mt-6 space-y-3 text-sm text-text-secondary">
          {/* Event Dates */}
          <div className="flex items-start gap-2 flex-wrap">
            <TablerCalendarEventFilled className="h-4 w-4 text-primary mt-0.5" />
            <div>
              <strong>Event:</strong>
              <div className="flex flex-col sm:flex-row sm:gap-4">
                <span>
                  <strong className="text-text-primary">Starts:</strong>{" "}
                  {formatDate(start_time)}
                </span>
                <span>
                  <strong className="text-text-primary">Ends:</strong>{" "}
                  {formatDate(end_time)}
                </span>
              </div>
            </div>
          </div>

          {/* Application Dates */}
          <div className="flex items-start gap-2 flex-wrap">
            <TablerFileSpark className="h-4 w-4 text-primary mt-0.5" />
            <div>
              <strong>Applications:</strong>
              <div className="flex flex-col sm:flex-row sm:gap-4">
                <span>
                  <strong className="text-text-primary">Opens:</strong>{" "}
                  {formatDate(application_open)}
                </span>
                <span>
                  <strong className="text-text-primary">Closes:</strong>{" "}
                  {formatDate(application_close)}
                </span>
              </div>
            </div>
          </div>

          {/* Optional Dates */}
          {decision_release && (
            <p className="flex items-center gap-2">
              <TablerMailOpened className="h-4 w-4 text-primary" />
              <strong>Decision Release:</strong> {formatDate(decision_release)}
            </p>
          )}

          {rsvp_deadline && (
            <p className="flex items-center gap-2">
              <TablerTicket className="h-4 w-4 text-primary" />
              <strong>RSVP Deadline:</strong> {formatDate(rsvp_deadline)}
            </p>
          )}

          {/* Location */}
          {location && (
            <p className="flex items-center gap-2 flex-wrap">
              <TablerLocation className="h-4 w-4 text-primary" />
              <strong>Location:</strong>{" "}
              {location_url ? (
                <a
                  href={location_url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-primary underline"
                >
                  {location}
                </a>
              ) : (
                location
              )}
            </p>
          )}

          {/* Website */}
          {website_url && (
            <p className="flex items-center gap-2 flex-wrap">
              <TablerWorld className="h-4 w-4 text-primary" />
              <strong>Website:</strong>{" "}
              <a
                href={website_url}
                target="_blank"
                rel="noopener noreferrer"
                className="text-primary underline break-all"
              >
                {website_url}
              </a>
            </p>
          )}
        </div>
      </div>
    </Modal>
  );
};

export default EventDetailsModal;
