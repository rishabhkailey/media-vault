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

export function getMonthlyMedia(
  mediaList: Array<Media>,
  mediaDateGetter: (media: Media) => Date,
): Array<MonthlyMedia> {
  const monthlyMediaList: Array<MonthlyMedia> = [];
  const getKey = (month: number, year: number) => {
    return `${month} ${year}`;
  };
  const parseKey = (key: string) => {
    const [month, year] = key.split(" ");
    return { month: Number(month), year: Number(year) };
  };
  const monthlyMediaMap = new Map<string, Array<Media>>();
  mediaList.forEach((media) => {
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
    monthlyMediaMap.get(key)?.push(media);
  });
  monthlyMediaMap.forEach((mediaList, key) => {
    const { month, year } = parseKey(key);
    monthlyMediaList.push({
      year,
      month,
      media: mediaList,
    });
  });
  monthlyMediaList.sort((a, b) => {
    return (
      new Date(b.year, b.month, 1).getTime() -
      new Date(a.year, a.month, 1).getTime()
    );
  });
  return monthlyMediaList;
}

export function getDailyMedia(
  mediaList: Array<Media>,
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
  const dailyMediaMap = new Map<string, Array<Media>>();
  mediaList.forEach((media) => {
    let key: string;
    const date = mediaDateGetter(media);
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
    dailyMediaMap.get(key)?.push(media);
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
  dailyMediaList.sort((a, b) => b.date - a.date);
  return dailyMediaList;
}
