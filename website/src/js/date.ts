export const UNKNOWN_DATE = new Date("01/01/0100");

// not 0 indexed
export const daysLong = [
  "invalid",
  "Sunday",
  "Monday",
  "Tuesday",
  "Wednesday",
  "Thursday",
  "Friday",
  "Saturday",
];

// not 0 indexed
export const daysShort = ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"];

// 0 indexed
export const monthLong = [
  "January",
  "February",
  "March",
  "April",
  "May",
  "June",
  "July",
  "August",
  "September",
  "October",
  "November",
  "December",
];

// 0 indexed
export const monthShort = [
  "Jan",
  "Feb",
  "Mar",
  "Apr",
  "May",
  "June",
  "July",
  "Aug",
  "Sept",
  "Oct",
  "Nov",
  "Dec",
];

export function getMonthlyMediaIndex(
  mediaList: Array<Media>,
  mediaDateGetter: (media: Media) => Date,
) {
  const monthlyMediaList: Array<MonthlyMedia> = [];
  const getKey = (month: number, year: number) => {
    return `${month} ${year}`;
  };
  const parseKey = (key: string) => {
    const [month, year] = key.split(" ");
    return { month: Number(month), year: Number(year) };
  };
  const monthlyMediaMap = new Map<string, Array<IndexMedia>>();
  mediaList.forEach((media, index) => {
    let key: string;
    const date = mediaDateGetter(media);
    try {
      key = getKey(date.getMonth(), date.getFullYear());
    } catch (err) {
      key = getKey(UNKNOWN_DATE.getMonth(), UNKNOWN_DATE.getFullYear());
    }
    if (monthlyMediaMap.get(key) === undefined) {
      monthlyMediaMap.set(key, []);
    }
    monthlyMediaMap.get(key)?.push({ media, index });
  });
  monthlyMediaMap.forEach((mediaList, key) => {
    const { month, year } = parseKey(key);
    monthlyMediaList.push({
      year,
      month,
      media: mediaList,
      indexOffset: 0,
    });
  });
  monthlyMediaList.sort((a, b) => {
    if (a.year !== b.year) {
      return (b.year - a.year) * 12;
    }
    if (a.month !== b.month) {
      return a.month - b.month;
    }
    return 0;
  });
  return monthlyMediaList;
}

export function getDailyMediaIndex(
  mediaList: Array<IndexMedia>,
  mediaDateGetter: (media: Media) => Date,
) {
  const dailyMediaList: Array<DailyMedia> = [];
  const getKey = (day: number, date: number, month: number, year: number) => {
    return `${day} ${date} ${month} ${year}`;
  };
  const parseKey = (key: string) => {
    const [day, date, month, year] = key.split(" ");
    return {
      day: Number(day),
      date: Number(date),
      month: Number(month),
      year: Number(year),
    };
  };
  const dailyMediaMap = new Map<string, Array<IndexMedia>>();
  mediaList.forEach((indexMedia) => {
    let key: string;
    const date = mediaDateGetter(indexMedia.media);
    try {
      key = getKey(
        date.getDay(),
        date.getDate(),
        date.getMonth(),
        date.getFullYear(),
      );
    } catch (err) {
      key = getKey(
        UNKNOWN_DATE.getDay(),
        UNKNOWN_DATE.getDate(),
        UNKNOWN_DATE.getMonth(),
        UNKNOWN_DATE.getFullYear(),
      );
    }
    if (dailyMediaMap.get(key) === undefined) {
      dailyMediaMap.set(key, []);
    }
    dailyMediaMap.get(key)?.push(indexMedia);
  });
  dailyMediaMap.forEach((mediaList, key) => {
    const { day, date, month, year } = parseKey(key);
    dailyMediaList.push({
      day,
      date,
      year,
      month,
      media: mediaList,
    });
  });
  dailyMediaList.sort((a, b) => (a.date > b.date ? -1 : 1));
  return dailyMediaList;
}
