export default {
    name: 'targets',
    delimiters: ["[[","]]"],
    data: function () {
        return {
            selectedIndex: 1,
            messages: '',
            targets:  [],
            nameTargets: '',
            selectedTargets: null,
            newTarget: {},
            sortedBy: "CreatedDate",
            sorting: {
                CreatedDate: true,
                LastExtractionDate: true,
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
        let multiToggleScript = document.createElement('script')
        multiToggleScript.setAttribute('src', 'https://unpkg.com/vue-taggable-select@latest')
        document.head.appendChild(multiToggleScript)
        const topAppBarElement = mdc.dataTable.MDCDataTable.attachTo(document.querySelector('.mdc-data-table'));
        const button = mdc.ripple.MDCRipple.attachTo(document.querySelector('.mdc-icon-button'));
        let styleElem = document.createElement('style');
        styleElem.textContent = `
            .overflow-x-scroll {
                overflow-x: hidden !important;
            }
            .taggableselectfield {
                max-width: 35%;
            }
            .column-sort {
                font-size: 16px;
                vertical-align: -3px;
            }`
        document.head.appendChild(styleElem);
    },
    created () {
        this.fetchTargets()
    },
    methods: {
        fetchTargets: function() {
            this.$http.get('/targets').then(function(response) {
                this.targets = response.data.Targets,
                this.nameTargets = response.data.NameTargets
            }).catch(function(error) {
                console.log(error)
            });
        },
        createTarget: function() {
            this.$http.put('/targets', {
                "selectedTargets": this.selectedTargets
                }).then(function(response) {
                    this.messages = response.data.Messages
                    this.fetchTargets()
                    this.newTarget = {}
            }).catch(function(error) {
                console.log(error)
            });
        },
        deleteTarget: function(index) {
            payload = 
            this.$http.delete('/targets', {"Name": this.targets[index].Name}).then(
                function(response) {
                    this.messages = response.data.Messages
                    this.fetchTargets()
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
            <h1>
                Add a new target
            </h1>
            <p v-for="(message, index) in messages">
                <span v-html="message"></span>
            </p>
            <div class="taggableselectfield">
                <span><br>Select an existing target or add a new one.</span>
                <vue-taggable-select
                    v-model="selectedTargets"
                    :options="nameTargets"
                    class="multiselect"
                    placeholder="Targets"
                    :taggable="true"
                    :max-results="1000"
                >
                </vue-taggable-select><br>
                <button class="mdc-button mdc-button--raised" v-on:click="createTarget">
                    <div class="mdc-button__ripple"></div>
                    <i class="material-icons mdc-button__icon" aria-hidden="true">check</i>
                    <span class="mdc-button__label">Add target</span>
                </button>
                <button class="mdc-button mdc-button--raised" v-on:click="csvExport">
                    <div class="mdc-button__ripple"></div>
                    <i class="material-icons mdc-button__icon" aria-hidden="true">arrow_downward</i>
                    <span class="mdc-button__label">Export table</span>
                </button>
            </div>
            <div>
                <h1>
                    My targets
                </h1>
                <div class="mdc-data-table">
                    <table class="mdc-data-table__table" aria-label="Created Targets">
                        <thead>
                            <tr class="mdc-data-table__header-row">
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
                                    <a class="column-header" @click="sortRows('Name')">
                                        Target
                                        <i v-if="sortedBy === 'Name' && sorting['Name'] === true" class="material-icons column-sort">keyboard_arrow_up</i>
                                        <i v-if="sortedBy === 'Name' && sorting['Name'] === false" class="material-icons column-sort">keyboard_arrow_down</i>
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
                                        Actual Jobs
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
                                <td class="mdc-data-table__cell">[[ target.CreatedDate ]]</td>
                                <td class="mdc-data-table__cell">[[ target.LastExtractionDate ]]</td>
                                <td class="mdc-data-table__cell">[[ target.Name ]]</td>
                                <td class="mdc-data-table__cell">[[ target.JobsAll ]]</td>
                                <td class="mdc-data-table__cell">[[ target.JobsNow ]]</td>
                                <td class="mdc-data-table__cell">[[ target.Opened ]]</td>
                                <td class="mdc-data-table__cell">[[ target.Closed ]]</td>
                                <td class="mdc-data-table__cell">
                                    <button 
                                        class="material-icons mdc-top-app-bar__action-item mdc-icon-button" 
                                        aria-label="Clear" 
                                        v-on:click="deleteTarget(index)">clear
                                    </button>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>`,
};