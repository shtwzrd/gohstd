import {Component} from 'angular2/core';
import {Http, Response, HTTP_PROVIDERS} from 'angular2/http';
import {NgForm} from 'angular2/common';
import {AuthService} from './auth.service';
import {Registration} from './registration';

@Component({
    selector: 'gohst-signup',
    templateUrl: 'app/signup.component.html',
    styleUrls: ['app/signup.component.css'],
    providers: [HTTP_PROVIDERS]
})
export class SignupComponent {
    model = new Registration("","","");

    constructor(public authService :AuthService, private http :Http) { }

    submitRegistration() {
        this.http.post('/api/users/register', JSON.stringify(this.model))
            .map(res => res.status)
            .subscribe(
                data => {
                    if (data == 200 || data == 201) {
                        this.authService.authenticate(
                            this.model.username, this.model.password);
                    }
                },
                err => alert(err)
            );
    }
}
