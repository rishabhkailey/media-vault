const Chacha20 = require("ts-chacha20").Chacha20;
let key = "01234567890123456789012345678901";
let nonce = "012345678901";
let message =
  "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijkl";
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
var mergedEncryptedMessage = new Uint8Array(
  encryptedMessage.length + encryptedMessage2.length
);
mergedEncryptedMessage.set(encryptedMessage);
mergedEncryptedMessage.set(encryptedMessage2, encryptedMessage.length);

let multiDecryptedMessage = new Chacha20(
  textEncoder.encode(key),
  textEncoder.encode(nonce)
).decrypt(mergedEncryptedMessage);
console.log(textDecoder.decode(multiDecryptedMessage));

// half correct encrypted message (first half we will use normal message instead of encrypted message)
// we were successfuly able to decrypt the second half without the correct first half :)
mergedEncryptedMessage = new Uint8Array(
  message.length + encryptedMessage2.length
);
mergedEncryptedMessage.set(textEncoder.encode(message));
mergedEncryptedMessage.set(encryptedMessage2, encryptedMessage.length);

multiDecryptedMessage = new Chacha20(
  textEncoder.encode(key),
  textEncoder.encode(nonce)
).decrypt(mergedEncryptedMessage);
console.log(textDecoder.decode(multiDecryptedMessage));

// custom code
// decrypt with counter (will directly decrypt encryptedMessage2)
// to check how counter work counter
// let temp = new Chacha20(textEncoder.encode(key), textEncoder.encode(nonce));
// console.log(temp.counter);
// console.log(textDecoder.decode(temp.decrypt(encryptedMessage)));
// console.log(encryptedMessage.length, temp.counter);

// let decryptedMessage2 = new Chacha20(
//   textEncoder.encode(key),
//   textEncoder.encode(nonce),
//   1, // 64 bytes this is also like a block :(
//   1
// ).decrypt(encryptedMessage2.slice(1));
// console.log(textDecoder.decode(decryptedMessage2), decryptedMessage2.length);

{
  let encryptor = new Chacha20(
    textEncoder.encode(key),
    textEncoder.encode(nonce)
  );
  const encryptedMessage = encryptor.encrypt(textEncoder.encode(message));
  let chunkSize = 10;
  let useless = textEncoder.encode(
    "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
  );
  for (let i = 0; i < encryptedMessage.length; i += chunkSize) {
    let start = i;
    let end = start + chunkSize;
    let counter = Math.floor(i / 64);
    let byteCounter = i % 64;
    let decryptor = new Chacha20(
      textEncoder.encode(key),
      textEncoder.encode(nonce),
      counter
    );
    // set the internal byte counter
    if (byteCounter !== 0) {
      decryptor.decrypt(useless.slice(0, byteCounter));
    }
    let decryptedmessage = decryptor.decrypt(
      encryptedMessage.slice(start, end)
    );
    console.log(textDecoder.decode(decryptedmessage));
  }
}

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

// even this has issues
// this is a stream cipher encryption so it will not be able to decrypt older (i.e. already decrypted) chunk.
// try counter while creating chcha20 instance for encrypting/decrypting data not starting from the start
// https://crypto.stackexchange.com/questions/70094/does-chacha20-counter-actually-increment-through-iterations

// if we want to decrypt a random think can we just append it with some random value and decrypt
