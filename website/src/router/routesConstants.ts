import { base64UrlEncode } from "@/js/utils";
import type { LocationQueryRaw, RouteLocationRaw } from "vue-router";

export const ABOUT_ROUTE_NAME = "about";
export const ERROR_SCREEN_ROUTE_NAME = "errorscreen";
export const NOT_FOUND_ROUTE_NAME = "NotFound";
export const HOME_ROUTE_NAME = "Home";
export const MEDIA_PREVIEW_ROUTE_NAME = "MediaPreview";
export const SEARCH_ROUTE_NAME = "search";
export const SEARCH_MEDIA_PREVIEW_ROUTE_NAME = "SearchMediaPreview";
export const ALBUMS_ROUTE_NAME = "Albums";
export const ALBUM_ROUTE_NAME = "Album";
export const ALBUM_MEDIA_PREVIEW_ROUTE_NAME = "AlbumMediaPreview";
export const ENTER_ENCRYPTION_KEY_ROUTE_NAME = "encryptionKey";
export const INITIAL_SETUP_ROUTE_NAME = "initialSetup";
export const PKCE_ROUTE_NAME = "PKCE";

export function aboutRoute(): RouteLocationRaw {
  return {
    name: ABOUT_ROUTE_NAME,
  };
}

export function errorScreenRoute(
  title: string,
  error: any,
  returnUri?: string,
): RouteLocationRaw {
  const query: LocationQueryRaw = {
    title: title,
  };

  if (returnUri !== undefined) {
    query.return_uri = returnUri;
  }

  if (typeof error === "string") {
    query.message = error;
  }
  if (error instanceof Error) {
    query.message = error.message + "\n" + error.stack;
  }

  return {
    name: ERROR_SCREEN_ROUTE_NAME,
    query: query,
  };
}

export function homeRoute(): RouteLocationRaw {
  return {
    name: HOME_ROUTE_NAME,
  };
}

export function albumsRoute(): RouteLocationRaw {
  return {
    name: ALBUMS_ROUTE_NAME,
  };
}

export function albumMediaPreviewRoute(
  index: number,
  mediaId: number,
  albumId: number,
  thumbnailClickLocation?: ThumbnailClickLocation,
): RouteLocationRaw {
  return {
    name: ALBUM_MEDIA_PREVIEW_ROUTE_NAME,
    params: {
      index: index,
      media_id: mediaId,
      album_id: albumId,
    },
    hash: `#${base64UrlEncode(thumbnailClickLocation)}`,
  };
}

export function albumRoute(albumId: number): RouteLocationRaw {
  return {
    name: ALBUM_ROUTE_NAME,
    params: {
      album_id: albumId,
    },
  };
}

export function searchRoute(query: string): RouteLocationRaw {
  return {
    name: SEARCH_ROUTE_NAME,
    params: {
      query: query,
    },
  };
}

export function searchMediaPreviewRoute(
  index: number,
  mediaId: number,
  query: string,
  thumbnailClickLocation?: ThumbnailClickLocation,
): RouteLocationRaw {
  return {
    name: SEARCH_MEDIA_PREVIEW_ROUTE_NAME,
    params: {
      index: index,
      media_id: mediaId,
      query: query,
    },
    hash: `#${base64UrlEncode(thumbnailClickLocation)}`,
  };
}

export function mediaPreviewRoute(
  index: number,
  mediaId: number,
  thumbnailClickLocation?: ThumbnailClickLocation,
): RouteLocationRaw {
  return {
    name: MEDIA_PREVIEW_ROUTE_NAME,
    params: {
      index: index,
      media_id: mediaId,
    },
    hash: `#${base64UrlEncode(thumbnailClickLocation)}`,
  };
}
