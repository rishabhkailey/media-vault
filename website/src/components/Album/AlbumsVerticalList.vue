<script setup lang="ts">
import { useAlbumStore } from "@/piniaStore/album";
import { albumRoute, albumsRoute } from "@/router/routesConstants";
import { storeToRefs } from "pinia";
import { computed } from "vue";
import { useRouter } from "vue-router";

const router = useRouter();
const albumStore = useAlbumStore();
const { albums, allAlbumsLoaded } = storeToRefs(albumStore);
const { loadMoreAlbums } = albumStore;

const maxNumberOfAlbums = 4;
const albumsSubSlice = computed<Array<Album>>(() => {
  return albums.value.length > maxNumberOfAlbums
    ? albums.value.slice(0, 4)
    : albums.value;
});
</script>
<template>
  <v-list-group value="Albums">
    <template v-slot:activator="{ props, isOpen }">
      <v-list-item
        title="Albums"
        prepend-icon="mdi-image-album"
        color=""
        @click="
          () => {
            router.push(albumsRoute());
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
    <v-infinite-scroll
      style="width: 100%; height: 100%"
      :items="albums"
      @load="
        ({ done }) => {
          loadMoreAlbums()
            .then((status) => {
              done(status);
            })
            .catch((_) => {
              done('error');
            });
        }
      "
    >
      <template #error> failed to load data from server </template>
      <template #empty> No more albums </template>
      <template #default>
        <v-list-item
          v-for="album in albumsSubSlice"
          :key="album.id"
          :title="album.name"
          prepend-icon="mdi-image"
          :value="album.name"
          :to="albumRoute(album.id)"
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
          :to="albumsRoute()"
          color="primary"
        ></v-list-item>
        <div
          v-if="!(allAlbumsLoaded || albums.length >= 5)"
          class="d-flex flex-row justify-center"
        ></div>
      </template>
    </v-infinite-scroll>
  </v-list-group>
</template>
