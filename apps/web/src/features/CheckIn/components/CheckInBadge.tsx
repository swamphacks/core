import { Badge } from "@/components/ui/Badge";

export function CheckInBadge({ isCheckedIn }: { isCheckedIn: boolean }) {
  if (isCheckedIn) {
    return (
      <Badge className="bg-green-100 text-green-800 border border-green-300">
        Checked In
      </Badge>
    );
  }

  return (
    <Badge className="bg-yellow-100 text-yellow-800 border border-yellow-300">
      Not Checked In
    </Badge>
  );
}

type EventRole = "admin" | "attendee" | "applicant" | "staff";

export function CheckInEventRoleBadge({ role }: { role?: EventRole | null }) {
  if (!role) {
    return (
      <Badge className="bg-red-100 text-red-800 border border-red-300">
        No Event Role
      </Badge>
    );
  }

  switch (role) {
    case "admin":
      return (
        <Badge className="bg-purple-100 text-purple-800 border border-purple-300">
          Admin
        </Badge>
      );
    case "attendee":
      return (
        <Badge className="bg-green-100 text-green-800 border border-green-300">
          Attendee
        </Badge>
      );
    case "applicant":
      return (
        <Badge className="bg-gray-100 text-gray-800 border border-gray-300">
          Admin
        </Badge>
      );
    case "staff":
      return (
        <Badge className="bg-blue-100 text-blue-800 border border-blue-300">
          Admin
        </Badge>
      );
  }
}
