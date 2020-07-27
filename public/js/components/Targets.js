export default {
    name: 'targetsT',
    delimiters: ["[[","]]"],
    data: function () {
        return {
            selectedIndex: 1,
            messages: '',
            targets:  [],
            nameTargets: '',
            selectedTargets: null,
            newTarget: {}
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
            }`
        document.head.appendChild(styleElem);
    },
    created () {
        this.fetchTargets()
    },
    methods: {
        fetchTargets: function() {
            this.$http.get('/testTargets').then(function(response) {
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
                    this.targets = response.data.Targets
                    this.newTarget = {}
            }).catch(function(error) {
                console.log(error)
            });
        },
        deleteTarget: function(index) {
            payload = {"Name": this.targets[index].Name}
            this.$http.put('/targets/remove', payload).then(function(response) {
                this.messages = response.data.Messages
                this.targets.splice(index, 1)
            }).catch(function(error) {
                console.log(error)
            });
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
                >
                </vue-taggable-select><br>
                <button class="mdc-button mdc-button--raised" v-on:click="createTarget">
                    <div class="mdc-button__ripple"></div>
                    <i class="material-icons mdc-button__icon" aria-hidden="true">check</i>
                    <span class="mdc-button__label">Add target</span>
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
                                <th class="mdc-data-table__header-cell" role="columnheader" scope="col">CreatedAt</th>
                                <th class="mdc-data-table__header-cell" role="columnheader" scope="col">Target</th>
                                <th class="mdc-data-table__header-cell" role="columnheader" scope="col">All Jobs</th>
                                <th class="mdc-data-table__header-cell" role="columnheader" scope="col">Actual Jobs</th>
                                <th class="mdc-data-table__header-cell" role="columnheader" scope="col">Opened Last 7 Days</th>
                                <th class="mdc-data-table__header-cell" role="columnheader" scope="col">Closed Last 7 Days</th>
                            </tr>
                        </thead>
                        <tbody class="mdc-data-table__content">
                            <tr v-for="(target, index) in targets" class="mdc-data-table__row">
                                <td class="mdc-data-table__cell">[[ target.CreatedDate ]]</td>
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