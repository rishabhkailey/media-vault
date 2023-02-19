const Chacha20 = require("ts-chacha20").Chacha20;
let key = "01234567890123456789012345678901";
let nonce = "012345678901";
let message = "abcdefghijklmnopqrstuvwxyz";
let message2 = "secret message";
const textEncoder = new TextEncoder();
const textDecoder = new TextDecoder();
console.log(Chacha20);

let encryptor = new Chacha20(
  textEncoder.encode(key),
  textEncoder.encode(nonce)
);
const encryptedMessage = encryptor.encrypt(textEncoder.encode(message));
const encryptedMessage2 = encryptor.decrypt(textEncoder.encode(message2));

console.log(textDecoder.decode(encryptedMessage));

// decrypt full
const decryptedMessage = new Chacha20(
  textEncoder.encode(key),
  textEncoder.encode(nonce)
).decrypt(encryptedMessage);
console.log(textDecoder.decode(decryptedMessage));


// half message
const halfDecryptedMessage = new Chacha20(
  textEncoder.encode(key),
  textEncoder.encode(nonce)
).decrypt(encryptedMessage.slice(0, encryptedMessage.length / 2));
console.log(textDecoder.decode(halfDecryptedMessage));

// byte by byte
let i = 0;
let decryptor = new Chacha20(
  textEncoder.encode(key),
  textEncoder.encode(nonce)
);
while (i < encryptedMessage.length) {
  let decryptedChar = decryptor.decrypt(encryptedMessage.slice(i, i + 1));
  console.log(textDecoder.decode(decryptedChar));
  i++;
}

// multiple input combined

var mergedEncryptedMessage = new Uint8Array(encryptedMessage.length + encryptedMessage2.length);
mergedEncryptedMessage.set(encryptedMessage);
mergedEncryptedMessage.set(encryptedMessage2, encryptedMessage.length);

let multiDecryptedMessage = new Chacha20(
  textEncoder.encode(key),
  textEncoder.encode(nonce)
).decrypt(mergedEncryptedMessage);
console.log(textDecoder.decode(multiDecryptedMessage));

// random/different nonce (doesn't work)
// i = 0;
// nonce = "890101234567";
// decryptor = new Chacha20(
//   textEncoder.encode(key),
//   textEncoder.encode(nonce)
// );
// while (i < encryptedMessage.length) {
//   let decryptedChar = decryptor.decrypt(encryptedMessage.slice(i, i + 1));
//   console.log(textDecoder.decode(decryptedChar));
//   i++;
// }

// benifit over libsodium
// we can use it like a stream can work like byte by byte but the chacha20 instance should be same for 1 stream.
