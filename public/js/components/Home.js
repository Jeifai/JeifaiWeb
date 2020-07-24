export default {
    name: 'home',
    delimiters: ["[[","]]"],
    data: function () {
        return {
            selectedIndex: 0,
        }
    },
    mounted() {
        this.$parent.selectedIndex = this.selectedIndex;
    },
    template: `
        <div>
            <h1>
                Your Home
            </h1>
        </div>`,
};