export default {
    name: 'Targets',
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