// export const onmessage = function (event) {
//   console.log("message received", event)
// };

self.addEventListener("message", (message) => {
  console.log("message received", message);
  postMessage("sending response" + message.data);
});

self.addEventListener("install", function (event) {
  console.log("install");
});

self.addEventListener("activate", function (event) {
  console.log("Claiming control");
  return self.clients.claim();
});

self.addEventListener("fetch", (event: any) => {
  if (
    event?.request?.method !== "GET" ||
    (typeof event?.request?.url === "string" &&
      !new URL(event.request.url).pathname.startsWith("/v1"))
  ) {
    return;
  }
  console.log("fetch event", event);
  event.respondWith(
    fetch(event.request).then(function (response) {
      console.log(response);
      return response;
      // if(response.url.match(".mp4")){
      //   console.log(event);
      //   responseCloned = response.clone();
      //   responseCloned.arrayBuffer().then(
      //     buffer =>{
      //       let length = 100
      //       view = new Uint8Array(buffer,0,length)
      //       for(i=0,j=length - 1; i<j; i++,j--){
      //           view[i] = view[i] ^ view[j]
      //           view[j] = view[j] ^ view[i]
      //           view[i] = view[i] ^ view[j]
      //       }
      //     }
      //   )
      // }
      // return responseCloned;
    })
  );
});
