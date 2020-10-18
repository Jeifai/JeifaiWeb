export default {
    name: 'jobs',
    delimiters: ["[[","]]"],
    data: function () {
        return {
            selectedIndex: 0,
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
          'page_title' : 'JobsExplorer',
          'page_path': '/JobsExplorer'
        });
    },
    created () {
        this.serveMetabaseJobs()
    },
    methods: {
        serveMetabaseJobs: function() {
            this.$http.get('/serveMetabaseJobs').then(function(response) {
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