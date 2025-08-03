import { useQuery } from "@tanstack/react-query";
import z from "zod";

export const StaffUserSchema = z.object({
  userId: z.uuid(),
  name: z.string().nullable(),
  email: z.email(),
  image: z.url().nullable(),
  assignedAt: z.iso.datetime(),
  eventRole: z.enum(["STAFF", "ADMIN"]),
});

export type StaffUser = z.infer<typeof StaffUserSchema>;

export function getEventStaffUsersQueryKey(eventId: string) {
  return ["event", eventId, "staff-users"] as const;
}

const mockData: StaffUser[] = [
  {
    userId: "c2107e85-efa8-4dda-b43c-9df306d289c4",
    name: "Alexander Wang",
    email: "alexanderwang@ufl.edu",
    image: null,
    assignedAt: "2025-08-03T18:45:30.123Z",
    eventRole: "ADMIN",
  },
  {
    userId: "a93f70d2-5f91-4e3c-b4de-96d785f51b99",
    name: null,
    email: "jessica@example.com",
    image: "https://example.com/images/jessica.png",
    assignedAt: "2024-12-15T09:30:00.000Z",
    eventRole: "STAFF",
  },
  {
    userId: "1d9843e9-0a97-4b8e-99c3-c67c1de7a9f7",
    name: "Brian Lee",
    email: "brian.lee@example.org",
    image: "https://example.org/images/brian.jpg",
    assignedAt: "2023-11-01T14:22:10.456Z",
    eventRole: "ADMIN",
  },
  {
    userId: "0f6c4787-c6c1-4e27-b1f0-5d3d3f6f11ae",
    name: "Sophia Kim",
    email: "sophia.kim@domain.com",
    image: null,
    assignedAt: "2025-01-20T07:15:45.000Z",
    eventRole: "ADMIN",
  },
  {
    userId: "f2d4159c-37db-4f7d-80fc-8468a890a3bc",
    name: "Michael Johnson",
    email: "michael.johnson@company.net",
    image: "https://company.net/profiles/michael.jpg",
    assignedAt: "2024-06-10T19:30:00.999Z",
    eventRole: "STAFF",
  },
  {
    userId: "bd5a0f1e-64bf-4c3d-9193-fd7a6b561a1b",
    name: null,
    email: "natalie@webmail.com",
    image: "https://webmail.com/avatars/natalie.png",
    assignedAt: "2023-09-12T12:00:00.000Z",
    eventRole: "STAFF",
  },
  {
    userId: "72de2c68-1bcb-4e3b-825d-44d5b58e3e4d",
    name: "James Carter",
    email: "james.carter@mailservice.com",
    image: null,
    assignedAt: "2024-02-25T16:45:30.123Z",
    eventRole: "STAFF",
  },
  {
    userId: "3a7d564b-8e25-4b0d-8c14-0b6e5b11c1d3",
    name: "Emma Thompson",
    email: "emma.thompson@education.edu",
    image: "https://education.edu/photos/emma.jpg",
    assignedAt: "2025-03-03T08:00:00.000Z",
    eventRole: "STAFF",
  },
  {
    userId: "c86f2b24-83d9-4de7-9f77-0b5048e6a7e5",
    name: null,
    email: "oliver@mail.com",
    image: null,
    assignedAt: "2023-12-31T23:59:59.999Z",
    eventRole: "STAFF",
  },
  {
    userId: "9e1a8b02-d8b4-4d20-bb27-b58b8f3c98b2",
    name: "Isabella Martinez",
    email: "isabella.martinez@service.org",
    image: "https://service.org/avatars/isabella.png",
    assignedAt: "2024-05-15T11:11:11.111Z",
    eventRole: "STAFF",
  },
];

function mockFetchWithTimeout<T>(data: T, delayMs: number = 0): Promise<T> {
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve(data);
    }, delayMs);
  });
}

export function useEventStaffUsers(eventId: string) {
  async function fetchEventStaffUsers() {
    return await mockFetchWithTimeout(mockData, 1_500);
  }

  return useQuery({
    queryKey: getEventStaffUsersQueryKey(eventId),
    queryFn: fetchEventStaffUsers,
    staleTime: 1000 * 60 * 5, // 5 minutes
  });
}
