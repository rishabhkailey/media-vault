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
}

interface Album {
  id: number;
  name: string;
  media_count: number;
  thumbnail_url: string;
  updated_at: Date;
  created_at: Date;
}
