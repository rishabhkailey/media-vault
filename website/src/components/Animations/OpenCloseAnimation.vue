<script setup lang="ts">
import { ref, onMounted } from "vue";
import type { StyleValue } from "vue";

const props = withDefaults(
  defineProps<{
    // element around which the animation will start
    sourceElementSelector: string;
    // element around which the animation will end
    targetElementSelector: string;
    // on changing id animation will be played again
    _key: string;
    style: StyleValue;
    class: string;
    timeMs?: number;
  }>(),
  {
    timeMs: 150,
    class: "",
  },
);

const containerRef = ref<HTMLElement | undefined>(undefined);
const animationClass = ref<"" | "close-animation" | "open-animation">("");
type Matrix = {
  x: number;
  y: number;
  scaleX: number;
  scaleY: number;
  skewX: number;
  skewY: number;
};

function getTrasformMatrixCssValue(matrix: Matrix): string {
  return `matrix(
      ${matrix.scaleX},
      0,
      0,
      ${matrix.scaleY},
      ${matrix.x},
      ${matrix.y}
    )`;
}

function calculateStartAndEndMatrix(): {
  startMatrix: Matrix | undefined;
  endMatrix: Matrix | undefined;
} {
  let startMatrix: Matrix | undefined = undefined;
  let endMatrix: Matrix | undefined = undefined;
  let sourceElementRect = document
    .querySelector(props.sourceElementSelector)
    ?.getBoundingClientRect();
  let targetElementRect = document
    .querySelector(props.targetElementSelector)
    ?.getBoundingClientRect();
  if (sourceElementRect === undefined || targetElementRect === undefined) {
    return {
      startMatrix,
      endMatrix,
    };
  }
  startMatrix = {
    x:
      sourceElementRect.x -
      (targetElementRect.width - sourceElementRect.width) / 2,
    y:
      sourceElementRect.y -
      (targetElementRect.height - sourceElementRect.height) / 2,
    scaleX: sourceElementRect.width / targetElementRect.width,
    scaleY: sourceElementRect.height / targetElementRect.height,
    skewX: 0,
    skewY: 0,
  };
  endMatrix = {
    x: targetElementRect.x,
    y: targetElementRect.y,
    scaleX: 1,
    scaleY: 1,
    skewX: 0,
    skewY: 0,
  };
  return {
    startMatrix,
    endMatrix,
  };
}

function animate(
  element: HTMLElement,
  startMatrix: Matrix,
  endMatrix: Matrix,
): Animation {
  return element.animate(
    [
      {
        transform: getTrasformMatrixCssValue(startMatrix),
      },
      {
        transform: getTrasformMatrixCssValue(endMatrix),
      },
    ],
    {
      duration: props.timeMs,
      fill: "both",
      easing: "ease-in",
    },
  );
}

let animateStartMatrix: Matrix | undefined = undefined;
let animateEndMatrix: Matrix | undefined = undefined;

async function openAnimation() {
  let targetElement = document.querySelector(props.targetElementSelector);
  let { startMatrix, endMatrix } = calculateStartAndEndMatrix();
  animateStartMatrix = startMatrix;
  animateEndMatrix = endMatrix;
  if (
    animateStartMatrix === undefined ||
    animateEndMatrix === undefined ||
    targetElement == undefined
  ) {
    return;
  }
  await animate(
    targetElement as HTMLElement,
    animateStartMatrix,
    animateEndMatrix,
  ).finish;
}

// todo check vuetify activators
async function closeAnimation() {
  let targetElement = document.querySelector(props.targetElementSelector);
  if (
    animateStartMatrix === undefined ||
    animateEndMatrix === undefined ||
    targetElement == undefined
  ) {
    return;
  }
  await animate(
    targetElement as HTMLElement,
    animateEndMatrix,
    animateStartMatrix,
  ).finished;
}
onMounted(() => {
  openAnimation();
});
</script>
<template>
  <div
    :class="`${props.class} ${animationClass}`"
    :style="style"
    ref="containerRef"
  >
    <slot name="default" :closeAnimation="closeAnimation"> </slot>
  </div>
</template>
