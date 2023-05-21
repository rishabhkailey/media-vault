import type { AxiosResponse } from "axios";
import axios from "axios";
import type { User } from "oidc-client-ts";
import { defineStore } from "pinia";
import { ref } from "vue";

export const useAuthStore = defineStore("auth", () => {
  const authenticated = ref(false);
  const userName = ref("");
  const email = ref("");
  const accessToken = ref("");
  const idToken = ref("");

  function setUserInfo(user: User) {
    console.log(user);
    if (user.profile.email === undefined) {
      throw new Error("Invalid user details. email missing.");
    }
    userName.value = user.profile.email;
    email.value = user.profile.email;
    accessToken.value = user.access_token;
    if (user.id_token) {
      idToken.value = user.id_token;
    }
    authenticated.value = true;
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

  function loadUserInfo(): Promise<boolean> {
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
    setUserInfo,
    setAuthenticated,
    setTokens,
    logOut,
    loadUserInfo,
    reset,
  };
});
