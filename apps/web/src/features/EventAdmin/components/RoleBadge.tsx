import { Badge } from "@/components/ui/Badge";

interface Props {
  role: "admin" | "staff" | "attendee" | "applicant";
}

const RoleBadge = ({ role }: Props) => {
  switch (role) {
    case "admin":
      return (
        <Badge size="md" className="bg-red-700 text-text-main px-4">
          Admin
        </Badge>
      );
    case "staff":
      return (
        <Badge size="md" className="bg-teal-600 text-text-main px-4">
          Staff
        </Badge>
      );
    case "attendee":
      return (
        <Badge size="md" className="bg-blue-600 text-text-main px-4">
          Attendee
        </Badge>
      );
    case "applicant":
      return (
        <Badge size="md" className="bg-yellow-600 text-text-main px-4">
          Applicant
        </Badge>
      );
  }
};

export default RoleBadge;
