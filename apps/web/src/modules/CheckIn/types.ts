export interface DashboardRow {
  applicationId: string;
  userId: string;
  name: string;
  email: string;
  status: string;
  checkedInAt: string | null;
  nfcLinked: boolean;
  rfid: string | null;
  signupSource: string;
  standbyArrival: string | null;
  standbyPosition: number | null;
  redemptions: { redeemableId: string; name: string; amount: number }[];
}

export interface Dashboard {
  capacity: number;
  outstandingOffers: number;
  unreservedSpots: number;
  hackathonId: string;
  confirmed: number;
  confirmedCheckedIn: number;
  dayOfSignups: number;
  standbyWaiting: number;
  rows: DashboardRow[];
}
