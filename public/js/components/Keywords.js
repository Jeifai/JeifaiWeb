export default {
    name: 'keywords',
    delimiters: ["[[","]]"],
    data: function () {
        return {
            keywords: [],
            targets: [],
        }
    },
    mounted() {
        let styleElem = document.createElement('style');
        styleElem.textContent = `
            .row {
                display: flex;
            }
            .column {
                flex: 50%;
            }
            .scrollable {
                overflow-y: scroll;
                height:36vh;
            }`
        document.head.appendChild(styleElem);
    },
    created () {
        this.fetchKeywords();
        this.fetchTargets();
    },
    methods: {
        fetchKeywords: function() {
            this.$http.get('/keywords').then(function(response) {
                this.keywords = response.data.Keywords
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
    template: `
        <div>
            <div class="row">
                <div class="column">
                    <div>
                        <h1>
                            My Keywords
                        </h1><br>
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
                                    <tr v-for="(row, index) in keywords" class="mdc-data-table__row">
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