export default {
    name: 'targets',
    delimiters: ["[[","]]"],
    data: function () {
        return {
            selectedIndex: 1,
        }
    },
    mounted() {
        this.$parent.selectedIndex = this.selectedIndex;
    },
    template: `
        <div>
            <h1>
                Your Targets
            </h1>
        </div>`,
};