import { describe, expect, it, test } from "vitest";
import { proccessThumbnailConstraints } from "@/utils/thumbnail";
import type { thumbnailConstraints } from "@/utils/thumbnail";

const RESOLUTION_16_X_9 = 16 / 9;
const RESOLUTION_9_X_16 = 9 / 16;

function safetyCheck(
  offset: point,
  thumbnailResolution: WidthHeight,
  sourceResolution: WidthHeight,
  inputResolution: WidthHeight
) {
  return (
    offset.x >= 0 &&
    offset.x <= inputResolution.width &&
    offset.y >= 0 &&
    offset.y <= inputResolution.height &&
    thumbnailResolution.height > 0 &&
    thumbnailResolution.width > 0 &&
    sourceResolution.height > 0 &&
    sourceResolution.width > 0 &&
    offset.x + sourceResolution.width <= inputResolution.width &&
    offset.y + sourceResolution.height <= inputResolution.height
  );
}

describe("thumbnail constraint tests", () => {
  test("constraint.maxHeightWidth = source.height = source.width", () => {
    const constraint: thumbnailConstraints = {
      maxHeightWidth: 300,
    };
    const inputResolution: WidthHeight = {
      width: 300,
      height: 300,
    };
    const expected: {
      offset: point;
      destinationResolution: WidthHeight;
      sourceResolution: WidthHeight;
    } = {
      offset: {
        x: 0,
        y: 0,
      },
      destinationResolution: {
        width: 300,
        height: 300,
      },
      sourceResolution: inputResolution,
    };
    const received = proccessThumbnailConstraints(constraint, inputResolution);
    expect(
      safetyCheck(
        received.offset,
        received.destinationResolution,
        received.sourceResolution,
        inputResolution
      )
    ).toStrictEqual(true);
    expect(
      received.destinationResolution.width <= constraint.maxHeightWidth &&
        received.destinationResolution.height <= constraint.maxHeightWidth
    ).toBe(true);
    expect(received).toStrictEqual(expected);
  });

  // no resolution change only scale change
  test("constraint.maxHeightWidth < source.height = source.width", () => {
    const constraint: thumbnailConstraints = {
      maxHeightWidth: 200,
    };
    const inputResolution: WidthHeight = {
      width: 300,
      height: 300,
    };
    const expected: {
      offset: point;
      destinationResolution: WidthHeight;
      sourceResolution: WidthHeight;
    } = {
      offset: {
        x: 0,
        y: 0,
      },
      destinationResolution: {
        width: 200,
        height: 200,
      },
      sourceResolution: inputResolution,
    };
    const received = proccessThumbnailConstraints(constraint, inputResolution);
    expect(
      safetyCheck(
        received.offset,
        received.destinationResolution,
        received.sourceResolution,
        inputResolution
      )
    ).toStrictEqual(true);
    expect(
      received.destinationResolution.width <= constraint.maxHeightWidth &&
        received.destinationResolution.height <= constraint.maxHeightWidth
    ).toBe(true);
    expect(received).toStrictEqual(expected);
  });

  test("verticle long image", () => {
    const constraint: thumbnailConstraints = {
      maxHeightWidth: 200,
    };
    const inputResolution: WidthHeight = {
      width: 300,
      height: 3000,
    };
    const expected: {
      offset: point;
      destinationResolution: WidthHeight;
      sourceResolution: WidthHeight;
    } = {
      offset: {
        x: 0,
        y: 1233,
      },
      destinationResolution: {
        width: Math.floor((9 / 16) * 200),
        height: 200,
      },
      sourceResolution: {
        width: 300,
        height: Math.floor((16 / 9) * 300),
      },
    };
    const received = proccessThumbnailConstraints(constraint, inputResolution);
    expect(
      safetyCheck(
        received.offset,
        received.destinationResolution,
        received.sourceResolution,
        inputResolution
      )
    ).toStrictEqual(true);
    expect(
      received.destinationResolution.width <= constraint.maxHeightWidth &&
        received.destinationResolution.height <= constraint.maxHeightWidth
    ).toBe(true);
    expect(
      Math.floor(
        Math.abs(
          // eslint-disable-next-line prettier/prettier
        (received.destinationResolution.width / received.destinationResolution.height) -
            RESOLUTION_9_X_16
        )
      )
    ).toBe(0);
    expect(
      Math.floor(
        Math.abs(
          // eslint-disable-next-line prettier/prettier
        (received.sourceResolution.width / received.sourceResolution.height) -
            RESOLUTION_9_X_16
        )
      )
    ).toBe(0);
    expect(received).toStrictEqual(expected);
  });

  test("horizontal long image", () => {
    const constraint: thumbnailConstraints = {
      maxHeightWidth: 200,
    };
    const inputResolution: WidthHeight = {
      width: 3000,
      height: 300,
    };
    const expected: {
      offset: point;
      destinationResolution: WidthHeight;
      sourceResolution: WidthHeight;
    } = {
      offset: {
        x: 1233,
        y: 0,
      },
      destinationResolution: {
        width: 200,
        height: 112,
      },
      sourceResolution: {
        width: 533,
        height: 300,
      },
    };
    const received = proccessThumbnailConstraints(constraint, inputResolution);
    expect(
      safetyCheck(
        received.offset,
        received.destinationResolution,
        received.sourceResolution,
        inputResolution
      )
    ).toStrictEqual(true);
    expect(
      received.destinationResolution.width <= constraint.maxHeightWidth &&
        received.destinationResolution.height <= constraint.maxHeightWidth
    ).toBe(true);
    // resolution check
    expect(
      Math.floor(
        Math.abs(
          // eslint-disable-next-line prettier/prettier
        (received.destinationResolution.width / received.destinationResolution.height)  -
            RESOLUTION_16_X_9
        )
      )
    ).toBe(0);
    expect(
      Math.floor(
        Math.abs(
          // eslint-disable-next-line prettier/prettier
        (received.sourceResolution.width / received.sourceResolution.height)  -
            RESOLUTION_16_X_9
        )
      )
    ).toBe(0);
    expect(received).toStrictEqual(expected);
  });

  test("horizontal long image but aspect ratio < 16/9", () => {
    const constraint: thumbnailConstraints = {
      maxHeightWidth: 200,
    };
    const inputResolution: WidthHeight = {
      width: 1204,
      height: 856,
    };
    const expected: {
      offset: point;
      destinationResolution: WidthHeight;
      sourceResolution: WidthHeight;
    } = {
      offset: {
        x: 1233,
        y: 0,
      },
      destinationResolution: {
        width: 200,
        height: 112,
      },
      sourceResolution: {
        width: 533,
        height: 300,
      },
    };
    const received = proccessThumbnailConstraints(constraint, inputResolution);
    expect(
      safetyCheck(
        received.offset,
        received.destinationResolution,
        received.sourceResolution,
        inputResolution
      )
    ).toStrictEqual(true);
    expect(
      received.destinationResolution.width <= constraint.maxHeightWidth &&
        received.destinationResolution.height <= constraint.maxHeightWidth
    ).toBe(true);
    // resolution check
    expect(
      Math.floor(
        Math.abs(
          // eslint-disable-next-line prettier/prettier
        (received.destinationResolution.width / received.destinationResolution.height)  -
            RESOLUTION_16_X_9
        )
      )
    ).toBe(0);
    expect(
      Math.floor(
        Math.abs(
          // eslint-disable-next-line prettier/prettier
        (received.sourceResolution.width / received.sourceResolution.height)  -
            RESOLUTION_16_X_9
        )
      )
    ).toBe(0);
    // expect(received).toStrictEqual(expected);
  });
});
