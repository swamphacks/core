import { api } from "@/lib/ky";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";

export type CampaignStatus =
  | "draft"
  | "scheduled"
  | "sending"
  | "sent"
  | "failed";

export type CampaignFormat = "text" | "html";

export type RecipientType =
  | "admins"
  | "staff"
  | "visitors"
  | "accepted_applicants"
  | "rejected_applicants"
  | "waitlisted_applicants"
  | "interest_subscribers";

/**
 * Two shapes on purpose: responses are snake_case (sqlc row struct), request
 * bodies are camelCase (huma handler structs). Merging them breaks writes.
 */
export interface EmailCampaign {
  id: string;
  hackathon_id: string;
  title: string;
  description: string | null;
  subject: string;
  body: string;
  format: CampaignFormat;
  recipient_types: RecipientType[];
  status: CampaignStatus;
  scheduled_at: string | null;
  sent_at: string | null;
  last_error: string | null;
  created_by_user_id: string | null;
  updated_by_user_id: string | null;
  created_at: string;
  updated_at: string;
}

export interface CreateCampaignRequest {
  title: string;
  subject: string;
  body: string;
  format: CampaignFormat;
  recipientTypes: RecipientType[];
  description?: string;
  scheduledAt?: string;
}

export type UpdateCampaignRequest = Partial<CreateCampaignRequest>;

const campaignKeys = {
  all: (hackathonId: string) => ["emailCampaigns", hackathonId] as const,
  detail: (hackathonId: string, campaignId: string) =>
    ["emailCampaigns", hackathonId, campaignId] as const,
};

export function useEmailCampaigns(hackathonId: string) {
  return useQuery({
    queryKey: campaignKeys.all(hackathonId),
    queryFn: () =>
      api
        .get<EmailCampaign[]>("email/campaigns", {
          searchParams: { hackathonId },
        })
        .json(),
    enabled: Boolean(hackathonId),
  });
}

export function useEmailCampaign(hackathonId: string, campaignId: string) {
  return useQuery({
    queryKey: campaignKeys.detail(hackathonId, campaignId),
    queryFn: () =>
      api
        .get<EmailCampaign>(`email/campaigns/${campaignId}`, {
          searchParams: { hackathonId },
        })
        .json(),
    enabled: Boolean(hackathonId && campaignId),
  });
}

export function useCreateEmailCampaign(hackathonId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateCampaignRequest) =>
      api
        .post<EmailCampaign>("email/campaigns", {
          json: { ...data, hackathonId },
        })
        .json(),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: campaignKeys.all(hackathonId),
      });
    },
  });
}

export function useUpdateEmailCampaign(hackathonId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      campaignId,
      data,
    }: {
      campaignId: string;
      data: UpdateCampaignRequest;
    }) =>
      api
        .patch<EmailCampaign>(`email/campaigns/${campaignId}`, {
          searchParams: { hackathonId },
          json: data,
        })
        .json(),
    onSuccess: (campaign) => {
      queryClient.invalidateQueries({
        queryKey: campaignKeys.all(hackathonId),
      });
      queryClient.invalidateQueries({
        queryKey: campaignKeys.detail(hackathonId, campaign.id),
      });
    },
  });
}

/**
 * Irreversible: queues one email per resolved recipient. No request body.
 * Returns status "sent" (meaning queued, not delivered) plus sent_at.
 * Rejects with 400 unless the campaign is still draft or scheduled.
 */
export function useSendEmailCampaign(hackathonId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (campaignId: string) =>
      api
        .post<EmailCampaign>(`email/campaigns/${campaignId}/send`, {
          searchParams: { hackathonId },
        })
        .json(),
    onSuccess: (campaign) => {
      queryClient.invalidateQueries({
        queryKey: campaignKeys.all(hackathonId),
      });
      queryClient.invalidateQueries({
        queryKey: campaignKeys.detail(hackathonId, campaign.id),
      });
    },
  });
}

/**
 * Responds 200 with an empty body, so never call .json() here — it would throw.
 * Only draft, scheduled, and failed are deletable; sent and sending reject 400.
 */
export function useDeleteEmailCampaign(hackathonId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (campaignId: string) => {
      await api.delete(`email/campaigns/${campaignId}`, {
        searchParams: { hackathonId },
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: campaignKeys.all(hackathonId),
      });
    },
  });
}
