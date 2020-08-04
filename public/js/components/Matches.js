export default {
    name: 'matches',
    delimiters: ["[[","]]"],
    data: function () {
        return {
            matches: [],
            selectedIndex: 3,
            filter: '',
            sortedBy: "CreatedDate",
            sorting: {
                CreatedDate: true,
                Target: false,
                Title: false,
                KeywordText: false
            }
        }
    },
    mounted() {
        this.$parent.selectedIndex = this.selectedIndex
        const textFilter= mdc.textField.MDCTextField.attachTo(document.getElementById("Filter"));
        let styleElem = document.createElement('style');
        styleElem.textContent = `
            .column-sort {
                font-size: 16px;
                vertical-align: -3px;
            }`
        document.head.appendChild(styleElem);
    },
    created () {
        this.fetchMatches()
    },
    methods: {
        fetchMatches: function() {
            this.$http.get('/matches').then(function(response) {
                this.matches = response.data.Data;
            }).catch(function(error) {
                console.log(error)
            });
        },
        select: function(raw) {
            window.open(raw.Url, "_blank");
        },
        sortRows: function(column) {
            this.sortedBy = column
            if (column == "CreatedDate") {
                if (this.sorting[column]) {
                    this.matches.sort((a,b) => (new Date(a[column]) - new Date(b[column])))
                    this.sorting[column] = false
                } else {
                    this.matches.sort((b,a) => (new Date(a[column]) - new Date(b[column])))
                    this.sorting[column] = true
                }
            } else {
                if (this.sorting[column]) {
                    this.matches.sort((a,b) => (a[column] > b[column]) ? 1 : ((b[column] > a[column]) ? -1 : 0))
                    this.sorting[column] = false
                } else {
                    this.matches.sort((b,a) => (a[column] > b[column]) ? 1 : ((b[column] > a[column]) ? -1 : 0))
                    this.sorting[column] = true
                }
            }
        }
    },
    computed: {
        filteredRows() {
            return this.matches.filter(row => {
                const CreatedDate = row.CreatedDate.toString().toLowerCase();
                const Target = row.Target.toString().toLowerCase();
                const Title = row.Title.toString().toLowerCase();
                const searchTerm = this.filter.toLowerCase();
                return (
                    CreatedDate.includes(searchTerm) || Target.includes(searchTerm) || Title.includes(searchTerm)
                );
            });
        }
    },
    template: `
        <div>
            <h1>
                Your matches user
            </h1>
            <div>
                <label id="Filter" for="Filter" class="mdc-text-field mdc-text-field--filled mdc-text-field--with-leading-icon">
                    <i class="material-icons mdc-text-field__icon mdc-text-field__icon--leading" tabindex="0" role="button">text_fields</i>
                    <input class="mdc-text-field__input" type="text" aria-labelledby="Filter" v-model="filter">
                    <span class="mdc-floating-label" id="Filter">Filter by any field</span>
                    <div class="mdc-line-ripple"></div>
                </label>
            </div>
            <table class="mdc-data-table__table" aria-label="Created Keywords">
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
                            <a class="column-header" @click="sortRows('KeywordText')">
                                Keyword
                                <i v-if="sortedBy === 'KeywordText' && sorting['KeywordText'] === true" class="material-icons column-sort">keyboard_arrow_up</i>
                                <i v-if="sortedBy === 'KeywordText' && sorting['KeywordText'] === false" class="material-icons column-sort">keyboard_arrow_down</i>
                            </a>
                        </th>
                        <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                            <a class="column-header" @click="sortRows('Target')">
                                Target
                                <i v-if="sortedBy === 'Target' && sorting['Target'] === true" class="material-icons column-sort">keyboard_arrow_up</i>
                                <i v-if="sortedBy === 'Target' && sorting['Target'] === false" class="material-icons column-sort">keyboard_arrow_down</i>
                            </a>
                        </th>
                        <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                            <a class="column-header" @click="sortRows('Title')">
                                Job Title
                                <i v-if="sortedBy === 'Title' && sorting['Title'] === true" class="material-icons column-sort">keyboard_arrow_up</i>
                                <i v-if="sortedBy === 'Title' && sorting['Title'] === false" class="material-icons column-sort">keyboard_arrow_down</i>
                            </a>
                        </th>
                        <th class="mdc-data-table__header-cell" role="columnheader" scope="col">Job Url</th>
                    </tr>
                </thead>
                <tbody class="mdc-data-table__content">
                    <tr v-for="(row, index) in filteredRows" class="mdc-data-table__row">
                        <td class="mdc-data-table__cell" v-html="row.CreatedDate"></td>
                        <td class="mdc-data-table__cell" v-html="row.KeywordText"></td>
                        <td class="mdc-data-table__cell" v-html="row.Target"></td>
                        <td class="mdc-data-table__cell" v-html="row.Title"></td>
                        <td class="mdc-data-table__cell">
                            <button 
                                class="material-icons mdc-top-app-bar__action-item mdc-icon-button" 
                                aria-label="Open"
                                target="_blank" 
                                rel="noopener"
                                v-on:click="select(row)">open_in_new
                            </button>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>`,
};