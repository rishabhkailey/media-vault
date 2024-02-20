import { useAuthStore } from "@/piniaStore/auth";
import axios from "axios";

// todo
function validateMedia(media: Media): Media {
  return media;
}

export async function getSingleMediaById(id: number): Promise<Media> {
  const { accessToken } = useAuthStore();
  const response = await axios.get<Media>(`/v1/media/${id}`, {
    headers: {
      Authorization: `Bearer ${accessToken}`,
    },
  });
  if (response.status !== 200) {
    throw new Error(`${response.status} status code`);
  }
  return validateMedia(response.data);
}
