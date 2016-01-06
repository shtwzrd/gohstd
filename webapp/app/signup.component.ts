import {Component} from 'angular2/core';
import {NgForm} from 'angular2/common';
import {AuthService} from './auth.service';
import {Registration} from './registration';

@Component({
    selector: 'gohst-signup',
    templateUrl: 'app/signup.component.html',
    styleUrls: ['app/signup.component.css']
})
export class SignupComponent {
    model = new Registration("","","");

    constructor(public authService :AuthService) { }

    submitRegistration() {
        window.fetch('/api/users/register', {
            method: 'POST',
            body: JSON.stringify(this.model)
        }).then((response :any) => {
                if (response.status == 201) {
                    this.authService.authenticate(this.model.username, this.model.password);
                } else {
                    alert(response.body);
                }
        });
    }
}
