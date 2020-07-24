const list = mdc.list.MDCList.attachTo(document.querySelector('.mdc-list'));
list.singleSelection = true;
list.selectedIndex = 0;

import Home from './components/Home.js';
import Targets from './components/Targets.js';
import Keywords from './components/Keywords.js';
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
        path: '/keywords',
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
            list.selectedIndex = val;
        }
    }
})

app.$mount('#app');