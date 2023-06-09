import { UNKNOWN_DATE } from "@/utils/date";
import axios from "axios";
import { defineStore } from "pinia";
import { ref } from "vue";
import { useAlbumMediaStore } from "./albumMedia";

// we will need lock or something else
// to prevent duplicates if the same request is called twice

// todo better function names to differentiate delete from backend or just remove from local state
export const useAlbumStore = defineStore("album", () => {
  const albumMediaStore = useAlbumMediaStore();

  const nextPageNumber = ref(1);
  const albums = ref<Array<Album>>([]);
  const allAlbumsLoaded = ref(false);

  function getAlbumByID(accessToken: string, albumID: number): Promise<Album> {
    return new Promise<Album>((resolve, reject) => {
      const album = albums.value.find((album) => album.id === albumID);
      if (album !== undefined) {
        resolve(album);
        return;
      }
      axios
        .get<Album>(`/v1/album/${albumID}`, {
          headers: {
            Authorization: `Bearer ${accessToken}`,
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
      console.log(a1.updated_at, a2.updated_at);
      return a1.updated_at > a2.updated_at ? -1 : 1;
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
          console.log(m1.updated_at, m2.updated_at);
          return m1.updated_at > m2.updated_at ? -1 : 1;
        });
      // todo remove duplicates
      // let finalAlbums: Array<Album>;
      // if (prepend) {
      //   finalAlbums = [...newAlbums, ...albums.value];
      // } else {
      //   finalAlbums = [...albums.value, ...newAlbums];
      // }
      albums.value = sortAndRemoveDuplicateAlbums([
        ...albums.value,
        ...newAlbums,
      ]);
    }
  }

  function loadMoreAlbums(accessToken: string): Promise<boolean> {
    return new Promise<boolean>((resolve, reject) => {
      axios
        .get<Array<Album>>(
          `/v1/albums?page=${nextPageNumber.value}&perPage=30&order=updated_at&sort=desc`,
          {
            headers: {
              Authorization: `Bearer ${accessToken}`,
            },
          }
        )
        .then((response) => {
          console.log(response);
          if (response.status == 200) {
            addAlbumsInLocalState(response.data);
            nextPageNumber.value += 1;
            allAlbumsLoaded.value = response.data.length == 0;
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
  // create album
  function createAlbum(
    accessToken: string,
    name: string,
    thumbnailUrl: string
  ): Promise<boolean> {
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
              Authorization: `Bearer ${accessToken}`,
            },
          }
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

  function deleteAlbum(accessToken: string, albumID: number): Promise<boolean> {
    return new Promise<boolean>((resolve, reject) => {
      axios
        .delete(`/v1/album/${albumID}`, {
          headers: {
            Authorization: `Bearer ${accessToken}`,
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
    accessToken: string,
    albumIDs: Array<number>
  ): Promise<Array<number>> {
    const failedAlbumIDs = new Set<number>();
    const successAlbumIDs = new Set<number>();

    for (const albumID of albumIDs) {
      try {
        const resp = await axios.delete(`/v1/album/${albumID}`, {
          headers: {
            Authorization: `Bearer ${accessToken}`,
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

  async function updateAlbumByID(
    accessToken: string,
    albumID: number
  ): Promise<boolean> {
    const album = await getAlbumByID(accessToken, albumID);
    try {
      removeAlbumByIDFromLocalState(albumID);
      // this will get and add the album to list
      await getAlbumByID(accessToken, albumID);
    } catch (err) {
      addAlbumsInLocalState([album]);
      throw err;
    }
    return true;
  }

  function removeMediaFromAlbum(
    accessToken: string,
    albumID: number,
    mediaIDs: Array<number>
  ): Promise<boolean> {
    return new Promise<boolean>((resolve, reject) => {
      axios
        .delete<{ media_ids: Array<number> }>(`/v1/album/${albumID}/media`, {
          headers: {
            Authorization: `Bearer ${accessToken}`,
          },
          data: { media_ids: mediaIDs },
        })
        .then((res) => {
          if (res.status != 200) {
            reject(new Error("non 200 status code"));
            return;
          }
          updateAlbumByID(accessToken, albumID);
          albumMediaStore.removeMediaByIDsFromLocalState(res.data.media_ids);
          resolve(true);
        })
        .catch((err) => {
          reject(err);
        });
    });
  }
  function addMediaToAlbum(
    accessToken: string,
    albumID: number,
    mediaIDs: Array<number>
  ): Promise<boolean> {
    return new Promise<boolean>((resolve, reject) => {
      axios
        .post<{ media_ids: Array<number> }>(
          `/v1/album/${albumID}/media`,
          {
            media_ids: mediaIDs,
          },
          {
            headers: {
              Authorization: `Bearer ${accessToken}`,
            },
          }
        )
        .then((res) => {
          if (res.status != 200) {
            reject(new Error("non 200 status code"));
            return;
          }
          updateAlbumByID(accessToken, albumID);
          resolve(true);
          return;
        })
        .catch((err) => {
          reject(err);
          return;
        });
    });
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
