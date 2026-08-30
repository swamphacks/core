import { Badge } from "@/components/ui/Badge";
import { Button } from "@/components/ui/Button";
import { Menu, MenuItem } from "@/components/ui/Menu";
import { Modal } from "@/components/ui/Modal";
import { MenuTrigger } from "react-aria-components";
import TablerDots from "~icons/tabler/dots-vertical";
import { campaignMenuItemClasses } from "./menuStyles";
import {
  useDeleteEmailCampaign,
  useUnscheduleEmailCampaign,
} from "./hooks/useEmailCampaigns";
import { Card } from "@/components/ui/Card";
import { cn } from "@/utils/cn";
import { useState } from "react";
import { format } from "date-fns";
import type { CampaignStatus, EmailCampaign } from "./hooks/useEmailCampaigns";

const STATUS_STYLES: Record<CampaignStatus, string> = {
  draft: "bg-zinc-100 border-zinc-400 text-zinc-600",
  scheduled: "bg-[#fef9c2] border-[#d08700] text-[#d08700]",
  sending: "bg-blue-100 border-blue-700 text-blue-700",
  sent: "bg-[#dcfce7] border-[#016630] text-[#016630]",
  failed: "bg-red-100 border-red-700 text-red-700",
};

const STATUS_LABELS: Record<CampaignStatus, string> = {
  draft: "Draft",
  scheduled: "Scheduled",
  sending: "Sending",
  sent: "Sent",
  failed: "Failed",
};

function formatDay(value: string) {
  return format(new Date(value), "yyyy-MM-dd");
}

interface CampaignCardProps {
  campaign: EmailCampaign;
  hackathonId: string;
  isSelected: boolean;
  onSelect: (campaign: EmailCampaign) => void;
  onReschedule: (campaign: EmailCampaign) => void;
}

export function CampaignCard({
  campaign,
  hackathonId,
  isSelected,
  onSelect,
  onReschedule,
}: CampaignCardProps) {
  const unscheduleCampaign = useUnscheduleEmailCampaign(hackathonId);
  const deleteCampaign = useDeleteEmailCampaign(hackathonId);
  const [confirmingDelete, setConfirmingDelete] = useState(false);

  // Matches the API: sent and sending campaigns stay on the record.
  const canDelete = campaign.status !== "sent" && campaign.status !== "sending";
  return (
    <Card
      role="button"
      tabIndex={0}
      onClick={() => onSelect(campaign)}
      onKeyDown={(e) => {
        if (e.key === "Enter" || e.key === " ") {
          e.preventDefault();
          onSelect(campaign);
        }
      }}
      className={cn(
        "focus-visible:ring-button-primary w-full max-w-none cursor-pointer rounded-[6px] border-[1.5px] border-[#d1d5dc] bg-transparent px-[17px] py-5 shadow-none transition-colors hover:border-zinc-400 focus-visible:ring-2 focus-visible:outline-none sm:max-w-none dark:border-zinc-700 dark:hover:border-zinc-500",
        isSelected && "border-[#2b7fff]",
      )}
    >
      <div className="flex items-start justify-between gap-3">
        <div className="flex min-w-0 items-center gap-[17px]">
          <h2 className="truncate text-xl leading-7 font-medium text-zinc-900 dark:text-zinc-50">
            {campaign.title}
          </h2>
          <Badge
            className={cn(
              "h-[21px] shrink-0 rounded-[18px] border px-2 py-0.5 text-xs leading-4 font-medium",
              STATUS_STYLES[campaign.status],
            )}
          >
            {STATUS_LABELS[campaign.status]}
          </Badge>
        </div>

        {/* The whole card is clickable, so keep menu interactions from opening it. */}
        <div
          onClick={(e) => e.stopPropagation()}
          onKeyDown={(e) => e.stopPropagation()}
          role="presentation"
          className="-mt-3 -mr-3.5 shrink-0"
        >
          <MenuTrigger>
            <Button
              variant="unstyled"
              size="auto"
              aria-label={`Actions for ${campaign.title}`}
              // Not variant="icon": that carries pressed:bg-black/10, and
              // MenuTrigger holds isPressed for as long as the menu is open, so
              // the tint never cleared. Own the states explicitly instead.
              className="cursor-pointer rounded-sm p-0.5 text-zinc-500 not-aria-expanded:hovered:bg-black/[7%] dark:text-zinc-400 dark:not-aria-expanded:hovered:bg-white/10"
            >
              <TablerDots className="size-4" />
            </Button>
            <Menu
              placement="bottom end"
              className="overflow-visible"
              popoverClassName="duration-75 min-w-[120px]"
            >
              {campaign.status === "scheduled" ? (
                <MenuItem
                  className={campaignMenuItemClasses}
                  onAction={() => onReschedule(campaign)}
                >
                  Reschedule
                </MenuItem>
              ) : null}
              {campaign.status === "scheduled" ? (
                <MenuItem
                  className={campaignMenuItemClasses}
                  onAction={() => unscheduleCampaign.mutate(campaign.id)}
                >
                  Unschedule
                </MenuItem>
              ) : null}
              {canDelete ? (
                <MenuItem
                  className={campaignMenuItemClasses}
                  onAction={() => setConfirmingDelete(true)}
                >
                  Delete
                </MenuItem>
              ) : null}
            </Menu>
          </MenuTrigger>

          <Modal
            size="md"
            isOpen={confirmingDelete}
            onOpenChange={setConfirmingDelete}
            title="Delete this campaign?"
          >
            <div className="flex flex-col gap-4">
              <p className="text-sm leading-5 text-zinc-700 dark:text-zinc-300">
                <strong>{campaign.title}</strong> will be permanently removed.
                This cannot be undone.
              </p>
              <div className="flex justify-end gap-2">
                <Button
                  variant="secondary"
                  size="sm"
                  onPress={() => setConfirmingDelete(false)}
                  isDisabled={deleteCampaign.isPending}
                >
                  Cancel
                </Button>
                <Button
                  variant="danger"
                  size="sm"
                  isDisabled={deleteCampaign.isPending}
                  onPress={async () => {
                    await deleteCampaign.mutateAsync(campaign.id);
                    setConfirmingDelete(false);
                  }}
                >
                  {deleteCampaign.isPending ? "Deleting..." : "Delete"}
                </Button>
              </div>
            </div>
          </Modal>
        </div>
      </div>

      {campaign.description && (
        <p className="mt-0.5 line-clamp-2 text-sm leading-5 text-zinc-900 dark:text-zinc-100">
          {campaign.description}
        </p>
      )}

      <p className="mt-0.5 text-sm leading-5 text-zinc-600 dark:text-zinc-400">
        Created {formatDay(campaign.created_at)}
        {campaign.status === "scheduled" && campaign.scheduled_at
          ? ` · Scheduled for ${formatDay(campaign.scheduled_at)}`
          : null}
        {campaign.status === "sent" && campaign.sent_at
          ? ` · Sent ${formatDay(campaign.sent_at)}`
          : null}
      </p>

      {campaign.status === "failed" && campaign.last_error && (
        <p className="mt-2 text-sm text-red-600">{campaign.last_error}</p>
      )}
    </Card>
  );
}
