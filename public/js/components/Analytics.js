export default {
    name: 'analytics',
    delimiters: ["[[","]]"],
    data: function () {
        return {
            selectedIndex: 4,
            jobs: [],
            jobsTotal: [],
            jobsTotalMinY: '',
            companyInfo: '',
            employeesTrend: '',
            employeesTotalMinY: '',
            jobTitlesWords: '',
            jobsTitlesMaxCount: '',
            targets: [],
            selectedTarget: null,
            chartOptions: {
                elements: {
                    point:{
                        radius: 0.8
                    }
                },
                scales: {
                    yAxes: [{
                        gridLines: {
                            display: false
                        },
                        ticks: {
                            callback: function(value, index, values) {
                                if (index === values.length - 1) return Math.min.apply(this, values);
                                else if (index === 0) return Math.max.apply(this,  values);
                                else return '';
                            }
                        }
                    }]
                },
                legend: {
                    display: true,
                }
            }
        }
    },
    mounted() {
        this.$parent.selectedIndex = this.selectedIndex;
        let multiToggleScript = document.createElement('script')
        multiToggleScript.setAttribute('src', 'https://unpkg.com/vue-single-select@latest')
        document.head.appendChild(multiToggleScript)
        let styleElem = document.createElement('style');
        styleElem.textContent = `
            .topSide {
                display: flex;
                align-items: center;
            }
            .singleselectfield {
                max-width: 15%;
                padding-right: 2%;
            }
            .topSideHeadline {
                font-size: 10px;
                vertical-align: sub;
            }
            .logoStyle {
                float: left;
                height: 50px;
                width: auto;
            }
            .topSideContent {
                float: left;
                padding-left: 70px;

            }
            .topSideContentValue {
                font-size:30px;
                vertical-align: super;
            }`
        document.head.appendChild(styleElem);
    },
    created () {
        this.fetchTargets()
    },
    watch: {
        selectedTarget: function(val) {
            if (this.selectedTarget !== null) {
                this.fetchData();
            }
        }
    },
    methods: {
        fetchTargets: function() {
            this.$http.get('/analytics/getTargets').then(function(response) {
                this.targets = response.data.Targets
            }).catch(function(error) {
                console.log(error)
            });
        },
        fetchData: function() {
            this.$http.get('/analytics/target/' + this.selectedTarget).then(function(response) {
                this.companyInfo = response.data.CompanyInfo;
                this.jobsTotal = [{name: 'TotalJobs', data: response.data.Jobs.CountTotal}];
                this.jobs = [
                    {name: 'jobsClosed', data: response.data.Jobs.CountClosed, color: "#ffadad"},
                    {name: 'jobsCreated', data: response.data.Jobs.CountCreated, color: "#caffbf"}
                ];
                this.jobsTotalMinY = response.data.Jobs.CountTotalMinY;
                this.employeesTrend = response.data.EmployeesTrend.CountEmployees;
                this.employeesTrend = [{name: 'EmployeesTrend', data: response.data.EmployeesTrend.CountEmployees}];
                this.employeesTotalMinY = response.data.EmployeesTrend.CountEmployeesMinY;
                this.jobTitlesWords = [{name: 'JobTitlesWords', data: response.data.JobTitlesWords.Words}];
                this.jobsTitlesMaxCount = response.data.JobTitlesWords.MaxCount;

            }).catch(function(error) {
                console.log(error)
            });
        },
    },
    template: `
        <div>
            <div class="topSide">
                <div class="singleselectfield">
                    <vue-single-select
                        v-model="selectedTarget"
                        :options="targets"
                        placeholder="Select a target"
                        :max-results="100000">
                    </vue-single-select>
                </div>
                <div v-if="selectedTarget">
                    <img class="logoStyle" v-bind:src="'/static/images/targets/' + selectedTarget + '.png'" v-bind:alt="selectedTarget">
                    <div class="topSideContent" div v-for="(value, name) in companyInfo">
                        <div class="topSideHeadline">[[name]]</div>
                        <div class="topSideContentValue">[[value]]</div>
                    </div>
                </div>
            </div>
            <div v-if="selectedTarget">
                <area-chart :data="jobsTotal" :min="jobsTotalMinY" width="95%" height="8%" :library="chartOptions" :colors="['#a0c4ff']"></area-chart><br>
                <area-chart :data="jobs" width="95%" height="8%" :library="chartOptions"></area-chart>
                <area-chart :data="employeesTrend" :min="employeesTotalMinY" width="95%" height="8%" :library="chartOptions" :colors="['#ffc6ff']"></area-chart>
                <column-chart :data="jobTitlesWords" :max="jobsTitlesMaxCount" width="95%" height="8%" :library="chartOptions" :colors="['#ffd6a5']"></column-chart>
            </div>
        </div>`,
};