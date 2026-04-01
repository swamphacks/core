import { Heading } from "react-aria-components";
import TablerBell from "~icons/tabler/bell";
import TeamInvitationCard from "./TeamInvitationCard";

export default function TeamInvitationSection() {
  return (
    <section>
      <Heading className="text-2xl text-text-main flex flex-row items-center gap-2 mb-4">
        <TablerBell className="w-5 h-5" /> Invitations
      </Heading>

      <div className="flex flex-col gap-4">
        <TeamInvitationCard />
      </div>
    </section>
  );
}
