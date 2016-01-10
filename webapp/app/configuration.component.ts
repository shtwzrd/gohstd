import {Component,OnChanges, SimpleChange,} from 'angular2/core';
import {NgForm} from 'angular2/common';
import {Configuration, NewConfiguration} from './configuration';
import {FILE_UPLOAD_DIRECTIVES, FileUploader} from 'ng2-file-upload';


@Component({
  selector: 'gohst-configuration',
  templateUrl: 'app/configuration.component.html'
})
export class ConfigurationComponent {
	model = new NewConfiguration("");
	
	submitConfiguration(){
		console.log(this.model);
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