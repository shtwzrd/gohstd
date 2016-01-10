import {Component,OnChanges, SimpleChange,} from 'angular2/core';
import {NgForm} from 'angular2/common';
import {Http, Response, HTTP_PROVIDERS} from 'angular2/http';
import {Configuration, NewConfiguration} from './configuration';
import {FILE_UPLOAD_DIRECTIVES, FileUploader} from 'ng2-file-upload';
import {AuthService} from './auth.service';

@Component({
  selector: 'gohst-configuration',
  templateUrl: 'app/configuration.component.html'
})
export class ConfigurationComponent {
	model = new NewConfiguration("");
	
    constructor(private http :Http, private auth :AuthService){
     
    }
	
	submitConfiguration(){
        if (!this.auth.authenticated) {
            alert("You must log in before submitting your configuration.");
            return;
        }
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