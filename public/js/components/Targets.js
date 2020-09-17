export default {
    name: 'targets',
    delimiters: ["[[","]]"],
    data: function () {
        return {
            selectedIndex: 1,
            targets:  [],
            sortedBy: "CreatedDate",
            sorting: {
                CreatedDate: true,
                LastExtractionDate: false,
                Employees: false,
                Name: false,
                JobsAll: false,
                JobsNow: false,
                Opened: false,
                Closed: false
            }
        }
    },
    mounted() {
        this.$parent.selectedIndex = this.selectedIndex
        let styleElem = document.createElement('style');
        styleElem.textContent = `
            .input-row {
                display: flex;
            }`
    },
    created () {
        this.fetchTargets()
    },
    methods: {
        fetchTargets: function() {
            this.$http.get('/targets/analytic').then(function(response) {
                this.targets = response.data.Targets;
            }).catch(function(error) {
                console.log(error)
            });
        },
        sortRows: function(column) {
            this.sortedBy = column
            if (column == "CreatedDate" || column == "LastExtractionDate" ) {
                if (this.sorting[column]) {
                    this.targets.sort((a,b) => (new Date(a[column]) - new Date(b[column])))
                    this.sorting[column] = false
                } else {
                    this.targets.sort((b,a) => (new Date(a[column]) - new Date(b[column])))
                    this.sorting[column] = true
                }
            } else {
                if (this.sorting[column]) {
                    this.targets.sort((a,b) => (a[column] > b[column]) ? 1 : ((b[column] > a[column]) ? -1 : 0))
                    this.sorting[column] = false
                } else {
                    this.targets.sort((b,a) => (a[column] > b[column]) ? 1 : ((b[column] > a[column]) ? -1 : 0))
                    this.sorting[column] = true
                }
            }
        },
        csvExport: function() {
            this.$parent.csvExport(this.targets, "targets.csv")
        }
    },
    template: `
        <div>
            <div>
                <div class="initial-row">
                    <h1>
                        My targets
                        <button
                            v-on:click="csvExport"
                            class="material-icons mdc-top-app-bar__action-item mdc-icon-button" 
                            aria-label="Export">arrow_downward
                        </button>
                    </h1>
                </div>
                <div class="mdc-data-table">
                    <table class="mdc-data-table__table" aria-label="Created Targets">
                        <thead>
                            <tr class="mdc-data-table__header-row">
                                <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                    <a class="column-header" @click="sortRows('Name')">
                                        Target
                                        <i v-if="sortedBy === 'Name' && sorting['Name'] === true" class="material-icons column-sort">keyboard_arrow_up</i>
                                        <i v-if="sortedBy === 'Name' && sorting['Name'] === false" class="material-icons column-sort">keyboard_arrow_down</i>
                                    </a>
                                </th>
                                <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                    <a class="column-header" @click="sortRows('CreatedDate')">
                                        CreatedAt
                                        <i v-if="sortedBy === 'CreatedDate' && sorting['CreatedDate'] === true" class="material-icons column-sort">keyboard_arrow_up</i>
                                        <i v-if="sortedBy === 'CreatedDate' && sorting['CreatedDate'] === false" class="material-icons column-sort">keyboard_arrow_down</i>
                                    </a>
                                </th>
                                <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                    <a class="column-header" @click="sortRows('LastExtractionDate')">
                                        LastExtractionAt
                                        <i v-if="sortedBy === 'LastExtractionDate' && sorting['LastExtractionDate'] === true" class="material-icons column-sort">keyboard_arrow_up</i>
                                        <i v-if="sortedBy === 'LastExtractionDate' && sorting['LastExtractionDate'] === false" class="material-icons column-sort">keyboard_arrow_down</i>
                                    </a>
                                </th>
                                <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                    <a class="column-header" @click="sortRows('Employees')">
                                        Employees
                                        <i v-if="sortedBy === 'Employees' && sorting['Employees'] === true" class="material-icons column-sort">keyboard_arrow_up</i>
                                        <i v-if="sortedBy === 'Employees' && sorting['Employees'] === false" class="material-icons column-sort">keyboard_arrow_down</i>
                                    </a>
                                </th>
                                <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                    <a class="column-header" @click="sortRows('JobsAll')">
                                        All Jobs
                                        <i v-if="sortedBy === 'JobsAll' && sorting['JobsAll'] === true" class="material-icons column-sort">keyboard_arrow_up</i>
                                        <i v-if="sortedBy === 'JobsAll' && sorting['JobsAll'] === false" class="material-icons column-sort">keyboard_arrow_down</i>
                                    </a>
                                </th>
                                <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                    <a class="column-header" @click="sortRows('JobsNow')">
                                        Current Jobs
                                        <i v-if="sortedBy === 'JobsNow' && sorting['JobsNow'] === true" class="material-icons column-sort">keyboard_arrow_up</i>
                                        <i v-if="sortedBy === 'JobsNow' && sorting['JobsNow'] === false" class="material-icons column-sort">keyboard_arrow_down</i>
                                    </a>
                                </th>
                                <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                    <a class="column-header" @click="sortRows('Opened')">
                                        Opened Last 7 Days
                                        <i v-if="sortedBy === 'Opened' && sorting['Opened'] === true" class="material-icons column-sort">keyboard_arrow_up</i>
                                        <i v-if="sortedBy === 'Opened' && sorting['Opened'] === false" class="material-icons column-sort">keyboard_arrow_down</i>
                                    </a>
                                </th>
                                <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                    <a class="column-header" @click="sortRows('Closed')">
                                        Closed Last 7 Days
                                        <i v-if="sortedBy === 'Closed' && sorting['Closed'] === true" class="material-icons column-sort">keyboard_arrow_up</i>
                                        <i v-if="sortedBy === 'Closed' && sorting['Closed'] === false" class="material-icons column-sort">keyboard_arrow_down</i>
                                    </a>
                                </th>
                            </tr>
                        </thead>
                        <tbody class="mdc-data-table__content">
                            <tr v-for="(target, index) in targets" class="mdc-data-table__row">
                                <td class="mdc-data-table__cell">[[ target.Name ]]</td>
                                <td class="mdc-data-table__cell">[[ target.CreatedDate ]]</td>
                                <td class="mdc-data-table__cell">[[ target.LastExtractionDate ]]</td>
                                <td class="mdc-data-table__cell">[[ target.Employees ]]</td>
                                <td class="mdc-data-table__cell">[[ target.JobsAll ]]</td>
                                <td class="mdc-data-table__cell">[[ target.JobsNow ]]</td>
                                <td class="mdc-data-table__cell">[[ target.Opened ]]</td>
                                <td class="mdc-data-table__cell">[[ target.Closed ]]</td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>`,
};