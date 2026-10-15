export interface DashboardRow {
  applicationId: string;
  userId: string;
  name: string;
  email: string;
  status: string;
  checkedInAt: string | null;
  rfid: string | null;
  signupSource: string;
  standbyArrival: string | null;
  redemptions: { redeemableId: string; name: string; amount: number }[];
}

export interface Dashboard {
  hackathonId: string;
  confirmed: number;
  confirmedCheckedIn: number;
  dayOfSignups: number;
  standbyWaiting: number;
  rows: DashboardRow[];
}
