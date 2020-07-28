export default {
    name: 'keywords',
    delimiters: ["[[","]]"],
    data: function () {
        return {
            selectedIndex: 2,
            messages: '',
            targets: [],
            utks: [],
            newKeyword: {},
            selectedTargets: null,
            filter: '',
            checkAll: false,
            checks: []
        }
    },
    mounted() {
        this.$parent.selectedIndex = this.selectedIndex
        let multiSelectScript = document.createElement('script')
        multiSelectScript.setAttribute('src', 'https://unpkg.com/vue-simple-multi-select@latest')
        document.head.appendChild(multiSelectScript)
        const topAppBarElement = mdc.dataTable.MDCDataTable.attachTo(document.querySelector('.mdc-data-table'))
        const textKeyword = mdc.textField.MDCTextField.attachTo(document.getElementById("Keyword"))
        const textFilter= mdc.textField.MDCTextField.attachTo(document.getElementById("Filter"))
        let styleElem = document.createElement('style');
        styleElem.textContent = `
            .multiselectfield {
                max-width: 35%;
            }
            .removeSelected {
                --mdc-theme-primary: #ea5a3d;
                --mdc-theme-secondary: #ea5a3d;
            }`
        document.head.appendChild(styleElem);
    },
    created () {
        this.fetchKeywords()
    },
    methods: {
        fetchKeywords: function() {
            this.$http.get('/keywords').then(function(response) {
                this.targets = response.data.Targets,
                this.utks = response.data.Utks
            }).catch(function(error) {
                console.log(error)
            });
        },
        createKeyword: function() {
            this.$http.put('/keywords', {
                "selectedTargets": this.selectedTargets,
                "newKeyword": this.newKeyword
                }).then(function(response) {
                    this.messages = response.data.Messages
                    this.utks = response.data.Utks;
            }).catch(function(error) {
                console.log(error)
            });
        },
        deleteUtks: function() {
            var payload = new Array();
            if (this.checks.length > 0) {
                for (var i = 0; i < this.checks.length; i++) {
                    payload.push({
                        "TargetName": this.filteredRows[this.checks[i]].TargetName,
                        "KeywordText": this.filteredRows[this.checks[i]].KeywordText
                    })
                }
                console.log(payload)
                payload_string = JSON.stringify(payload);
                this.$http.put('/keywords/remove', payload_string).then(function(response) {
                    this.messages = response.data.Messages
                    this.utks = response.data.Utks;
                    this.checks = [];
                    this.checkAll = false;
                }).catch(function(error) {
                    console.log(error)
                });
            }
        },
        selectAllTargets: function() {
            this.selectedTargets = this.targets;
        },
        selectAll: function() {
            this.checks = [];
            if (!this.checkAll) {
                for (var i = 0; i < this.filteredRows.length; i++) {
                    this.checks.push(i)
                }
            }
        }
    },
    computed: {
        filteredRows() {
            if (this.utks != null) {
                return this.utks.filter(row => {
                    const CreatedDate = row.CreatedDate.toString().toLowerCase();
                    const KeywordText = row.KeywordText.toString().toLowerCase();
                    const TargetName = row.TargetName.toString().toLowerCase();
                    const searchTerm = this.filter.toLowerCase();
                    return (
                        CreatedDate.includes(searchTerm) || 
                        KeywordText.includes(searchTerm) || 
                        TargetName.includes(searchTerm)
                    );
                });
            }
        }
    },
    template: `
        <div>
            <div class="addkeyword"> 
                <h1>
                    Add a new keyword
                </h1>
                <p v-for="(message, index) in messages">
                    <span v-html="message"></span>
                </p>
                <div>
                    <label id="Keyword" for="Keyword" class="mdc-text-field mdc-text-field--filled mdc-text-field--with-leading-icon">
                        <i class="material-icons mdc-text-field__icon mdc-text-field__icon--leading" tabindex="0" role="button">text_fields</i>
                        <input class="mdc-text-field__input" type="text" aria-labelledby="Keyword" v-model="newKeyword.Text">
                        <span class="mdc-floating-label" id="Keyword">Keyword</span>
                        <div class="mdc-line-ripple"></div>
                    </label>
                </div>
                <div class="multiselectfield">
                    <span><br>Keyword's targets</span>
                    <vue-multi-select
                        v-model="selectedTargets"
                        :options="targets"
                        class="multiselect"
                        placeholder="Targets"
                    >
                    </vue-multi-select><br>
                    <button class="mdc-button mdc-button--raised" v-on:click="selectAllTargets">
                        <span class="mdc-button__label">Select all</span>
                    </button>
                    <button class="mdc-button mdc-button--raised" v-on:click="createKeyword">
                        <div class="mdc-button__ripple"></div>
                        <i class="material-icons mdc-button__icon" aria-hidden="true">check</i>
                        <span class="mdc-button__label">Add keyword</span>
                    </button>
                </div>
            </div>
            <div>
                <h1>
                    My keywords
                </h1>
                <div>
                    <label id="Filter" for="Filter" class="mdc-text-field mdc-text-field--filled mdc-text-field--with-leading-icon">
                        <i class="material-icons mdc-text-field__icon mdc-text-field__icon--leading" tabindex="0" role="button">text_fields</i>
                        <input class="mdc-text-field__input" type="text" aria-labelledby="Filter" v-model="filter">
                        <span class="mdc-floating-label" id="Filter">Filter by any field</span>
                        <div class="mdc-line-ripple"></div>
                    </label>
                </div><br>
                <div>
                    <button class="mdc-button mdc-button--raised removeSelected" v-on:click="deleteUtks">
                        <div class="mdc-button__ripple"></div>
                        <i class="material-icons mdc-button__icon" aria-hidden="true">clear</i>
                        <span class="mdc-button__label">Remove selected</span>
                    </button>
                </div><br>
                <div class="mdc-data-table">
                    <table class="mdc-data-table__table" aria-label="Created Keywords">
                        <thead>
                            <tr class="mdc-data-table__header-row">
                                <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                    <input type="checkbox" v-model="checkAll" @click="selectAll">
                                </th>
                                <th class="mdc-data-table__header-cell" role="columnheader" scope="col">CreatedAt</th>
                                <th class="mdc-data-table__header-cell" role="columnheader" scope="col">Keyword</th>
                                <th class="mdc-data-table__header-cell" role="columnheader" scope="col">Target</th>
                            </tr>
                        </thead>
                        <tbody class="mdc-data-table__content">
                            <tr v-for="(row, index) in filteredRows" class="mdc-data-table__row">
                                <td class="mdc-data-table__cell">
                                    <input type="checkbox" v-model="checks" :value="index">
                                </td>
                                <td class="mdc-data-table__cell" v-html="row.CreatedDate"></td>
                                <td class="mdc-data-table__cell" v-html="row.KeywordText"></td>
                                <td class="mdc-data-table__cell" v-html="row.TargetName"></td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>`,
};