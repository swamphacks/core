import { type UserContext } from "./user";
import { type ErrorResponse } from "./error";

export interface AuthUserResponse {
  user: UserContext | null;
  error: ErrorResponse | null;
}
