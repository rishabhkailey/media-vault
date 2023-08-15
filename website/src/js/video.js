// import MP4Box from "mp4box";
// type scaleImageCallback = (originalImage: Blob, scaledImage: Blob) => void;
// type onErrorCallback = (error: Error) => void;

// export const generateVideoThumbnail: (
//   originalImage: Blob,
//   onSuccess: scaleImageCallback,
//   onError: onErrorCallback
// ) => void = (videoFile, onSuccess, onError) => {
//   const mp4boxfile = MP4Box.createFile();
//   let fileInfoCompleted = false;
//   mp4boxfile.onError = (e: any) => {
//     fileInfoCompleted = true;
//     console.log(e);
//   };
//   mp4boxfile.onReady = (info: any) => {
//     fileInfoCompleted = true;
//     console.log("mp4box is ready.", info);
//   };
//   mp4boxfile.onMoovStart = () => {
//     fileInfoCompleted = true;
//     console.log("Starting to receive File Information");
//   };
//   const reader = videoFile.stream().getReader();
//   let buffer = new ArrayBuffer(4000);
//   let bytesReceived = 0;
//   let offset = 0;

//   while (offset < buffer.byteLength && !fileInfoCompleted) {
//     // reader.read(new Uint8Array(buffer, offset, buffer.byteLength - offset))
//     reader.read().then(function processBytes({ done, value }): any {
//       // Result objects contain two properties:
//       // done  - true if the stream has already given all its data.
//       // value - some data. Always undefined when done is true.

//       if (done) {
//         // There is no more data in the stream
//         return;
//       }

//       mp4boxfile.appendBuffer(value);

//       buffer = value.buffer;
//       offset += value.byteLength;
//       bytesReceived += value.byteLength;

//       // Read some more, and call this function again
//       return reader.read().then(processBytes);
//       // .read(new Uint8Array(buffer, offset, buffer.byteLength - offset))
//     });
//   }
// };


export const getVideoThumbnail = async (video) => {
  if (window.MediaStreamTrackProcessor) {
    let stopped = false;
    const track = await getVideoTrack();
    // eslint-disable-next-line no-undef
    const processor = new MediaStreamTrackProcessor(track);
    const reader = processor.readable.getReader();
    readChunk();
    
    function readChunk() {
      reader.read().then(async ({ done, value }) => {
        if (value) {
          const bitmap = await createImageBitmap(value);
          const index = frames.length;
          frames.push(bitmap);
          select.append(new Option("Frame #" + (index + 1), index));
          value.close();
        }
        if (!done && !stopped) {
          readChunk();
        } else {
          select.disabled = false;
        }
      });
    }
    button.onclick = (evt) => stopped = true;
    button.textContent = "stop";
  } else {
    console.error("your browser doesn't support this API yet");
  }
};

select.onchange = (evt) => {
  const frame = frames[select.value];
  canvas.width = frame.width;
  canvas.height = frame.height;
  ctx.drawImage(frame, 0, 0);
};

async function getVideoTrack() {
  const video = document.createElement("video");
  video.crossOrigin = "anonymous";
  video.src = "https://upload.wikimedia.org/wikipedia/commons/a/a4/BBH_gravitational_lensing_of_gw150914.webm";
  document.body.append(video);
  await video.play();
  const [track] = video.captureStream().getVideoTracks();
  video.onended = (evt) => track.stop();
  return track;
}

// https://jsfiddle.net/yzorbwg0/13/