export default {
    name: 'keywords',
    delimiters: ["[[","]]"],
    data: function () {
        return {
            selectedIndex: 1,
            keywords:  [],
            sortedBy: "CreatedDate",
            sorting: {
                CreatedDate: true,
                Name: false,
                CountResults: false,
                CountResultsPerc: false,
                CountResultsDay: false,
                CountTargets: false,
                TotalMatches: false,
                LastWeekMatches: false,
                AvgMatchesDay: false,
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
        this.fetchKeywords()
    },
    methods: {
        fetchKeywords: function() {
            this.$http.get('/keywords/analytic').then(function(response) {
                this.keywords = response.data.Keywords;
            }).catch(function(error) {
                console.log(error)
            });
        },
        sortRows: function(column) {
            this.sortedBy = column
            if (column == "CreatedDate" || column == "LastExtractionDate" ) {
                if (this.sorting[column]) {
                    this.keywords.sort((a,b) => (new Date(a[column]) - new Date(b[column])))
                    this.sorting[column] = false
                } else {
                    this.keywords.sort((b,a) => (new Date(a[column]) - new Date(b[column])))
                    this.sorting[column] = true
                }
            } else {
                if (this.sorting[column]) {
                    this.keywords.sort((a,b) => (a[column] > b[column]) ? 1 : ((b[column] > a[column]) ? -1 : 0))
                    this.sorting[column] = false
                } else {
                    this.keywords.sort((b,a) => (a[column] > b[column]) ? 1 : ((b[column] > a[column]) ? -1 : 0))
                    this.sorting[column] = true
                }
            }
        },
        csvExport: function() {
            this.$parent.csvExport(this.keywords, "keywords.csv")
        }
    },
    template: `
        <div>
            <div>
                <div class="initial-row">
                    <h1>
                        My keywords
                        <button
                            v-on:click="csvExport"
                            class="material-icons mdc-top-app-bar__action-item mdc-icon-button" 
                            aria-label="Export">arrow_downward
                        </button>
                    </h1>
                </div>
                <div class="mdc-data-table">
                    <table class="mdc-data-table__table" aria-label="Your Keywords">
                        <thead>
                            <tr class="mdc-data-table__header-row">
                                <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                    <a class="column-header" @click="sortRows('CreatedDate')">
                                        CreatedAt
                                        <i v-if="sortedBy === 'CreatedDate' && sortedBy['CreatedDate'] === true" class="material-icons column-sort">keyboard_arrow_up</i>
                                        <i v-if="sortedBy === 'CreatedDate' && sortedBy['CreatedDate'] === false" class="material-icons column-sort">keyboard_arrow_down</i>
                                    </a>
                                </th>
                                <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                    <a class="column-header" @click="sortRows('Name')">
                                        Keyword
                                        <i v-if="sortedBy === 'Name' && sortedBy['Name'] === true" class="material-icons column-sort">keyboard_arrow_up</i>
                                        <i v-if="sortedBy === 'Name' && sortedBy['Name'] === false" class="material-icons column-sort">keyboard_arrow_down</i>
                                    </a>
                                </th>
                                <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                    <a class="column-header" @click="sortRows('CountResults')">
                                        Results
                                        <i v-if="sortedBy === 'CountResults' && sortedBy['CountResults'] === true" class="material-icons column-sort">keyboard_arrow_up</i>
                                        <i v-if="sortedBy === 'CountResults' && sortedBy['CountResults'] === false" class="material-icons column-sort">keyboard_arrow_down</i>
                                    </a>
                                </th>
                                <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                    <a class="column-header" @click="sortRows('CountResultsPerc')">
                                        Results %
                                        <i v-if="sortedBy === 'CountResultsPerc' && sortedBy['CountResultsPerc'] === true" class="material-icons column-sort">keyboard_arrow_up</i>
                                        <i v-if="sortedBy === 'CountResultsPerc' && sortedBy['CountResultsPerc'] === false" class="material-icons column-sort">keyboard_arrow_down</i>
                                    </a>
                                </th>
                                <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                    <a class="column-header" @click="sortRows('CountResultsDay')">
                                        Results / Day
                                        <i v-if="sortedBy === 'CountResultsDay' && sortedBy['CountResultsDay'] === true" class="material-icons column-sort">keyboard_arrow_up</i>
                                        <i v-if="sortedBy === 'CountResultsDay' && sortedBy['CountResultsDay'] === false" class="material-icons column-sort">keyboard_arrow_down</i>
                                    </a>
                                </th>
                                <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                    <a class="column-header" @click="sortRows('CountTargets')">
                                        Targets
                                        <i v-if="sortedBy === 'CountTargets' && sortedBy['CountTargets'] === true" class="material-icons column-sort">keyboard_arrow_up</i>
                                        <i v-if="sortedBy === 'CountTargets' && sortedBy['CountTargets'] === false" class="material-icons column-sort">keyboard_arrow_down</i>
                                    </a>
                                </th>
                                <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                    <a class="column-header" @click="sortRows('TotalMatches')">
                                        Matches
                                        <i v-if="sortedBy === 'TotalMatches' && sortedBy['TotalMatches'] === true" class="material-icons column-sort">keyboard_arrow_up</i>
                                        <i v-if="sortedBy === 'TotalMatches' && sortedBy['TotalMatches'] === false" class="material-icons column-sort">keyboard_arrow_down</i>
                                    </a>
                                </th>
                                <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                    <a class="column-header" @click="sortRows('LastWeekMatches')">
                                        7 Days Matches
                                        <i v-if="sortedBy === 'LastWeekMatches' && sortedBy['LastWeekMatches'] === true" class="material-icons column-sort">keyboard_arrow_up</i>
                                        <i v-if="sortedBy === 'LastWeekMatches' && sortedBy['LastWeekMatches'] === false" class="material-icons column-sort">keyboard_arrow_down</i>
                                    </a>
                                </th>
                                <th class="mdc-data-table__header-cell" role="columnheader" scope="col">
                                    <a class="column-header" @click="sortRows('LastWeekMatches')">
                                        Matches / Day
                                        <i v-if="sortedBy === 'AvgMatchesDay' && sortedBy['AvgMatchesDay'] === true" class="material-icons column-sort">keyboard_arrow_up</i>
                                        <i v-if="sortedBy === 'AvgMatchesDay' && sortedBy['AvgMatchesDay'] === false" class="material-icons column-sort">keyboard_arrow_down</i>
                                    </a>
                                </th>
                            </tr>
                        </thead>
                        <tbody class="mdc-data-table__content">
                            <tr v-for="(row, index) in keywords" class="mdc-data-table__row">
                                <td class="mdc-data-table__cell" v-html="row.CreatedDate"></td>
                                <td class="mdc-data-table__cell" v-html="row.Name"></td>
                                <td class="mdc-data-table__cell" v-html="row.CountResults"></td>
                                <td class="mdc-data-table__cell" v-html="row.CountResultsPerc"></td>
                                <td class="mdc-data-table__cell" v-html="row.CountResultsDay"></td>
                                <td class="mdc-data-table__cell" v-html="row.CountTargets"></td>
                                <td class="mdc-data-table__cell" v-html="row.TotalMatches"></td>
                                <td class="mdc-data-table__cell" v-html="row.LastWeekMatches"></td>
                                <td class="mdc-data-table__cell" v-html="row.AvgMatchesDay"></td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>`,
};