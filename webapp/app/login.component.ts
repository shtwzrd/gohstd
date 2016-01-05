import {Component} from 'angular2/core';
import {NgForm} from 'angular2/common';
import {AuthService} from './auth.service';
import {Credentials} from './credentials';

@Component({
    selector: 'gohst-login',
    templateUrl: 'app/login.component.html',
    styleUrls: ['app/login.component.css'],
    providers: [AuthService]
})
export class LoginComponent {
    loggedIn :boolean;
    model = new Credentials("","");

    constructor(private authService :AuthService) {
        this.loggedIn = false;
    }

    submitLogin() {
        this.authService.authenticate(this.model.username, this.model.password)
            .then((success :any) => {
                if (!success) {
                    this.loggedIn = false;
                    alert("Username and password do not match");
                } else {
                    this.loggedIn = true;
                }
            });
    }
}
