export default {
    name: 'test',
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
    },
    created () {
        this.getMetabaseUrl()
    },
    methods: {
        getMetabaseUrl: function() {
            this.$http.get('/getMetabaseUrl').then(function(response) {
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
                height="680"
                allowtransparency
            ></iframe>
        </div>`,
};