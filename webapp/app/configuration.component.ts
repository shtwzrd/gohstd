import {Component, OnChanges, SimpleChange} from 'angular2/core';
import {NgForm} from 'angular2/common';
import {AuthService} from './auth.service';
import 'rxjs/add/operator/map';

@Component({
    selector: 'gohst-configuration',
    templateUrl: 'app/configuration.component.html'
})
export class ConfigurationComponent {
    private file :File;
    constructor(private auth :AuthService){
    }

    submitConfiguration() {
        if (!this.auth.authenticated) {
            alert("You must log in before submitting your configuration.");
            return;
        }

        var basic = this.auth.getHeader().get('Authorization');
        var xhr = new XMLHttpRequest();
        xhr.setRequestHeader('Authorization', basic);
        xhr.setRequestHeader('Content-Type', 'image/png');
        xhr.open('POST', 'api/users/' + this.auth.username + '/picture');

        var file    = (<HTMLInputElement>document.querySelector('input[type=file]')).files[0];
        var reader  = new FileReader();

        reader.onload = function (evt) {
            console.log(evt.target);
            xhr.send(evt.target);
            alert("upload completed");
        }
        reader.readAsBinaryString(file);
    }

    picturePreview(){
        var preview = document.querySelector('img');
        var file    = (<HTMLInputElement>document.querySelector('input[type=file]')).files[0];
        var reader  = new FileReader();

        reader.onloadend = function () {
            (<HTMLImageElement>preview).src = reader.result;
        }

        if (file) {
            reader.readAsDataURL(file);
        } else {
            (<HTMLImageElement>preview).src = "";
        }
    }
}
