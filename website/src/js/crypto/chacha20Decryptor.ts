import { Chacha20 } from "ts-chacha20";

const padding = new TextEncoder().encode(
  "00000000000000000000000000000000000000000000000000000000000000000",
);

export function newChaCha20Decryptor(
  encryptionKey: string,
  nonce: string,
  offset: number,
): Chacha20 {
  const counter = Math.floor(offset / 64);
  const byteCounter = offset % 64;
  const textEncoder = new TextEncoder();
  const decryptor = new Chacha20(
    textEncoder.encode(encryptionKey),
    textEncoder.encode(nonce),
    counter,
  );
  // set the internal byte counter
  if (byteCounter !== 0) {
    decryptor.decrypt(padding.slice(0, byteCounter));
  }
  return decryptor;
}

export const newDecryptTransformer: (
  decryptor: Chacha20,
) => TransformStream<Uint8Array, Uint8Array> = (decryptor) =>
  new TransformStream<Uint8Array, Uint8Array>({
    start() {},
    transform(chunk, controller) {
      const decryptedChunk = decryptChunk(chunk, decryptor);
      controller.enqueue(decryptedChunk);
    },
    flush() {},
  });

export const decryptChunk: (
  input: Uint8Array,
  decryptor: Chacha20,
) => Uint8Array = (input, decryptor) => {
  const decrypted = decryptor.decrypt(input);
  return decrypted;
};
