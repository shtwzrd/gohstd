import {Component} from 'angular2/core';
import {AuthService} from './auth.service'

@Component({
    selector: 'gohst-app',
    templateUrl: 'gohst/app.component.html',
    providers: [AuthService]
})

export class AppComponent { };
