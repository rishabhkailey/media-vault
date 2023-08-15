// todo try mediainfo.js
export function fileType(file: File): string {
  const extension = file.name.split(".").pop()?.toLocaleLowerCase();
  switch (extension) {
    case "png":
      return "image/png";
    case "jpg":
      return "image/jpeg";
    case "jpeg":
      return "image/jpeg";
    case "mp4":
      return "video/mp4";
    case "webm":
      return "video/webm";
  }
  return "unknown";
}
