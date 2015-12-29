import {Http, HTTP_PROVIDERS} from 'angular2/http';
export class AuthService {
    constructor(
        private _backend: BackendService,
        private _logger: Logger) { }

    private _username:string = "";
    private _password:string = "";
    private _authenticated:bool = false;

    authenticate(username, password :string) :bool {
        var headers = new Headers();
        headers.append('Authorization', 'Basic ' + btoa(username + ':' + password));
        http.get('/api/users/login');
    }

    needsCredentials() :boolean {}
    getAuthHeader() :string {}
    getHeroes() {
        this._backend.getAll(Hero).then( (heroes:Hero[]) => {
            this._logger.log(`Fetched ${heroes.length} heroes.`);
            this._heroes.push(...heroes); // fill cache
        });
        return this._heroes;
    }
}
