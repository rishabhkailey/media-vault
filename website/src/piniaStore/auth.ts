import { refreshSessionWithResourceServer } from "@/js/api/user";
import type { AxiosResponse } from "axios";
import axios from "axios";
import type { User } from "oidc-client-ts";
import { defineStore } from "pinia";
import { computed, ref } from "vue";
import { useErrorsStore } from "./errors";

export const useAuthStore = defineStore("auth", () => {
  const { appendError } = useErrorsStore();
  const authenticated = ref(false);
  const userName = ref("");
  const email = ref("");
  const accessToken = ref("");
  const idToken = ref("");
  const expireAt = ref(Date.now() / 1000);
  const expired = computed(() => Date.now() / 1000 > expireAt.value);

  function setUserAuthInfo(user: User) {
    if (user.profile.email !== undefined) {
      email.value = user.profile.email;
    } else {
      email.value = user.profile.sub;
    }
    if (user.profile.name !== undefined) {
      userName.value = user.profile.name;
    }
    accessToken.value = user.access_token;
    if (user.id_token) {
      idToken.value = user.id_token;
    }
    authenticated.value = true;
    if (user.expires_at !== undefined) {
      expireAt.value = user.expires_at;
    } else {
      // default 24 hours
      expireAt.value = Date.now() / 1000 + 24 * 60 * 60;
    }
    refreshSessionWithResourceServer(accessToken.value).catch((err) => {
      appendError(
        "failed to refresh session with resource server",
        `if facing any issues please refresh the page. error message: ${err}`,
        -1,
      );
    });
  }

  function setAuthenticated(_authenticated: boolean) {
    authenticated.value = _authenticated;
  }

  function setTokens(_accessToken: string, _idToken: string) {
    accessToken.value = _accessToken;
    idToken.value = _idToken;
  }

  function logOut(): Promise<boolean> {
    return new Promise((resolve, reject) => {
      axios
        .post("/v1/logout")
        .then(() => {
          reset();
          resolve(true);
        })
        .catch((err) => {
          console.debug(err);
          reject(err);
        });
    });
  }

  function loadUserAuthInfo(): Promise<boolean> {
    return new Promise((resolve, reject) => {
      axios
        .get("/v1/userinfo")
        .then((res: AxiosResponse) => {
          try {
            const _userName: string = res.data?.userName;
            const _email: string = res.data?.email;
            if (_userName.length == 0 && _email.length == 0) {
              throw new Error("Invalid Response");
            }
            authenticated.value = true;
            userName.value = _email;
            email.value = _email;
            resolve(true);
          } catch (_) {
            reject("Invalid Response");
          }
        })
        .catch((err: any) => {
          console.debug(err);
          if (axios.isAxiosError(err) && err.response?.status === 401) {
            setAuthenticated(false);
            // resolve or reject?
            reject(new Error("Unauthorized"));
            return;
          }
          reject("Something went wrong");
        });
    });
  }

  function reset() {
    authenticated.value = false;
    userName.value = "";
    email.value = "";
    accessToken.value = "";
  }

  return {
    authenticated,
    userName,
    email,
    accessToken,
    expired,
    setUserAuthInfo,
    setAuthenticated,
    setTokens,
    logOut,
    loadUserAuthInfo,
    reset,
  };
});
