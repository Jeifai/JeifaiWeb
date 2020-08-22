export default {
    name: 'keywords',
    delimiters: ["[[","]]"],
    data: function () {
        return {
            selectedIndex: 2,
            messages: '',
            targets: [],
            keywordsInfo: [],
            utks: [],
            newKeyword: {},
            selectedTargets: null,
            filter: '',
            checkAll: false,
            checks: [],
            sortedByKeywords: "CreatedDate",
            sortingKeywords: {
                CreatedDate: true,
                Name: false,
                CountTargets: false,
                TotalMatches: false,
                LastWeekMatches: false,
                AvgMatchesDay: false,
            },
            sortedByCombinations: "CreatedDate",
            sortingCombinations: {
                CreatedDate: true,
                KeywordText: false,
                TargetName: false
            }
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
                max-width: 80%;
            }
            .removeSelected {
                --mdc-theme-primary: #ea5a3d;
                --mdc-theme-secondary: #ea5a3d;
            }
            .row {
                display: flex;
            }
            .column {
                flex: 50%;
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
                this.utks = response.data.Utks,
                this.keywordsInfo = response.data.KeywordsInfo
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
                    this.fetchKeywords()
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
                this.$http.delete('/keywords', JSON.stringify(payload)).then(function(response) {
                    this.messages = response.data.Messages
                    this.checks = []
                    this.checkAll = false
                    this.fetchKeywords()
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
        },
        unselectAllTargets: function() {
            this.selectedTargets = [];
        },
        sortRowsKeywords: function(column) {
            this.sortedByKeywords = column
            if (column == "CreatedDate") {
                if (this.sortingKeywords[column]) {
                    this.keywordsInfo.sort((a,b) => (new Date(a[column]) - new Date(b[column])))
                    this.sortingKeywords[column] = false
                } else {
                    this.utks.sort((b,a) => (new Date(a[column]) - new Date(b[column])))
                    this.sortingKeywords[column] = true
                }
            } else {
                if (this.sortingKeywords[column]) {
                    this.keywordsInfo.sort((a,b) => (a[column] > b[column]) ? 1 : ((b[column] > a[column]) ? -1 : 0))
                    this.sortingKeywords[column] = false
                } else {
                    this.keywordsInfo.sort((b,a) => (a[column] > b[column]) ? 1 : ((b[column] > a[column]) ? -1 : 0))
                    this.sortingKeywords[column] = true
                }
            }
        },
        sortRowsCombinations: function(column) {
            this.sortedByCombinations = column
            if (column == "CreatedDate") {
                if (this.sortingCombinations[column]) {
                    this.utks.sort((a,b) => (new Date(a[column]) - new Date(b[column])))
                    this.sortingCombinations[column] = false
                } else {
                    this.utks.sort((b,a) => (new Date(a[column]) - new Date(b[column])))
                    this.sortingCombinations[column] = true
                }
            } else {
                if (this.sortingCombinations[column]) {
                    this.utks.sort((a,b) => (a[column] > b[column]) ? 1 : ((b[column] > a[column]) ? -1 : 0))
                    this.sortingCombinations[column] = false
                } else {
                    this.utks.sort((b,a) => (a[column] > b[column]) ? 1 : ((b[column] > a[column]) ? -1 : 0))
                    this.sortingCombinations[column] = true
                }
            }
        },
        csvExport: function() {
            this.$parent.csvExport(this.filteredRows, "keywords.csv")
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
            <div class="keywords">
                <h1>
                    My Keywords
                </h1>
                <table class="mdc-data-table__table" aria-label="Your Keywords">
                    <thead>
                        <tr class="mdc-data-table__header-row">
                            <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                <a class="column-header" @click="sortRowsKeywords('CreatedDate')">
                                    CreatedAt
                                    <i v-if="sortedByKeywords === 'CreatedDate' && sortingKeywords['CreatedDate'] === true" class="material-icons column-sort">keyboard_arrow_up</i>
                                    <i v-if="sortedByKeywords === 'CreatedDate' && sortingKeywords['CreatedDate'] === false" class="material-icons column-sort">keyboard_arrow_down</i>
                                </a>
                            </th>
                            <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                <a class="column-header" @click="sortRowsKeywords('Name')">
                                    Keyword Text
                                    <i v-if="sortedByKeywords === 'Name' && sortingKeywords['Name'] === true" class="material-icons column-sort">keyboard_arrow_up</i>
                                    <i v-if="sortedByKeywords === 'Name' && sortingKeywords['Name'] === false" class="material-icons column-sort">keyboard_arrow_down</i>
                                </a>
                            </th>
                            <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                <a class="column-header" @click="sortRowsKeywords('CountTargets')">
                                    Count Targets
                                    <i v-if="sortedByKeywords === 'CountTargets' && sortingKeywords['CountTargets'] === true" class="material-icons column-sort">keyboard_arrow_up</i>
                                    <i v-if="sortedByKeywords === 'CountTargets' && sortingKeywords['CountTargets'] === false" class="material-icons column-sort">keyboard_arrow_down</i>
                                </a>
                            </th>
                            <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                <a class="column-header" @click="sortRowsKeywords('TotalMatches')">
                                    Total Matches
                                    <i v-if="sortedByKeywords === 'TotalMatches' && sortingKeywords['TotalMatches'] === true" class="material-icons column-sort">keyboard_arrow_up</i>
                                    <i v-if="sortedByKeywords === 'TotalMatches' && sortingKeywords['TotalMatches'] === false" class="material-icons column-sort">keyboard_arrow_down</i>
                                </a>
                            </th>
                            <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                <a class="column-header" @click="sortRowsKeywords('LastWeekMatches')">
                                    Last 7 Days Matches
                                    <i v-if="sortedByKeywords === 'LastWeekMatches' && sortingKeywords['LastWeekMatches'] === true" class="material-icons column-sort">keyboard_arrow_up</i>
                                    <i v-if="sortedByKeywords === 'LastWeekMatches' && sortingKeywords['LastWeekMatches'] === false" class="material-icons column-sort">keyboard_arrow_down</i>
                                </a>
                            </th>
                            <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                <a class="column-header" @click="sortRowsKeywords('LastWeekMatches')">
                                    Avg Matches / Day
                                    <i v-if="sortedByKeywords === 'AvgMatchesDay' && sortingKeywords['AvgMatchesDay'] === true" class="material-icons column-sort">keyboard_arrow_up</i>
                                    <i v-if="sortedByKeywords === 'AvgMatchesDay' && sortingKeywords['AvgMatchesDay'] === false" class="material-icons column-sort">keyboard_arrow_down</i>
                                </a>
                            </th>
                        </tr>
                    </thead>
                    <tbody class="mdc-data-table__content">
                        <tr v-for="(row, index) in keywordsInfo" class="mdc-data-table__row">
                            <td class="mdc-data-table__cell" v-html="row.CreatedDate"></td>
                            <td class="mdc-data-table__cell" v-html="row.Name"></td>
                            <td class="mdc-data-table__cell" v-html="row.CountTargets"></td>
                            <td class="mdc-data-table__cell" v-html="row.TotalMatches"></td>
                            <td class="mdc-data-table__cell" v-html="row.LastWeekMatches"></td>
                            <td class="mdc-data-table__cell" v-html="row.AvgMatchesDay"></td>
                        </tr>
                    </tbody>
                </table>
            </div>
            <div class="row">
                <div class="column">
                    <div class="addkeyword"> 
                        <h1>
                            Create combination
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
                        <button class="mdc-button mdc-button--raised" v-on:click="selectAllTargets">
                            <span class="mdc-button__label">Select all</span>
                        </button>
                        <button class="mdc-button mdc-button--raised" v-on:click="unselectAllTargets">
                            <span class="mdc-button__label">Unselect all</span>
                        </button>
                        <button class="mdc-button mdc-button--raised" v-on:click="createKeyword">
                            <div class="mdc-button__ripple"></div>
                            <i class="material-icons mdc-button__icon" aria-hidden="true">check</i>
                            <span class="mdc-button__label">Add keyword</span>
                        </button>
                        <div class="multiselectfield">
                            <vue-multi-select
                                v-model="selectedTargets"
                                :options="targets"
                                class="multiselect"
                                placeholder="Selected targets"
                                :max-results="1000"
                            >
                            </vue-multi-select><br>
                        </div>
                    </div>
                </div>
                <div class="column">
                    <div class="combinations">
                        <h1>
                            My Combinations
                        </h1>
                        <div>
                            <label id="Filter" for="Filter" class="mdc-text-field mdc-text-field--filled mdc-text-field--with-leading-icon">
                                <i class="material-icons mdc-text-field__icon mdc-text-field__icon--leading" tabindex="0" role="button">text_fields</i>
                                <input class="mdc-text-field__input" type="text" aria-labelledby="Filter" v-model="filter">
                                <span class="mdc-floating-label" id="Filter">Filter by any field</span>
                                <div class="mdc-line-ripple"></div>
                            </label>
                        </div>
                        <div>
                            <button class="mdc-button mdc-button--raised removeSelected" v-on:click="deleteUtks">
                                <div class="mdc-button__ripple"></div>
                                <i class="material-icons mdc-button__icon" aria-hidden="true">clear</i>
                                <span class="mdc-button__label">Remove selected</span>
                            </button>
                            <button class="mdc-button mdc-button--raised" v-on:click="csvExport">
                                <div class="mdc-button__ripple"></div>
                                <i class="material-icons mdc-button__icon" aria-hidden="true">arrow_downward</i>
                                <span class="mdc-button__label">Export table</span>
                            </button>
                        </div><br>
                        <div class="mdc-data-table">
                            <table class="mdc-data-table__table" aria-label="Created Keywords">
                                <thead>
                                    <tr class="mdc-data-table__header-row">
                                        <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                            <input type="checkbox" v-model="checkAll" @click="selectAll">
                                        </th>
                                        <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                            <a class="column-header" @click="sortRowsCombinations('CreatedDate')">
                                                CreatedAt
                                                <i v-if="sortedByCombinations === 'CreatedDate' && sortingCombinations['CreatedDate'] === true" class="material-icons column-sort">keyboard_arrow_up</i>
                                                <i v-if="sortedByCombinations === 'CreatedDate' && sortingCombinations['CreatedDate'] === false" class="material-icons column-sort">keyboard_arrow_down</i>
                                            </a>
                                        </th>
                                        <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                            <a class="column-header" @click="sortRowsCombinations('KeywordText')">
                                                Keyword
                                                <i v-if="sortedByCombinations === 'KeywordText' && sortingCombinations['KeywordText'] === true" class="material-icons column-sort">keyboard_arrow_up</i>
                                                <i v-if="sortedByCombinations === 'KeywordText' && sortingCombinations['KeywordText'] === false" class="material-icons column-sort">keyboard_arrow_down</i>
                                            </a>
                                        </th>
                                        <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                            <a class="column-header" @click="sortRowsCombinations('TargetName')">
                                                Target
                                                <i v-if="sortedByCombinations === 'TargetName' && sortingCombinations['TargetName'] === true" class="material-icons column-sort">keyboard_arrow_up</i>
                                                <i v-if="sortedByCombinations === 'TargetName' && sortingCombinations['TargetName'] === false" class="material-icons column-sort">keyboard_arrow_down</i>
                                            </a>
                                        </th>
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
                </div>
            </div>
        </div>`,
};