import { useAuthStore } from "@/piniaStore/auth";
import { revokeSession } from "@/js/api/user";
import { useErrorsStore } from "@/piniaStore/errors";
import { logOutUsingUserManager } from "@/js/auth";
import { aboutRoute } from "@/router/routesConstants";
// import { useRouter } from "vue-router";
import { router } from "@/router/index";

// should only be called from inside the vue component where pinia store is expected to be initialized
export async function logOut(): Promise<any> {
  // const router = useRouter();
  const authStore = useAuthStore();
  const { appendError } = useErrorsStore();

  try {
    await logOutUsingUserManager();
  } catch (err) {
    appendError(
      "logout failed",
      `auth server retuned error response - ${err}`,
      -1,
    );
    return false;
  }

  try {
    await revokeSession;
  } catch (err) {
    appendError(
      "logout failed",
      `resource server retuned error response - ${err}`,
      -1,
    );
  }

  authStore.reset();
  try {
    const err = await router.push(aboutRoute());
    if (err instanceof Error) {
      throw err;
    }
  } catch (err) {
    appendError("redirect to about page failed", `error - ${err}`, -1);
  }
  return true;
}
