export default {
    name: 'manage',
    delimiters: ["[[","]]"],
    data: function () {
        return {
            macroPivot: '',
            jobs: [],
            sortedByJobs: "CreatedDate",
            sortingJobs: {
                CreatedDate: true,
                TargetName: false,
                KeywordText: false,
                Title: false,
                Location: false,
                Url: false,
                IsSaved: false,
            },
            utksKeyword: {},
            utksTarget: {},
            allKeywords: [],
            userKeywords: [],
            inputKeyword: '',
            checksKeywords: [],
            sortedByKeywords: "CreatedDate",
            sortingKeywords: {
                CreatedDate: true,
                Text: false,
            },
            allTargets: [],
            userTargets: [],
            inputTarget: '',
            checksTargets: [],
            sortedByTargets: "CreatedDate",
            sortingTargets: {
                CreatedDate: true,
                Name: false,
            },
            autoCompleteStyle : {
                vueSimpleSuggest: "",
                inputWrapper: "",
                defaultInput : "form-control",
                suggestions: "suggestions-style",
                suggestItem: "list-group-item"
            },
            message: '',
            messageLoading: '',
        }
    },
    mounted() {
        this.$parent.selectedIndex = this.selectedIndex;
        mdc.textField.MDCTextField.attachTo(document.getElementById("KeywordsField"));
        mdc.textField.MDCTextField.attachTo(document.getElementById("TargetsField"));

        let styleElem = document.createElement('style');
        styleElem.textContent = `
            .row {
                display: flex;
            }
            .input-row {
                display: flex;
            }
            .column-match {
                padding-left: 3%;
                padding-right: 3%;
                display: flex;
                flex-direction: column;
                justify-content: space-between;
            }
            .combinations-button-div {
                height: 80%;
            }
            .combinations-button {
                top: 80%;
            }
            .message-div {
                height: 20%;
                text-align: center;"
                background-color: green;
            }
            .message-text {
                color: #1D3557;
                font-size: 25px;
            }
            .arrowBold {
                font-weight: bold !important;
                color: black !important;
            }
            .scrollable {
                overflow-y: scroll;
                height:30vh;
            }
            table.tableBodyScroll tbody {
              display: block;
              max-height: 35vh;
              overflow-y: scroll;
            }
            table.tableBodyScroll thead, table.tableBodyScroll tbody tr {
              display: table;
              table-layout: fixed;
              width: 98%;
              text-align: center;
            }
            .loader-message {
                width: 40px;
                height: 40px;
                border: 8px solid #f3f3f3;
                border-radius: 50%;
                border-top: 6px solid #457B9D;
                border-bottom: 6px solid #457B9D;
                margin: 0 auto;
                animation: spin 0.5s linear infinite;
            }
            @keyframes spin {
                0% { transform: rotate(0deg); }
                100% { transform: rotate(360deg); }
            }`
        document.head.appendChild(styleElem);
    },
    created: async function() {
        this.fetchUserTargetsKeywords();
        this.fetchUserKeywords();
        this.fetchAllKeywords();
        this.fetchUserTargets();
        this.fetchAllTargets();
        this.fetchResults();
    },
    methods: {
        fetchResults: function() {
            this.$http.get('/jobs').then(function(response) {
                if (response.data.Jobs !== null) {
                    this.jobs = response.data.Jobs;
                }
            }).catch(function(error) {
                console.log(error);
            });
        },
        fetchUserTargetsKeywords: function() {
            this.$http.get('/utks').then(function(response) {
                if (response.data.Utks !== null) {
                    this.utksKeyword = response.data.Utks[0];
                    this.utksTarget = response.data.Utks[1];
                }
            }).catch(function(error) {
                console.log(error);
            });
        },
        createUserTargetsKeywords: function() {
            this.messageLoading = true;
            var keywords = [];
            var targets = [];
            for (var i = 0; i < this.checksKeywords.length; i++) {
                keywords.push(this.userKeywords[this.checksKeywords[i]].Text);
            };
            for (var q = 0; q < this.checksTargets.length; q++) {
                targets.push(this.userTargets[this.checksTargets[q]].Name);
            };
            this.$http.put('/utks', {
                "macroPivot": this.macroPivot,
                "keywords": keywords,
                "targets": targets
            }).then(function(response) {
                    this.fetchUserTargetsKeywords();
                    this.fetchResults();
                    this.messageLoading = false;
                    this.message = response.data.Message;
                    setTimeout(() => this.message = '', 2000);
            }).catch(function(error) {
                console.log(error)
            });
        },  
        fetchUserKeywords: function() {
            this.$http.get('/keywords/user').then(function(response) {
                this.userKeywords = response.data.Keywords;
            }).catch(function(error) {
                console.log(error);
            });
        },
        fetchAllKeywords: function() {
            this.$http.get('/keywords/all').then(function(response) {
                this.allKeywords = response.data.Keywords;
            }).catch(function(error) {
                console.log(error);
            });
        },
        createKeyword: function() {
            this.messageLoading = true;
            this.$http.put('/keywords/' + this.inputKeyword).then(
                function(response) {
                    this.fetchUserKeywords();
                    this.fetchAllKeywords();
                    this.fetchUserTargetsKeywords();
                    this.inputKeyword = '';
                    this.messageLoading = false;
                    this.message = response.data.Message;
                    setTimeout(() => this.message = '', 2000);
            }).catch(function(error) {
                console.log(error);
            });
        },
        deleteKeyword: function(index) {
            if (confirm("Are you sure you want to delete the keyword?")) {
                this.messageLoading = true;
                this.$http.delete('/keywords/' + this.userKeywords[index].Text).then(
                    function(response) {
                        this.fetchUserKeywords();
                        this.fetchAllKeywords();
                        this.fetchUserTargetsKeywords();
                        this.fetchResults();
                        this.messageLoading = false;
                        this.message = response.data.Message;
                        setTimeout(() => this.message = '', 2000);
                }).catch(function(error) {
                    console.log(error)
                });
            }
        },
        selectAllKeywords: function() {
            if (this.checksKeywords.length <  this.userKeywords.length) {
                for (var i = 0; i < this.userKeywords.length; i++) {
                    this.checksKeywords.push(i)
                }
                return;
            }
            this.checksKeywords = [];
        },
        sortRowsKeywords: function(column) {
            this.sortedByKeywords = column
            if (column == "CreatedDate") {
                if (this.sortingKeywords[column]) {
                    this.userKeywords.sort((a,b) => (new Date(a[column]) - new Date(b[column])))
                    this.sortingKeywords[column] = false
                } else {
                    this.userKeywords.sort((b,a) => (new Date(a[column]) - new Date(b[column])))
                    this.sortingKeywords[column] = true
                }
            } else {
                if (this.sortingKeywords[column]) {
                    this.userKeywords.sort((a,b) => (a[column] > b[column]) ? 1 : ((b[column] > a[column]) ? -1 : 0))
                    this.sortingKeywords[column] = false
                } else {
                    this.userKeywords.sort((b,a) => (a[column] > b[column]) ? 1 : ((b[column] > a[column]) ? -1 : 0))
                    this.sortingKeywords[column] = true
                }
            }
        },
        fetchUserTargets: function() {
            this.$http.get('/targets/user').then(function(response) {
                this.userTargets = response.data.Targets;
            }).catch(function(error) {
                console.log(error);
            });
        },
        fetchAllTargets: function() {
            this.$http.get('/targets/all').then(function(response) {
                this.allTargets = response.data.Targets;
            }).catch(function(error) {
                console.log(error);
            });
        },
        createTarget: function() {
            this.messageLoading = true;
            this.$http.put('/targets/' + this.inputTarget).then(
                function(response) {
                    this.fetchUserTargets();
                    this.fetchAllTargets();
                    this.fetchUserTargetsKeywords();
                    this.inputTarget = '';
                    this.messageLoading = false;
                    this.message = response.data.Message;
                    setTimeout(() => this.message = '', 2000);
            }).catch(function(error) {
                console.log(error);
            });
        },
        deleteTarget: function(index) {
            if (confirm("Are you sure you want to delete the target?")) {
                this.messageLoading = true;
                this.$http.delete('/targets/' + this.userTargets[index].Name).then(
                    function(response) {
                        this.fetchUserTargets();
                        this.fetchAllTargets();
                        this.fetchUserTargetsKeywords();
                        this.fetchResults();
                        this.messageLoading = false;
                        this.message = response.data.Message;
                        setTimeout(() => this.message = '', 2000);
                }).catch(function(error) {
                    console.log(error)
                });
            }
        },
        selectAllTargets: function() {
            if (this.checksTargets.length <  this.userTargets.length) {
                for (var i = 0; i < this.userTargets.length; i++) {
                    this.checksTargets.push(i)
                }
                return;
            }
            this.checksTargets = [];
        },
        sortRowsTargets: function(column) {
            this.sortedByTargets = column
            if (column == "CreatedDate") {
                if (this.sortingTargets[column]) {
                    this.userTargets.sort((a,b) => (new Date(a[column]) - new Date(b[column])))
                    this.sortingTargets[column] = false
                } else {
                    this.userTargets.sort((b,a) => (new Date(a[column]) - new Date(b[column])))
                    this.sortingTargets[column] = true
                }
            } else {
                if (this.sortingTargets[column]) {
                    this.userTargets.sort((a,b) => (a[column] > b[column]) ? 1 : ((b[column] > a[column]) ? -1 : 0))
                    this.sortingTargets[column] = false
                } else {
                    this.userTargets.sort((b,a) => (a[column] > b[column]) ? 1 : ((b[column] > a[column]) ? -1 : 0))
                    this.sortingTargets[column] = true
                }
            }
        },
        sortCheckboxes: function(arr_values, arr_checks) {
            var temp_result_present = [];
            var temp_result_not_present = [];
            var temp_checkbox = [];
            var matches = 0
            for (var i = 0; i < arr_values.length; i++) {
                var is_present = false;
                for (var x = 0; x < arr_checks.length; x++) {
                    if (i === arr_checks[x]) {
                        is_present = true;
                        temp_result_present.push(arr_values[i]);
                        temp_checkbox.push(matches);
                        matches++;
                    }
                }
                if (!is_present) {
                    temp_result_not_present.push(arr_values[i]);
                }
            }
            return [
                temp_result_present.concat(temp_result_not_present),
                temp_checkbox
            ];
        },
        updateCheckboxes: function(pivotOn, index) {

            // INSTANCIATE MACRO PIVOT
            if (this.macroPivot == '' && this.checksKeywords.length == 1) {
                this.macroPivot = 'keywords';
            } else if (this.macroPivot == '' && this.checksTargets.length == 1) {
                this.macroPivot = 'targets';
            }

            // DO NOTHING IF THE USER IS PLAYING AROUND
            if (this.macroPivot == 'targets' && pivotOn == 'keywords' && this.checksTargets.length == 1) return;
            if (this.macroPivot == 'keywords' && pivotOn == 'targets' && this.checksKeywords.length == 1) return;

            // RESET INITIAL CONDITIONS
            if (pivotOn == 'keywords' && this.checksKeywords.length == 0) {
                this.checksTargets = [];
                this.macroPivot = '';
                document.getElementById('table-targets').scrollTop = 0;
                return;
            }
            if (pivotOn == 'targets' && this.checksTargets.length == 0) {
                this.checksKeywords = [];
                this.macroPivot = '';
                document.getElementById('table-keywords').scrollTop = 0;
                return;
            }

            // FILL CHECKBOXES BASED ON DICT
            if (pivotOn == 'keywords' && this.checksKeywords.length == 1) {
                var selectedKeyword = this.userKeywords[index].Text;
                var keywordTargets = this.utksKeyword[selectedKeyword];
                if (keywordTargets !== null) {
                    for (var i = 0; i < this.userTargets.length; i++) {
                        if (keywordTargets.includes(this.userTargets[i].Name)) {
                            this.checksTargets.push(i);
                        }
                    }
                    var sorted_elem = this.sortCheckboxes(this.userTargets, this.checksTargets);
                    this.userTargets = sorted_elem[0];
                    this.checksTargets = sorted_elem[1];
                }
            }
            if (pivotOn == 'targets' && this.checksTargets.length == 1) {
                var selectedTarget = this.userTargets[index].Name;
                var targetKeywords = this.utksTarget[selectedTarget];
                if (targetKeywords !== null) {
                    for (var i = 0; i < this.userKeywords.length; i++) {
                        if (targetKeywords.includes(this.userKeywords[i].Text)) {
                            this.checksKeywords.push(i);
                        }
                    }
                    var sorted_elem = this.sortCheckboxes(this.userKeywords, this.checksKeywords);
                    this.userKeywords = sorted_elem[0];
                    this.checksKeywords = sorted_elem[1];
                }
            }
        },
        sortRowsJobs: function(column) {
            this.sortedByJobs = column
            if (column == "CreatedDate") {
                if (this.sortingJobs[column]) {
                    this.jobs.sort((a,b) => (new Date(a[column]) - new Date(b[column])))
                    this.sortingJobs[column] = false
                } else {
                    this.jobs.sort((b,a) => (new Date(a[column]) - new Date(b[column])))
                    this.sortingJobs[column] = true
                }
            } else {
                if (this.sortingJobs[column]) {
                    this.jobs.sort((a,b) => (a[column] > b[column]) ? 1 : ((b[column] > a[column]) ? -1 : 0))
                    this.sortingJobs[column] = false
                } else {
                    this.jobs.sort((b,a) => (a[column] > b[column]) ? 1 : ((b[column] > a[column]) ? -1 : 0))
                    this.sortingJobs[column] = true
                }
            }
        },
        select: function(row) {
            window.open(row.Url, "_blank");
        },
        save: function(row, event) {
            this.$http.put('/favourite/' + row.Id + '/' + event.target.checked).then(function(response) {
                    this.message = response.data.Message;
                    setTimeout(() => this.message = '', 2000);
            }).catch(function(error) {
                console.log(error)
            });
        },
        arrowDefine: function (sortedBy, column) {
            return {
                'arrowBold': sortedBy == column,
            }
        }
    },
    computed: {
        checkAllKeywords() {
            if (this.checksKeywords.length == this.userKeywords.length) {
                return true;
            } else {
                return false;
            }
        },
        checkAllTargets() {
            if (this.checksTargets.length == this.userTargets.length) {
                return true;
            } else {
                return false;
            }
        },
        filteredRows() {
            var string_keywords = "";
            for (var i = 0; i < this.checksKeywords.length; i++) {
                string_keywords = string_keywords + this.userKeywords[this.checksKeywords[i]].Text.toLowerCase();
            }
            var string_targets = "";
            for (var i = 0; i < this.checksTargets.length; i++) {
                string_targets = string_targets + this.userTargets[this.checksTargets[i]].Name.toLowerCase();
            }
            return this.jobs.filter(row => {
                const KeywordText = row.KeywordText.toString().toLowerCase();
                const TargetName = row.TargetName.toString().toLowerCase();
                if (this.checksKeywords.length == 0 && this.checksTargets.length == 0) {
                    return true;
                } else {
                    return string_keywords.includes(KeywordText) && string_targets.includes(TargetName);
                }
            })
        },
    },
    template: `
        <div>
            <div class="row">
                <div>
                    <div>
                        <h3>
                            My Keywords
                        </h3>
                        <div class="input-row">
                            <vue-simple-suggest
                                v-model="inputKeyword"
                                :list="allKeywords"
                                :styles="autoCompleteStyle"
                                :destyled=true
                                :filter-by-query="true">
                                <label id="KeywordsField" for="Filter" class="mdc-text-field mdc-text-field--filled mdc-text-field--with-leading-icon">
                                    <i class="material-icons mdc-text-field__icon mdc-text-field__icon--leading" tabindex="0" role="button">text_fields</i>
                                    <input class="mdc-text-field__input" type="text">
                                    <span class="mdc-floating-label">Add keyword</span>
                                    <div class="mdc-line-ripple"></div>
                                </label>
                            </vue-simple-suggest>
                            <button
                                v-on:click="createKeyword"
                                class="material-icons mdc-top-app-bar__action-item mdc-icon-button" 
                                aria-label="Add">add
                            </button>
                        </div>
                        <div class="mdc-data-table scrollable" id="table-keywords">
                            <table class="mdc-data-table__table" aria-label="Created Keywords">
                                <thead>
                                    <tr class="mdc-data-table__header-row">
                                        <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                            <input
                                                v-if="checksTargets.length && (macroPivot != 'keywords')"
                                                type="checkbox" 
                                                v-model="checkAllKeywords"
                                                @click="selectAllKeywords"
                                                :disabled="macroPivot == 'keywords' && checksKeywords.length == 1"
                                            >
                                        </th>
                                        <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                            <a class="column-header" @click="sortRowsKeywords('CreatedDate')">
                                                CreatedDate
                                                <i v-if="sortingKeywords['CreatedDate'] === true" class="material-icons column-sort" v-bind:class="{'arrowBold': sortedByKeywords == 'CreatedDate'}">
                                                    keyboard_arrow_up
                                                </i>
                                                <i v-if="sortingKeywords['CreatedDate'] === false" class="material-icons column-sort" v-bind:class="{'arrowBold': sortedByKeywords == 'CreatedDate'}">
                                                    keyboard_arrow_down
                                                </i>
                                            </a>
                                        </th>
                                        <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                            <a class="column-header" @click="sortRowsKeywords('Text')">
                                                Keyword
                                                <i v-if="sortingKeywords['Text'] === true" class="material-icons column-sort" v-bind:class="{'arrowBold': sortedByKeywords == 'Text'}">
                                                    keyboard_arrow_up
                                                </i>
                                                <i v-if="sortingKeywords['Text'] === false" class="material-icons column-sort" v-bind:class="{'arrowBold': sortedByKeywords == 'Text'}">
                                                    keyboard_arrow_down
                                                </i>
                                            </a>
                                        </th>
                                    </tr>
                                </thead>
                                <tbody class="mdc-data-table__content">
                                    <tr v-for="(row, index) in userKeywords" class="mdc-data-table__row">
                                        <td class="mdc-data-table__cell">
                                            <input
                                                type="checkbox"
                                                v-model="checksKeywords"
                                                :value="index"
                                                @change="updateCheckboxes('keywords', index)"
                                                :disabled="macroPivot == 'keywords' && checksKeywords[0] != index"
                                            >
                                        </td>
                                        <td class="mdc-data-table__cell" v-html="row.CreatedDate"></td>
                                        <td class="mdc-data-table__cell" v-html="row.Text"></td>
                                        <td class="mdc-data-table__cell">
                                            <button
                                                v-on:click="deleteKeyword(index)"
                                                class="material-icons mdc-top-app-bar__action-item mdc-icon-button" 
                                                aria-label="Clear">clear
                                            </button>
                                        </td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
                <div class="column-match">
                    <div class="combinations-button-div">
                        <button
                            v-on:click="createUserTargetsKeywords"
                            class="mdc-button mdc-button--raised combinations-button"
                            :disabled="checksKeywords.length == 0 && checksTargets.length == 0">
                            <div class="mdc-button__ripple"></div>
                            <span class="mdc-button__label">Save combinations</span>
                        </button>
                    </div>
                    <div v-if="messageLoading" class="message-div loader-message"></div>
                    <div v-if="message" class="message-div">
                        <span v-html="message" class="message-text"></span>
                    </div>
                </div>
                <div>
                    <div>
                        <h3>
                            My Targets
                        </h3>
                        <div class="input-row">
                            <vue-simple-suggest
                                v-model="inputTarget"
                                :list="allTargets"
                                :styles="autoCompleteStyle"
                                :destyled=true
                                :filter-by-query="true">
                                <label id="TargetsField" for="Filter" class="mdc-text-field mdc-text-field--filled mdc-text-field--with-leading-icon">
                                    <i class="material-icons mdc-text-field__icon mdc-text-field__icon--leading" tabindex="0" role="button">text_fields</i>
                                    <input class="mdc-text-field__input" type="text">
                                    <span class="mdc-floating-label">Add target</span>
                                    <div class="mdc-line-ripple"></div>
                                </label>
                            </vue-simple-suggest>
                            <button
                                v-on:click="createTarget"
                                class="material-icons mdc-top-app-bar__action-item mdc-icon-button" 
                                aria-label="Add">add
                            </button>
                        </div>
                        <div class="mdc-data-table scrollable" id="table-targets">
                            <table class="mdc-data-table__table" aria-label="Created Targets">
                                <thead>
                                    <tr class="mdc-data-table__header-row">
                                        <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                            <input
                                                v-if="checksKeywords.length && (macroPivot != 'targets')"
                                                type="checkbox"
                                                v-model="checkAllTargets"
                                                @click="selectAllTargets"
                                                :disabled="macroPivot == 'targets' && checksTargets.length == 1"
                                            >
                                        </th>
                                        <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                            <a class="column-header" @click="sortRowsTargets('CreatedDate')">
                                                CreatedDate
                                                <i v-if="sortingTargets['CreatedDate'] === true" class="material-icons column-sort" v-bind:class="{'arrowBold': sortedByTargets == 'CreatedDate'}">
                                                    keyboard_arrow_up
                                                </i>
                                                <i v-if="sortingTargets['CreatedDate'] === false" class="material-icons column-sort" v-bind:class="{'arrowBold': sortedByTargets == 'CreatedDate'}">
                                                    keyboard_arrow_down
                                                </i>
                                            </a>
                                        </th>
                                        <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                            <a class="column-header" @click="sortRowsTargets('Name')">
                                                Target
                                                <i v-if="sortingTargets['Name'] === true" class="material-icons column-sort" v-bind:class="{'arrowBold': sortedByTargets == 'Name'}">
                                                    keyboard_arrow_up
                                                </i>
                                                <i v-if="sortingTargets['Name'] === false" class="material-icons column-sort" v-bind:class="{'arrowBold': sortedByTargets == 'Name'}">
                                                    keyboard_arrow_down
                                                </i>
                                            </a>
                                        </th>
                                    </tr>
                                </thead>
                                <tbody class="mdc-data-table__content">
                                    <tr v-for="(row, index) in userTargets" class="mdc-data-table__row">
                                        <td class="mdc-data-table__cell">
                                            <input
                                                type="checkbox"
                                                v-model="checksTargets"
                                                :value="index"
                                                @change="updateCheckboxes('targets', index)"
                                                :disabled="macroPivot == 'targets' && checksTargets[0] != index"
                                            >
                                        </td>
                                        <td class="mdc-data-table__cell" v-html="row.CreatedDate"></td>
                                        <td class="mdc-data-table__cell" v-html="row.Name"></td>
                                        <td class="mdc-data-table__cell">
                                            <button
                                                v-on:click="deleteTarget(index)"
                                                class="material-icons mdc-top-app-bar__action-item mdc-icon-button" 
                                                aria-label="Clear">clear
                                            </button>
                                        </td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
            </div><br>
            <div class="row">
                <table class="mdc-data-table__table tableBodyScroll" aria-label="Results" style="">
                    <thead>
                        <tr class="mdc-data-table__header-row">
                            <th class="mdc-data-table__header-cell" role="columnheader" scope="col" style="width:7%;white-space:nowrap;">Job Url</th>
                            <th class="mdc-data-table__header-cell" role="columnheader" scope="col" style="width:7%;white-space:nowrap;">Save</th>
                            <th class="mdc-data-table__header-cell" role="columnheader" scope="col" style="width:11%;white-space:nowrap;">
                                <a class="column-header" @click="sortRowsJobs('CreatedDate')">
                                    CreatedAt
                                    <i v-if="sortingJobs['CreatedDate'] === true" class="material-icons column-sort" v-bind:class="{'arrowBold': sortedByJobs == 'CreatedDate'}">
                                        keyboard_arrow_up
                                    </i>
                                    <i v-if="sortingJobs['CreatedDate'] === false" class="material-icons column-sort" v-bind:class="{'arrowBold': sortedByJobs == 'CreatedDate'}">
                                        keyboard_arrow_down
                                    </i>
                                </a>
                            </th>
                            <th class="mdc-data-table__header-cell" role="columnheader" scope="col" style="width:11%;white-space:nowrap;">
                                <a class="column-header" @click="sortRowsJobs('KeywordText')">
                                    Keyword
                                    <i v-if="sortingJobs['KeywordText'] === true" class="material-icons column-sort" v-bind:class="{'arrowBold': sortedByJobs == 'KeywordText'}">
                                        keyboard_arrow_up
                                    </i>
                                    <i v-if="sortingJobs['KeywordText'] === false" class="material-icons column-sort" v-bind:class="{'arrowBold': sortedByJobs == 'KeywordText'}">
                                        keyboard_arrow_down
                                    </i>
                                </a>
                            </th>
                            <th class="mdc-data-table__header-cell" role="columnheader" scope="col" style="width:13%;white-space:nowrap;">
                                <a class="column-header" @click="sortRowsJobs('TargetName')">
                                    Target
                                    <i v-if="sortingJobs['TargetName'] === true" class="material-icons column-sort" v-bind:class="{'arrowBold': sortedByJobs == 'TargetName'}">
                                        keyboard_arrow_up
                                    </i>
                                    <i v-if="sortingJobs['TargetName'] === false" class="material-icons column-sort" v-bind:class="{'arrowBold': sortedByJobs == 'TargetName'}">
                                        keyboard_arrow_down
                                    </i>
                                </a>
                            </th>
                            <th class="mdc-data-table__header-cell" role="columnheader" scope="col" style="width:19%;white-space:nowrap;">
                                <a class="column-header" @click="sortRowsJobs('Location')">
                                    Location
                                    <i v-if="sortingJobs['Location'] === true" class="material-icons column-sort" v-bind:class="{'arrowBold': sortedByJobs == 'Location'}">
                                        keyboard_arrow_up
                                    </i>
                                    <i v-if="sortingJobs['Location'] === false" class="material-icons column-sort" v-bind:class="{'arrowBold': sortedByJobs == 'Location'}">
                                        keyboard_arrow_down
                                    </i>
                                </a>
                            </th>
                            <th class="mdc-data-table__header-cell" role="columnheader" scope="col" style="width:32%;white-space:nowrap;">
                                <a class="column-header" @click="sortRowsJobs('Title')">
                                    Title
                                    <i v-if="sortingJobs['Title'] === true" class="material-icons column-sort" v-bind:class="{'arrowBold': sortedByJobs == 'Title'}">
                                        keyboard_arrow_up
                                    </i>
                                    <i v-if="sortingJobs['Title'] === false" class="material-icons column-sort" v-bind:class="{'arrowBold': sortedByJobs == 'Title'}">
                                        keyboard_arrow_down
                                    </i>
                                </a>
                            </th>
                        </tr>
                    </thead>
                    <tbody class="mdc-data-table__content">
                        <tr v-for="(row, index) in filteredRows" class="mdc-data-table__row">
                            <td class="mdc-data-table__cell" style="width:7%;white-space:nowrap;">
                                <button 
                                    class="material-icons mdc-top-app-bar__action-item mdc-icon-button" 
                                    aria-label="Open"
                                    target="_blank" 
                                    rel="noopener"
                                    v-on:click="select(row)">open_in_new
                                </button>
                            </td>
                            <td class="mdc-data-table__cell" style="width:7%;white-space:nowrap;">
                                <input
                                    type="checkbox"
                                    v-model="row.IsSaved"
                                    @change="save(row,$event)"
                                >
                            </td>
                            <td class="mdc-data-table__cell" v-html="row.CreatedDate" style="width:11%;white-space:nowrap;"></td>
                            <td class="mdc-data-table__cell" v-html="row.KeywordText" style="width:11%;white-space:nowrap;"></td>
                            <td class="mdc-data-table__cell comment" v-html="row.TargetName" style="width:13%;white-space:nowrap;"></td>
                            <td class="mdc-data-table__cell comment" v-html="row.Location" style="width:19%;white-space:nowrap;" v-bind:title="row.Location"></td>
                            <td class="mdc-data-table__cell" v-html="row.Title" style="width:32%;white-space:nowrap;" v-bind:title="row.Title"></td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>`,
};