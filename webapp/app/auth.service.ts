import {Http, Headers} from 'angular2/http';
import {Injectable} from 'angular2/core';

@Injectable()
export class AuthService {
    private authenticated :boolean;
    public username :string;
    private header :Headers;

    constructor(private http :Http) {
        this.authenticated = false;
    }

    // authenticate returns a Promise, whose result will contain a boolean
    // value; true if authentication was successful and false otherwise
    authenticate(username :string, password :string) :Promise<boolean> {
        var promise = new Promise((resolve, reject) => {
            var h = new Headers();
            h.append('Authorization', 'Basic ' + btoa(username + ':' + password));
            this.http.get('/api/users/login', {headers: h})
                .map(res => res.status)
                .subscribe(data => {
                    if (data == 200) {
                        this.header = h;
                        this.username = username;
                        this.authenticated = true;
                        resolve(true);
                    } else {
                        resolve(false);
                    }
                });
        });
        return promise;
    }

    getHeader() :Headers {
        if (!this.authenticated) {
            console.log('Error: not authenticated.')
        } else {
            return this.header;
        }
    }
}
