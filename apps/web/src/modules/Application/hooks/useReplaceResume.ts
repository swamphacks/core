import { api } from "@/lib/ky";
import { useMutation } from "@tanstack/react-query";

// Replaces the resume on an already-submitted application without touching any
// of the question responses. Mirrors the multipart upload used by the submit flow.
export async function replaceResume(resume: File): Promise<void> {
  const formData = new FormData();
  formData.append("resume", resume);

  await api.put(`application/resume`, {
    body: formData,
  });
}

export function useReplaceResume() {
  return useMutation({
    mutationFn: (resume: File) => replaceResume(resume),
  });
}
