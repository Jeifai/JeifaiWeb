import Manage from './components/Manage.js';
import Watch from './components/Watch.js';
import Targets from './components/Targets.js';
import Analytics from './components/Analytics.js';
import Profile from './components/Profile.js';

const router = new VueRouter({
    mode: 'hash',
    routes: [
    {
        path: '/',
        component: Manage
    },
    {
        path: '/watch',
        component: Watch
    },
    {
        path: '/targets',
        component: Targets
    },
    {
        path: '/analytics',
        component: Analytics
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
            if (typeof(list) !== 'undefined') {
                list.selectedIndex = val;
            }
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

        let styleElem = document.createElement('style');
        styleElem.textContent = `
            .column-sort {
                font-size: 16px;
                vertical-align: -3px;
                color: #BEBEBE;
            }`
        document.head.appendChild(styleElem);
    }
})

app.$mount('#app');