```ts
@ = v-on
// :=v-bind??
:key="js" for evaluating js inside string
v-model for 2 way bindings (use this for changing data)
v-bind for making html element attribute dynamic
:=v-bind

// dynamic html element classes
:class="{
    btn: true,
    flex-1: true
}"

// static + dynamic class use both class and :class
<div class="static-class" :class="dynamic-class>

// we can also have class as array instead of objects
:class="['btn', 'flex-1']"

// computed properties
import {ref, computed} from "vue"
a = ref("")
charCount = computed(()=> a.value.length)

// ref vs reactive

```


vitest
```bash
```

https://vueschool.io/courses/rapid-testing-with-vitest
https://vueschool.io/courses/rapid-development-with-vite
https://vueschool.io/courses/pinia-the-enjoyable-vue-store
https://vueschool.io/courses/typescript-with-vue-js-3
https://vueschool.io/courses/vue-router-4-for-everyone
https://vueschool.io/courses/typescript-fundamentals