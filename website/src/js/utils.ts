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
  defaultValue: string,
): string {
  try {
    const queryParam = query[key];
    if (queryParam && typeof queryParam === "string") {
      return queryParam;
    } else {
      return defaultValue;
    }
  } catch (err) {
    return defaultValue;
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
