export function calculateThumbnailResolution(
  constraints: thumbnailConstraints,
  imageResolution: WidthHeight,
): {
  offset: point;
  destinationResolution: WidthHeight;
  sourceResolution: WidthHeight;
} {
  if (constraints.preserveAspectRatio) {
    return calculateThumbnailResolutionWithOriginalAspectRatio(
      constraints,
      imageResolution,
    );
  }
  return calculate_16x9_9x16_resolution(constraints, imageResolution);
}

function calculateThumbnailResolutionWithOriginalAspectRatio(
  constraints: thumbnailConstraints,
  imageResolution: WidthHeight,
): {
  offset: point;
  destinationResolution: WidthHeight;
  sourceResolution: WidthHeight;
} {
  const imageAspectRatio = imageResolution.width / imageResolution.height;
  if (imageResolution.width > imageResolution.height) {
    const width = Math.min(imageResolution.width, constraints.maxHeightWidth);
    const height = Math.floor(width / imageAspectRatio);
    return {
      offset: {
        x: 0,
        y: 0,
      },
      destinationResolution: {
        width: width,
        height: height,
      },
      sourceResolution: imageResolution,
    };
  } else {
    const height = Math.min(imageResolution.height, constraints.maxHeightWidth);
    const width = Math.floor(height * imageAspectRatio);
    return {
      offset: {
        x: 0,
        y: 0,
      },
      destinationResolution: {
        width: width,
        height: height,
      },
      sourceResolution: imageResolution,
    };
  }
}

function calculate_16x9_9x16_resolution(
  constraints: thumbnailConstraints,
  imageResolution: WidthHeight,
): {
  offset: point;
  destinationResolution: WidthHeight;
  sourceResolution: WidthHeight;
} {
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
    destinationResolution.height = maxDimention;
    destinationResolution.width =
      destinationResolution.height * requiredAspectRatio;
  } else if (imageResolution.width > imageResolution.height) {
    requiredAspectRatio = 16 / 9;
    // we don't want thumbnail bigger than image
    maxDimention = Math.min(maxDimention, imageResolution.width);
    destinationResolution.width = maxDimention;
    destinationResolution.height =
      destinationResolution.width / requiredAspectRatio;
  } else {
    requiredAspectRatio = 1;
    maxDimention = Math.min(maxDimention, imageResolution.width);
    destinationResolution.width = maxDimention;
    destinationResolution.height = maxDimention;
  }
  const imageAspectRatio: number =
    imageResolution.width / imageResolution.height;
  if (imageAspectRatio > requiredAspectRatio) {
    // image is longer in the horizontal direction then reqruired
    // we will be ignoring horizontal sides
    // source image height will remain same only width will be cropped
    sourceResolution.height = imageResolution.height;
    sourceResolution.width = Math.floor(
      sourceResolution.height * requiredAspectRatio,
    );

    // offset.y = actualImageWidth - actualImageWidthused / 2
    // eslint-disable-next-line prettier/prettier
    offset.x = Math.abs((imageResolution.width - (imageResolution.height * requiredAspectRatio))/2);
    offset.y = 0;
  } else {
    // source image width will remain same only height will be cropped
    sourceResolution.width = imageResolution.width;
    sourceResolution.height = Math.floor(
      sourceResolution.width / requiredAspectRatio,
    );
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
}
