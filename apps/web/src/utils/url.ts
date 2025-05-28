const DEV_BASE_API_URL = "http://localhost:8080/api";
const PROD_BASE_API_URL = "https://api.swamphacks.com";

export const getBaseUrl = () => {
  return import.meta.env.DEV ? DEV_BASE_API_URL : PROD_BASE_API_URL;
};

export const APP_URL = {
  VERIFY_AUTH: `${getBaseUrl()}/auth/me`,
};
