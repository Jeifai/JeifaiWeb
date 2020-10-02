import Jobs from './components/Jobs.js';
import Company from './components/Company.js';

const router = new VueRouter({
    mode: 'hash',
    routes: [
    {
        path: '/',
        component: Jobs
    },
    {
        path: '/company',
        component: Company
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
                color: #2E353B;
            }`
        document.head.appendChild(styleElem);
    }
})

app.$mount('#app');