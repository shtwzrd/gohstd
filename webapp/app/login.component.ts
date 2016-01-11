import {Component, OnInit} from 'angular2/core';
import {NgForm} from 'angular2/common';
import {AuthService} from './auth.service';
import {Credentials} from './credentials';

@Component({
    selector: 'gohst-login',
    templateUrl: 'app/login.component.html',
    styleUrls: ['app/login.component.css']
})
export class LoginComponent implements OnInit {
    model = new Credentials("","");

    constructor(public authService :AuthService) { }

    submitLogin() {
        this.authService.authenticate(this.model.username, this.model.password)
		.then((success :any) => {
			if (!success) {
				alert("Username and password do not match");
			}
			else {
					if ((<HTMLInputElement>document.getElementById('rmb')).checked){
						localStorage.setItem('gohstusername', this.model.username);
					}
					else {
						localStorage.clear(	);
					}
			}
		});
	}
	showSignup(){
		(<HTMLElement>document.getElementsByClassName('container')[1]).style.display = "block";
	}
	
	getPreferedUsername(){
		var output = "";
		if (window.localStorage) {
		   if (localStorage.length) {
				this.model.username=localStorage.getItem(localStorage.key(0));
				(<HTMLInputElement>document.getElementById('rmb')).checked=true;
		   } else {
			  output = 'There is no data stored for this domain.';
		   }
		} else {
		   output = 'Your browser does not support local storage.'
		}
		console.log(output);
	}
	
	ngOnInit() {
		this.getPreferedUsername()
	}
}
