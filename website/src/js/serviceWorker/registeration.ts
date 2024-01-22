import decryptWorker from "@/worker/dist/bundle.js?url";
// import decryptWorker from "@/worker/decrypt?url";

export function updateOrRegisterServiceWorker(): Promise<ServiceWorker> {
  return new Promise<ServiceWorker>((resolve, reject) => {
    navigator.serviceWorker.getRegistration().then((registration) => {
      if (registration === undefined) {
        console.debug("worker not registered. registering new worker");
        registerServiceWorker()
          .then((serviceWorker) => {
            console.debug("new worker registered");
            resolve(serviceWorker);
            return;
          })
          .catch((err) => {
            console.debug("new worker registeration failed");
            reject(err);
            return;
          });
      } else {
        console.debug("worker already registered. updating the worker.");
        updateServiceWorker(registration)
          .then((serviceWorker) => {
            console.debug("updated worker");
            resolve(serviceWorker);
            return;
          })
          .catch((err) => {
            console.debug("update failed");
            reject(err);
            return;
          });
      }
    });
  });
}

function updateServiceWorker(
  registration: ServiceWorkerRegistration,
): Promise<ServiceWorker> {
  return new Promise<ServiceWorker>((resolve, reject) => {
    // if hard reload or service worker url changed
    if (
      navigator.serviceWorker.controller === null ||
      !navigator.serviceWorker.controller.scriptURL.endsWith(decryptWorker)
    ) {
      console.log(
        "looks like a hard reload. unregistering the existing worker and try to reregister worker",
      );
      // https://github.com/rishabhkailey/media-service/issues/2
      // https://stackoverflow.com/a/66816077
      return registration
        .unregister()
        .then((unregister) => {
          console.log(unregister);
          return registerServiceWorker();
        })
        .catch((err) => {
          console.log("worker unregisteration failed");
          reject(err);
        });
    }

    registration
      .update()
      .then(() => {
        console.log("updated");
        if (registration.active === null) {
          reject(new Error("got null service worker after update"));
          return;
        }
        resolve(registration.active);
        return;
      })
      .catch((err) => {
        reject(err);
        return;
      });
  });
}

// https://github.com/jimmywarting/StreamSaver.js/blob/master/mitm.html#L39
function registerServiceWorker(): Promise<ServiceWorker> {
  return new Promise<ServiceWorker>((resolve, reject) => {
    console.log("registering new worker");
    if ("serviceWorker" in navigator) {
      // unregister the existing service worker
      navigator.serviceWorker
        .register(decryptWorker, {
          scope: "./",
          type: "module",
        })
        .then((swReg) => {
          if (swReg.active !== null) {
            console.debug("Service Worker registsered");
            resolve(swReg.active);
            return;
          }
          const swRegTmp = swReg.installing || swReg.waiting;
          if (swRegTmp === null) {
            reject(new Error("got null service worker registration"));
            return;
          }
          let callback: () => void;
          console.debug("waiting for Service Worder to registser");
          swRegTmp.addEventListener(
            "statechange",
            (callback = () => {
              if (swRegTmp.state === "activated") {
                console.debug("Service Worder registed and active");
                swRegTmp.removeEventListener("statechange", callback);
                resolve(swRegTmp);
              }
            }),
          );
        })
        .catch((err) => {
          reject(err);
          return;
        });
    }
  });
}
