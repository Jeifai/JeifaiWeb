export default {
    name: 'analytics',
    delimiters: ["[[","]]"],
    data: function () {
        return {
            selectedIndex: 4,
            jobs: [],
        }
    },
    mounted() {
        this.$parent.selectedIndex = this.selectedIndex;
    },
    created () {
        this.fetchJobs()
    },
    methods: {
        fetchJobs: function() {
            this.$http.get('/analytics').then(function(response) {
                this.jobs = {};
                for (var key in response.data.Jobs) {
                    this.jobs[response.data.Jobs[key].Date] = response.data.Jobs[key].Count;
                }
            }).catch(function(error) {
                console.log(error)
            });
        }
    },
    template: `
        <div>
            <line-chart :data="jobs"></line-chart>
        </div>`,
};