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
        this.http.get('/api/posts')
            .map((res :Response) => res.json())
            .subscribe((posts :Array<Post>) => this.posts = posts);
    }

    submitPost() {
        if (!this.auth.authenticated) {
            alert("You must log in before submitting a post.");
            return;
        }
        this.http.post('/api/users/' + this.auth.username + '/posts',
                       JSON.stringify(this.model),
                       {headers: this.auth.getHeader()})
            .map(res => res.status)
            .subscribe(
                data => {
                    if (data == 200) {
                        this.model.message = ""
                        this.model.title = ""
                        this.getPosts()
                    }
                },
                err => alert(err)
            );
    }

    ngOnInit() {
        this.getPosts();
    }
}
