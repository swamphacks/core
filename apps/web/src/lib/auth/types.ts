// based on the RFC: https://datatracker.ietf.org/doc/html/rfc6749
// certain providers will have their own unique params too
export type BaseOAuthParams = {
  response_type: "code";
  scope: string;
  client_id: string;
  redirect_uri: string;
  state?: string;
};

export type OauthProvider = {
  id: string;
  authorization: {
    url: string;
    scopes: string;
    clientId: string;
  };
};

export type OauthState = {
  nonce: string;
  providerId: string;
  returnTo?: string;
};

export type AuthConfig = {
  providers: OauthProvider[];
  redirect_uri: string;
};
