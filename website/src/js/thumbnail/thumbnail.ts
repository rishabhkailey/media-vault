import { getFileType } from "@/js/utils";
import { calculateThumbnailResolution } from "./resolution";

export function generateThumbnailAsArrayBuffer(
  file: File,
  constraints: thumbnailConstraints,
): Promise<{ thumbnail: Uint8Array; resolution: WidthHeight }> {
  return new Promise((resolve, reject) => {
    generateThumbnail(file, constraints)
      .then(({ thumbnail, resolution }) => {
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
          resolve({ thumbnail: new Uint8Array(arrayBuffer), resolution });
          return;
        };
        fileReader.readAsArrayBuffer(thumbnail);
      })
      .catch((err) => {
        reject(err);
        return;
      });
  });
}

function generateThumbnail(
  file: File,
  constraints: thumbnailConstraints,
): Promise<{ thumbnail: Blob; resolution: WidthHeight }> {
  if (getFileType(file).startsWith("image")) {
    return generateImageThumbnail(file, constraints);
  }
  if (getFileType(file).startsWith("video")) {
    return getnerateVideoThumbnail(file, constraints);
  }
  return new Promise((_, reject) => {
    return reject(new Error("unsupported file type"));
  });
}

function generateImageThumbnail(
  originalImage: Blob,
  constraints: thumbnailConstraints,
): Promise<{ thumbnail: Blob; resolution: WidthHeight }> {
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
          calculateThumbnailResolution(constraints, {
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
          destinationResolution.height,
        );
        // todo directry use canvas.toBlob
        fetch(canvas.toDataURL("image/jpeg"))
          .then((res) => {
            res
              .blob()
              .then((scaledImage: Blob) => {
                resolve({
                  thumbnail: scaledImage,
                  resolution: destinationResolution,
                });
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

// test cases
// 16/9 > maxHeightWidth
// 9/16 > maxHeightWidth
// horizontal long image > maxHeightWidth
// verticle long image > maxHeightWidth
// 16/9 < maxHeightWidth
// 9/16 < maxHeightWidth
// horizontal long image < maxHeightWidth
// verticle long image < maxHeightWidth

// returns
// offset - offset for the image according to the source resolution
// sourceResolution - cropped thumbnail but full resolution
// thumbnailResolution - cropped thumbnail with required resolution

/*
// cases for thumbnail resolution
if image is horizontal long
  thumbnail will also be long 16/9
  height will be set to max dimension

if image is verically long
  thumnail will also be vertically long 9/16
  height will be set to max dimension


  
// cases for source resolution
if thumbnail resolution > source resolution
  sourcewidth will remain same
  height will be cropped

if thumbnail resolution < source resolution
  sourceheight will not change
  widht will be cropped
*/

function getnerateVideoThumbnail(
  file: File,
  constraints: thumbnailConstraints,
): Promise<{ thumbnail: Blob; resolution: WidthHeight }> {
  return new Promise((resolve, reject) => {
    const video = document.createElement("video");

    video.src = URL.createObjectURL(file);
    video.onerror = (e) => {
      reject(e);
    };
    video.onloadedmetadata = (event) => {
      console.debug("metadata", event);
      console.debug("video duration", video.duration);
      // console.debug(video.fastSeek);
    };
    let requestVideoFrameCallbackCalled = false;
    video.oncanplay = async (event) => {
      if (requestVideoFrameCallbackCalled) {
        return;
      }
      requestVideoFrameCallbackCalled = true;
      console.debug("canplay", event);
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
      console.debug(now, metadata);
      requestVideoFrameCallbackCalled = true;
      let thumbnail: Blob | undefined;
      let resolution: WidthHeight | undefined;
      try {
        // video.height, width will be the html element width, height
        // which will always be 0, 0 as we are not rendering it on display
        ({ thumbnail, resolution } = await generateThumbnailOfCurrentFrame(
          video,
          {
            width: metadata.width,
            height: metadata.height,
          },
          constraints,
        ));
      } catch (err) {
        reject(err);
        video.remove();
        return;
      }
      if (thumbnail === undefined) {
        reject("got empty blob");
        return;
      }
      resolve({ thumbnail, resolution });
      return;
    });
  });
}

function generateThumbnailOfCurrentFrame(
  video: HTMLVideoElement,
  resolution: WidthHeight,
  constraints: thumbnailConstraints,
): Promise<{ thumbnail: Blob; resolution: WidthHeight }> {
  return new Promise((resolve, reject) => {
    const canvas = document.createElement("canvas");
    const ctx: CanvasRenderingContext2D | null = canvas.getContext("2d");
    if (!ctx) {
      reject(new Error("canvas.GetContext returned null`"));
      return;
    }
    const { destinationResolution, sourceResolution, offset } =
      calculateThumbnailResolution(constraints, {
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
      destinationResolution.height,
    );
    canvas.toBlob((thumbnail) => {
      if (thumbnail === null) {
        reject(new Error("got null"));
        return;
      }
      resolve({ thumbnail, resolution: destinationResolution });
      return;
    }, "image/jpeg");
  });
}
