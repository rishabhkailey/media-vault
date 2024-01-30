import { UNKNOWN_DATE } from "@/js/date";
import axios from "axios";
import { defineStore, storeToRefs } from "pinia";
import { ref } from "vue";
import { useAuthStore } from "./auth";

// todo we will need lock or something else
// to prevent duplicates if the same request is called twice

export const useMediaStore = defineStore("media", () => {
  const nextPageNumber = ref(1);
  const mediaList = ref<Array<Media>>([]);
  const allMediaLoaded = ref(false);
  // todo move this to different store
  // config store?
  const orderByUploadDateKey = "uploaded_at";
  const orderByDateKey = "date";
  const orderBy = ref(orderByUploadDateKey);
  const { accessToken } = storeToRefs(useAuthStore());

  function orderByUploadDate() {
    if (orderBy.value !== orderByUploadDateKey) {
      reset();
    }
    orderBy.value = orderByUploadDateKey;
  }
  function orderByMediaDate() {
    if (orderBy.value !== orderByDateKey) {
      reset();
    }
    orderBy.value = orderByDateKey;
  }
  function getMediaDateAccordingToOrderBy(media: Media): Date {
    if (orderBy.value === orderByUploadDateKey) {
      return media.uploaded_at;
    }
    return media.date;
  }

  function reset() {
    nextPageNumber.value = 1;
    mediaList.value = [];
    allMediaLoaded.value = false;
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
          if (typeof media.uploaded_at === "string") {
            try {
              media.uploaded_at = new Date(media.uploaded_at);
            } catch (err) {
              media.uploaded_at = UNKNOWN_DATE;
            }
          }
          return media;
        })
        .sort((m1, m2) => {
          return m1.date > m2.date ? -1 : 1;
        });
      mediaList.value = sortAndRemoveDuplicateMedia([
        ...mediaList.value,
        ...newMediaList,
      ]);
    }
  }

  function sortAndRemoveDuplicateMedia(mediaList: Array<Media>): Array<Media> {
    if (mediaList.length <= 1) {
      return mediaList;
    }
    let compare: (m1: Media, m2: Media) => number;
    switch (orderBy.value) {
      case "uploaded_at": {
        compare = (m1: Media, m2: Media) => {
          return m1.uploaded_at > m2.uploaded_at ? -1 : 1;
        };
        break;
      }
      default: {
        compare = (m1: Media, m2: Media) => {
          return m1.date > m2.date ? -1 : 1;
        };
        break;
      }
    }
    mediaList = mediaList.sort((m1, m2) => compare(m1, m2));
    const uniqueMediaList: Array<Media> = [mediaList[0]];
    let previousMediaID = mediaList[0].id;
    for (let index = 1; index < mediaList.length; index++) {
      if (mediaList[index].id === previousMediaID) {
        continue;
      }
      uniqueMediaList.push(mediaList[index]);
      previousMediaID = mediaList[index].id;
    }
    return uniqueMediaList;
  }

  function loadMoreMedia(): Promise<boolean> {
    return new Promise<boolean>((resolve, reject) => {
      let url = `/v1/mediaList?order=${orderBy.value}&sort=desc&per_page=30`;
      if (mediaList.value.length !== 0) {
        const lastMedia = mediaList.value[mediaList.value.length - 1];
        const lastDate =
          getMediaDateAccordingToOrderBy(lastMedia).toISOString();
        url = `/v1/mediaList?order=${orderBy.value}&sort=desc&per_page=30&last_media_id=${lastMedia.id}&last_date=${lastDate}`;
      }
      axios
        .get<Array<Media>>(url, {
          headers: {
            Authorization: `Bearer ${accessToken.value}`,
          },
        })
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

  function loadAllMediaForDate(date: Date): Promise<boolean> {
    return new Promise((resolve, reject) => {
      let lastMedia = mediaList.value[mediaList.value.length - 1];
      if (lastMedia.date > date) {
        resolve(true);
        return;
      }
      while (lastMedia.date > date && !allMediaLoaded.value) {
        loadMoreMedia()
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
      (media) => !mediaIDs.includes(media.id),
    );
  }

  function deleteMedia(mediaID: number): Promise<boolean> {
    return new Promise<boolean>((resolve, reject) => {
      axios
        .delete(`/v1/media/${mediaID}`, {
          headers: {
            Authorization: `Bearer ${accessToken.value}`,
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
    mediaIDs: Array<number>,
  ): Promise<Array<number>> {
    const failedMediaIDs = new Set<number>();
    const successMediaIDs = new Set<number>();

    for (const mediaID of mediaIDs) {
      try {
        const resp = await axios.delete(`/v1/media/${mediaID}`, {
          headers: {
            Authorization: `Bearer ${accessToken.value}`,
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
    appendMedia,
    orderByMediaDate,
    orderByUploadDate,
    getMediaDateAccordingToOrderBy,
  };
});
