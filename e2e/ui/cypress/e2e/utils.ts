export function sortMediaByUploadDateDesc(mediaList: Array<Media>): Array<Media> {
  return mediaList.sort((m1, m2) => new Date(m2.uploaded_at).getTime() - new Date(m1.uploaded_at).getTime())
}