import config from "@/config";
import ky from "ky";

export const api = ky.create({
  prefixUrl: config.BASE_API_URL,
  credentials: "include",
});
