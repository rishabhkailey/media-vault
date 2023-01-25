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

* https://vueschool.io/courses/rapid-testing-with-vitest
* https://vueschool.io/courses/rapid-development-with-vite
* https://vueschool.io/courses/pinia-the-enjoyable-vue-store
* https://vueschool.io/courses/typescript-with-vue-js-3
* https://vueschool.io/courses/vue-router-4-for-everyone
* https://vueschool.io/courses/typescript-fundamentals

### go through
* https://vuejs.org/api/options-state.html
* 


### store (global state)
mutation-types.js
```js
export const SOME_MUTATION = 'SOME_MUTATION'
```

```js
import { SOME_MUTATION } from './mutation-types'

store = createStore({
    state: {
        count: 10,
        left: 80
    }, 
    getters: {
        total (state) {
            return state.count + state.left
        },
        // using getters inside getters
        totalMinusTen(state, getters) {
            return getters.total - 10
        },
        // method style 
        totalMinusX: (state, getters) => (x) => {
            return getters.total - x
        }
    },
    mutations: {
        incrementCount (state) {
            state.count++
        },
        incrementCountBy (state, payload) {
            state.count += payload.amount
        },
        // mutation names from constant
        [SOME_MUTATION] (state) {
        }
    },
    // mutations are synchronous for ansynchronous we should use actions
    // use for api calls 
    actions: {
        increment (context) {
          context.commit('incrementCount')
        },
        lilElegent({commit}) {
            commit('incrementCount')
        },
        incrementAsync ({ commit }) {
            setTimeout(() => {
                commit('increment')
            }, 1000)
        },
        incrementBy (context, payload) {
          context.commit('incrementCount', payload)
        },
        incrementByAsync ({ commit }, payload) {
            setTimeout(() => {
                commit('increment', payload)
            }, 1000)
        },
        // action can(should?) return promise to tell when it is complete
        actionA ({ commit }) {
            return new Promise((resolve, reject) => {
                setTimeout(() => {
                    commit('someMutation')
                    resolve()
                }, 1000)
            })
        },
        // action in action + usage of promise returned
        actionB ({ dispatch, commit }) {
            return dispatch('actionA').then(() => {
                commit('someOtherMutation')
            })
        },
        // usage of async/await (instead of returning promise we can also do it like this)
        async actionA ({ commit }) {
            commit('gotData', await getData())
        },

    }

})

// inside component

// usage of state
// using state inside component
computed: {
    count() {
        return this.$.store.state.count
    }
}
// or using map state
computed: mapState({
    // arrow functions can make the code very succinct!
    count: state => state.count,
    // to access local state with `this`, a normal function must be used
    countPlusLocalState (state) {
        return state.count + this.localCount
    }
})

computed: mapState([
  'count',
  'left',
])

computed: {
  localComputed () { },
  ...mapState({})
}

// usage of getters
this.$store.getters.total
// with local computed
computed: {
    ...mapGetters([
        'total',
        'totalMinusTen'
    ])
}
// or with different name
computed: {
    ...mapGetters({
        newTotal: 'total'
    })
}

// usage of mutations
store.commit('incrementCount')
store.commit('incrementCount', {amount: 10})
store.commit({
  type: 'incrementCount',
  amount: 10
})

methods: {
    ...mapMutations([
        'incrementCount', // map `this.increment()` to `this.$store.commit('increment')`
        'incrementCountBy' // map `this.incrementBy(amount)` to `this.$store.commit('incrementBy', amount)`
    ]),
    ...mapMutations({
        add: 'incrementCount' // map `this.add()` to `this.$store.commit('increment')`
    })
}

// action usage
store.dispatch('increment')
store.dispatch('incrementBy', {
    amount: 10
})
store.dispatch({
    type: 'incrementByAsync',
    amount: 10
})

methods: {
    ...mapActions([
        'increment', // map `this.increment()` to `this.$store.dispatch('increment')`
        // `mapActions` also supports payloads:
        'incrementBy' // map `this.incrementBy(amount)` to `this.$store.dispatch('incrementBy', amount)`
    ]),
    ...mapActions({
        add: 'increment' // map `this.add()` to `this.$store.dispatch('increment')`
    })
}

// action returning promise
store.dispatch('actionA').then(() => {

})

// modules - https://vuex.vuejs.org/guide/modules.html
```


### Provide / inject
Provide used by parent to define properties
Inject used by any children or children of children to get the value defined by parent

todo - https://vuex.vuejs.org/guide/modules.html

https://next.vuetifyjs.com/en/components/navigation-drawers/
https://next.vuetifyjs.com/en/components/app-bars/