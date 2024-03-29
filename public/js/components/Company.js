export default {
    name: 'company',
    delimiters: ["[[","]]"],
    data: function () {
        return {
            selectedIndex: 4,
            metabase: '',
        }
    },
    mounted() {
        this.$parent.selectedIndex = this.selectedIndex;
        let multiToggleScript = document.createElement('script')
        let styleElem = document.createElement('style');
        styleElem.textContent = ``
        document.head.appendChild(styleElem);
        gtag('config', 'UA-180812973-1', {
          'page_title' : 'CompanyExplorer',
          'page_path': '/CompanyExplorer'
        });
    },
    created () {
        this.serveMetabaseCompany()
    },
    methods: {
        serveMetabaseCompany: function() {
            this.$http.get('/serveMetabaseCompany').then(function(response) {
                this.metabase = response.data.Metabase
            }).catch(function(error) {
                console.log(error)
            });
        },
    },
    template: `
        <div>
            <iframe
                v-bind:src="metabase"
                frameborder="0"
                width="100%"
                height="680px"
                allowtransparency
            ></iframe>
        </div>`,
};