// based on the RFC: https://datatracker.ietf.org/doc/html/rfc6749
// certain providers will have their own unique params too
type BaseOAuthParams = {
  response_type: "code";
  scope: string;
  client_id: string;
  redirect_uri: string;
  state?: string;
};

type OauthProvider = {
  authorization: {
    url: string;
    scopes: string;
    clientId: string;
  };
};

const discordOauth: OauthProvider = {
  authorization: {
    url: "https://discord.com/oauth2/authorize",
    scopes: "identify email",
    clientId: import.meta.env["SWAMPHACKS_DISCORD_OAUTH_CLIENT_ID"],
  },
};

const buildOauthParams = (
  oauthProvider: OauthProvider,
  redirectUri: string,
): BaseOAuthParams => ({
  response_type: "code",
  scope: oauthProvider.authorization.scopes,
  client_id: oauthProvider.authorization.clientId,
  redirect_uri: redirectUri,
});

type SignInParams = {
  provider: "discord";
  redirect_uri: string;
};

const signInWithOauth = ({ provider, redirect_uri }: SignInParams) => {
  let oauthProvider: OauthProvider;

  if (provider === "discord") {
    oauthProvider = discordOauth;
  } else {
    throw new Error(
      "Unsupported OAuth provider. Only Discord OAuth is currently supported.",
    );
  }

  const url = new URL(oauthProvider.authorization.url);

  const params = buildOauthParams(oauthProvider, redirect_uri);

  for (const [param, value] of Object.entries(params)) {
    url.searchParams.set(param, value);
  }

  window.location.href = url.toString();
};

export default {
  signInWithOauth,
};
