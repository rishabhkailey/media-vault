<script setup lang="ts">
import LazyMediaThumbnailsPreview from "@/components/LazyMediaThumbnailsPreview.vue";
import { LOAD_MORE_MEDIA_ACTION } from "@/store/modules/media";
import { computed } from "vue";
import { useStore } from "vuex";
const store = useStore();
const mediaList = computed<Array<Media>>(() => store.getters.mediaList);
const allMediaLoaded = computed<boolean>(() => store.getters.allMediaLoaded);
const loadMoreMedia = () => store.dispatch(LOAD_MORE_MEDIA_ACTION);
</script>

<template>
  <LazyMediaThumbnailsPreview
    :media-list="mediaList"
    :all-media-loaded="allMediaLoaded"
    :load-more-media="loadMoreMedia"
  />
</template>
