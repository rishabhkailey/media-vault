import type { Module } from "vuex";
import type { AxiosError, AxiosResponse } from "axios";

import axios from "axios";

type AuthModuleState = {
  authenticated: Boolean;
  userName: string;
  email: string;
};

// todo define root state type
const authModule: Module<AuthModuleState, any> = {
  state: {
    authenticated: false,
    userName: "",
    email: "",
  },
  mutations: {
    setUserInfo(state, payload: { userName: string; email: string }) {
      if (!!payload?.userName && typeof payload.userName === "string") {
        state.userName = payload.userName;
      }
      if (!!payload?.email && typeof payload.email === "string") {
        state.email = payload.email;
      }
    },
    setAuthenticated(state, payload: { authenticated: string }) {
      if (
        !!payload?.authenticated &&
        typeof payload.authenticated === "boolean"
      ) {
        state.authenticated = payload.authenticated;
      }
    },
  },
  actions: {
    setLoggedInUserInfo({ commit }) {
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
            if (axios.isAxiosError(err) && err.status === 401) {
              reject("Unauthorized");
              return;
            }
            reject("Something went wrong!");
          });
      });
    },
  },
  // type of this?
  getters: {
    authenticated (state) {
      return state.authenticated
    },
    userName (state) {
      return state.userName
    },
    email (state) {
      return state.email
    }
  },
};

export default authModule;
