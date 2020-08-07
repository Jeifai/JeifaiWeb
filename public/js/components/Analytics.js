export default {
    name: 'analytics',
    delimiters: ["[[","]]"],
    data: function () {
        return {
            selectedIndex: 4,
            jobs: [],
            jobsTotal: [],
            hasJobs: false,
            targets: [],
            selectedTarget: null,
        }
    },
    mounted() {
        this.$parent.selectedIndex = this.selectedIndex;
        let multiToggleScript = document.createElement('script')
        multiToggleScript.setAttribute('src', 'https://unpkg.com/vue-single-select@latest')
        document.head.appendChild(multiToggleScript)
        let styleElem = document.createElement('style');
        styleElem.textContent = `
            .singleselectfield {
                max-width: 50%;
            }`
        document.head.appendChild(styleElem);
    },
    created () {
        this.fetchTargets()
    },
    methods: {
        fetchTargets: function() {
            this.$http.get('/analytics/getTargets').then(function(response) {
                this.targets = response.data.Targets
            }).catch(function(error) {
                console.log(error)
            });
        },
        fetchJobs: function() {
            this.$http.get('/analytics/getJobsPerDayPerTarget/' + this.selectedTarget).then(function(response) {
                var jobsClosed = {};
                var jobsCreated = {};
                this.jobsTotal = {};
                for (var key in response.data.Jobs) {
                    jobsClosed[response.data.Jobs[key].Date] = response.data.Jobs[key].CountClosed;
                    jobsCreated[response.data.Jobs[key].Date] = response.data.Jobs[key].CountCreated;
                    this.jobsTotal[response.data.Jobs[key].Date] = response.data.Jobs[key].CountTotal;
                }
                if (Object.keys(this.jobsTotal).length > 0) {
                    this.hasJobs = true;
                }
                this.jobs = [
                    {name: 'jobsClosed', data: jobsClosed, color: "#585858"},
                    {name: 'jobsCreated', data: jobsCreated, color: "#00CC66"}
                ]
            }).catch(function(error) {
                console.log(error)
            });
        },
    },
    template: `
        <div>
            <div class="taggableselectfield">
                <span><br>Select a target.</span>
                <div class="singleselectfield">
                    <vue-single-select
                        v-model="selectedTarget"
                        :options="targets"
                        placeholder="Targets"
                        :max-results="100000">
                    </vue-single-select>
                </div>
                <button class="mdc-button mdc-button--raised" v-on:click="fetchJobs">
                    <div class="mdc-button__ripple"></div>
                    <i class="material-icons mdc-button__icon" aria-hidden="true">check</i>
                    <span class="mdc-button__label">Explore target</span>
                </button>
            </div><br>
            <div v-if="hasJobs">
                <area-chart :data="jobsTotal" width="95%" height="250px"></area-chart><br>
                <column-chart :data="jobs" width="95%" height="250px"></column-chart>
            </div>
        </div>`,
};