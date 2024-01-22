import { timestamp } from "../utils";

const channelName = "encryption-key-channel";

export class EncryptionKeyChannelClient {
  private _channel: BroadcastChannel;
  private _encryptionKey: string;
  constructor(encryptionKey: string) {
    this._channel = new BroadcastChannel(channelName);
    this._encryptionKey = encryptionKey;
    this._channel.addEventListener("message", (event) => {
      console.log(
        `${timestamp()}: client: want encryption key message received`
      );
      if (event.data.wantEncryptionKey !== true) {
        return;
      }
      this._channel.postMessage({
        encryptionKey: this._encryptionKey,
      });
    });
  }
  set encryptionKey(encryptionKey: string) {
    this._encryptionKey = encryptionKey;

    this._channel.postMessage({
      encryptionKey,
    });
  }
}

export class EncryptionKeyChannelWorker {
  private _requestEncryptionKeyInProgress: boolean;
  private _channel: BroadcastChannel;
  private _encryptionKey: string;
  constructor() {
    this._channel = new BroadcastChannel(channelName);
    this._encryptionKey = "";
    this._requestEncryptionKeyInProgress = false;
    this._channel.addEventListener("message", (event) => {
      if (
        event.data?.encryptionKey &&
        typeof event.data.encryptionKey === "string" &&
        event.data.encryptionKey.length !== 0
      ) {
        console.log(
          `${timestamp()}: worker: set encryption key message received`
        );
        this._encryptionKey = event.data.encryptionKey;
        return;
      }
    });
  }
  get encryptionKey(): string {
    return this._encryptionKey;
  }
  public async requestEncryptionKey(): Promise<string> {
    if (this.encryptionKey.length !== 0) {
      return this.encryptionKey;
    }
    if (this._requestEncryptionKeyInProgress === false) {
      console.debug(`${timestamp()}: requesting encryption key from client`);
      this._channel.postMessage({
        wantEncryptionKey: true,
      });
      this._requestEncryptionKeyInProgress = true;
    }
    try {
      // make a function for wait separate
      // in ms
      const interval = 100;
      const timeout = 10000;

      let time = 0;
      while (time < timeout) {
        if (this.encryptionKey.length !== 0) {
          return this.encryptionKey;
        }
        await new Promise((r) => setTimeout(r, interval));
        time += interval;
      }
    } catch (err) {
      throw new Error(`${timestamp()}: failed to watch encryptionKey update`);
    } finally {
      this._requestEncryptionKeyInProgress = false;
    }
    throw new Error(
      `${timestamp()}: timedout waiting for encryption key from clients`
    );
  }
}
