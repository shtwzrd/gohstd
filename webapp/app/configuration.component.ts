import {Component} from 'angular2/core';
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
		alert(this.model.picture);
	}

}