import ky from "ky";
import appConfig from "@/config";

export const api = ky.create({ prefixUrl: appConfig.BASE_API_URL });
