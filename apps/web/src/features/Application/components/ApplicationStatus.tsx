import Loading from "@/components/Loading";
import { useApplication } from "@/features/Application/hooks/useApplication";
import { Heading } from "react-aria-components";
import TablerUserCode from "~icons/tabler/user-code";
import TablerUsersGroup from "~icons/tabler/users-group";

interface ApplicationStatusProps {
  eventId: string;
}

export default function ApplicationStatus({ eventId }: ApplicationStatusProps) {
  const application = useApplication(eventId);

  const renderStatus = () => {
    if (!application.data) return null;

    const status = application.data["status"]["application_status"];

    switch (status) {
      case "submitted":
      case "under_review":
        return (
          <div>
            <div className="mb-5">
              <p className="text-lg">
                Your application is currently{" "}
                <b className="underline text-badge-text-under-review">
                  under review
                </b>
                .
              </p>
              <p className="text-lg">
                If you have any questions or concerns, please contact us at{" "}
                <span className="underline">tech@swamphacks.com</span>
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
          <HackerProfile
            data={JSON.parse(atob(application.data["application"]))}
          />

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
      )}
    </main>
  );
}

function HackerProfile({ data }: { data: any }) {
  return (
    <div className="border border-input-border rounded-md p-3 w-fit max-w-150">
      <div className="mb-3">
        <p className="text-lg flex gap-2 items-center">
          <TablerUserCode /> Hacker Profile
        </p>
      </div>

      <div className="space-y-3">
        <HackerProfileField
          label="Name"
          value={`${data["firstName"]} ${data["lastName"]}`}
        />
        <HackerProfileField label="School" value={data["school"]} />

        <HackerProfileField
          label="Major(s)"
          value={data["majors"].join(", ")}
        />

        <HackerProfileField
          label="Graduation Year"
          value={data["graduationYear"]}
        />

        <HackerProfileField label="Shirt Size" value={data["shirtSize"]} />

        <HackerProfileField label="Email" value={data["preferredEmail"]} />

        <HackerProfileField
          label="University Email"
          value={data["preferredEmail"]}
        />

        <HackerProfileField label="Github" value={data["github"]} />

        <HackerProfileField label="LinkedIn" value={data["linkedin"]} />
      </div>
    </div>
  );
}

function HackerProfileField({
  label,
  value,
}: {
  label: string;
  value: string;
}) {
  return (
    <div>
      <span>{label}:</span>
      <p className="text-text-secondary">{value}</p>
    </div>
  );
}
