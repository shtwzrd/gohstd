import {Component, OnInit} from 'angular2/core';
import {NgIf} from 'angular2/common';
import {Router} from 'angular2/router';
import {Hero} from './hero';
import {LoginComponent} from './login.component';
import {HeroService} from './hero.service';
import {AuthService} from './auth.service';

@Component({
    selector: 'gohst-about',
    templateUrl: 'app/about.component.html',
    styleUrls: ['app/about.component.css'],
    directives: [LoginComponent, NgIf]
})
export class AboutComponent implements OnInit {
    public heroes: Hero[] = [];

    constructor(private _heroService: HeroService, private _router: Router, private _authService: AuthService) { }

    ngOnInit() {
        this._heroService.getHeroes().then(heroes => this.heroes = heroes.slice(1,5));
    }

    gotoDetail(hero: Hero) {
        this._router.navigate(['HeroDetail', { id: hero.id }]);
    }
}
