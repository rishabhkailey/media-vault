import { UNKNOWN_DATE } from "@/utils/date";
import axios from "axios";
import { defineStore } from "pinia";
import { ref } from "vue";

export const useSearchStore = defineStore("search", () => {
  const nextPageNumber = ref(1);
  const mediaList = ref<Array<Media>>([]);
  const allMediaLoaded = ref(false);
  const query = ref("");

  function reset() {
    nextPageNumber.value = 1;
    mediaList.value = [];
    allMediaLoaded.value = false;
    query.value = "";
  }

  function appendMedia(_mediaList: Array<Media>) {
    if (_mediaList.length > 0) {
      const newMediaList = _mediaList
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
      mediaList.value = [...mediaList.value, ...newMediaList];
    }
  }

  function loadMoreSearchResults(
    accessToken: string,
    _query: string
  ): Promise<boolean> {
    console.log(query.value, _query, query.value !== _query);
    return new Promise<boolean>((resolve, reject) => {
      if (query.value !== _query) {
        console.log("query changed resetting media list");
        reset();
      }
      query.value = _query;
      axios
        .get<Array<Media>>(
          `/v1/search?query=${_query}&page=${nextPageNumber.value}&perPage=30&order=date&sort=desc`,
          {
            headers: {
              Authorization: `Bearer ${accessToken}`,
            },
          }
        )
        .then((response) => {
          console.log(response);
          if (response.status == 200) {
            appendMedia(response.data);
            nextPageNumber.value += 1;
            allMediaLoaded.value = response.data.length === 0;
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
  }

  function setQuery(_query: string) {
    if (query.value === _query) {
      return;
    }
    query.value = _query;
    reset();
  }

  return {
    nextPageNumber,
    mediaList,
    allMediaLoaded,
    loadMoreSearchResults,
    reset,
    setQuery,
  };
});
