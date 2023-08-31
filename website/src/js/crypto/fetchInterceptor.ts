import { type IRequestRange, getRequestRange } from "../request";
import {
  newChaCha20Decryptor,
  newDecryptTransformer,
} from "./chacha20Decryptor";

function getNonceFromUrl(url: string): string {
  let nonce = "";
  if (url.indexOf("/v1/media/") !== -1) {
    nonce = url.substring(url.indexOf("/v1/media/") + "/v1/media/".length);
  }
  if (url.indexOf("/v1/thumbnail/") !== -1) {
    nonce = url.substring(
      url.indexOf("/v1/thumbnail/") + "/v1/thumbnail/".length
    );
  }
  console.log("nonce", nonce);
  if (nonce.length === 0) {
    throw new Error("invalid url");
  }
  nonce = nonce.substring(0, 12);
  return nonce;
}

// todo better name
export async function internalFetch(
  req: Request,
  encryptionKeyGetter: () => Promise<string>
) {
  const range: IRequestRange | undefined = getRequestRange(req.headers);
  const nonce = getNonceFromUrl(req.url);
  console.log("nonce", nonce);
  const offset = range !== undefined ? range.start : 0;
  try {
    const encryptionKey = await encryptionKeyGetter();
    const decryptor = newChaCha20Decryptor(encryptionKey, nonce, offset);
    const res = await fetch(req);
    // return;
    const encryptedStream = res.body;
    if (encryptedStream === null) {
      return new Response("res");
    }
    const decryptedStream = encryptedStream.pipeThrough(
      newDecryptTransformer(decryptor)
    );
    return new Response(decryptedStream, {
      headers: res.headers,
      status: res.status,
      statusText: res.statusText,
    });
  } catch (err) {
    console.log(err);
    return new Response(undefined, {
      status: 500,
    });
  }
}
