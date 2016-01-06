import {Component} from 'angular2/core';
import {NgForm} from 'angular2/common';
import {AuthService} from './auth.service';
import {Credentials} from './credentials';

@Component({
    selector: 'gohst-login',
    templateUrl: 'app/login.component.html',
    styleUrls: ['app/login.component.css']
})
export class LoginComponent {
    model = new Credentials("","");

    constructor(public authService :AuthService) { }

    submitLogin() {
        this.authService.authenticate(this.model.username, this.model.password)
            .then((success :any) => {
                if (!success) {
                    alert("Username and password do not match");
                }
            });
    }
}
