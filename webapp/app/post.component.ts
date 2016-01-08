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
    public posts: Post[]

    model = new Post("","");
    constructor(public authService :AuthService){
        this.getPosts();
    }

    getPosts(){
        window.fetch('/api/posts',	{
            headers: {
                'Authorization': this.authService.getAuthHeader()
            }
        }).then((response :any) => {
            this.posts = response.body;
        }).catch((response :any) => {
            alert(response);
        });
    }

    ngOnInit(){
        this.getPosts();
    }

	  submitPost(){
		    alert("okthatskm8");
	  }
}
