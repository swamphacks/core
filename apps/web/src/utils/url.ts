const DEV_BASE_API_URL = "http://localhost:8080/v1";
const PROD_BASE_API_URL = "https://api.swamphacks.com";

export const getBaseUrl = () => {
  return import.meta.env.DEV ? DEV_BASE_API_URL : PROD_BASE_API_URL;
};

// TODO: find a better way to handle this
export const APP_URL = {
  VERIFY_AUTH: `${getBaseUrl()}/auth/me`,
  REDIRECT_URI: `${getBaseUrl()}/auth/callback`,
  REDIRECT_FALLBACK: `${getBaseUrl()}/dashboard`,
};
