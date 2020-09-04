export default {
    name: 'keywords',
    delimiters: ["[[","]]"],
    data: function () {
        return {
            macroPivot: '',
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
            messagesKeywords: '',
            allTargets: [],
            userTargets: [],
            inputTarget: '',
            checksTargets: [],
            sortedByTargets: "CreatedDate",
            sortingTargets: {
                CreatedDate: true,
                Name: false,
            },
            messagesTargets: '',
            autoCompleteStyle : {
                vueSimpleSuggest: "",
                inputWrapper: "",
                defaultInput : "form-control",
                suggestions: "suggestions-style",
                suggestItem: "list-group-item"
            },
        }
    },
    mounted() {
        this.$parent.selectedIndex = this.selectedIndex
        mdc.textField.MDCTextField.attachTo(document.getElementById("KeywordsFilter"));
        mdc.textField.MDCTextField.attachTo(document.getElementById("TargetsFilter"));

        let styleElem = document.createElement('style');
        styleElem.textContent = `
            .row {
                display: flex;
            }
            .column {
                flex: 50%;
            }
            .input-row {
                display: flex;
            }
            .suggestions-style {
                position: absolute;
                z-index: 1000;
                text-align: center;
                width: 30%;
                background-color: rgba(245, 245, 245, 0.8);
            }
            .hover {
                background-color: #007bff;
                color: #fff;
            }
            .scrollable {
                overflow-y: scroll;
                height:36vh;
            }`
        document.head.appendChild(styleElem);
    },
    created: async function() {
        this.fetchUserTargetsKeywords();
        this.fetchUserKeywords();
        this.fetchAllKeywords();
        this.fetchUserTargets();
        this.fetchAllTargets();
    },
    methods: {
        fetchUserTargetsKeywords: function() {
            this.$http.get('/targets/keywords').then(function(response) {
                this.utksKeyword = response.data.Utks[0];
                this.utksTarget = response.data.Utks[1];
            }).catch(function(error) {
                console.log(error);
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
            this.$http.put('/keywords/' + this.inputKeyword).then(
                function(response) {
                    this.messagesKeywords = response.data.Messages;
                    this.fetchUserKeywords();
                    this.fetchAllKeywords();
                    this.fetchUserTargetsKeywords();
                    this.inputKeyword = '';
                    setTimeout(() => this.messagesKeywords = '', 2000);
            }).catch(function(error) {
                console.log(error);
            });
        },
        deleteKeyword: function(index) {
            if (confirm("Are you sure you want to delete the keyword?")) {
                this.$http.delete('/keywords/' + this.filteredKeywords[index].Text).then(
                    function(response) {
                        this.messagesKeywords = response.data.Messages;
                        this.fetchUserKeywords();
                        this.fetchAllKeywords();
                        this.fetchUserTargetsKeywords();
                        setTimeout(() => this.messagesKeywords = '', 2000);
                }).catch(function(error) {
                    console.log(error)
                });
            }
        },
        selectAllKeywords: function() {
            if (this.checksKeywords.length <  this.filteredKeywords.length) {
                for (var i = 0; i < this.filteredKeywords.length; i++) {
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
                    this.filteredKeywords.sort((a,b) => (new Date(a[column]) - new Date(b[column])))
                    this.sortingKeywords[column] = false
                } else {
                    this.filteredKeywords.sort((b,a) => (new Date(a[column]) - new Date(b[column])))
                    this.sortingKeywords[column] = true
                }
            } else {
                if (this.sortingKeywords[column]) {
                    this.filteredKeywords.sort((a,b) => (a[column] > b[column]) ? 1 : ((b[column] > a[column]) ? -1 : 0))
                    this.sortingKeywords[column] = false
                } else {
                    this.filteredKeywords.sort((b,a) => (a[column] > b[column]) ? 1 : ((b[column] > a[column]) ? -1 : 0))
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
            this.$http.put('/targets/' + this.inputTarget).then(
                function(response) {
                    this.messagesTargets = response.data.Messages;
                    this.fetchUserTargets();
                    this.fetchAllTargets();
                    this.fetchUserTargetsKeywords();
                    this.inputTarget = '';
                    setTimeout(() => this.messagesTargets = '', 2000);
            }).catch(function(error) {
                console.log(error);
            });
        },
        deleteTarget: function(index) {
            if (confirm("Are you sure you want to delete the target?")) {
                this.$http.delete('/targets/' + this.filteredTargets[index].Name).then(
                    function(response) {
                        this.messagesTargets = response.data.Messages;
                        this.fetchUserTargets();
                        this.fetchAllTargets();
                        this.fetchUserTargetsKeywords();
                        setTimeout(() => this.messagesTargets = '', 2000);
                }).catch(function(error) {
                    console.log(error)
                });
            }
        },
        selectAllTargets: function() {
            if (this.checksTargets.length <  this.filteredTargets.length) {
                for (var i = 0; i < this.filteredTargets.length; i++) {
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
                    this.filteredTargets.sort((a,b) => (new Date(a[column]) - new Date(b[column])))
                    this.sortingTargets[column] = false
                } else {
                    this.filteredTargets.sort((b,a) => (new Date(a[column]) - new Date(b[column])))
                    this.sortingTargets[column] = true
                }
            } else {
                if (this.sortingTargets[column]) {
                    this.filteredTargets.sort((a,b) => (a[column] > b[column]) ? 1 : ((b[column] > a[column]) ? -1 : 0))
                    this.sortingTargets[column] = false
                } else {
                    this.filteredTargets.sort((b,a) => (a[column] > b[column]) ? 1 : ((b[column] > a[column]) ? -1 : 0))
                    this.sortingTargets[column] = true
                }
            }
        },
        sortKeywordsCheckboxes: function() {
            this.checksKeywords.sort((b,a) => (a > b) ? 1 : ((b > a) ? -1 : 0))
        },
        sortCheckboxes: function(arr_filtered, arr_checks, arr_user) {
            var temp_result_present = [];
            var temp_result_not_present = [];
            var temp_checkbox = [];
            var matches = 0
            for (var i = 0; i < arr_filtered.length; i++) {
                var is_present = false;
                for (var x = 0; x < arr_checks.length; x++) {
                    if (i === arr_checks[x]) {
                        is_present = true;
                        temp_result_present.push(arr_filtered[i]);
                        temp_checkbox.push(matches);
                        matches++;
                    }
                }
                if (!is_present) {
                    temp_result_not_present.push(arr_filtered[i]);
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
                return;
            }
            if (pivotOn == 'targets' && this.checksTargets.length == 0) {
                this.checksKeywords = [];
                this.macroPivot = '';
                return;
            }

            // FILL CHECKBOXES BASED ON DICT
            if (pivotOn == 'keywords' && this.checksKeywords.length == 1) {
                var selectedKeyword = this.filteredKeywords[index].Text;
                var keywordTargets = this.utksKeyword[selectedKeyword];
                for (var i = 0; i < this.filteredTargets.length; i++) {
                    if (keywordTargets.includes(this.filteredTargets[i].Name)) {
                        this.checksTargets.push(i);
                    }
                }
                var sorted_elem = this.sortCheckboxes(this.filteredTargets, this.checksTargets, this.userTargets);
                this.userTargets = sorted_elem[0];
                this.checksTargets = sorted_elem[1];
            }
            if (pivotOn == 'targets' && this.checksTargets.length == 1) {
                var selectedTarget = this.filteredTargets[index].Name;
                var targetKeywords = this.utksTarget[selectedTarget];
                for (var i = 0; i < this.filteredKeywords.length; i++) {
                    if (targetKeywords.includes(this.filteredKeywords[i].Text)) {
                        this.checksKeywords.push(i);
                    }
                }
                var sorted_elem = this.sortCheckboxes(this.filteredKeywords, this.checksKeywords, this.userKeywords);
                this.userKeywords = sorted_elem[0];
                this.checksKeywords = sorted_elem[1];
            }
        },

    },
    computed: {
        filteredKeywords() {
            return this.userKeywords.filter(row => {
                const Text = row.Text.toString().toLowerCase();
                const searchTerm = this.inputKeyword.toLowerCase();
                return (
                    Text.includes(searchTerm)
                );
            });
        },
        filteredTargets() {
            return this.userTargets.filter(row => {
                const Name = row.Name.toString().toLowerCase();
                const searchTerm = this.inputTarget.toLowerCase();
                return (
                    Name.includes(searchTerm)
                );
            });
        },
        checkAllKeywords() {
            if (this.checksKeywords.length == this.filteredKeywords.length) {
                return true;
            } else {
                return false;
            }
        },
        checkAllTargets() {
            if (this.checksTargets.length == this.filteredTargets.length) {
                return true;
            } else {
                return false;
            }
        },
    },
    template: `
        <div>
            <div class="row">
                <div class="column">
                    <div>
                        <h1>
                            My Keywords
                        </h1><br>
                        <p v-for="(message, index) in messagesKeywords">
                            <span v-html="message"></span>
                        </p>
                        <div class="input-row">
                            <vue-simple-suggest
                                v-model="inputKeyword"
                                :list="allKeywords"
                                :styles="autoCompleteStyle"
                                :destyled=true
                                :filter-by-query="true">
                                <label id="KeywordsFilter" for="Filter" class="mdc-text-field mdc-text-field--filled mdc-text-field--with-leading-icon">
                                    <i class="material-icons mdc-text-field__icon mdc-text-field__icon--leading" tabindex="0" role="button">text_fields</i>
                                    <input class="mdc-text-field__input" type="text" aria-labelledby="Filter">
                                    <span class="mdc-floating-label" id="KeywordsFilter">Filter or add a new keyword</span>
                                    <div class="mdc-line-ripple"></div>
                                </label>
                            </vue-simple-suggest>
                            <button
                                v-on:click="createKeyword"
                                class="material-icons mdc-top-app-bar__action-item mdc-icon-button" 
                                aria-label="Add">add
                            </button>
                        </div>
                        <div class="mdc-data-table scrollable">
                            <table class="mdc-data-table__table" aria-label="Created Keywords">
                                <thead>
                                    <tr class="mdc-data-table__header-row">
                                        <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                            <input
                                                v-if="checksTargets.length"
                                                type="checkbox" 
                                                v-model="checkAllKeywords"
                                                @click="selectAllKeywords"
                                                :disabled="macroPivot == 'keywords' && checksKeywords.length == 1"
                                            >
                                        </th>
                                        <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                            <a class="column-header" @click="sortRowsKeywords('CreatedDate')">
                                                CreatedDate
                                                <i v-if="sortedByKeywords === 'CreatedDate' && sortingKeywords['CreatedDate'] === true" class="material-icons column-sort">
                                                    keyboard_arrow_up
                                                </i>
                                                <i v-if="sortedByKeywords === 'CreatedDate' && sortingKeywords['CreatedDate'] === false" class="material-icons column-sort">
                                                    keyboard_arrow_down
                                                </i>
                                            </a>
                                        </th>
                                        <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                            <a class="column-header" @click="sortRowsKeywords('Text')">
                                                Keyword
                                                <i v-if="sortedByKeywords === 'Text' && sortingKeywords['Text'] === true" class="material-icons column-sort">
                                                    keyboard_arrow_up
                                                </i>
                                                <i v-if="sortedByKeywords === 'Text' && sortingKeywords['Text'] === false" class="material-icons column-sort">
                                                    keyboard_arrow_down
                                                </i>
                                            </a>
                                        </th>
                                    </tr>
                                </thead>
                                <tbody class="mdc-data-table__content">
                                    <tr v-for="(row, index) in filteredKeywords" class="mdc-data-table__row">
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
                <div class="column">
                    <div>
                        <h1>
                            My Targets
                        </h1><br>
                        <p v-for="(message, index) in messagesTargets">
                            <span v-html="message"></span>
                        </p>
                        <div class="input-row">
                            <vue-simple-suggest
                                v-model="inputTarget"
                                :list="allTargets"
                                :styles="autoCompleteStyle"
                                :destyled=true
                                :filter-by-query="true">
                                <label id="TargetsFilter" for="Filter" class="mdc-text-field mdc-text-field--filled mdc-text-field--with-leading-icon">
                                    <i class="material-icons mdc-text-field__icon mdc-text-field__icon--leading" tabindex="0" role="button">text_fields</i>
                                    <input class="mdc-text-field__input" type="text" aria-labelledby="Filter">
                                    <span class="mdc-floating-label" id="TargetsFilter">Filter or add a new target</span>
                                    <div class="mdc-line-ripple"></div>
                                </label>
                            </vue-simple-suggest>
                            <button
                                v-on:click="createTarget"
                                class="material-icons mdc-top-app-bar__action-item mdc-icon-button" 
                                aria-label="Add">add
                            </button>
                        </div>
                        <div class="mdc-data-table scrollable">
                            <table class="mdc-data-table__table" aria-label="Created Targets">
                                <thead>
                                    <tr class="mdc-data-table__header-row">
                                        <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                            <input
                                                v-if="checksKeywords.length"
                                                type="checkbox"
                                                v-model="checkAllTargets"
                                                @click="selectAllTargets"
                                                :disabled="macroPivot == 'targets' && checksTargets.length == 1"
                                            >
                                        </th>
                                        <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                            <a class="column-header" @click="sortRowsTargets('CreatedDate')">
                                                CreatedDate
                                                <i v-if="sortedByTargets === 'CreatedDate' && sortingTargets['CreatedDate'] === true" class="material-icons column-sort">
                                                    keyboard_arrow_up
                                                </i>
                                                <i v-if="sortedByTargets === 'CreatedDate' && sortingTargets['CreatedDate'] === false" class="material-icons column-sort">
                                                    keyboard_arrow_down
                                                </i>
                                            </a>
                                        </th>
                                        <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                            <a class="column-header" @click="sortRowsTargets('Name')">
                                                Target
                                                <i v-if="sortedByTargets === 'Name' && sortingTargets['Name'] === true" class="material-icons column-sort">
                                                    keyboard_arrow_up
                                                </i>
                                                <i v-if="sortedByTargets === 'Name' && sortingTargets['Name'] === false" class="material-icons column-sort">
                                                    keyboard_arrow_down
                                                </i>
                                            </a>
                                        </th>
                                    </tr>
                                </thead>
                                <tbody class="mdc-data-table__content">
                                    <tr v-for="(row, index) in filteredTargets" class="mdc-data-table__row">
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
            </div>
        </div>`,
};