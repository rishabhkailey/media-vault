import { UserManager, type User, WebStorageStateStore } from "oidc-client-ts";
import { v4 } from "uuid";

export interface InternalState {
  internalRedirectPath: string;
  internalRedirectQuery: string;
  nonce: string;
}

export const userManager = new UserManager({
  authority: import.meta.env.VITE_AUTH_SERVICE_URL,
  metadataUrl: import.meta.env.VITE_AUTH_SERVICE_DISCOVERY_ENDPOINT,
  client_id: import.meta.env.VITE_AUTH_SERVICE_CLIENT_ID,
  redirect_uri: window.location.origin + "/pkce",
  response_type: "code",
  scope: "openid profile email user",
  post_logout_redirect_uri: window.location.origin,
  // silent_redirect_uri: window.location.origin + "/static/silent-renew.html",
  accessTokenExpiringNotificationTimeInSeconds: 10,
  automaticSilentRenew: true,
  // if true it removes the nonce
  filterProtocolClaims: false,
  loadUserInfo: true,
  stateStore: new WebStorageStateStore({ store: window.localStorage }),
  userStore: new WebStorageStateStore({ store: window.localStorage }),
});

export function signinUsingUserManager(
  userManager: UserManager,
  redirectToHome: boolean,
) {
  const nonce = v4();
  const state: InternalState = {
    internalRedirectPath: "/",
    internalRedirectQuery: "",
    nonce: nonce,
  };
  if (!redirectToHome) {
    state.internalRedirectPath = location.pathname;
    state.internalRedirectQuery = location.search;
  }
  userManager
    .signinRedirect({
      nonce: nonce,
      state: state,
    })
    .then((response) => {
      console.log(response);
    })
    .catch((err) => {
      console.log(err);
    });
}

export async function handlePostLoginUsingUserManager(
  userManager: UserManager,
): Promise<User> {
  return new Promise((resolve, reject) => {
    userManager
      .signinRedirectCallback()
      .then((user: User) => {
        userManager.storeUser(user);
        console.log(user);
        // store token in store or store userManager in store?
        // todo validate internal redirect
        try {
          const internalState = user.state as InternalState;
          if (
            internalState.nonce.length === 0 ||
            user.profile.nonce?.length === 0 ||
            internalState.nonce.length !== user.profile.nonce?.length
          ) {
            throw new Error("Nonce mismatch");
          }
          resolve(user);
        } catch (err) {
          reject(err);
          return;
        }
      })
      .catch((err) => {
        reject(err);
        return;
      });
  });
}
