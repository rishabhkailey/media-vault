import axios from "axios";
import { defineStore, storeToRefs } from "pinia";
import { ref } from "vue";
import { useAuthStore } from "./auth";
import bcrypt from "bcryptjs";
import { error } from "console";

interface UserInfo {
  id: string;
  prefered_timezone: string;
  storage_usage: number;
  encryption_key_checksum: string;
}

export const useUserInfoStore = defineStore("userInfo", () => {
  const { accessToken } = storeToRefs(useAuthStore());
  const userStoreLoaded = ref(false);

  const preferedTimezone = ref("");
  const storageUsage = ref(0);
  const initRequired = ref(false);
  const encryptionKeyChecksum = ref("");

  const encryptionKey = ref("");
  // of 32 length
  const usableEncryptionKey = ref("");
  const encryptionKeyValidated = ref(false);

  function loadUserInfoIfRequred(): Promise<boolean> {
    if (userStoreLoaded.value === false) {
      return loadUserInfo();
    }
    return new Promise<boolean>((resolve) => resolve(true));
  }

  function loadUserInfo(): Promise<boolean> {
    return new Promise<boolean>((resolve, reject) => {
      axios
        .get<UserInfo>("/v1/user-info", {
          headers: {
            Authorization: `Bearer ${accessToken.value}`,
          },
        })
        .then((res) => {
          if (res.status == 200) {
            preferedTimezone.value = res.data.prefered_timezone;
            storageUsage.value = res.data.storage_usage;
            encryptionKeyChecksum.value = res.data.encryption_key_checksum;
            userStoreLoaded.value = true;
            resolve(true);
            return;
          }
          reject(new Error("non 200 and 404 response from server"));
          return;
        })
        .catch((err) => {
          if (axios.isAxiosError(err)) {
            if (err.response?.status == 404) {
              console.log("init required");
              initRequired.value = true;
              storageUsage.value = 0;
              preferedTimezone.value = "";
              resolve(false);
              return;
            }
          }
          reject(err);
          return;
        });
    });
  }

  function postUserInfo(
    _preferedTimezone: string,
    encryptionKey: string
  ): Promise<boolean> {
    return new Promise<boolean>((resolve, reject) => {
      const hash = bcrypt.hashSync(encryptionKey, 8);
      axios
        .post<UserInfo>(
          "/v1/user-info",
          {
            prefered_timezone: _preferedTimezone,
            encryption_key_checksum: hash,
          },
          {
            headers: {
              Authorization: `Bearer ${accessToken.value}`,
            },
          }
        )
        .then((res) => {
          if (res.status != 200) {
            reject(new Error("non 200 status"));
            return;
          }
          preferedTimezone.value = res.data.prefered_timezone;
          storageUsage.value = res.data.storage_usage;
          initRequired.value = false;
          resolve(true);
          return;
        })
        .catch((err) => {
          if (axios.isAxiosError(err)) {
            if (err.response?.status == 409) {
              reject(new Error("info of user already exist"));
              return;
            }
          }
          reject(err);
          return;
        });
    });
  }

  function validateEncryptionKey(_encryptionKey: string): boolean {
    if (!bcrypt.compareSync(_encryptionKey, encryptionKeyChecksum.value)) {
      encryptionKey.value = "";
      encryptionKeyValidated.value = false;
      return false;
    }
    encryptionKey.value = _encryptionKey;
    encryptionKeyValidated.value = true;
    setUsableEncryptionKey(_encryptionKey);
    return true;
  }

  function setUsableEncryptionKey(encryptionKey: string) {
    if (encryptionKey.length === 0) {
      return;
    }
    while (encryptionKey.length < 32) {
      encryptionKey += encryptionKey;
    }
    usableEncryptionKey.value = encryptionKey.substring(0, 32);
  }
  return {
    preferedTimezone,
    initRequired,
    usableEncryptionKey,
    encryptionKeyValidated,
    loadUserInfo,
    loadUserInfoIfRequred,
    postUserInfo,
    validateEncryptionKey,
  };
});

// todo
// what is more secure
// 1. using a secure hashing algo, and store the hash on server. after user auth is done client can request the hash from server to verify the encryption key entered by the user
// 2. use a less secure hashing algo like md5, store the hash on server. client can not get the hash, client will need to send the hash encryption key entered by the user to the server in order to verify it
