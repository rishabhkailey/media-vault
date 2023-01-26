import type { Module } from "vuex";
import type { AxiosResponse } from "axios";
import axios from "axios";

import {
  SET_LOGGED_IN_USERINFO_ACTION,
  LOGOUT_ACTION,
} from "@/store/actions-type";

type AuthModuleState = {
  authenticated: Boolean;
  userName: string;
  email: string;
};

// todo define root state type
export const authModule: Module<AuthModuleState, any> = {
  state: {
    authenticated: false,
    userName: "",
    email: "",
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
              resolve("Unauthorized");
              return;
            }
            reject("Something went wrong");
          });
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
  },
};
