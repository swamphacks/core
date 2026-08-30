import { Badge } from "@/components/ui/Badge";
import { Card } from "@/components/ui/Card";
import { cn } from "@/utils/cn";
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
  isSelected: boolean;
  onSelect: (campaign: EmailCampaign) => void;
}

export function CampaignCard({
  campaign,
  isSelected,
  onSelect,
}: CampaignCardProps) {
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
