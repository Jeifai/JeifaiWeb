import Home from './components/Home.js';
import Targets from './components/Targets.js';
import Keywords from './components/Keywords.js';
import Matches from './components/Matches.js';

const router = new VueRouter({
    mode: 'hash',
    routes: [
    {
        path: '/homeT',
        component: Home
    },
    {
        path: '/targetsT',
        component: Targets
    },
    {
        path: '/keywordsT',
        component: Keywords
    },
    {
        path: '/matchesT',
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
    },
    mounted() {
        const list = mdc.list.MDCList.attachTo(document.querySelector('.mdc-list'));
        list.singleSelection = true;
        list.selectedIndex = this.selectedIndex;
    }
})

app.$mount('#app');