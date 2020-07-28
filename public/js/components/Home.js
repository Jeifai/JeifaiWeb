export default {
    name: 'home',
    delimiters: ["[[","]]"],
    data: function () {
        return {
            selectedIndex: 0,
            homeInfo: [],
        }
    },
    mounted() {
        this.$parent.selectedIndex = this.selectedIndex;
    },
    created () {
        this.fetchHome()
    },
    methods: {
        fetchHome: function() {
            this.$http.get('/getHome').then(function(response) {
                this.homeInfo = response.data.Home;
            }).catch(function(error) {
                console.log(error)
            });
        }
    },
    template: `
        <div>
            <div class="mdc-typography">
                <h2 class="mdc-typography--headline3">Hi [[homeInfo.UserName]]~~~</h2>
                <span class="mdc-typography--body1">
                    <b>[[homeInfo.CountTargets]]</b> targets are being monitored by <b>[[homeInfo.CountKeywords]]</b> keywords so far.<br>
                    In total you have created <b>[[homeInfo.CountTargetsKeywords]]</b> combinations between targets and keywords.<br>
                    Those triggers have been able to generate <b>[[homeInfo.CountMatchesLast7Days]]</b> matches in the last 7 days.<br><br>
                    Currently in total your targets have <b>[[homeInfo.CountOpenPositions]]</b> open job positions on their career pages.<br>
                    <b>[[homeInfo.CountResultsLast7Days]]</b> have been published within 7 days, 
                    while <b>[[homeInfo.CountClosedPositionLast7Days]]</b> have been closed in the last 7 days.
                </span><br><br>
                <span class="mdc-typography--body2">
                    If you want to see other details on this homepage, just let us know :)
                </span>
            </div>
        </div>`,
};