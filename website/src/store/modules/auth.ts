import type { Module } from "vuex";
import type { AxiosResponse } from "axios";
import axios from "axios";

import {
  SET_LOGGED_IN_USERINFO_ACTION,
  LOGOUT_ACTION,
  SET_USERINFO_ACTION,
  RESET_USERINFO_ACTION,
} from "@/store/actions-type";
import type { User } from "oidc-client-ts";

type AuthModuleState = {
  authenticated: Boolean;
  userName: string;
  email: string;
  accessToken: string;
  idToken: string;
};

// todo define root state type
export const authModule: Module<AuthModuleState, any> = {
  state: {
    authenticated: false,
    userName: "",
    email: "",
    accessToken: "",
    idToken: "",
  },
  mutations: {
    setUserInfo(state, payload: { userName: string; email: string }) {
      if (
        payload?.userName !== undefined &&
        typeof payload.userName === "string"
      ) {
        state.userName = payload.userName;
      }
      if (payload?.email !== undefined && typeof payload.email === "string") {
        state.email = payload.email;
      }
    },
    setAuthenticated(state, payload: { authenticated: boolean }) {
      if (
        payload?.authenticated !== undefined &&
        typeof payload.authenticated === "boolean"
      ) {
        state.authenticated = payload.authenticated;
        console.log(state);
      }
    },
    setTokens(state, payload: { idToken: string; accessToken: string }) {
      state.accessToken = payload.accessToken;
      state.idToken = payload.idToken;
    },
  },
  actions: {
    [LOGOUT_ACTION]({ commit }) {
      return new Promise((resolve, reject) => {
        axios
          .post("/v1/logout")
          .then(() => {
            commit("setAuthenticated", {
              authenticated: false,
            });
            commit("setUserInfo", {
              email: "",
              userName: "",
            });
            resolve("Success!");
          })
          .catch((err) => {
            console.debug(err);
            reject("Something went wrong");
          });
      });
    },
    [SET_USERINFO_ACTION]({ commit }, payload: User) {
      console.log(payload);
      // todo contants for mutations/commit? mutations are not used outside the store logic
      return new Promise((resolve, reject) => {
        if (payload.profile?.email === undefined) {
          reject("invalid User object");
          return;
        }
        commit("setAuthenticated", {
          authenticated: true,
        });
        commit("setUserInfo", {
          email: payload.profile.email,
          userName: payload.profile.email,
        });
        commit("setTokens", {
          idToken: payload.id_token,
          accessToken: payload.access_token,
        });
        resolve(true);
      });
    },
    [SET_LOGGED_IN_USERINFO_ACTION]({ commit }) {
      return new Promise((resolve, reject) => {
        axios
          .get("/v1/userinfo")
          .then((res: AxiosResponse) => {
            try {
              const userName: string = res.data?.userName;
              const email: string = res.data?.email;
              if (userName.length == 0 && email.length == 0) {
                throw new Error("Invalid Response");
              }
              commit("setAuthenticated", {
                authenticated: true,
              });
              commit("setUserInfo", {
                email,
                userName,
              });
              resolve("Success!");
            } catch (_) {
              reject("Invalid Response");
            }
          })
          .catch((err: any) => {
            console.debug(err);
            if (axios.isAxiosError(err) && err.response?.status === 401) {
              commit("setAuthenticated", {
                authenticated: false,
              });
              // resolve or reject?
              resolve("Unauthorized");
              return;
            }
            reject("Something went wrong");
          });
      });
    },
    [RESET_USERINFO_ACTION]({ commit }) {
      return new Promise((resolve) => {
        commit("setUserInfo", {
          userName: "",
          email: "",
        });
        commit("setAuthenticated", {
          authenticated: false,
        });
        resolve(true);
      });
    },
  },
  // type of this?
  getters: {
    authenticated(state) {
      return state.authenticated;
    },
    userName(state) {
      return state.userName;
    },
    email(state) {
      return state.email;
    },
    accessToken(state) {
      return state.accessToken;
    },
  },
};
