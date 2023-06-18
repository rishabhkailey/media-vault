const channelName = "encryption-key-channel";

export class EncryptionKeyChannelClient {
  private _channel: BroadcastChannel;
  private _encryptionKey: string;
  constructor(encryptionKey: string) {
    this._channel = new BroadcastChannel(channelName);
    this._encryptionKey = encryptionKey;
    this._channel.addEventListener("message", (event) => {
      console.log("client: want encryption key message received");
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
  private _channel: BroadcastChannel;
  private _encryptionKey: string;
  constructor() {
    this._channel = new BroadcastChannel(channelName);
    this._encryptionKey = "";
    this._channel.addEventListener("message", (event) => {
      if (
        event.data?.encryptionKey &&
        typeof event.data.encryptionKey === "string" &&
        event.data.encryptionKey.length !== 0
      ) {
        console.log("worker: set encryption key message received");
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
    this._channel.postMessage({
      wantEncryptionKey: true,
    });
    // in ms
    const interval = 100;
    const timeout = 10000;

    let time = 0;
    while (time < timeout) {
      if (this.encryptionKey.length !== 0) {
        await new Promise((r) => setTimeout(r, interval));
        return this.encryptionKey;
      }
      time += interval;
    }
    throw new Error("timedout waiting for encryption key from clients");
  }
}
