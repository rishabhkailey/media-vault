<script setup lang="ts">
import InfiniteScroller from "@/components/InfiniteScroller/InfiniteScroller.vue";
import { useAlbumStore } from "@/piniaStore/album";
import { useAuthStore } from "@/piniaStore/auth";
import { storeToRefs } from "pinia";
import { computed } from "vue";
import { useRouter } from "vue-router";
import { useDisplay } from "vuetify";
const display = useDisplay();

const router = useRouter();
const props = defineProps<{
  modelValue: boolean;
}>();
const emit = defineEmits<{
  (e: "update:modelValue", value: boolean): void;
}>();
const smallDisplay = computed(
  () => display.mobile.value || display.smAndDown.value
);

const albumStore = useAlbumStore();
const { albums, allAlbumsLoaded } = storeToRefs(albumStore);
const { loadMoreAlbums } = albumStore;

const maxNumberOfAlbums = 4;
const albumsSubSlice = computed<Array<Album>>(() => {
  return albums.value.length > maxNumberOfAlbums
    ? albums.value.slice(0, 4)
    : albums.value;
});

const authStore = useAuthStore();
const { accessToken } = storeToRefs(authStore);
</script>

<template>
  <v-navigation-drawer
    :model-value="props.modelValue"
    @update:model-value="
      (newValue) => {
        emit('update:modelValue', newValue);
      }
    "
    :permanent="!smallDisplay"
    :temporary="smallDisplay"
    :rounded="false"
    elevation="2"
  >
    <v-list nav>
      <v-list-item
        prepend-icon="mdi-home"
        title="Home"
        value="Home"
        :to="{
          name: 'Home',
        }"
        :exact="true"
        color="primary"
      ></v-list-item>
      <v-list-group value="Albums">
        <!-- todo move to separate component in navigationBar directory -->
        <template v-slot:activator="{ props, isOpen }">
          <v-list-item
            title="Albums"
            prepend-icon="mdi-image-album"
            color=""
            @click="
              () => {
                router.push({
                  name: 'Albums',
                });
              }
            "
            value="Albums"
          >
            <template #append>
              <v-icon
                :icon="isOpen ? 'mdi-menu-up' : 'mdi-menu-down'"
                v-bind="props"
                @click.stop.prevent="() => {}"
              />
            </template>
          </v-list-item>
        </template>
        <v-list-item
          v-for="album in albumsSubSlice"
          :key="album.id"
          :title="album.name"
          prepend-icon="mdi-image"
          :value="album.name"
          :to="{
            name: 'Album',
            params: {
              album_id: album.id,
            },
          }"
          color="primary"
        >
          <template #prepend>
            <div class="mr-3">
              <v-img
                style="width: 1.5em; height: 1.5em"
                cover
                :src="album.thumbnail_url"
                v-if="album.thumbnail_url.length > 0"
              >
                <template #error>
                  <v-icon icon="mdi-image" />
                </template>
              </v-img>
              <v-icon v-else icon="mdi-image" />
            </div>
          </template>
        </v-list-item>
        <v-list-item
          v-if="albums.length > maxNumberOfAlbums"
          title="View all albums"
          prepend-icon="mdi-arrow-right-thin"
          :to="{
            name: 'Albums',
          }"
          color="primary"
        ></v-list-item>
        <div
          v-if="!(allAlbumsLoaded || albums.length >= 5)"
          class="d-flex flex-row justify-center"
        >
          <InfiniteScroller
            v-if="!allAlbumsLoaded"
            :on-threshold-reach="loadMoreAlbums"
            :threshold="0.1"
            :min-height="100"
            :min-width="100"
            :root-margin="10"
          >
            <v-progress-circular indeterminate></v-progress-circular>
          </InfiniteScroller>
        </div>
      </v-list-group>
    </v-list>
  </v-navigation-drawer>
</template>
