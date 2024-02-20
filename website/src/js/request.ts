export interface IRequestRange {
  unit: string;
  start: number;
  end: number;
}

export interface IResponseRange {
  unit: string;
  start: number;
  end: number;
  size: number;
}

export function getRequestRange(headers: Headers): IRequestRange | undefined {
  const rangeHeader = headers.get("Range");
  if (rangeHeader === null) {
    return undefined;
  }
  return parseRequestRangeHeader(rangeHeader);
}

export function parseRequestRangeHeader(range: string): IRequestRange {
  const parts = range.split("=");
  if (parts.length != 2) {
    throw new Error("Invalid range header " + range);
  }
  const unit = parts[0];
  const rangePart = parts[1];
  let end: number = -1;
  if (
    rangePart.split("-").length === 2 &&
    rangePart.split("-")[1].length !== 0
  ) {
    end = Number(rangePart.split("-")[1]);
  }
  const start = Number(rangePart.split("-")[0]);
  return {
    unit,
    start,
    end,
  };
}

export function parseResponseRangeHeader(range: string): IResponseRange {
  let parts = range.split(" ");
  if (parts.length != 2) {
    throw new Error("Invalid range header " + range);
  }
  const unit = parts[0];
  // rangeSizePart = 200-1000/67589
  const rangeSizePart = parts[1];
  parts = rangeSizePart.split("/");
  if (parts.length != 2) {
    throw new Error("Invalid range header " + range);
  }
  const size = Number(parts[1]);
  // rangePart = 200-1000
  const rangePart = parts[0];
  if (rangePart.split("-").length != 2) {
    throw new Error("Invalid range header " + range);
  }
  const start = Number(rangePart.split("-")[0]);
  const end = Number(rangePart.split("-")[1]);
  return {
    unit,
    start,
    end,
    size,
  };
}
