import type { LocationQuery } from "vue-router";

// todo throw error?
export function base64UrlEncode(data: any) {
  try {
    if (data === undefined) {
      return undefined;
    }
    return window.btoa(JSON.stringify(data));
  } catch (err) {
    console.error(err);
    return "";
  }
}

// todo throw error?
export function base64UrlDecode(data: string) {
  try {
    if (data === undefined) {
      return undefined;
    }
    return JSON.parse(window.atob(data));
  } catch (err) {
    console.error(err);
    return "";
  }
}

export function timestamp(): string {
  const date = new Date();
  return `${date.getHours().toString().padStart(2, "0")}:${date
    .getMinutes()
    .toString()
    .padStart(2, "0")}:${date.getSeconds().toString().padStart(2, "0")}.${date
    .getMilliseconds()
    .toString()
    .padEnd(3, "0")}`;
}

export function getQueryParamStringValue(
  query: LocationQuery,
  key: string,
): string | undefined {
  try {
    const queryParam = query[key];
    if (queryParam && typeof queryParam === "string") {
      return queryParam;
    } else {
      return undefined;
    }
  } catch (err) {
    return undefined;
  }
}

export function getQueryParamNumberValue(
  query: LocationQuery,
  key: string,
): number | undefined {
  try {
    const queryParam = query[key];
    if (
      queryParam &&
      typeof queryParam === "string" &&
      !isNaN(Number(queryParam))
    ) {
      return Number(queryParam);
    } else {
      return undefined;
    }
  } catch (err) {
    return undefined;
  }
}

export class PromiseTimeoutError extends Error {}

export function promiseTimeout(
  promise: Promise<any>,
  timeoutMS: number,
): Promise<any> {
  let timer: NodeJS.Timeout;
  return Promise.race([
    promise,
    new Promise((_, reject) => {
      timer = setTimeout(
        reject,
        timeoutMS,
        new PromiseTimeoutError("timeout waiting for promise"),
      );
    }),
  ]).finally(() => {
    clearTimeout(timer);
  });
}

export function getFileType(file: File): string {
  const extension = file.name.split(".").pop()?.toLocaleLowerCase();
  switch (extension) {
    case "png":
      return "image/png";
    case "jpg":
      return "image/jpeg";
    case "jpeg":
      return "image/jpeg";
    case "webp":
      return "image/webp";
    case "mp4":
      return "video/mp4";
    case "webm":
      return "video/webm";
  }
  return "unknown";
}

// useful for creating keys for components to rerender them on route change
export function generateQueryParamsKey(
  query: LocationQuery,
  ...keys: Array<string>
): string | undefined {
  try {
    let componentKey = "";
    keys.forEach((key) => {
      componentKey += `${key}=${getQueryParamStringValue(query, key)};`;
    });
    return componentKey;
  } catch (err) {
    console.log(err);
    return "";
  }
}
