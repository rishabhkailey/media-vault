import { signinUsingUserManager } from "@/js/auth";
import { useAuthStore } from "@/piniaStore/auth";
import { useUserInfoStore } from "@/piniaStore/userInfo";
import { storeToRefs } from "pinia";
import { userManager } from "@/js/auth";
import type {
  LocationQueryRaw,
  NavigationGuard,
  RouteLocationNamedRaw,
} from "vue-router";
import { updateOrRegisterServiceWorker } from "@/js/serviceWorker/registeration";
import { promiseTimeout } from "@/js/utils";

const aboutNavigationPath: RouteLocationNamedRaw = {
  name: "about",
};

function errorScreenRoute(
  title: string,
  error: any,
  returnUri?: string,
): RouteLocationNamedRaw {
  const query: LocationQueryRaw = {
    title: title,
  };

  if (returnUri !== undefined) {
    query.return_uri = returnUri;
  }

  if (typeof error === "string") {
    query.message = error;
  }
  if (error instanceof Error) {
    query.message = error.message + "\n" + error.stack;
  }

  return {
    name: "errorscreen",
    query: query,
  };
}

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
    // ignore error and user will remain unauthenticated
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
    return errorScreenRoute("Failed to load user info", err, to.fullPath);
  }
  if (userInfoInitRequired.value) {
    return {
      name: "initialSetup",
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
    return errorScreenRoute("Failed to load user info", err, to.fullPath);
  }
  console.debug(
    `encryption key gaurd, ${to.fullPath}`,
    `encryptionKeyValidated = ${encryptionKeyValidated.value}`,
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
    console.group("service worker registeration");
    await promiseTimeout(updateOrRegisterServiceWorker(), 5 * 1000);
    console.groupEnd();
  } catch (err) {
    let errorMessage = "service worker registeration failed.";
    if (err instanceof Error) {
      errorMessage += " " + err.message;
    }
    console.groupEnd();
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
