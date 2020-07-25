export default {
    name: 'keywordsT',
    delimiters: ["[[","]]"],
    data: function () {
        return {
            selectedIndex: 2,
        }
    },
    mounted() {
        this.$parent.selectedIndex = this.selectedIndex;
    },
    template: `
        <div>
            <h1>
                Your Keywords
            </h1>
        </div>`,
};