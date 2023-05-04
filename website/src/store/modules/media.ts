import type { Module } from "vuex";
import axios from "axios";
import { UNKNOWN_DATE } from "@/utils/date";

// mutations
const SET_MEDIA_LIST = "setLoggedInUserInfo";
const SET_NEXT_PAGE_NUMBER = "setNextPageNumber";
const SET_LOADED_ALL_MEDIA = "setLoadedAllMedia";
const ADD_MEDIA_TO_LIST = "addMediaToList";

// actions
export const LOAD_MORE_MEDIA_ACTION = "loadMoreMedia";

type MediaModuleState = {
  nextPageNumber: number;
  mediaList: Array<Media>;
  allMediaLoaded: boolean;
};

// todo define root state type
export const mediaModule: Module<MediaModuleState, any> = {
  state: {
    nextPageNumber: 1,
    mediaList: [],
    allMediaLoaded: false,
  },
  mutations: {
    [SET_NEXT_PAGE_NUMBER](state, payload: number) {
      if (payload > 0) {
        state.nextPageNumber = payload;
      }
    },
    [SET_MEDIA_LIST](state, payload: Array<Media>) {
      if (payload.length > 0) {
        state.mediaList = payload;
      }
    },
    [ADD_MEDIA_TO_LIST](state, payload: Array<Media>) {
      if (payload.length > 0) {
        const newMediaList = payload
          .map((media) => {
            if (typeof media.date === "string") {
              try {
                media.date = new Date(media.date);
              } catch (err) {
                media.date = UNKNOWN_DATE;
              }
            }
            return media;
          })
          .sort((m1, m2) => {
            return m1.date > m2.date ? -1 : 1;
          });
        state.mediaList = [...state.mediaList, ...newMediaList];
      }
    },
    [SET_LOADED_ALL_MEDIA](state, payload: boolean) {
      console.trace(payload);
      state.allMediaLoaded = payload;
    },
  },
  actions: {
    [LOAD_MORE_MEDIA_ACTION]({ commit, getters, rootGetters }) {
      return new Promise<boolean>((resolve, reject) => {
        const accessToken: string = rootGetters.accessToken;
        if (accessToken.length === 0) {
          reject(new Error("empty access token"));
          return;
        }
        axios
          .get<Array<Media>>(
            `/v1/mediaList?page=${getters.nextPageNumber}&perPage=30&order=date&sort=desc`,
            {
              headers: {
                Authorization: `Bearer ${accessToken}`,
              },
            }
          )
          .then((response) => {
            console.log(response);
            if (response.status == 200) {
              commit(ADD_MEDIA_TO_LIST, response.data);
              commit(SET_NEXT_PAGE_NUMBER, getters.nextPageNumber + 1);
              commit(SET_LOADED_ALL_MEDIA, response.data.length === 0);
              resolve(true);
              return;
            }
            reject(new Error("non 200 status code"));
            return;
          })
          .catch((err) => {
            console.log(err);
            reject(err);
          });
      });
    },
  },
  // type of this?
  getters: {
    nextPageNumber(state) {
      return state.nextPageNumber;
    },
    mediaList(state) {
      console.log(state);
      return state.mediaList;
    },
    allMediaLoaded(state) {
      return state.allMediaLoaded;
    },
  },
};
