import { fileType } from "./file";

export function generateThumbnailAsArrayBuffer(
  file: File,
  constraints: thumbnailConstraints
): Promise<Uint8Array> {
  return new Promise((resolve, reject) => {
    generateThumbnail(file, constraints)
      .then((blob) => {
        const fileReader = new FileReader();
        fileReader.onload = function (event) {
          if (
            event.target?.result === undefined ||
            event.target?.result === null ||
            typeof event.target?.result === "string" // todo check if string is valid return
          ) {
            reject("unable to convert blob to array buffer");
            return;
          }
          const arrayBuffer: ArrayBuffer = event.target.result;
          resolve(new Uint8Array(arrayBuffer));
          return;
        };
        fileReader.readAsArrayBuffer(blob);
      })
      .catch((err) => {
        reject(err);
        return;
      });
  });
}

export function generateThumbnail(
  file: File,
  constraints: thumbnailConstraints
): Promise<Blob> {
  if (fileType(file).startsWith("image")) {
    return generateImageThumbnail(file, constraints);
  }
  if (fileType(file).startsWith("video")) {
    return getnerateVideoThumbnail(file, constraints);
  }
  return new Promise((_, reject) => {
    return reject(new Error("unsupported file type"));
  });
}

function generateImageThumbnail(
  originalImage: Blob,
  constraints: thumbnailConstraints
): Promise<Blob> {
  return new Promise((resolve, reject) => {
    const fileReader: FileReader = new FileReader();
    const canvas = document.createElement("canvas");
    const ctx: CanvasRenderingContext2D | null = canvas.getContext("2d");
    if (!ctx) {
      reject(new Error("canvas.GetContext returned null`"));
      return;
    }
    fileReader.onload = (e: any) => {
      if (
        e?.targe?.result === undefined &&
        typeof e.target.result !== "string"
      ) {
        reject(new Error("file reader returned undefined/invalid data type"));
        return;
      }
      // canvas.src = e.target.result;
      const image: HTMLImageElement = document.createElement("img");
      image.src = e.target.result;
      image.onload = () => {
        const { destinationResolution, sourceResolution, offset } =
          proccessThumbnailConstraints(constraints, {
            width: image.width,
            height: image.height,
          });
        canvas.width = destinationResolution.width;
        canvas.height = destinationResolution.height;
        ctx.drawImage(
          image,
          offset.x,
          offset.y,
          sourceResolution.width,
          sourceResolution.height,
          0,
          0,
          destinationResolution.width,
          destinationResolution.height
        );
        // todo directry use canvas.toBlob
        fetch(canvas.toDataURL("image/jpeg"))
          .then((res) => {
            res
              .blob()
              .then((scaledImage: Blob) => {
                resolve(scaledImage);
              })
              .catch((err) => {
                reject(err);
              });

            return;
          })
          .catch((err) => {
            reject(err);
          });
      };
      image.onerror = () => {
        reject(new Error("image load failed"));
        return;
      };
    };
    fileReader.onerror = () => {
      reject(new Error("error in fileReader"));
      return;
    };
    fileReader.readAsDataURL(originalImage);
    return;
  });
}

export type thumbnailConstraints = {
  maxHeightWidth: number;
};
// test cases
// 16/9 > maxHeightWidth
// 9/16 > maxHeightWidth
// horizontal long image > maxHeightWidth
// verticle long image > maxHeightWidth
// 16/9 < maxHeightWidth
// 9/16 < maxHeightWidth
// horizontal long image < maxHeightWidth
// verticle long image < maxHeightWidth
const proccessThumbnailConstraints = (
  constraints: thumbnailConstraints,
  imageResolution: WidthHeight
) => {
  // sourceResolution + offset will be used to crop
  // check diagram here - https://developer.mozilla.org/en-US/docs/Web/API/CanvasRenderingContext2D/drawImage
  const sourceResolution: WidthHeight = {
    height: 0,
    width: 0,
  };
  const offset: point = {
    x: 0,
    y: 0,
  };

  // source image will be scaled to destinationResolution
  const destinationResolution: WidthHeight = {
    height: 0,
    width: 0,
  };
  let requiredAspectRatio: number = 1;
  let maxDimention = constraints.maxHeightWidth;
  if (imageResolution.height > imageResolution.width) {
    requiredAspectRatio = 9 / 16;
    // we don't want thumbnail bigger than image
    maxDimention = Math.min(maxDimention, imageResolution.height);
  } else if (imageResolution.width > imageResolution.height) {
    requiredAspectRatio = 16 / 9;
    // we don't want thumbnail bigger than image
    maxDimention = Math.min(maxDimention, imageResolution.width);
  } else {
    requiredAspectRatio = 1;
    maxDimention = Math.min(maxDimention, imageResolution.width);
  }
  const imageAspectRatio: number =
    imageResolution.width / imageResolution.height;
  if (imageAspectRatio > requiredAspectRatio) {
    // image is longer in the horizontal direction then reqruired
    // we will be ignoring horizontal si
    destinationResolution.width = maxDimention;
    destinationResolution.height =
      destinationResolution.width / requiredAspectRatio;

    // source image height will remain same only width will be cropped
    sourceResolution.height = imageResolution.height;
    sourceResolution.width = sourceResolution.height * requiredAspectRatio;

    // offset.y = actualImageWidth - actualImageWidthused / 2
    // eslint-disable-next-line prettier/prettier
    offset.x = Math.abs((imageResolution.width - (imageResolution.height * requiredAspectRatio))/2);
    offset.y = 0;
  } else {
    destinationResolution.height = maxDimention;
    destinationResolution.width =
      destinationResolution.height * requiredAspectRatio;

    // source image width will remain same only height will be cropped
    sourceResolution.width = imageResolution.width;
    sourceResolution.height = sourceResolution.width / requiredAspectRatio;
    // offset.y = actualImageHeight - actualImageHeightUsed / 2
    // eslint-disable-next-line prettier/prettier
    offset.y = Math.abs((imageResolution.height - (imageResolution.width / requiredAspectRatio))/2);
    offset.x = 0;
  }
  offset.x = Math.floor(offset.x);
  offset.y = Math.floor(offset.y);
  destinationResolution.width = Math.floor(destinationResolution.width);
  destinationResolution.height = Math.floor(destinationResolution.height);
  return { offset, destinationResolution, sourceResolution };
};
function getnerateVideoThumbnail(
  file: File,
  constraints: thumbnailConstraints
): Promise<Blob> {
  return new Promise((resolve, reject) => {
    const video = document.createElement("video");

    video.src = URL.createObjectURL(file);
    video.onerror = (e) => {
      reject(e);
    };
    video.onloadedmetadata = (event) => {
      console.log("metadata", event);
      console.log(video.duration);
      // console.log(video.fastSeek);
    };
    let requestVideoFrameCallbackCalled = false;
    let thumbnailGenerated = false;
    video.oncanplay = async (event) => {
      if (requestVideoFrameCallbackCalled) {
        return;
      }
      requestVideoFrameCallbackCalled = true;
      console.log("canplay", event);
      // we want the thumbnail to be in the first minute
      // for video longer than 1 minute = 30, for video less than 1 minute = duration/2
      const thumbnailTime = Math.min(30, Math.floor(video.duration / 2));
      // video.fastSeek(thumbnailTime);
      video.currentTime = thumbnailTime;
      video.volume = 0;
      await video.play();
      await video.pause();
    };
    video.requestVideoFrameCallback(async (now, metadata) => {
      console.log(now, metadata);
      thumbnailGenerated = true;
      requestVideoFrameCallbackCalled = true;
      let thumbnail: Blob | undefined;
      // thumbnailGenerated = true;
      try {
        // video.height, width will be the html element width, height
        // which will always be 0, 0 as we are not rendering it on display
        thumbnail = await generateThumbnailOfCurrentFrame(
          video,
          {
            width: metadata.width,
            height: metadata.height,
          },
          constraints
        );
      } catch (err) {
        reject(err);
        video.remove();
        return;
      }
      console.log("thumbnail", thumbnail);
      if (thumbnail === undefined) {
        reject("got empty blob");
        return;
      }
      resolve(thumbnail);
      return;
    });
  });
}

function generateThumbnailOfCurrentFrame(
  video: HTMLVideoElement,
  resolution: WidthHeight,
  constraints: thumbnailConstraints
): Promise<Blob> {
  return new Promise((resolve, reject) => {
    const canvas = document.createElement("canvas");
    const ctx: CanvasRenderingContext2D | null = canvas.getContext("2d");
    if (!ctx) {
      reject(new Error("canvas.GetContext returned null`"));
      return;
    }
    const { destinationResolution, sourceResolution, offset } =
      proccessThumbnailConstraints(constraints, {
        width: resolution.width,
        height: resolution.height,
      });
    canvas.width = destinationResolution.width;
    canvas.height = destinationResolution.height;
    ctx.drawImage(
      video,
      offset.x,
      offset.y,
      sourceResolution.width,
      sourceResolution.height,
      0,
      0,
      destinationResolution.width,
      destinationResolution.height
    );
    canvas.toBlob((thumbnail) => {
      if (thumbnail === null) {
        reject(new Error("got null"));
        return;
      }
      resolve(thumbnail);
      return;
    }, "image/jpeg");
  });
}
