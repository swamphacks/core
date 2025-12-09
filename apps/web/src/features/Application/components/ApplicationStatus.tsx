import Loading from "@/components/Loading";
import { Button } from "@/components/ui/Button";
import { useMyApplication } from "@/features/Application/hooks/useMyApplication";
import { EventBadge } from "@/features/Event/components/EventBadge";
import { Heading } from "react-aria-components";
import TablerUserCode from "~icons/tabler/user-code";
import TablerUsersGroup from "~icons/tabler/users-group";
import TablerDownload from "~icons/tabler/download";
import { api } from "@/lib/ky";
import { EventButton } from "@/features/Event/components/EventButton";

interface ApplicationStatusProps {
  eventId: string;
}

export default function ApplicationStatus({ eventId }: ApplicationStatusProps) {
  const application = useMyApplication(eventId);

  const renderStatus = () => {
    if (!application.data) return null;

    const status = application.data["status"]["application_status"];

    switch (status) {
      case "submitted":
      case "under_review":
        return (
          <div>
            <div className="mb-5">
              <p className="text-lg flex items-center gap-2">
                Your application is{" "}
                <EventBadge status="underReview" className="text-sm" />
              </p>
              <p className="my-1 text-text-secondary">
                If you have any questions or concerns, please contact us at{" "}
                <span className="underline hover:text-text-main">
                  <a href="mailto:tech@swamphacks.com">tech@swamphacks.com</a>
                </span>
              </p>
            </div>
          </div>
        );
      case "rejected":
        return (
          <div>
            <div className="mb-1">
              <p className="text-lg flex items-center gap-2">
                Your application is{" "}
                <EventBadge status="rejected" className="text-sm" />
              </p>
              <p className="my-1 text-text-secondary">
                We couldn't accomodate all the applications this year. If you
                are still interested, please join the waitlist!{" "}
              </p>
            </div>
            <div>
              <EventButton
                className="w-1/2 mt-4 mb-5"
                status={status}
                eventId={eventId}
              />
            </div>
          </div>
        );
      case "accepted":
        return (
          <div>
            <div className="mb-5">
              <p className="text-lg flex items-center gap-2">
                Your application is{" "}
                <EventBadge status="accepted" className="text-sm" />
              </p>
              <p className="my-1 text-text-secondary">
                Congratulations! We would love to see you at Swamphacks XI{" "}
              </p>
            </div>
          </div>
        );
      default:
        break;
    }

    return null;
  };

  if (application.isLoading) {
    return <Loading />;
  }

  if (!application.data || !application.data["application"]) {
    return (
      <main>
        <Heading className="text-1xl lg:text-2xl font-semibold mb-6">
          Application Status
        </Heading>

        <p>
          Unable to load application status. Please contact us through Discord
          or tech@swamphacks.com
        </p>
      </main>
    );
  }

  return (
    <main>
      <Heading className="text-xl lg:text-2xl font-semibold mb-6">
        Application Status
      </Heading>

      {renderStatus()}

      {application.data && application.data["application"] && (
        <div className="flex gap-3 flex-wrap">
          <div className="max-w-100 w-fit">
            <HackerProfile
              applicationData={JSON.parse(
                atob(application.data["application"]),
              )}
              eventId={eventId}
            />
          </div>

          <div className="max-w-100">
            <div className="border border-input-border rounded-md p-3 w-fit max-w-150">
              <div className="mb-3">
                <p className="text-lg flex gap-2 items-center">
                  <TablerUsersGroup /> Team Info
                </p>
              </div>

              <p className="text-text-secondary">
                Coming soon! You will be notified when team formation begins.
              </p>
            </div>
          </div>
        </div>
      )}
    </main>
  );
}

interface HackerProfileProps {
  applicationData: any;
  eventId: string;
}

function HackerProfile({ applicationData, eventId }: HackerProfileProps) {
  return (
    <div className="border border-input-border rounded-md p-3">
      <div className="mb-3">
        <p className="text-lg flex gap-2 items-center">
          <TablerUserCode /> Hacker Profile
        </p>
      </div>

      <div className="space-y-3">
        <HackerProfileField
          label="Name"
          value={`${applicationData["firstName"]} ${applicationData["lastName"]}`}
        />
        <HackerProfileField label="School" value={applicationData["school"]} />

        <HackerProfileField
          label="Major(s)"
          value={
            Array.isArray(applicationData["majors"])
              ? applicationData["majors"].join(", ")
              : applicationData["majors"].split(",").join(", ")
          }
        />

        <HackerProfileField
          label="Graduation Year"
          value={applicationData["graduationYear"]}
        />

        <HackerProfileField
          label="Shirt Size"
          value={applicationData["shirtSize"]}
        />

        <HackerProfileField
          label="Preferred Email"
          value={applicationData["preferredEmail"]}
        />

        <HackerProfileField
          label="University Email"
          value={applicationData["universityEmail"]}
        />

        <HackerProfileField
          isUrl
          label="Github"
          value={applicationData["github"]}
        />

        <HackerProfileField
          isUrl
          label="LinkedIn"
          value={applicationData["linkedin"]}
        />

        <div className="mt-8">
          <DownloadResume
            eventId={eventId}
            hackerFullname={`${applicationData["firstName"]} ${applicationData["lastName"]}`}
          />
        </div>
      </div>
    </div>
  );
}

interface HackerProfileFieldProps {
  label: string;
  value: string;
  isUrl?: boolean;
}

function HackerProfileField({
  label,
  value,
  isUrl = false,
}: HackerProfileFieldProps) {
  return (
    <div>
      <span>{label}:</span>
      {isUrl ? (
        <p className="truncate">
          <a
            href={value}
            target="_blank"
            className="text-text-secondary underline"
          >
            {value}
          </a>
        </p>
      ) : (
        <p className="text-text-secondary text-wrap">{value}</p>
      )}
    </div>
  );
}

interface DownloadResumeProps {
  eventId: string;
  hackerFullname: string;
}

function DownloadResume({ eventId, hackerFullname }: DownloadResumeProps) {
  const handleDownload = async () => {
    const res = await api.get(`events/${eventId}/application/download-resume`);
    const presignedUrl: string = await res.json();

    const fileRes = await fetch(presignedUrl);
    const blob = await fileRes.blob();

    const url = URL.createObjectURL(blob);

    const link = document.createElement("a");
    link.href = url;
    link.download = `${hackerFullname}_Resume.pdf`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);

    URL.revokeObjectURL(url);
  };

  return (
    <Button onClick={handleDownload} variant="secondary" className="w-full">
      <TablerDownload /> Download resume
    </Button>
  );
}
