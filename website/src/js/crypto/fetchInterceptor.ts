import { getNonceFromFileName } from "./utils";
import { type IRequestRange, getRequestRange } from "../request";
import {
  newChaCha20Decryptor,
  newDecryptTransformer,
} from "./chacha20Decryptor";

function getNonceFromUrl(url: string): string {
  const fileName = new RegExp(`/v1/file/(?<fileName>[^/]+)(/thumbnail)?$`).exec(
    url,
  )?.groups?.fileName;
  if (fileName === undefined) {
    throw new Error("invalid url");
  }
  return getNonceFromFileName(fileName);
}

// todo better name
export async function internalFetch(
  req: Request,
  encryptionKeyGetter: () => Promise<string>,
) {
  console.debug(`[internalFetch] ${req.url}`);
  const range: IRequestRange | undefined = getRequestRange(req.headers);
  const offset = range !== undefined ? range.start : 0;
  try {
    const nonce = getNonceFromUrl(req.url);
    const encryptionKey = await encryptionKeyGetter();
    const decryptor = newChaCha20Decryptor(encryptionKey, nonce, offset);
    const res = await fetch(req);
    // return;
    const encryptedStream = res.body;
    if (encryptedStream === null) {
      return new Response(undefined, {
        status: 500,
      });
    }
    const decryptedStream = encryptedStream.pipeThrough(
      newDecryptTransformer(decryptor),
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
