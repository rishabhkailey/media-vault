type WidthHeight = {
  width: number;
  height: number;
};

type point = {
  x: number;
  y: number;
};

interface Album {
  id: number;
  name: string;
  media_count: number;
  thumbnail_url: string;
  updated_at: Date;
  created_at: Date;
}

interface Media {
  id: number;
  name: string;
  type: string;
  date: Date;
  uploaded_at: Date;
  size: number;
  thumbnail: boolean;
  url: string;
  thumbnail_url: string;
  thumbnail_aspect_ratio: number;
  // todo
  // getSortBy(): string];
  // getDateAccordingToSortBy(): Date;
}

interface AlbumMedia extends Media {
  added_at: Date;
}

interface IndexMedia {
  media: Media;
  index: number;
}

interface MonthlyMedia {
  month: number;
  year: number;
  media: Array<IndexMedia>;
  indexOffset: number;
}

interface DailyMedia {
  month: number;
  day: number;
  year: number;
  date: number;
  media: Array<IndexMedia>;
}

interface UserInfo {
  id: string;
  prefered_timezone: string;
  encryption_key_checksum: string;
  storage_usage: number;
}

interface ThumbnailClickLocation {
  x: number;
  y: number;
  top: number;
  left: number;
  width: number;
  height: number;
}

type LoadMoreMediaStatus = "ok" | "empty";
type LoadMoreMedia = () => Promise<LoadMoreMediaStatus>;

type LoadMoreAlbumStatus = "ok" | "empty";
type LoadMoreAlbum = () => Promise<LoadMoreAlbumStatus>;
