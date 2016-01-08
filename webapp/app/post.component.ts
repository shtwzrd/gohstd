import {Component, OnInit} from 'angular2/core';
import {NgForm} from 'angular2/common';
import {Post} from './post';
import {AuthService} from './auth.service'

@Component({
    selector: 'gohst-post',
    templateUrl: 'app/post.component.html',
    styleUrls: ['app/post.component.css']
})
export class PostComponent implements OnInit  {
	public posts : Post[]
	
	model = new Post("","");
	constructor(public authService :AuthService){}
	
	getPosts(){
		this.posts = [];
		
		window.fetch('/api/posts',	{
            headers: {
                 'Authorization': this.authService.getAuthHeader()
            }
		}).then((response :any) => {
			this.posts = JSON.parse (response.body);
		});
	}
	
	ngOnInit(){
		this.getPosts();
	}
	
	submitPost(){
		alert("okthatskm8");
	}
}
