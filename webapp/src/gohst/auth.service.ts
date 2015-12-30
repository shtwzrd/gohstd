import {Injectable} from 'angular2/core';
import {Http, Headers, HTTP_PROVIDERS} from 'angular2/http';

@Injectable()
export class AuthService {

    private _http :Http;
    private _username :string = "";
    private _password :string = "";
    private _authenticated :boolean = false;

    constructor(http :Http) {
        this._http = http;
    }

    authenticate(username, password :string) :boolean {
        var headers = new Headers();
        headers.append('Authorization', 'Basic ' + btoa(username + ':' + password));
        this._http.get('/api/users/login');
        return false;
    }

    needsCredentials() :boolean {
        return false;
    }
    getAuthHeader() :string {
        return ""
    }
}
