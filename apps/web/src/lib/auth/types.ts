export type User = {
  userId: string;
  name: string;
  image: string;
  onboarded: boolean;
  role: string;
};

// based on the RFC: https://datatracker.ietf.org/doc/html/rfc6749
// certain providers will have their own unique params too
export type BaseOAuthParams = {
  response_type: "code";
  scope: string;
  client_id: string;
  redirect_uri: string;
  state?: string;
};

export type OAuthProvider = {
  id: string;
  authorization: {
    url: string;
    scopes: string;
    clientId: string;
  };
};

export type OAuthState = {
  nonce: string;
  provider: string;
  redirect?: string;
};

export type AuthConfig = {
  providers: OAuthProvider[];
  redirectUri: string;
};
