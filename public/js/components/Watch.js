export default {
    name: 'watch',
    delimiters: ["[[","]]"],
    data: function () {
        return {
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
        }
    },
    mounted() {
        this.$parent.selectedIndex = this.selectedIndex;

        let styleElem = document.createElement('style');
        styleElem.textContent = `
            .row {
                display: flex;
            }
            .input-row {
                display: flex;
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
            table.tableBodyScroll tbody {
              display: block;
              max-height: 80vh;
              overflow-y: scroll;
            }
            table.tableBodyScroll thead, table.tableBodyScroll tbody tr {
              display: table;
              table-layout: fixed;
              width: 98%;
              text-align: center;`
        document.head.appendChild(styleElem);
    },
    created: async function() {
        this.fetchResults();
    },
    methods: {
        fetchResults: function() {
            this.$http.get('/jobs/favourite').then(function(response) {
                if (response.data.Jobs !== null) {
                    this.jobs = response.data.Jobs;
                }
            }).catch(function(error) {
                console.log(error);
            });
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
    template: `
        <div>
            <div class="row">
                <table class="mdc-data-table__table tableBodyScroll" aria-label="Results" style="">
                    <thead>
                        <tr class="mdc-data-table__header-row">
                            <th class="mdc-data-table__header-cell" role="columnheader" scope="col" style="width:9%;white-space:nowrap;">Job Url</th>
                            <th class="mdc-data-table__header-cell" role="columnheader" scope="col" style="width:8%;white-space:nowrap;">Save</th>
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
                                <a class="column-header" @click="sortRowsJobs('ClosedDate')">
                                    ClosedAt
                                    <i v-if="sortingJobs['ClosedDate'] === true" class="material-icons column-sort" v-bind:class="{'arrowBold': sortedByJobs == 'ClosedDate'}">
                                        keyboard_arrow_up
                                    </i>
                                    <i v-if="sortingJobs['ClosedDate'] === false" class="material-icons column-sort" v-bind:class="{'arrowBold': sortedByJobs == 'ClosedDate'}">
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
                            <th class="mdc-data-table__header-cell" role="columnheader" scope="col" style="width:14%;white-space:nowrap;">
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
                            <th class="mdc-data-table__header-cell" role="columnheader" scope="col" style="width:23%;white-space:nowrap;">
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
                        <tr v-for="(row, index) in jobs" class="mdc-data-table__row">
                            <td class="mdc-data-table__cell" style="width:9%;white-space:nowrap;">
                                <button 
                                    class="material-icons mdc-top-app-bar__action-item mdc-icon-button" 
                                    aria-label="Open"
                                    target="_blank" 
                                    rel="noopener"
                                    v-on:click="select(row)">open_in_new
                                </button>
                            </td>
                            <td class="mdc-data-table__cell" style="width:8%;white-space:nowrap;">
                                <input
                                    type="checkbox"
                                    v-model="row.IsSaved"
                                    @change="save(row,$event)"
                                >
                            </td>
                            <td class="mdc-data-table__cell" v-html="row.CreatedDate" style="width:11%;white-space:nowrap;"></td>
                            <td class="mdc-data-table__cell" v-html="row.ClosedDate" style="width:11%;white-space:nowrap;"></td>
                            <td class="mdc-data-table__cell" v-html="row.KeywordText" style="width:11%;white-space:nowrap;"></td>
                            <td class="mdc-data-table__cell comment" v-html="row.TargetName" style="width:13%;white-space:nowrap;"></td>
                            <td class="mdc-data-table__cell comment" v-html="row.Location" style="width:14%;white-space:nowrap;" v-bind:title="row.Location"></td>
                            <td class="mdc-data-table__cell" v-html="row.Title" style="width:23%;white-space:nowrap;" v-bind:title="row.Title"></td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>`,
};