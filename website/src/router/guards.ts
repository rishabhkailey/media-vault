import { signinUsingUserManager } from "@/js/auth";
import { useAuthStore } from "@/piniaStore/auth";
import { useUserInfoStore } from "@/piniaStore/userInfo";
import { storeToRefs } from "pinia";
import { userManager } from "@/js/auth";
import type { NavigationGuard, RouteLocationNamedRaw } from "vue-router";
import { updateOrRegisterServiceWorker } from "@/js/serviceWorker/registeration";

const aboutNavigationPath: RouteLocationNamedRaw = {
  name: "about",
};

// ensure user is logged in and user onboarding is also done
export const loginGaurd: NavigationGuard = async (to) => {
  console.debug(`login gaurd, ${to.fullPath}`);
  // we can not define stores globally in this file, as it will not work outside component setup
  // https://router.vuejs.org/guide/advanced/navigation-guards.html
  const userInfoStore = useUserInfoStore();
  const { initRequired: userInfoInitRequired } = storeToRefs(userInfoStore);
  const { loadUserInfoIfRequred } = userInfoStore;
  const authStore = useAuthStore();
  const { authenticated } = storeToRefs(authStore);

  if (userManager === undefined) {
    // todo on error redirect to error page with redirect uri?
    console.error("userManager undefined: ", userManager);
    return false;
  }
  // load user from local storage
  try {
    const user = await userManager.getUser();
    if (
      user !== null &&
      (user.expired === false || user.expired === undefined)
    ) {
      // usermanager automatically tries to renew the token before expiration
      // if already expired we are not setting user auth info
      authStore.setUserInfo(user);
    }
  } catch (err) {
    console.error(err);
  }

  if (!authenticated.value) {
    if (to.name === "Home") {
      return aboutNavigationPath;
    }
    signinUsingUserManager(userManager, false);
    return false;
  }

  try {
    await loadUserInfoIfRequred();
  } catch (err) {
    // todo: redirect and give option to clear cookies?
    console.error(err);
    return false;
  }
  console.log(userInfoInitRequired.value);
  if (userInfoInitRequired.value) {
    return {
      name: "onboarding",
    };
  }
  return true;
};

// ensure user's encryption key is validated and stored in store
export const encryptionKeyGaurd: NavigationGuard = async (to) => {
  const userInfoStore = useUserInfoStore();
  const { encryptionKeyValidated } = storeToRefs(userInfoStore);
  try {
    await userInfoStore.loadUserInfoIfRequred();
  } catch (err) {
    console.log("error: ", err);
    return false;
  }
  console.debug(
    `encryption key gaurd, ${to.fullPath}`,
    `encryptionKeyValidated = ${encryptionKeyValidated.value}`
  );
  if (encryptionKeyValidated.value == false) {
    return {
      name: "encryptionKey",
      query: {
        return_uri: to.fullPath,
      },
    };
  }
  return true;
};

export const serviceWrokerGaurd: NavigationGuard = async (to) => {
  try {
    await updateOrRegisterServiceWorker();
  } catch (err) {
    let errorMessage = "service worker registeration failed.";
    if (err instanceof Error) {
      errorMessage += " " + Error.toString();
    }
    return {
      name: "errorscreen",
      query: {
        title: "Service Worker registeration failed",
        message: errorMessage,
        return_uri: to.fullPath,
      },
    };
  }
  return true;
};
