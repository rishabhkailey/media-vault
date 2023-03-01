// fileReader is async so everything needs to be done via callbacks
type scaleImageCallback = (originalImage: Blob, scaledImage: Blob) => void;
type onErrorCallback = (error: Error) => void;
export const generateThumbnail: (
  originalImage: Blob,
  constraints: thumbnailConstraints,
  onSuccess: scaleImageCallback,
  onError: onErrorCallback
) => void = (originalImage, constraints, onSuccess, onError) => {
  const fileReader: FileReader = new FileReader();
  const canvas = document.createElement("canvas");
  const ctx: CanvasRenderingContext2D | null = canvas.getContext("2d");
  if (!ctx) {
    onError(new Error("canvas.GetContext returned null`"));
    return;
  }
  fileReader.onload = (e: any) => {
    if (e?.targe?.result === undefined && typeof e.target.result !== "string") {
      onError(new Error("file reader returned undefined/invalid data type"));
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
              onSuccess(originalImage, scaledImage);
            })
            .catch((err) => {
              onError(err);
            });

          return;
        })
        .catch((err) => {
          onError(err);
        });
    };
    image.onerror = () => {
      onError(new Error("image load failed"));
      return;
    };
  };
  fileReader.onerror = () => {
    onError(new Error("error in fileReader"));
    return;
  };
  fileReader.readAsDataURL(originalImage);
  return;
};

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
export const proccessThumbnailConstraints = (
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
