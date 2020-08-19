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
            minYEmployees: '',
            maxYEmployees: '',
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
                        }
                    }]
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
            }
            .middleCharts {
                display: flex;
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
                this.jobsTotal = response.data.Jobs.CountTotal;
                this.jobs = [
                    {name: 'jobsClosed', data: response.data.Jobs.CountClosed, color: "#585858"},
                    {name: 'jobsCreated', data: response.data.Jobs.CountCreated, color: "#00CC66"}
                ];
                this.jobsTotalMinY = response.data.Jobs.CountTotalMinY;




                /** 
                this.employeesTrend = {};
                for (var key in response.data.EmployeesTrend) {
                    jobsClosed[response.data.Jobs[key].Date] = response.data.Jobs[key].CountClosed;
                    this.employeesTrend[response.data.EmployeesTrend[key].Date] = response.data.EmployeesTrend[key].CountEmployees;
                }

                var tempEmployeesKeys = this.employeesTrend;
                var employeesKeys = Object.keys(tempEmployeesKeys).sort(function(a,b) { return tempEmployeesKeys[a] - tempEmployeesKeys[b];});
                this.minYEmployees = Math.floor(tempEmployeesKeys[employeesKeys[0]] * 0.96);
                this.maxYEmployees = Math.floor(tempEmployeesKeys[employeesKeys[employeesKeys.length - 1]] * 1.04);
                 */

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
                <area-chart :data="jobsTotal" :min="jobsTotalMinY"  width="95%" height="10%" :library="chartOptions"></area-chart><br>
                <div class="middleCharts">
                    <column-chart :data="jobs" height="250px"></column-chart>
                    <line-chart :min="minYEmployees" :max="maxYEmployees" :data="employeesTrend" height="250px"></line-chart>
                </div">
            </div>
        </div>`,
};