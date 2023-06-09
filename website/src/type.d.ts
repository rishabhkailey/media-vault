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
  size: number;
  thumbnail: boolean;
  url: string;
  thumbnail_url: string;
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
