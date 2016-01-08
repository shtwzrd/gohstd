import {Component} from 'angular2/core';
import {NgForm} from 'angular2/common';
import {Post} from './Post';

@Component({
    selector: 'gohst-post',
    templateUrl: 'app/post.component.html',
    styleUrls: ['app/post.component.css']
})
export class PostComponent {
	model = new Post("","");
	
	submitPost(){
		alert("okthatskm8");
	}
}
