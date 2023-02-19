const _sodium = require("libsodium-wrappers");

_sodium.ready.then(() => {
  const sodium = _sodium;

  let key = sodium.crypto_secretstream_xchacha20poly1305_keygen();

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

  let r = sodium.crypto_secretstream_xchacha20poly1305_pull(
    state_in,
    new Uint8Array([...c1, ...c2])
  );
  console.log(r);
  let [m, tag] = [sodium.to_string(r.message), r.tag];
  console.log(m);

  // let c1Part1 = c1.slice(0, c1.length / 2);
  // let c1Part2 = c1.slice(c1.length / 2);
  // let r1Part1 = sodium.crypto_secretstream_xchacha20poly1305_pull(state_in, c1Part1);
  // let [m1Part1, tag1Part1] = [sodium.to_string(r1Part1.message), r1Part1.tag];

  // let r1Part2 = sodium.crypto_secretstream_xchacha20poly1305_pull(state_in, c1Part2);
  // let [m1Part2, tag1Part2] = [sodium.to_string(r1Part2.message), r1Part2.tag];
  // console.log(m1Part1 + m1Part2);

  // let r1 = sodium.crypto_secretstream_xchacha20poly1305_pull(state_in, c1);
  // let [m1, tag1] = [sodium.to_string(r1.message), r1.tag];
  // console.log(m1);
  // let r2 = sodium.crypto_secretstream_xchacha20poly1305_pull(state_in, c2);
  // let [m2, tag2] = [sodium.to_string(r2.message), r2.tag];
  // console.log(m2);
});

// it requires exact amount of data to decrypt, can not be used as stream cipher
// and we cannot control chunk/range size requested by video player (or we can modify the request in service worker)?
