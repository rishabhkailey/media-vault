import { UNKNOWN_DATE } from "@/utils/date";
import axios from "axios";
import { defineStore } from "pinia";
import { ref } from "vue";

// we will need lock or something else
// to prevent duplicates if the same request is called twice

export const useMediaStore = defineStore("media", () => {
  const nextPageNumber = ref(1);
  const mediaList = ref<Array<Media>>([]);
  const allMediaLoaded = ref(false);

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

  function loadMoreMedia(accessToken: string): Promise<boolean> {
    return new Promise<boolean>((resolve, reject) => {
      axios
        .get<Array<Media>>(
          `/v1/mediaList?page=${nextPageNumber.value}&perPage=30&order=date&sort=desc`,
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
            allMediaLoaded.value = response.data.length == 0;
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

  function loadAllMediaForDate(
    accessToken: string,
    date: Date
  ): Promise<boolean> {
    return new Promise((resolve, reject) => {
      let lastMedia = mediaList.value[mediaList.value.length - 1];
      if (lastMedia.date > date) {
        resolve(true);
        return;
      }
      while (lastMedia.date > date && !allMediaLoaded.value) {
        loadMoreMedia(accessToken)
          .then(() => {
            lastMedia = mediaList.value[mediaList.value.length - 1];
          })
          .catch((err) => {
            reject(err);
            return;
          });
      }
      resolve(true);
      return;
    });
  }

  function removeMediaByID(mediaID: number) {
    mediaList.value = mediaList.value.filter((media) => media.id !== mediaID);
  }

  function removeMediaByIDs(mediaIDs: Array<number>) {
    mediaList.value = mediaList.value.filter(
      (media) => !mediaIDs.includes(media.id)
    );
  }

  function deleteMedia(accessToken: string, mediaID: number): Promise<boolean> {
    return new Promise<boolean>((resolve, reject) => {
      axios
        .delete(`/v1/media/${mediaID}`, {
          headers: {
            Authorization: `Bearer ${accessToken}`,
          },
        })
        .then((resp) => {
          if (resp.status === 200) {
            removeMediaByID(mediaID);
            resolve(true);
            return;
          }
          resolve(false);
          return;
        })
        .catch((err) => {
          reject(err);
        });
    });
  }

  // returns failed media ids
  // reject only in case of unexpected/unhandeled error
  async function deleteMultipleMedia(
    accessToken: string,
    mediaIDs: Array<number>
  ): Promise<Array<number>> {
    const failedMediaIDs = new Set<number>();
    const successMediaIDs = new Set<number>();

    for (const mediaID of mediaIDs) {
      try {
        const resp = await axios.delete(`/v1/media/${mediaID}`, {
          headers: {
            Authorization: `Bearer ${accessToken}`,
          },
        });
        if (resp.status !== 200) {
          failedMediaIDs.add(mediaID);
          continue;
        } else {
          successMediaIDs.add(mediaID);
        }
      } catch {
        failedMediaIDs.add(mediaID);
      }
    }
    removeMediaByIDs(Array.from(successMediaIDs));
    return Array.from(failedMediaIDs);
  }

  return {
    nextPageNumber,
    mediaList,
    allMediaLoaded,
    loadMoreMedia,
    loadAllMediaForDate,
    deleteMedia,
    deleteMultipleMedia,
  };
});
