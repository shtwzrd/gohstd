import {Component} from 'angular2/core';
import {AuthService} from './auth.service';

@Component({
    selector: 'gohst-login',
    templateUrl: 'app/login.component.html',
    styleUrls: ['app/login.component.css'],
    providers: [AuthService]
})

export class LoginComponent {

    constructor(private authService :AuthService) {}

    submitLogin(event :any, username :string, password :string) {
        var success = this.authService.authenticate(username, password);
        if (!success) {
            alert("Username and password do not match");
        }
        console.log(this.authService.authenticated);
    }
}
