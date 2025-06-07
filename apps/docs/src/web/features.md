
## Authentication
SwampHacks' authentication system is built using our own authentication library, with an API inspired by NextAuth.js.
Currently, it supports [Discord Oauth2](https://discord.com/developers/docs/topics/oauth2) for log in and registration.

### Implementation
Our auth library uses the Authorization Code Grant, which allows our backend server to retrieve an access code and exchange it for the user's access token.

The auth library can be initialized like so:

```typescript
const authClient = Auth({
  providers: [Discord],
  redirect_uri: <insert backend's url to handle oauth callback>,
});
```

The library is flexible enough to handle more providers other than Discord.
A provider has following type:

```typescript
type Provider = {
  id: string;
  authorization: {
    url: string;
    scopes: string;
    clientId: string;
  };
}
```

To create a new provider, the `createProvider` function in `providers.ts` must be called so that types are inferred correctly.

For example:

```typescript
const Google = createProvider({
  id: "google",
  authorization: {
    url: <Google Oauth url>,
    scopes: <Oauth scopes>,
    clientId: <obtain client id from Google developer portal>,
  },
});

```