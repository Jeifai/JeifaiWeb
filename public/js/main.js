import Home from './components/Home.js';
import Targets from './components/Targets.js';
import Keywords from './components/Keywords.js';
import Matches from './components/Matches.js';
import Profile from './components/Profile.js';

const router = new VueRouter({
    mode: 'hash',
    routes: [
    {
        path: '/',
        component: Home
    },
    {
        path: '/targets',
        component: Targets
    },
    {
        path: '/keywords',
        component: Keywords
    },
    {
        path: '/matches',
        component: Matches
    },
    {
        path: '/profile',
        component: Profile
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
    methods: {
        csvExport(arrData, fileName) {
            let csvContent = "data:text/csv;charset=utf-8,";
            csvContent += [
                Object.keys(arrData[0]).join(";"),
                ...arrData.map(item => Object.values(item).join(";"))
            ].join("\n").replace(/(^\[)|(\]$)/gm, "");
            const data = encodeURI(csvContent);
            const link = document.createElement("a");
            link.setAttribute("href", data);
            link.setAttribute("download", fileName);
            link.click();
        }
    },
    mounted() {
        const list = mdc.list.MDCList.attachTo(document.querySelector('.mdc-list'));
        list.singleSelection = true;
        list.selectedIndex = this.selectedIndex;
    }
})

app.$mount('#app');