import {Component, OnInit} from 'angular2/core';
import {NgForm} from 'angular2/common';
import {Post} from './post';
import {PostService} from './post.service'

@Component({
    selector: 'gohst-post',
    templateUrl: 'app/post.component.html',
    styleUrls: ['app/post.component.css'],
    providers: [PostService]
})
export class PostComponent implements OnInit {
    public posts: Post[]
    model = new NewPost("","");

    constructor(private postService :PostService){
        this.posts = [];
    }

    getPosts(){
        this.postService.getPosts().then(posts => this.posts = posts);
    }

    submitPost(){
		    alert("okthatskm8");
	  }

    ngOnInit(){
        this.getPosts();
    }
}
