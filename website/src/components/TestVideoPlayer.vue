<script setup lang="ts"></script>

<template>
  <video-player
    src="/v1/testGetVideoWithRange/test.mp4"
    controls
    :loop="true"
    :volume="0.6"
  />
  <!-- <video-player src="https://vjs.zencdn.net/v/oceans.mp4" /> -->
</template>
