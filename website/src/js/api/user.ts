import axios from "axios";

// todo access token required?
export async function revokeSession(): Promise<boolean> {
  const response = await axios.post("/v1/terminateSession");
  if (response.status !== 200) {
    throw new Error("terminate session failed");
  }
  return true;
}

export async function refreshSessionWithResourceServer(
  accessToken: string,
): Promise<boolean> {
  const response = await axios.post(
    "/v1/refreshSession",
    {},
    {
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    },
  );
  if (response.status !== 200) {
    throw new Error("terminate session failed");
  }
  return true;
}
