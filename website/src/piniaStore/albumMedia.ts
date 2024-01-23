import { UNKNOWN_DATE } from "@/js/date";
import axios from "axios";
import { defineStore, storeToRefs } from "pinia";
import { computed, ref } from "vue";
import { useAuthStore } from "./auth";

// we will need lock or something else
// to prevent duplicates if the same request is called twice

type OrderBySearchParam = "added_at" | "uploaded_at" | "date";
type SortyBySearchParam = "asc" | "desc";
interface IOrderByProperties {
  responseAttribute: keyof AlbumMedia;
  // for sorting in asc order
  // album1 > album2 => 1
  // album1 < album2 => -1
  // album1 = album2 => 0
  compare: (am1: AlbumMedia, am2: AlbumMedia) => number;

  orderBySearchParam: OrderBySearchParam;
  sortBySearchParam: SortyBySearchParam;
}

type UserFriendlyOrderBy =
  | "Newest date first"
  | "Oldest date first"
  | "Newest uploaded first"
  | "Oldest uploaded first"
  | "Newest Added first"
  | "Oldest Added first";

const userFriendlyOrderByToProperties = new Map<
  UserFriendlyOrderBy,
  IOrderByProperties
>();
userFriendlyOrderByToProperties.set("Newest date first", {
  compare: (a1, a2) => a1.date.getTime() - a2.date.getTime(),
  responseAttribute: "date",
  orderBySearchParam: "date",
  sortBySearchParam: "desc",
});
userFriendlyOrderByToProperties.set("Oldest date first", {
  compare: (a1, a2) => a1.date.getTime() - a2.date.getTime(),
  responseAttribute: "date",
  orderBySearchParam: "date",
  sortBySearchParam: "asc",
});
userFriendlyOrderByToProperties.set("Newest uploaded first", {
  compare: (a1, a2) => a1.uploaded_at.getTime() - a2.uploaded_at.getTime(),
  responseAttribute: "uploaded_at",
  orderBySearchParam: "uploaded_at",
  sortBySearchParam: "desc",
});
userFriendlyOrderByToProperties.set("Oldest uploaded first", {
  compare: (a1, a2) => a1.uploaded_at.getTime() - a2.uploaded_at.getTime(),
  responseAttribute: "uploaded_at",
  orderBySearchParam: "uploaded_at",
  sortBySearchParam: "asc",
});
userFriendlyOrderByToProperties.set("Newest Added first", {
  compare: (a1, a2) => a1.added_at.getTime() - a2.added_at.getTime(),
  responseAttribute: "added_at",
  orderBySearchParam: "added_at",
  sortBySearchParam: "desc",
});
userFriendlyOrderByToProperties.set("Newest Added first", {
  compare: (a1, a2) => a1.added_at.getTime() - a2.added_at.getTime(),
  responseAttribute: "added_at",
  orderBySearchParam: "added_at",
  sortBySearchParam: "asc",
});

export const useAlbumMediaStore = defineStore("albumMedia", () => {
  const mediaList = ref<Array<AlbumMedia>>([]);
  // todo add lastMediaLoadFailed so we don't keep retrying failed request
  const allMediaLoaded = ref(false);
  const albumID = ref(0);
  const { accessToken } = storeToRefs(useAuthStore());
  const lastMediaId = ref<null | number>(null);
  const orderBy = ref<UserFriendlyOrderBy>("Newest Added first");
  const orderByProperties = computed<IOrderByProperties>(() => {
    const properties = userFriendlyOrderByToProperties.get(orderBy.value);
    if (properties === undefined) {
      throw new Error(
        `"${orderBy.value} "invalid order by value. unable to get properties for the selected order by.`,
      );
    }
    return properties;
  });

  // todo date getter?
  // will return the date according to which the media is sorted?

  function reset() {
    console.debug("reset album media store");
    mediaList.value = [];
    allMediaLoaded.value = false;
    albumID.value = 0;
    lastMediaId.value = null;
  }
  function setAlbumID(_albumID: number) {
    if (albumID.value != _albumID) {
      reset();
    }
    albumID.value = _albumID;
  }
  function appendMedia(_mediaList: Array<AlbumMedia>) {
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
      mediaList.value = [...mediaList.value, ...newMediaList];
    }
  }

  function loadMoreMedia(): Promise<boolean> {
    const perPage = 30;
    return new Promise<boolean>((resolve, reject) => {
      const url = new URL(
        `/v1/album/${albumID.value}/media`,
        import.meta.env.VITE_BASE_URL,
      );
      url.searchParams.append("per_page", perPage.toString());
      url.searchParams.append(
        "order",
        orderByProperties.value.orderBySearchParam,
      );
      url.searchParams.append(
        "sort",
        orderByProperties.value.sortBySearchParam,
      );

      if (lastMediaId.value !== null) {
        url.searchParams.append("last_media_id", lastMediaId.value.toString());
      }
      axios
        .get<Array<AlbumMedia>>(url.toString(), {
          headers: {
            Authorization: `Bearer ${accessToken.value}`,
          },
        })
        .then((response) => {
          console.log(response);
          if (response.status == 200) {
            if (response.data.length > 0) {
              lastMediaId.value = response.data[response.data.length - 1].id;
              appendMedia(response.data);
            }
            allMediaLoaded.value =
              response.data.length == 0 || response.data.length < perPage;
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

  function loadAllMediaUntil(date: Date): Promise<boolean> {
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

  function removeMediaByIDsFromLocalState(mediaIDs: Array<number>) {
    mediaList.value = mediaList.value.filter(
      (media) => !mediaIDs.includes(media.id),
    );
  }

  return {
    mediaList,
    allMediaLoaded,
    albumID,
    reset,
    setAlbumID,
    loadMoreMedia,
    loadAllMediaUntil,
    removeMediaByIDsFromLocalState,
  };
});
