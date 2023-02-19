<script setup lang="ts">
import { ref, onMounted } from "vue";
import decryptWorker from "@/worker/decrypt?url";
import sodium from "libsodium-wrappers";
let soidumReady: boolean = false;
const video = ref<HTMLVideoElement | undefined>(undefined);
let worker: Worker | undefined = undefined;
onMounted(() => {
  sodium.ready
    .then(() => {
      soidumReady = true;
      // encryption and decryption
      // todo check best bractices like using random/different nonce every time
      // https://doc.libsodium.org/advanced/stream_ciphers/chacha20
      // https://www.npmjs.com/package/libsodium-wrappers?activeTab=readme
      let password = "berrysikure";
      if (password.length == 0) {
        return;
      }
      while (password.length < 32) {
        password += password;
      }
      password = password.substring(0, 32);
      const encoder = new TextEncoder();
      let key = encoder.encode(password);
      console.log(sodium);
      console.log(
        sodium.crypto_secretstream_xchacha20poly1305_keygen().length,
        key.length
      );
      let res = sodium.crypto_secretstream_xchacha20poly1305_init_push(key);
      let [state_out, header] = [res.state, res.header];
      let c1 = sodium.crypto_secretstream_xchacha20poly1305_push(
        state_out,
        sodium.from_string("message 1"),
        null,
        sodium.crypto_secretstream_xchacha20poly1305_TAG_MESSAGE
      );
      let c2 = sodium.crypto_secretstream_xchacha20poly1305_push(
        state_out,
        sodium.from_string("message 2"),
        null,
        sodium.crypto_secretstream_xchacha20poly1305_TAG_FINAL
      );
      let state_in = sodium.crypto_secretstream_xchacha20poly1305_init_pull(
        header,
        key
      );
      let r1 = sodium.crypto_secretstream_xchacha20poly1305_pull(state_in, c1);
      let [m1, tag1] = [sodium.to_string(r1.message), r1.tag];
      let r2 = sodium.crypto_secretstream_xchacha20poly1305_pull(state_in, c2);
      let [m2, tag2] = [sodium.to_string(r2.message), r2.tag];

      console.log(m1);
      console.log(m2);
    })
    .catch((err) => {
      console.log(err);
    });
  if ("serviceWorker" in navigator) {
    navigator.serviceWorker
      .register(decryptWorker, {
        scope: "./",
        type: "module",
      })
      .then((registration) => {
        if (registration.installing) {
          console.log("Service worker installing");
        } else if (registration.waiting) {
          console.log("Service worker installed");
        } else if (registration.active) {
          console.log("Service worker active");
        }
      })
      .catch((error) => {
        console.error(`Registration failed with ${error}`);
      });
  }
  fetch("/v1/userinfo");
  // if (window.Worker) {
  //   worker = new decryptWorker();
  //   console.log(worker);
  //   worker.postMessage("Hi");
  //   fetch("/v1/userinfo");
  // }
});
</script>

<template>
  <video ref="video"></video>
</template>
// custom stream source //
https://stackoverflow.com/questions/70502425/javascript-how-to-modify-the-current-response-in-service-worker
// in aes encryption block size is 16 bytes, so we may need to ignore remainder
and update the range while decrypting accordingly // let's use chach20 instead
it is a stream ciphor (bit by bit) encryption not block cipher
https://security.stackexchange.com/questions/173758/can-i-use-streaming-decryption-by-every-algorithm
// https://www.npmjs.com/package/chacha //
https://www.npmjs.com/package/ts-chacha20 //
https://www.npmjs.com/package/libsodium //
https://www.npmjs.com/package/vue-upload-component?activeTab=readme
