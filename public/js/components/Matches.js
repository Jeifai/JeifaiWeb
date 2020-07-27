export default {
    name: 'matches',
    delimiters: ["[[","]]"],
    data: function () {
        return {
            matches: [],
            selectedIndex: 3,
            filter: '',
        }
    },
    mounted() {
        this.$parent.selectedIndex = this.selectedIndex
        const textFilter= mdc.textField.MDCTextField.attachTo(document.getElementById("Filter"));
    },
    created () {
        this.fetchMatches()
    },
    methods: {
        fetchMatches: function() {
            this.$http.get('/testMatch').then(function(response) {
                this.matches = response.data.Data;
            }).catch(function(error) {
                console.log(error)
            });
        },
        select: function(raw) {
            window.open(raw.Url, "_blank");
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
                        <th class="mdc-data-table__header-cell" role="columnheader" scope="col">CreatedAt</th>
                        <th class="mdc-data-table__header-cell" role="columnheader" scope="col">Target</th>
                        <th class="mdc-data-table__header-cell" role="columnheader" scope="col">Job Title</th>
                        <th class="mdc-data-table__header-cell" role="columnheader" scope="col">Job Url</th>
                    </tr>
                </thead>
                <tbody class="mdc-data-table__content">
                    <tr v-for="(row, index) in filteredRows" class="mdc-data-table__row">
                        <td class="mdc-data-table__cell" v-html="row.CreatedDate"></td>
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