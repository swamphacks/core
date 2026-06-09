import type { components } from "./schema";

export type UserContext = components["schemas"]["UserContext"];
export type Application = components["schemas"]["HackerApplication"];
export type Hackathon = components["schemas"]["PublicHackathon"];
export type StaffHackathon = components["schemas"]["Hackathon"];
export type ApplicationStats = components["schemas"]["ApplicationStatistics"];
export type HackathonStaff = Array<components["schemas"]["User"]>;
export type AssignedApplications = Array<
  components["schemas"]["AssignedApplication"]
>;
export type ApplicationReviewDetails =
  components["schemas"]["ApplicationReviewDetails"];
export type ApplicationAutoDecisionRequest = Array<
  components["schemas"]["ListAutoDecisionRequestsRow"]
>;
