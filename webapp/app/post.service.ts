import {Post} from './post';
import {Injectable} from 'angular2/core';
import {AuthService} from './auth.service';
@Injectable()
export class PostService {

    constructor(public authService :AuthService) { }

    getPosts() :Promise<Post[]> {
        return new Promise<Post[]>(resolve => {
            window.fetch('/api/posts',	{
                headers: {
                    'Authorization': this.authService.getAuthHeader()
                }
            });
        });
    }
}
