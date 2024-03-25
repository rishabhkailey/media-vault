export const PDF_TYPE = "application/pdf";
export const JPEG_TYPE = "image/jpeg";
export const PNG_TYPE = "image/png";
export const WEBP_TYPE = "image/webp";
export const GIF_TYPE = "image/gif";
export const MP4_TYPE = "video/mp4";
export const WEBM_TYPE = "video/webm";
export const MOV_TYPE = "video/quicktime";
export const AVI_TYPE = "video/x-msvideo";
export const MP3_TYPE = "audio/mp3";
export const WAV_TYPE = "audio/wav";
export const OGG_TYPE = "audio/ogg";

export function isImage(contentType: string): boolean {
  return [JPEG_TYPE, PNG_TYPE, WEBP_TYPE, GIF_TYPE].includes(contentType);
}

export function isVideo(contentType: string): boolean {
  return [MP4_TYPE, WEBM_TYPE, MOV_TYPE, AVI_TYPE].includes(contentType);
}

export function isVideoFormatSupported(contentType: string): boolean {
  const videoElement = document.createElement("video");
  return videoElement.canPlayType(contentType) !== "";
}

export function isPdf(contentType: string): boolean {
  return contentType === PDF_TYPE;
}

export function isAudio(contentType: string): boolean {
  return [MP3_TYPE, WAV_TYPE, OGG_TYPE].includes(contentType);
}

export function getFileType(file: File): string {
  const extension = file.name.split(".").pop()?.toLocaleLowerCase();
  switch (extension) {
    case "png":
      return PNG_TYPE;
    case "jpg":
      return JPEG_TYPE;
    case "jpeg":
      return JPEG_TYPE;
    case "webp":
      return WEBP_TYPE;
    case "gif":
      return GIF_TYPE;
    case "mp4":
      return MP4_TYPE;
    case "webm":
      return WEBM_TYPE;
    case "pdf":
      return PDF_TYPE;
    case "mp3":
      return MP3_TYPE;
    case "wav":
      return WAV_TYPE;
    case "ogg":
      return OGG_TYPE;
    case "mov":
      return MOV_TYPE;
    case "avi":
      return AVI_TYPE;
  }
  return "unknown";
}
