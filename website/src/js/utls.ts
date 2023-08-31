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
