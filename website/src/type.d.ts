type WidthHeight = {
  width: number;
  height: number;
};

type point = {
  x: number;
  y: number;
}

interface Media {
  name: string;
  type: string;
  date: Date;
  size: number;
  thumbnail: boolean;
  url: string;
  thumbnail_url: string;
}
