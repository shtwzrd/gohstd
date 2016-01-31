import {Component} from 'angular2/core';
import {RouteConfig, ROUTER_DIRECTIVES, ROUTER_PROVIDERS} from 'angular2/router';
import {AboutComponent} from './about.component';
import {LoginComponent} from './login.component';
import {SignupComponent} from './signup.component';
import {AuthService} from './auth.service';

@Component({
  selector: 'gohst-app',
  template: `
    <h1>{{title}}</h1>
    <nav>
      <a [routerLink]="['/About']">About</a>
    </nav>
        <gohst-login>
        </gohst-login>
        <gohst-signup>
        </gohst-signup>
    <router-outlet></router-outlet>
  `,
  styleUrls: ['app/app.component.css'],
  directives: [ROUTER_DIRECTIVES, LoginComponent, SignupComponent],
  providers: [AuthService]
})
@RouteConfig([
  {path: '/about', as: 'About', component: AboutComponent, useAsDefault: true}
])
export class AppComponent {
  public title = 'gohst';
}
