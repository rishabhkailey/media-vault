import type { Module } from "vuex";
import axios from "axios";
import { UNKNOWN_DATE } from "@/utils/date";

export const MODULE_NAME = "searchModule";

// getters
const NEXT_PAGE_NUMBER = "nextPageNumber";
const SEARCH_RESULTS = "searchResults";
const ALL_SEARCH_RESULTS_LOADED = "allSearchResultsLoaded";
const SEARCH_QUERY = "query";

// mutations
const SET_SEARCH_RESULTS = "setLoggedInUserInfo";
const SET_NEXT_PAGE_NUMBER = "setNextPageNumber";
const SET_LOADED_ALL_SEARCH_RESULTS = "setLoadedAllMedia";
const SET_SEARCH_QUERY = "setSearchQuery";
const ADD_MEDIA_TO_SEARCH_RESULTS = "addMediaToList";
const RESET_SEARCH_RESULTS = "resetMediaList";
// actions
const _LOAD_MORE_SEARCH_RESULTS_ACTION = "loadMoreMedia";
const _RESET_SEARCH_RESULTS_ACTION = "resetMediaList";

// external actions
export const LOAD_MORE_SEARCH_RESULTS_ACTION = `${MODULE_NAME}/${_LOAD_MORE_SEARCH_RESULTS_ACTION}`;
export const RESET_SEARCH_RESULTS_ACTION = `${MODULE_NAME}/${_RESET_SEARCH_RESULTS_ACTION}`;

// external getters
export const NEXT_PAGE_NUMBER_GETTER = `${MODULE_NAME}/${NEXT_PAGE_NUMBER}`;
export const SEARCH_RESULTS_GETTER = `${MODULE_NAME}/${SEARCH_RESULTS}`;
export const ALL_SEARCH_RESULTS_LOADED_GETTER = `${MODULE_NAME}/${ALL_SEARCH_RESULTS_LOADED}`;
export const SEARCH_QUERY_GETTER = `${MODULE_NAME}/${SEARCH_QUERY}`;

interface LoadSearchResultsPayload {
  query: string;
  // todo add before and after dates
}
export type LoadSearchResults = LoadSearchResultsPayload;

type MediaModuleState = {
  [NEXT_PAGE_NUMBER]: number;
  [SEARCH_RESULTS]: Array<Media>;
  [ALL_SEARCH_RESULTS_LOADED]: boolean;
  [SEARCH_QUERY]: string;
};

// todo define root state type
export const searchModule: Module<MediaModuleState, any> = {
  namespaced: true,
  state: {
    [NEXT_PAGE_NUMBER]: 1,
    [SEARCH_RESULTS]: [],
    [ALL_SEARCH_RESULTS_LOADED]: false,
    [SEARCH_QUERY]: "",
  },
  mutations: {
    [SET_NEXT_PAGE_NUMBER](state, payload: number) {
      if (payload > 0) {
        state[NEXT_PAGE_NUMBER] = payload;
      }
    },
    [SET_SEARCH_RESULTS](state, payload: Array<Media>) {
      if (payload.length > 0) {
        state[SEARCH_RESULTS] = payload;
      }
    },
    [SET_SEARCH_QUERY](state, payload: string) {
      state[SEARCH_QUERY] = payload;
    },
    [ADD_MEDIA_TO_SEARCH_RESULTS](state, payload: Array<Media>) {
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
        state.searchResults = [...state.searchResults, ...newMediaList];
      }
    },
    [SET_LOADED_ALL_SEARCH_RESULTS](state, payload: boolean) {
      console.trace(payload);
      state[ALL_SEARCH_RESULTS_LOADED] = payload;
    },

    [RESET_SEARCH_RESULTS](state) {
      state[ALL_SEARCH_RESULTS_LOADED] = false;
      state[SEARCH_RESULTS] = [];
      state[NEXT_PAGE_NUMBER] = 1;
    },
  },
  actions: {
    [_RESET_SEARCH_RESULTS_ACTION]({ commit }) {
      commit(RESET_SEARCH_RESULTS);
    },
    [_LOAD_MORE_SEARCH_RESULTS_ACTION](
      { commit, getters, rootGetters, state },
      payload: LoadSearchResultsPayload
    ) {
      return new Promise<boolean>((resolve, reject) => {
        if (state.query !== payload.query) {
          console.log("query changed resetting media list");
          commit(RESET_SEARCH_RESULTS);
        }
        commit(SET_SEARCH_QUERY, payload.query);
        const accessToken: string = rootGetters.accessToken;
        if (accessToken.length === 0) {
          reject(new Error("empty access token"));
          return;
        }
        axios
          .get<Array<Media>>(
            `/v1/search?query=${payload.query}&page=${getters.nextPageNumber}`,
            {
              headers: {
                Authorization: `Bearer ${accessToken}`,
              },
            }
          )
          .then((response) => {
            console.log(response);
            if (response.status == 200) {
              commit(ADD_MEDIA_TO_SEARCH_RESULTS, response.data);
              commit(SET_NEXT_PAGE_NUMBER, getters.nextPageNumber + 1);
              commit(SET_LOADED_ALL_SEARCH_RESULTS, response.data.length === 0);
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
    searchResults(state) {
      console.log(state);
      return state.searchResults;
    },
    allSearchResultsLoaded(state) {
      return state.allSearchResultsLoaded;
    },
  },
};
