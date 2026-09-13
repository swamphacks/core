import { Button } from "@/components/ui/Button";
import { DatePicker } from "@/components/ui/DatePicker";
import { Menu, MenuItem } from "@/components/ui/Menu";
import { Modal } from "@/components/ui/Modal";
import { MenuTrigger } from "react-aria-components";
import {
  fromDate,
  getLocalTimeZone,
  now,
  toCalendarDate,
  toCalendarDateTime,
  type CalendarDate,
  type CalendarDateTime,
} from "@internationalized/date";
import TablerChevronDown from "~icons/tabler/chevron-down";
import { format as formatDate } from "date-fns";
import { Label } from "@/components/ui/Field";
import { cn } from "@/utils/cn";
import { MultiSelect } from "@/components/ui/MultiSelect";
import { TextField } from "@/components/ui/TextField";
import TablerChevronLeft from "~icons/tabler/chevron-left";
import { useState } from "react";
import { campaignMenuItemClasses } from "./menuStyles";
import {
  useCreateEmailCampaign,
  useScheduleEmailCampaign,
  useSendEmailCampaign,
  useUnscheduleEmailCampaign,
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

/** Why a campaign can no longer be edited, matching the API's canEditCampaign. */
const READ_ONLY_REASON: Partial<Record<EmailCampaign["status"], string>> = {
  sending: "This campaign is being sent right now, so it cannot be changed.",
  sent: "This campaign has already been sent. Sent campaigns are kept as a record and cannot be edited.",
  failed:
    "This campaign failed to send and can no longer be edited. Delete it and create a new one.",
};

interface CampaignFormPanelProps {
  hackathonId: string;
  campaign: EmailCampaign | null;
  openScheduleOnMount?: boolean;
  onClose: () => void;
}

export function CampaignFormPanel({
  hackathonId,
  campaign,
  openScheduleOnMount = false,
  onClose,
}: CampaignFormPanelProps) {
  const isEdit = campaign !== null;
  const isReadOnly =
    campaign !== null &&
    campaign.status !== "draft" &&
    campaign.status !== "scheduled";

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
  const sendCampaign = useSendEmailCampaign(hackathonId);
  const scheduleCampaign = useScheduleEmailCampaign(hackathonId);
  const unscheduleCampaign = useUnscheduleEmailCampaign(hackathonId);

  const [confirmingSend, setConfirmingSend] = useState(false);
  const [scheduling, setScheduling] = useState(openScheduleOnMount);

  // survives unscheduling (the API reads a null scheduledAt as "leave it alone"),
  // so a draft can still be carrying a stale timestamp.
  const [scheduleAt, setScheduleAt] = useState<CalendarDateTime | null>(
    campaign?.status === "scheduled" && campaign.scheduled_at
      ? toCalendarDateTime(
          fromDate(new Date(campaign.scheduled_at), getLocalTimeZone()),
        )
      : null,
  );

  const scheduleIsPast =
    scheduleAt !== null &&
    scheduleAt.toDate(getLocalTimeZone()).getTime() <= Date.now();

  const isSaving =
    createCampaign.isPending ||
    updateCampaign.isPending ||
    sendCampaign.isPending ||
    scheduleCampaign.isPending ||
    unscheduleCampaign.isPending;
  const saveError =
    createCampaign.error ??
    updateCampaign.error ??
    sendCampaign.error ??
    scheduleCampaign.error ??
    unscheduleCampaign.error;

  // The API rejects blank values, so mirror its requirements before sending.
  const canSave =
    !isReadOnly &&
    recipients.length > 0 &&
    formatValue !== null &&
    title.trim() !== "" &&
    subject.trim() !== "" &&
    body.trim() !== "";

  /** Creates or updates, returning the stored campaign so send/schedule can chain off it. */
  async function persistCampaign() {
    if (formatValue === null) return null;

    const payload = {
      title: title.trim(),
      subject: subject.trim(),
      body,
      format: formatValue,
      recipientTypes: recipients,
      description: description.trim() || undefined,
    };

    if (isEdit) {
      return updateCampaign.mutateAsync({
        campaignId: campaign.id,
        data: payload,
      });
    }
    return createCampaign.mutateAsync(payload);
  }

  async function handleSaveDraft() {
    if (await persistCampaign()) onClose();
  }

  // Send always saves first: a new campaign has no id yet, and an edited one
  // must go out with what is on screen rather than what was last stored.
  async function handleSendNow() {
    const stored = await persistCampaign();
    if (!stored) return;
    await sendCampaign.mutateAsync(stored.id);
    setConfirmingSend(false);
    onClose();
  }

  async function handleUnschedule() {
    if (!campaign) return;
    await unscheduleCampaign.mutateAsync(campaign.id);
    onClose();
  }

  async function handleSchedule() {
    if (!scheduleAt || scheduleIsPast) return;
    const stored = await persistCampaign();
    if (!stored) return;
    await scheduleCampaign.mutateAsync({
      campaignId: stored.id,
      scheduledAt: scheduleAt.toDate(getLocalTimeZone()).toISOString(),
    });
    setScheduling(false);
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
        {isReadOnly
          ? "Email Campaign"
          : isEdit
            ? "Edit an Email Campaign"
            : "Create an Email Campaign"}
      </h2>

      {isReadOnly && campaign && (
        <p className="rounded-[4px] border border-zinc-300 bg-zinc-200/60 px-3 py-2 text-sm leading-5 text-zinc-700 dark:border-zinc-700 dark:bg-zinc-800/60 dark:text-zinc-300">
          {READ_ONLY_REASON[campaign.status] ??
            "This campaign cannot be edited."}
          {campaign.status === "failed" && campaign.last_error
            ? ` Reason: ${campaign.last_error}.`
            : null}
        </p>
      )}

      <TextField
        label="Campaign Title"
        name="title"
        isRequired
        type="text"
        value={title}
        onChange={setTitle}
        isDisabled={isReadOnly}
        placeholder="Internal name, e.g. Welcome to SwampHacks XII"
      />

      <TextField
        label="Campaign Description"
        name="description"
        type="text"
        value={description}
        onChange={setDescription}
        isDisabled={isReadOnly}
        placeholder="Internal optional note for organizers"
        description="The title and description above are internal only, never sent to recipients."
      />

      {/* Everything below this line is delivered to recipients. */}
      <div className="my-2 border-t border-zinc-300 dark:border-zinc-600" />

      {/* MultiSelect has no disabled prop, so lock it from the outside. */}
      <div
        aria-disabled={isReadOnly}
        className={cn(isReadOnly && "pointer-events-none opacity-60")}
      >
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
      </div>

      <TextField
        label="Email Subject"
        name="subject"
        isRequired
        type="text"
        value={subject}
        onChange={setSubject}
        isDisabled={isReadOnly}
        placeholder="What recipients see in their inbox"
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
              disabled={isReadOnly}
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
          readOnly={isReadOnly}
          rows={14}
          style={{ height: 349 }}
          placeholder={
            formatValue === "html"
              ? "<h1>Hello</h1><p>Email contents</p>"
              : formatValue === "text"
                ? "Email contents"
                : "Choose a format above to start writing"
          }
          className={cn(
            "bg-input-bg text-text-main w-full resize-y rounded-[4px] border border-zinc-300 px-[13px] py-3 text-sm leading-5 outline-0 dark:border-zinc-700",
            isReadOnly &&
              "cursor-not-allowed resize-none bg-zinc-200/50 text-zinc-500 dark:bg-zinc-800/50 dark:text-zinc-400",
          )}
        />
      </div>

      {saveError && (
        <p className="text-sm text-red-600">
          {saveError instanceof Error && saveError.message
            ? saveError.message
            : "Could not save. Check the fields and try again."}
        </p>
      )}

      {campaign?.status === "scheduled" && campaign.scheduled_at && (
        <div className="flex flex-col gap-2 rounded-[4px] border border-[#d08700] bg-[#fef9c2] px-3 py-2">
          <p className="text-sm leading-5 text-[#d08700]">
            Scheduled to send{" "}
            {formatDate(
              new Date(campaign.scheduled_at),
              "EEEE, MMMM d 'at' h:mm a",
            )}
            .
          </p>
          <div className="flex gap-2">
            <Button
              variant="unstyled"
              size="auto"
              onPress={() => setScheduling(true)}
              isDisabled={isSaving}
              className="h-[26px] cursor-pointer rounded-[4px] border border-[#d08700] px-3 text-xs leading-4 font-medium text-[#d08700]"
            >
              Reschedule
            </Button>
            <Button
              variant="unstyled"
              size="auto"
              onPress={handleUnschedule}
              isDisabled={isSaving}
              className="h-[26px] cursor-pointer rounded-[4px] border border-[#d08700] bg-[#d08700] px-3 text-xs leading-4 font-medium text-white"
            >
              {unscheduleCampaign.isPending ? "Unscheduling..." : "Unschedule"}
            </Button>
          </div>
        </div>
      )}

      <div
        className={cn("flex items-center gap-2 pb-6", isReadOnly && "hidden")}
      >
        {/* Figma 878:95 — one 26px pill, label and chevron sharing the fill. */}
        <div
          className={cn(
            "inline-flex h-[26px] items-center rounded-[4px] bg-[#2b7fff]",
            (!canSave || isSaving) && "opacity-50",
          )}
        >
          <Button
            variant="unstyled"
            size="auto"
            onPress={() => setConfirmingSend(true)}
            isDisabled={!canSave || isSaving}
            className="h-full cursor-pointer rounded-l-[4px] pr-1 pl-3 text-xs leading-4 font-medium text-white"
          >
            {isSaving ? "Working..." : "Send"}
          </Button>

          <MenuTrigger>
            <Button
              variant="unstyled"
              size="auto"
              aria-label="More send options"
              isDisabled={!canSave || isSaving}
              className="h-full cursor-pointer rounded-r-[4px] pr-2.5 pl-1 text-white"
            >
              <TablerChevronDown className="size-2.5" />
            </Button>
            <Menu
              placement="bottom start"
              className="overflow-visible"
              popoverClassName="duration-75"
            >
              <MenuItem
                className={campaignMenuItemClasses}
                onAction={() => setScheduling(true)}
              >
                Schedule for later
              </MenuItem>
            </Menu>
          </MenuTrigger>
        </div>

        {/* Saving without sending is the low-stakes action, so it stays neutral. */}
        <Button
          variant="unstyled"
          size="auto"
          onPress={handleSaveDraft}
          isDisabled={!canSave || isSaving}
          className={cn(
            "h-[26px] cursor-pointer rounded-[4px] border border-zinc-300 bg-zinc-200 px-3 text-xs leading-4 font-medium text-zinc-700 transition-colors hover:bg-zinc-300 dark:border-zinc-700 dark:bg-zinc-800 dark:text-zinc-200 dark:hover:bg-zinc-700",
            (!canSave || isSaving) && "opacity-50",
          )}
        >
          Save for later
        </Button>
      </div>

      {/* Sending is irreversible, so name the audience before it goes out. */}
      <Modal
        size="md"
        isOpen={confirmingSend}
        onOpenChange={setConfirmingSend}
        title="Send this campaign?"
      >
        <div className="flex flex-col gap-4">
          <p className="text-sm leading-5 text-zinc-700 dark:text-zinc-300">
            This emails everyone in{" "}
            <strong>
              {recipients
                .map(
                  (r) =>
                    RECIPIENT_OPTIONS.find((o) => o.value === r)?.label ?? r,
                )
                .join(", ")}
            </strong>{" "}
            immediately. It cannot be undone or recalled.
          </p>
          <p className="text-sm leading-5 text-zinc-500 dark:text-zinc-400">
            Subject: {subject.trim() || "(none)"}
          </p>
          <div className="flex justify-end gap-2">
            <Button
              variant="secondary"
              size="sm"
              onPress={() => setConfirmingSend(false)}
              isDisabled={isSaving}
            >
              Cancel
            </Button>
            <Button
              variant="primary"
              size="sm"
              onPress={handleSendNow}
              isDisabled={isSaving}
            >
              {isSaving ? "Sending..." : "Send now"}
            </Button>
          </div>
        </div>
      </Modal>

      {/* The API accepts past times, so only offer future ones. */}
      <Modal
        size="md"
        isOpen={scheduling}
        onOpenChange={setScheduling}
        title="Schedule this campaign"
      >
        <div className="flex flex-col gap-4">
          {/* Date and time are one CalendarDateTime underneath, edited separately. */}
          <DatePicker
            label="Date"
            value={scheduleAt ? toCalendarDate(scheduleAt) : null}
            onChange={(date: CalendarDate | null) => {
              if (!date) return setScheduleAt(null);
              setScheduleAt(
                scheduleAt
                  ? scheduleAt.set({
                      year: date.year,
                      month: date.month,
                      day: date.day,
                    })
                  : toCalendarDateTime(date).set({ hour: 9, minute: 0 }),
              );
            }}
            minValue={toCalendarDate(now(getLocalTimeZone()))}
          />

          <div className="flex flex-col gap-1">
            <Label className="text-sm leading-5 font-normal text-zinc-600 dark:text-zinc-400">
              Time
            </Label>
            <input
              type="time"
              value={
                scheduleAt
                  ? `${String(scheduleAt.hour).padStart(2, "0")}:${String(scheduleAt.minute).padStart(2, "0")}`
                  : ""
              }
              onChange={(e) => {
                const [hour, minute] = e.target.value.split(":").map(Number);
                if (Number.isNaN(hour) || Number.isNaN(minute)) return;
                const base =
                  scheduleAt ?? toCalendarDateTime(now(getLocalTimeZone()));
                setScheduleAt(base.set({ hour, minute, second: 0 }));
              }}
              className="bg-input-bg border-input-border text-text-main w-fit rounded-sm border px-2.5 py-2 text-sm leading-5 outline-0"
            />
          </div>

          <p
            className={cn(
              "text-sm leading-5",
              scheduleIsPast
                ? "text-red-600"
                : "text-zinc-600 dark:text-zinc-400",
            )}
          >
            {!scheduleAt
              ? "Pick a date and time."
              : scheduleIsPast
                ? `${formatDate(scheduleAt.toDate(getLocalTimeZone()), "EEEE, MMMM d 'at' h:mm a")} is in the past. Pick a future time.`
                : `Sends ${formatDate(scheduleAt.toDate(getLocalTimeZone()), "EEEE, MMMM d 'at' h:mm a")}`}
          </p>

          {/* The panel is behind the modal, so failures have to surface here. */}
          {saveError && (
            <p className="text-sm leading-5 text-red-600">
              {saveError instanceof Error && saveError.message
                ? saveError.message
                : "Could not schedule. Try again."}
            </p>
          )}

          <div className="flex justify-end gap-2">
            <Button
              variant="secondary"
              size="sm"
              onPress={() => setScheduling(false)}
              isDisabled={isSaving}
            >
              Cancel
            </Button>
            <Button
              variant="primary"
              size="sm"
              onPress={handleSchedule}
              isDisabled={
                isSaving || scheduleAt === null || scheduleIsPast || !canSave
              }
            >
              {isSaving ? "Scheduling..." : "Schedule"}
            </Button>
          </div>
        </div>
      </Modal>
    </div>
  );
}
