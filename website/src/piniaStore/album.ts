import { UNKNOWN_DATE } from "@/js/date";
import axios from "axios";
import { defineStore, storeToRefs } from "pinia";
import { computed, ref } from "vue";
import { useAlbumMediaStore } from "./albumMedia";
import { useMediaStore } from "./media";
import { useAuthStore } from "./auth";
import { useConfigStore } from "./config";
// todo we will need lock or something else
// to prevent duplicates if the same request is called twice

type OrderBySearchParam = "created_at" | "updated_at";
type SortyBySearchParam = "asc" | "desc";
interface IOrderByProperties {
  responseAttribute: keyof Album;
  // for sorting in asc order
  // album1 > album2 => 1
  // album1 < album2 => -1
  // album1 = album2 => 0
  compare: (album1: Album, album2: Album) => number;
  orderBySearchParam: OrderBySearchParam;
  sortBySearchParam: SortyBySearchParam;
}

// todo
type UserFriendlyOrderBy =
  | "Newest date first"
  | "Oldest date first"
  | "Newest updated first"
  | "Oldest updated first";

const userFriendlyOrderByToProperties = new Map<
  UserFriendlyOrderBy,
  IOrderByProperties
>();
userFriendlyOrderByToProperties.set("Newest date first", {
  compare: (a1, a2) => a1.created_at.getTime() - a2.created_at.getTime(),
  responseAttribute: "created_at",
  orderBySearchParam: "created_at",
  sortBySearchParam: "desc",
});
userFriendlyOrderByToProperties.set("Oldest date first", {
  compare: (a1, a2) => a1.created_at.getTime() - a2.created_at.getTime(),
  responseAttribute: "created_at",
  orderBySearchParam: "created_at",
  sortBySearchParam: "asc",
});
userFriendlyOrderByToProperties.set("Newest updated first", {
  compare: (a1, a2) => a1.updated_at.getTime() - a2.updated_at.getTime(),
  responseAttribute: "updated_at",
  orderBySearchParam: "updated_at",
  sortBySearchParam: "desc",
});
userFriendlyOrderByToProperties.set("Oldest updated first", {
  compare: (a1, a2) => a1.updated_at.getTime() - a2.updated_at.getTime(),
  responseAttribute: "updated_at",
  orderBySearchParam: "updated_at",
  sortBySearchParam: "asc",
});

// todo better function names to differentiate delete from backend or just remove from local state
export const useAlbumStore = defineStore("album", () => {
  const albumMediaStore = useAlbumMediaStore();
  const { mediaList } = storeToRefs(useMediaStore());
  const { accessToken } = storeToRefs(useAuthStore());
  const nextPageNumber = ref(1);
  const lastAlbumId = ref<null | number>(null);
  const albums = ref<Array<Album>>([]);
  const allAlbumsLoaded = ref(false);
  const orderBy = ref<UserFriendlyOrderBy>("Newest updated first");
  const orderByProperties = computed<IOrderByProperties>(() => {
    const properties = userFriendlyOrderByToProperties.get(orderBy.value);
    if (properties === undefined) {
      throw new Error(
        `"${orderBy.value} "invalid order by value. unable to get properties for the selected order by.`,
      );
    }
    return properties;
  });
  function getAlbumByID(albumID: number): Promise<Album> {
    return new Promise<Album>((resolve, reject) => {
      const album = albums.value.find((album) => album.id === albumID);
      if (album !== undefined) {
        resolve(album);
        return;
      }
      axios
        .get<Album>(`/v1/album/${albumID}`, {
          headers: {
            Authorization: `Bearer ${accessToken.value}`,
          },
        })
        .then((res) => {
          if (res.status !== 200) {
            reject(new Error("non 200 status code"));
            return;
          }
          addAlbumsInLocalState([res.data]);
          resolve(res.data);
        })
        .catch((err) => {
          reject(err);
          return;
        });
    });
  }

  function sortAndRemoveDuplicateAlbums(albums: Array<Album>): Array<Album> {
    if (albums.length <= 1) {
      return albums;
    }
    albums = albums.sort((a1, a2) => {
      if (orderByProperties.value.sortBySearchParam === "asc") {
        return orderByProperties.value.compare(a1, a2);
      }
      return -1 * orderByProperties.value.compare(a1, a2);
    });

    const uniqueAlbums: Array<Album> = [albums[0]];
    let previousAlbumID = albums[0].id;
    for (let index = 1; index < albums.length; index++) {
      if (albums[index].id === previousAlbumID) {
        continue;
      }
      uniqueAlbums.push(albums[index]);
      previousAlbumID = albums[index].id;
    }
    // albums = albums.filter(
    //   (album, index, self) =>
    //     index === self.findIndex((_album) => _album.id == album.id)
    // );
    return uniqueAlbums;
  }

  function addAlbumsInLocalState(_albums: Array<Album>) {
    if (_albums.length > 0) {
      const newAlbums = _albums
        .map((album) => {
          if (typeof album.updated_at === "string") {
            try {
              album.updated_at = new Date(album.updated_at);
            } catch (err) {
              album.updated_at = UNKNOWN_DATE;
            }
          }
          if (typeof album.created_at === "string") {
            try {
              album.created_at = new Date(album.updated_at);
            } catch (err) {
              album.created_at = UNKNOWN_DATE;
            }
          }
          return album;
        })
        .sort((m1, m2) => {
          if (orderByProperties.value.sortBySearchParam === "asc") {
            return orderByProperties.value.compare(m1, m2);
          }
          return -1 * orderByProperties.value.compare(m1, m2);
        });
      albums.value = sortAndRemoveDuplicateAlbums([
        ...albums.value,
        ...newAlbums,
      ]);
    }
  }

  function loadMoreAlbums(): Promise<LoadMoreAlbumStatus> {
    const perPage = 30;
    return new Promise<LoadMoreAlbumStatus>((resolve, reject) => {
      // todo fix: new URL("/abc", "http://localhost:5173/v1").toString() -> "http://localhost:5173/abc"
      const url = new URL("/v1/albums", window.location.origin);
      url.searchParams.append("per_page", perPage.toString());
      url.searchParams.append(
        "order",
        orderByProperties.value.orderBySearchParam,
      );
      url.searchParams.append(
        "sort",
        orderByProperties.value.sortBySearchParam,
      );

      if (lastAlbumId.value !== null) {
        url.searchParams.append("last_album_id", lastAlbumId.value.toString());
      }
      axios
        .get<Array<Album>>(url.toString(), {
          headers: {
            Authorization: `Bearer ${accessToken.value}`,
          },
        })
        .then((response) => {
          if (response.status == 200) {
            addAlbumsInLocalState(response.data);
            if (response.data.length > 0) {
              lastAlbumId.value = response.data[response.data.length - 1].id;
            }
            nextPageNumber.value += 1;
            if (response.data.length == 0 || response.data.length < perPage) {
              allAlbumsLoaded.value = true;
              resolve("empty");
              return;
            }
            resolve("ok");
            return;
          }
          reject(new Error("non 200 status code"));
          return;
        })
        .catch((err) => {
          reject(err);
        });
    });
  }
  // create album
  function createAlbum(name: string, thumbnailUrl: string): Promise<boolean> {
    return new Promise<boolean>((resolve, reject) => {
      axios
        .post<Album>(
          "/v1/album",
          {
            name: name,
            thumbnail_url: thumbnailUrl,
          },
          {
            headers: {
              "Content-Type": "application/json",
              Authorization: `Bearer ${accessToken.value}`,
            },
          },
        )
        .then((res) => {
          if (res.status === 200) {
            resolve(true);
          }
          addAlbumsInLocalState([res.data]);
          reject(new Error(`${res.status} status code`));
        })
        .catch((err) => {
          reject(err);
        });
    });
  }

  function removeAlbumByIDFromLocalState(albumID: number) {
    albums.value = albums.value.filter((album) => album.id !== albumID);
  }

  function removeAlbumsByIDFromLocalState(albumIDs: Array<number>) {
    albums.value = albums.value.filter((album) => !albumIDs.includes(album.id));
  }

  function deleteAlbum(albumID: number): Promise<boolean> {
    return new Promise<boolean>((resolve, reject) => {
      axios
        .delete(`/v1/album/${albumID}`, {
          headers: {
            Authorization: `Bearer ${accessToken.value}`,
          },
        })
        .then((resp) => {
          if (resp.status === 200) {
            removeAlbumByIDFromLocalState(albumID);
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

  // returns failed album ids
  // reject only in case of unexpected/unhandeled error
  async function deleteMultipleAlbums(
    albumIDs: Array<number>,
  ): Promise<Array<number>> {
    const failedAlbumIDs = new Set<number>();
    const successAlbumIDs = new Set<number>();

    for (const albumID of albumIDs) {
      try {
        const resp = await axios.delete(`/v1/album/${albumID}`, {
          headers: {
            Authorization: `Bearer ${accessToken.value}`,
          },
        });
        if (resp.status !== 200) {
          failedAlbumIDs.add(albumID);
          continue;
        } else {
          successAlbumIDs.add(albumID);
        }
      } catch {
        failedAlbumIDs.add(albumID);
      }
    }
    removeAlbumsByIDFromLocalState(Array.from(successAlbumIDs));
    return Array.from(failedAlbumIDs);
  }

  async function updateAlbumByID(albumID: number): Promise<boolean> {
    const album = await getAlbumByID(albumID);
    try {
      removeAlbumByIDFromLocalState(albumID);
      // this will get and add the album to list
      await getAlbumByID(albumID);
    } catch (err) {
      addAlbumsInLocalState([album]);
      throw err;
    }
    return true;
  }

  function removeMediaFromAlbum(
    albumID: number,
    mediaIDs: Array<number>,
  ): Promise<boolean> {
    return new Promise<boolean>((resolve, reject) => {
      axios
        .delete<{ media_ids: Array<number> }>(`/v1/album/${albumID}/media`, {
          headers: {
            Authorization: `Bearer ${accessToken.value}`,
          },
          data: { media_ids: mediaIDs },
        })
        .then((res) => {
          if (res.status != 200) {
            reject(new Error("non 200 status code"));
            return;
          }
          updateAlbumByID(albumID);
          albumMediaStore.removeMediaByIDsFromLocalState(res.data.media_ids);
          resolve(true);
        })
        .catch((err) => {
          reject(err);
        });
    });
  }

  async function addMediaToAlbum(
    albumID: number,
    mediaIDs: Array<number>,
  ): Promise<boolean> {
    // will call api if not present in local store
    let album = await getAlbumByID(albumID);
    const res = await axios.post<{ media_ids: Array<number> }>(
      `/v1/album/${albumID}/media`,
      {
        media_ids: mediaIDs,
      },
      {
        headers: {
          Authorization: `Bearer ${accessToken.value}`,
        },
      },
    );
    if (album.media_count == 0) {
      const thumbnailUrl = getFirstThumbnail(mediaIDs);
      if (thumbnailUrl.length != 0) {
        album = await updateAlbumThumbnail(albumID, thumbnailUrl);
      }
    }
    if (res.status != 200) {
      throw new Error("non 200 status code");
    }
    updateAlbumByID(albumID);
    return true;
  }

  function updateAlbumThumbnail(
    albumID: number,
    thumbnailUrl: string,
  ): Promise<Album> {
    return new Promise<Album>((resolve, reject) => {
      axios
        .patch(
          `/v1/album/${albumID}`,
          {
            thumbnail_url: thumbnailUrl,
          },
          {
            headers: {
              Authorization: `Bearer ${accessToken.value}`,
            },
          },
        )
        .then((res) => {
          if (res.status != 200) {
            reject(new Error("non 200 status code"));
            return;
          }
          resolve(res.data);
          return;
        })
        .catch((err) => {
          reject(err);
          return;
        });
    });
  }

  function getFirstThumbnail(mediaIDs: Array<number>): string {
    let thumbnailUrl = "";
    for (const mediaID of mediaIDs) {
      for (const media of mediaList.value) {
        if (media.id === mediaID && media.thumbnail) {
          thumbnailUrl = media.thumbnail_url;
          return thumbnailUrl;
        }
      }
    }
    return "";
  }

  return {
    nextPageNumber,
    albums,
    allAlbumsLoaded,
    loadMoreAlbums,
    deleteAlbum,
    deleteMultipleAlbums,
    createAlbum,
    getAlbumByID,
    addMediaToAlbum,
    removeMediaFromAlbum,
  };
});
