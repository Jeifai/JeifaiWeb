const list = mdc.list.MDCList.attachTo(document.querySelector('.mdc-list'));
list.singleSelection = true;

import Home from './components/Home.js';
import Targets from './components/Targets.js';
import Matches from './components/Matches.js';

const router = new VueRouter({
  mode: 'history',
  routes: [
    {
        path: '/home',
        component: Home
    },
    {
        path: '/targets',
        component: Targets
    },
    {
        path: '/matches',
        component: Matches
    },
  ]
})

var app = new Vue({
    router,
    delimiters: ["[[","]]"],
    data() {
        return {
            selectedIndex: this.$router.currentRoute.selectedIndex
        }
    },
    watch: {
        selectedIndex: function(val) {
            list.selectedIndex = 2;
            alert(val);
        }
    }
})

app.$mount('#app');

console.log(app.data)
console.log(app.$data)
console.log(app.$data.selectedIndex)

// const list = mdc.list.MDCList.attachTo(document.querySelector('.mdc-list'));
// list.singleSelection = true;
// list.selectedIndex = app.$data.selectedIndex;
list.selectedIndex = 2;