export default {
    name: 'profile',
    delimiters: ["[[","]]"],
    data: function () {
        return {
            messages: '',
            fields: []
        }
    },
    mounted() {
        const textUserName = mdc.textField.MDCTextField.attachTo(document.getElementById("UserName"));
        const textFirstName = mdc.textField.MDCTextField.attachTo(document.getElementById("FirstName"));
        const textLastName = mdc.textField.MDCTextField.attachTo(document.getElementById("LastName"));
        const textEmail = mdc.textField.MDCTextField.attachTo(document.getElementById("Email"));
        const textDateOfBirth = mdc.textField.MDCTextField.attachTo(document.getElementById("DateOfBirth"));
        const textCountry = mdc.textField.MDCTextField.attachTo(document.getElementById("Country"));
        const textCity = mdc.textField.MDCTextField.attachTo(document.getElementById("City"));
        const textNewPassword = mdc.textField.MDCTextField.attachTo(document.getElementById("NewPassword"));
        const textRepeatNewPassword = mdc.textField.MDCTextField.attachTo(document.getElementById("RepeatNewPassword"));
        const textCurrentPassword = mdc.textField.MDCTextField.attachTo(document.getElementById("CurrentPassword"));
        const button = mdc.ripple.MDCRipple.attachTo(document.querySelector('.mdc-fab'));
        const select = mdc.select.MDCSelect.attachTo(document.querySelector('.mdc-select'));

        select.listen('MDCSelect:change', () => {
            this.setSelect(select.value);
        });
    },
    created () {
        this.fetchProfile()
    },
    methods: {
        fetchProfile: function() {
            this.$http.get('/profile').then(function(response) {
                this.fields = response.data.User;
                const select = mdc.select.MDCSelect.attachTo(document.querySelector('.mdc-select'));
                select.value = response.data.User.Gender;

                // Workaround https://github.com/material-components/material-components-web/issues/6290
                const edtTexts = [].map.call(document.querySelectorAll('.mdc-floating-label'), function(el) {
                    if (response.data.User[el.id] !== "") {
                        el.classList.add('mdc-floating-label--float-above');
                    }
                });
            }).catch(function(error) {
                console.log(error)
            });
        },
        updateProfile: function() {
            this.messages = '';
            this.$http.put('/profile', this.fields).then(function(response) {
                this.messages = response.data.Messages
            }).catch(function(error) {
                console.log(error)
            });
        },
        setSelect: function (value) {
            this.fields.Gender = value;
        },
    },
    template: `
        <div>
            <p v-for="(message, index) in messages">
                <span v-html="message"></span>
            </p>
            <label id="UserName" for="UserName" class="mdc-text-field mdc-text-field--filled mdc-text-field--with-leading-icon">
                <i class="material-icons mdc-text-field__icon mdc-text-field__icon--leading" tabindex="0" role="button">person</i>
                <input class="mdc-text-field__input" type="text" aria-labelledby="UserName" v-model="fields.UserName">
                <span class="mdc-floating-label" id="UserName">UserName</span>
                <div class="mdc-line-ripple"></div>
            </label>
            <label id="FirstName" for="FirstName" class="mdc-text-field mdc-text-field--filled mdc-text-field--with-leading-icon">
                <i class="material-icons mdc-text-field__icon mdc-text-field__icon--leading" tabindex="1" role="button">accessibility</i>
                <input class="mdc-text-field__input" type="text" aria-labelledby="FirstName" v-model="fields.FirstName">
                <span class="mdc-floating-label" id="FirstName">FirstName</span>
                <div class="mdc-line-ripple"></div>
            </label>
            <p></p>
            <label id="LastName" for="LastName" class="mdc-text-field mdc-text-field--filled mdc-text-field--with-leading-icon">
                <i class="material-icons mdc-text-field__icon mdc-text-field__icon--leading" tabindex="2" role="button">accessibility</i>
                <input class="mdc-text-field__input" type="text" aria-labelledby="LastName" v-model="fields.LastName">
                <span class="mdc-floating-label" id="LastName">LastName</span>
                <div class="mdc-line-ripple"></div>
            </label>
            <label id="Email" for="Email" class="mdc-text-field mdc-text-field--filled mdc-text-field--with-leading-icon">
                <i class="material-icons mdc-text-field__icon mdc-text-field__icon--leading" tabindex="3" role="button">alternate_email</i>
                <input class="mdc-text-field__input" type="text" aria-labelledby="Email" v-model="fields.Email">
                <span class="mdc-floating-label" id="Email">Email</span>
                <div class="mdc-line-ripple"></div>
            </label>
            <p></p>
            <label id="DateOfBirth" for="DateOfBirth" class="mdc-text-field mdc-text-field--filled mdc-text-field--with-leading-icon">
                <i class="material-icons mdc-text-field__icon mdc-text-field__icon--leading" tabindex="4" role="button">date_range</i>
                <input class="mdc-text-field__input" type="date" aria-labelledby="DateOfBirth" v-model="fields.DateOfBirth">
                <span class="mdc-floating-label" id="DateOfBirth">DateOfBirth</span>
                <div class="mdc-line-ripple"></div>
            </label>
            <label id="Country" for="Country" class="mdc-text-field mdc-text-field--filled mdc-text-field--with-leading-icon">
                <i class="material-icons mdc-text-field__icon mdc-text-field__icon--leading" tabindex="3" role="button">layers</i>
                <input class="mdc-text-field__input" type="text" aria-labelledby="Country" v-model="fields.Country">
                <span class="mdc-floating-label" id="Country">Country</span>
                <div class="mdc-line-ripple"></div>
            </label>
            <p></p>
            <label id="City" for="City" class="mdc-text-field mdc-text-field--filled mdc-text-field--with-leading-icon">
                <i class="material-icons mdc-text-field__icon mdc-text-field__icon--leading" tabindex="4" role="button">place</i>
                <input class="mdc-text-field__input" type="text" aria-labelledby="City" v-model="fields.City">
                <span class="mdc-floating-label" id="City">City</span>
                <div class="mdc-line-ripple"></div>
            </label>
            <div class="mdc-select mdc-select--filled mdc-select--with-leading-icon">
                <div class="mdc-select__anchor">
                    <i class="material-icons mdc-select__icon mdc-select__icon--leading" tabindex="4" role="button">adb</i>
                    <div class="mdc-select__selected-text"></div>
                    <span class="mdc-floating-label" id="Gender">Gender</span>
                    <div class="mdc-line-ripple"></div>
                </div>
                <div class="mdc-select__menu mdc-menu mdc-menu-surface">
                    <ul class="mdc-list">
                        <li class="mdc-list-item mdc-list-item--selected" data-value="" aria-selected="true"></li>
                        <li class="mdc-list-item" data-value="female">
                            <span class="mdc-list-item__text">
                                Female
                            </span>
                        </li>
                        <li class="mdc-list-item" data-value="male">
                            <span class="mdc-list-item__text">
                                Male
                            </span>
                        </li>
                    </ul>
                </div>
            </div>
            <p></p>
            <label id="NewPassword" for="NewPassword" class="mdc-text-field mdc-text-field--filled mdc-text-field--with-leading-icon">
                <i class="material-icons mdc-text-field__icon mdc-text-field__icon--leading" tabindex="3" role="button">security</i>
                <input class="mdc-text-field__input" type="password" aria-labelledby="NewPassword" v-model="fields.NewPassword">
                <span class="mdc-floating-label" id="NewPassword">NewPassword</span>
                <div class="mdc-line-ripple"></div>
            </label>
            <label id="RepeatNewPassword" for="RepeatNewPassword" class="mdc-text-field mdc-text-field--filled mdc-text-field--with-leading-icon">
                <i class="material-icons mdc-text-field__icon mdc-text-field__icon--leading" tabindex="3" role="button">security</i>
                <input class="mdc-text-field__input" type="password" aria-labelledby="RepeatNewPassword" v-model="fields.RepeatNewPassword">
                <span class="mdc-floating-label" id="RepeatNewPassword">RepeatNewPassword</span>
                <div class="mdc-line-ripple"></div>
            </label>
            <p></p>
            <label id="CurrentPassword" for="CurrentPassword" class="mdc-text-field mdc-text-field--filled mdc-text-field--with-leading-icon">
                <i class="material-icons mdc-text-field__icon mdc-text-field__icon--leading" tabindex="3" role="button">security</i>
                <input class="mdc-text-field__input" type="password" aria-labelledby="CurrentPassword" v-model="fields.CurrentPassword">
                <span class="mdc-floating-label" id="CurrentPassword">CurrentPassword</span>
                <div class="mdc-line-ripple"></div>
            </label>
            <button 
                class="mdc-fab app-fab--absolute"
                aria-label="Check"
                v-on:click="updateProfile">
                <div class="mdc-fab__ripple"></div>
                <span class="mdc-fab__icon material-icons">check</span>
            </button>
        </div>`,
};