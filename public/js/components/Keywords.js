export default {
    name: 'keywords',
    delimiters: ["[[","]]"],
    data: function () {
        return {
            userKeywords: [],
            filterUserKeywords: '',
            allkeywords: [],
            targets: [],
            chosen: '',
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
    created () {
        this.fetchUserKeywords();
        this.fetchAllKeywords();
        this.fetchTargets();
    },
    methods: {
        fetchUserKeywords: function() {
            this.$http.get('/keywords/user').then(function(response) {
                this.userKeywords = response.data.Keywords
            }).catch(function(error) {
                console.log(error)
            });
        },
        fetchAllKeywords: function() {
            this.$http.get('/keywords/all').then(function(response) {
                this.allkeywords = response.data.Keywords
            }).catch(function(error) {
                console.log(error)
            });
        },
        fetchTargets: function() {
            this.$http.get('/targets').then(function(response) {
                this.targets = response.data.Targets
            }).catch(function(error) {
                console.log(error)
            });
        },
    },
    computed: {
        applyFilterUserKeywords() {
            return this.userKeywords.filter(row => {
                const Text = row.Text.toString().toLowerCase();
                const searchTerm = this.filterUserKeywords.toLowerCase();
                return (
                    Text.includes(searchTerm)
                );
            });
        }
    },
    template: `
        <div>
            <div class="row">
                <div class="column">
                    <div>
                        <h1>
                            My Keywords
                        </h1><br>
                        <div class="input-row">
                            <vue-simple-suggest
                                v-model="chosen"
                                :list="allkeywords"
                                :styles="autoCompleteStyle"
                                :destyled=true
                                :filter-by-query="true">
                                <label id="KeywordsFilter" for="Filter" class="mdc-text-field mdc-text-field--filled mdc-text-field--with-leading-icon">
                                    <i class="material-icons mdc-text-field__icon mdc-text-field__icon--leading" tabindex="0" role="button">text_fields</i>
                                    <input class="mdc-text-field__input" type="text" aria-labelledby="Filter" v-model="filterUserKeywords">
                                    <span class="mdc-floating-label" id="KeywordsFilter">Filter or add a new keyword</span>
                                    <div class="mdc-line-ripple"></div>
                                </label>
                            </vue-simple-suggest>
                            <button 
                                class="material-icons mdc-top-app-bar__action-item mdc-icon-button" 
                                aria-label="Add">add
                            </button>
                        </div>
                        <div class="mdc-data-table scrollable">
                            <table class="mdc-data-table__table" aria-label="Created Keywords">
                                <thead>
                                    <tr class="mdc-data-table__header-row">
                                        <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                            <input type="checkbox">
                                        </th>
                                        <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                            <a class="column-header">
                                                CreatedDate
                                            </a>
                                        </th>
                                        <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                            <a class="column-header">
                                                Keyword
                                            </a>
                                        </th>
                                    </tr>
                                </thead>
                                <tbody class="mdc-data-table__content">
                                    <tr v-for="(row, index) in applyFilterUserKeywords" class="mdc-data-table__row">
                                        <td class="mdc-data-table__cell">
                                            <input type="checkbox">
                                        </td>
                                        <td class="mdc-data-table__cell" v-html="row.CreatedDate"></td>
                                        <td class="mdc-data-table__cell" v-html="row.Text"></td>
                                        <td class="mdc-data-table__cell">
                                            <button 
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
                        <div class="input-row">
                            <vue-simple-suggest
                                v-model="chosen"
                                :list="allkeywords"
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
                                class="material-icons mdc-top-app-bar__action-item mdc-icon-button" 
                                aria-label="Add">add
                            </button>
                        </div>
                        <div class="mdc-data-table scrollable">
                            <table class="mdc-data-table__table" aria-label="Created Targets">
                                <thead>
                                    <tr class="mdc-data-table__header-row">
                                        <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                            <input type="checkbox">
                                        </th>
                                        <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                            <a class="column-header">
                                                CreatedDate
                                            </a>
                                        </th>
                                        <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                            <a class="column-header">
                                                Target
                                            </a>
                                        </th>
                                    </tr>
                                </thead>
                                <tbody class="mdc-data-table__content">
                                    <tr v-for="(row, index) in targets" class="mdc-data-table__row">
                                        <td class="mdc-data-table__cell">
                                            <input type="checkbox">
                                        </td>
                                        <td class="mdc-data-table__cell" v-html="row.CreatedDate"></td>
                                        <td class="mdc-data-table__cell" v-html="row.Name"></td>
                                        <td class="mdc-data-table__cell">
                                            <button 
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