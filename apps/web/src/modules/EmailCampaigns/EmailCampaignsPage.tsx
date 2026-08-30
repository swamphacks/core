import { Button } from "@/components/ui/Button";
import { cn } from "@/utils/cn";
import { Spinner } from "@/components/ui/Spinner";
import { useHackathon } from "@/modules/Hackathon/hooks/useHackathon";
import { useState } from "react";
import { CampaignCard } from "./CampaignCard";
import { CampaignFormPanel } from "./CampaignFormPanel";
import {
  useEmailCampaigns,
  type EmailCampaign,
} from "./hooks/useEmailCampaigns";

export default function EmailCampaignsPage() {
  const { data: hackathon, isPending: hackathonPending } = useHackathon();
  const hackathonId = hackathon?.id ?? "";

  const {
    data: campaigns,
    isPending,
    isError,
  } = useEmailCampaigns(hackathonId);

  // null = panel closed, "new" = create, otherwise the campaign being viewed.
  const [selected, setSelected] = useState<EmailCampaign | "new" | null>(null);
  const isPanelOpen = selected !== null;

  // Set when the panel is opened from the card's Reschedule action, so it can
  // land straight on the date picker instead of the form.
  const [openScheduleOnMount, setOpenScheduleOnMount] = useState(false);

  function openCampaign(campaign: EmailCampaign) {
    setOpenScheduleOnMount(false);
    setSelected(campaign);
  }

  function rescheduleCampaign(campaign: EmailCampaign) {
    setOpenScheduleOnMount(true);
    setSelected(campaign);
  }

  if (hackathonPending || (hackathonId && isPending)) {
    return (
      <div className="flex justify-center p-12">
        <Spinner />
      </div>
    );
  }

  return (
    // One animated property drives the whole layout: the panel track grows from
    // 0 to 610px (586 panel + 24 gutter) and the list column follows for free,
    // so the two panes cannot drift out of sync.
    <div
      className={cn(
        "mx-auto grid w-full max-w-[1200px] transition-[grid-template-columns] duration-300 motion-reduce:transition-none",
        isPanelOpen
          ? "ease-out [grid-template-columns:minmax(0,1fr)_610px]"
          : "ease-in [grid-template-columns:minmax(0,1fr)_0px]",
      )}
    >
      <div className="min-w-0 px-8 py-6">
        {/* Content stays 511px wide and is centred in whatever space is left. */}
        <div className="mx-auto w-full max-w-[511px]">
          <header className="mb-6 flex flex-col gap-2">
            <h1 className="text-2xl leading-8 font-medium text-zinc-900 dark:text-zinc-50">
              Email Campaigns
            </h1>
            <div className="flex items-center justify-between gap-4">
              <p className="text-lg leading-6 text-zinc-600 dark:text-zinc-400">
                Create and manage email campaigns.
              </p>
              <Button
                variant="primary"
                size="auto"
                className="h-[26px] w-[112px] shrink-0 justify-between rounded-[4px] bg-[#2b7fff] px-3 text-xs leading-4 font-medium"
                onPress={() => {
                  setOpenScheduleOnMount(false);
                  setSelected("new");
                }}
              >
                <span>Create Email</span>
                <span aria-hidden>+</span>
              </Button>
            </div>
          </header>

          {isError && (
            <p className="text-red-600">
              Could not load campaigns. Refresh to try again.
            </p>
          )}

          {campaigns && campaigns.length === 0 && (
            <p className="text-text-secondary py-12 text-center">
              No campaigns yet. Create one to email a group of users.
            </p>
          )}

          <div className="flex flex-col gap-6">
            {campaigns?.map((campaign) => (
              <CampaignCard
                key={campaign.id}
                campaign={campaign}
                hackathonId={hackathonId}
                isSelected={selected !== "new" && selected?.id === campaign.id}
                onSelect={openCampaign}
                onReschedule={rescheduleCampaign}
              />
            ))}
          </div>
        </div>
      </div>

      {/* Inline panel, not an overlay: the list stays readable beside it. */}
      <aside
        aria-hidden={!isPanelOpen}
        inert={!isPanelOpen ? true : undefined}
        className={cn(
          "sticky top-6 my-6 h-[calc(100vh-3rem)] overflow-hidden",
          !isPanelOpen && "pointer-events-none",
        )}
      >
        {/* Fixed width so nothing reflows while the track animates. */}
        <div
          className={cn(
            "ml-6 h-full w-[586px] overflow-y-auto border border-zinc-300 bg-[#f0f0f0] transition-[transform,opacity] duration-300 motion-reduce:transition-none dark:border-zinc-700 dark:bg-neutral-900",
            isPanelOpen
              ? "translate-x-0 opacity-100 ease-out"
              : "translate-x-4 opacity-0 ease-in",
          )}
        >
          {selected !== null && (
            <CampaignFormPanel
              key={selected === "new" ? "new" : selected.id}
              hackathonId={hackathonId}
              campaign={selected === "new" ? null : selected}
              openScheduleOnMount={openScheduleOnMount}
              onClose={() => setSelected(null)}
            />
          )}
        </div>
      </aside>
    </div>
  );
}
