import {Component} from 'angular2/core';
import {RouteConfig, ROUTER_DIRECTIVES, ROUTER_PROVIDERS} from 'angular2/router';
import {HeroesComponent} from './heroes.component';
import {HeroDetailComponent} from './hero-detail.component';
import {AboutComponent} from './about.component';
import {LoginComponent} from './login.component';
import {SignupComponent} from './signup.component';
import {HeroService} from './hero.service';
import {AuthService} from './auth.service';

@Component({
  selector: 'gohst-app',
  template: `
    <h1>{{title}}</h1>
    <nav>
      <a [routerLink]="['/About']">About</a>
      <a [routerLink]="['/History']">History</a>
      <a [routerLink]="['/Posts']">Posts</a>
    </nav>
        <gohst-login>
        </gohst-login>
        <gohst-signup>
        </gohst-signup>
    <router-outlet></router-outlet>
  `,
  styleUrls: ['app/app.component.css'],
  directives: [ROUTER_DIRECTIVES, LoginComponent, SignupComponent],
  providers: [AuthService, HeroService]
})
@RouteConfig([
  {path: '/about', as: 'About', component: AboutComponent, useAsDefault: true},
  {path: '/history', as: 'History', component: HeroesComponent},
  {path: '/posts', as: 'Posts', component: HeroDetailComponent}
])
export class AppComponent {
  public title = 'gohst';
}
