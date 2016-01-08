import {Component, OnInit} from 'angular2/core';
import {NgForm} from 'angular2/common';
import {Http, Response, HTTP_PROVIDERS} from 'angular2/http';
import {Post, NewPost} from './post';
import {AuthService} from './auth.service';
import 'rxjs/add/operator/map';

@Component({
    selector: 'gohst-post',
    templateUrl: 'app/post.component.html',
    styleUrls: ['app/post.component.css'],
    providers: [HTTP_PROVIDERS]
})
export class PostComponent implements OnInit {
    public posts: Post[]
    model = new NewPost("","");

    constructor(private http :Http, private auth :AuthService){
        this.posts = [];
    }

    getPosts() {
        this.http.get('/api/posts', {headers: this.auth.getAuthHeader()})
            .map((res :Response) => res.json())
            .subscribe((posts :Array<Post>) => this.posts = posts);
    }

    submitPost() {
        alert("okthatskm8");
    }

    ngOnInit() {
        this.getPosts();
    }
}
