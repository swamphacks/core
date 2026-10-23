import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/ky";

type AttendanceWorkshop = {
  id: string;
  title: string;
  type: string;
  startTime: string;
  checkedInCount: number;
};

export function useWorkshopAttendance() {
  return useQuery({
    queryKey: ["workshop-attendance-options"],
    queryFn: () =>
      api.get("workshops/attendance-options").json<AttendanceWorkshop[]>(),
    refetchInterval: 10000,
  });
}
