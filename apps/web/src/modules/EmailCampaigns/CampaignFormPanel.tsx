import { Button } from "@/components/ui/Button";
import { Label } from "@/components/ui/Field";
import { cn } from "@/utils/cn";
import { MultiSelect } from "@/components/ui/MultiSelect";
import { TextField } from "@/components/ui/TextField";
import TablerChevronLeft from "~icons/tabler/chevron-left";
import { useState } from "react";
import {
  useCreateEmailCampaign,
  useUpdateEmailCampaign,
  type CampaignFormat,
  type EmailCampaign,
  type RecipientType,
} from "./hooks/useEmailCampaigns";

/** Only groups the API can actually resolve. Values match the recipient_type enum. */
const RECIPIENT_OPTIONS: { value: RecipientType; label: string }[] = [
  { value: "admins", label: "Admins" },
  { value: "staff", label: "Staff" },
  { value: "visitors", label: "Visitors" },
  { value: "accepted_applicants", label: "Accepted Applicants" },
  { value: "waitlisted_applicants", label: "Waitlisted Applicants" },
  { value: "rejected_applicants", label: "Rejected Applicants" },
  { value: "interest_subscribers", label: "Interest Subscribers" },
];

interface CampaignFormPanelProps {
  hackathonId: string;
  campaign: EmailCampaign | null;
  onClose: () => void;
}

export function CampaignFormPanel({
  hackathonId,
  campaign,
  onClose,
}: CampaignFormPanelProps) {
  const isEdit = campaign !== null;

  const [recipients, setRecipients] = useState<RecipientType[]>(
    campaign?.recipient_types ?? [],
  );
  const [title, setTitle] = useState(campaign?.title ?? "");
  const [subject, setSubject] = useState(campaign?.subject ?? "");
  const [description, setDescription] = useState(campaign?.description ?? "");
  const [body, setBody] = useState(campaign?.body ?? "");
  const [formatValue, setFormatValue] = useState<CampaignFormat | null>(
    campaign?.format ?? null,
  );

  const createCampaign = useCreateEmailCampaign(hackathonId);
  const updateCampaign = useUpdateEmailCampaign(hackathonId);

  const isSaving = createCampaign.isPending || updateCampaign.isPending;
  const saveError = createCampaign.error ?? updateCampaign.error;

  // The API rejects blank values, so mirror its requirements before sending.
  const canSave =
    recipients.length > 0 &&
    formatValue !== null &&
    title.trim() !== "" &&
    subject.trim() !== "" &&
    body.trim() !== "";

  async function handleSave() {
    if (formatValue === null) return;

    const payload = {
      title: title.trim(),
      subject: subject.trim(),
      body,
      format: formatValue,
      recipientTypes: recipients,
      description: description.trim() || undefined,
    };

    if (isEdit) {
      await updateCampaign.mutateAsync({
        campaignId: campaign.id,
        data: payload,
      });
    } else {
      await createCampaign.mutateAsync(payload);
    }
    onClose();
  }

  return (
    <div className="flex h-full flex-col gap-4 overflow-y-auto px-8 py-6 [&_input]:min-h-0 [&_input]:rounded-[4px] [&_input]:border-zinc-300 [&_input]:px-3 [&_input]:py-2 [&_input]:text-sm [&_input]:leading-5 [&_label]:text-sm [&_label]:leading-5 [&_label]:font-normal [&_label]:text-zinc-600 [&_label>span]:text-sm dark:[&_input]:border-zinc-700 dark:[&_label]:text-zinc-400">
      <button
        type="button"
        onClick={onClose}
        className="flex cursor-pointer items-center gap-1.5 self-start text-sm leading-5 text-[#2b7fff] hover:underline"
      >
        <TablerChevronLeft className="size-4" />
        Return
      </button>

      <h2 className="text-2xl leading-8 font-medium text-zinc-900 dark:text-zinc-50">
        {isEdit ? "Edit an Email Campaign" : "Create an Email Campaign"}
      </h2>

      <TextField
        label="Campaign Title"
        name="title"
        isRequired
        type="text"
        value={title}
        onChange={setTitle}
        placeholder="Internal name, e.g. Welcome to SwampHacks XII"
      />

      <TextField
        label="Campaign Description"
        name="description"
        type="text"
        value={description}
        onChange={setDescription}
        placeholder="Internal optional note for organizers"
      />

      {/* Everything below this line is delivered to recipients. */}
      <div className="my-2 border-t border-zinc-300 dark:border-zinc-600" />

      <MultiSelect
        name="recipientTypes"
        label="Email Recipients"
        isRequired
        options={RECIPIENT_OPTIONS}
        value={RECIPIENT_OPTIONS.filter((o) => recipients.includes(o.value))}
        onChange={(selected) =>
          setRecipients(selected.map((o) => o.value as RecipientType))
        }
      />

      <TextField
        label="Email Subject"
        name="subject"
        isRequired
        type="text"
        value={subject}
        onChange={setSubject}
        placeholder="What recipients see in their inbox"
        description="Shown to recipients. The title above is only for organizers."
      />

      <div className="flex flex-col gap-0.5">
        <Label
          id="campaign-format-label"
          isRequired
          className="text-sm leading-5 font-normal text-zinc-600 dark:text-zinc-400"
        >
          Email Format
        </Label>
        {/* Two options only, so a segmented control beats a dropdown. */}
        <div
          role="radiogroup"
          aria-labelledby="campaign-format-label"
          className="inline-flex w-fit gap-0.5 rounded-[4px] border border-zinc-300 p-0.5 dark:border-zinc-700"
        >
          {(
            [
              { value: "html", label: "HTML" },
              { value: "text", label: "Plain text" },
            ] as { value: CampaignFormat; label: string }[]
          ).map((option) => (
            <button
              key={option.value}
              type="button"
              role="radio"
              aria-checked={formatValue === option.value}
              onClick={() => setFormatValue(option.value)}
              className={cn(
                "cursor-pointer rounded-[3px] px-3 py-1 text-sm leading-5 transition-colors",
                formatValue === option.value
                  ? "bg-[#2b7fff] font-medium text-white"
                  : "text-zinc-600 hover:bg-zinc-200 dark:text-zinc-400 dark:hover:bg-zinc-800",
              )}
            >
              {option.label}
            </button>
          ))}
        </div>
      </div>

      <div className="flex flex-col gap-1">
        <Label
          isRequired
          className="text-sm leading-5 font-normal text-zinc-600 dark:text-zinc-400"
        >
          Email Body
        </Label>
        <textarea
          value={body}
          onChange={(e) => setBody(e.target.value)}
          rows={14}
          style={{ height: 349 }}
          placeholder={
            formatValue === "html"
              ? "<h1>Hello</h1><p>Email contents</p>"
              : formatValue === "text"
                ? "Email contents"
                : "Choose a format above to start writing"
          }
          className="bg-input-bg text-text-main w-full resize-y rounded-[4px] border border-zinc-300 px-[13px] py-3 text-sm leading-5 outline-0 dark:border-zinc-700"
        />
      </div>

      {saveError && (
        <p className="text-sm text-red-600">
          Could not save. Check the fields and try again.
        </p>
      )}

      <div className="flex gap-2 pb-6">
        <Button
          variant="primary"
          size="auto"
          className="h-[26px] rounded-[4px] bg-[#2b7fff] px-[18px] text-xs leading-4 font-medium"
          onPress={handleSave}
          isDisabled={!canSave || isSaving}
        >
          {isSaving ? "Saving..." : isEdit ? "Update" : "Save Draft"}
        </Button>
      </div>
    </div>
  );
}
