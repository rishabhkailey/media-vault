<script setup lang="ts">
import axios from "axios";
import { onMounted, ref } from "vue";

const props = withDefaults(
  defineProps<{
    page: number;
    perPage: number;
    accessToken: string;
  }>(),
  {
    page: 0,
    items: 30,
  }
);

interface media {
  name: string;
  type: string;
  date: Date;
  size: number;
  thumbnail: boolean;
  url: string;
  thumbnail_url: string;
}
const mediaList = ref<Array<media>>([]);

onMounted(() => {
  axios
    .get<Array<media>>(
      `/v1/mediaList?page=${props.page}&perPage=${props.perPage}`,
      {
        headers: {
          Authorization: `Bearer ${props.accessToken}`,
        },
      }
    )
    .then((response) => {
      console.log(response);
      if (response.status == 200) {
        mediaList.value = response.data;
      }
    })
    .catch((err) => {
      console.log(err);
    });
});
</script>

<template>
  <div :key="`${index}+${media.name}`" v-for="(media, index) in mediaList">
    <img v-if="media.thumbnail" :src="media.thumbnail_url" />
  </div>
</template>
