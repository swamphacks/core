import { type UserContext } from "./user";
import { type ErrorResponse } from "./error";

export interface AuthMeResponse {
  user: UserContext | null;
  error: ErrorResponse | null;
}
