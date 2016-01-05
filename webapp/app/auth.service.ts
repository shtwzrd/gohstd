export class AuthService {
    private _username :string = "";
    private _password :string = "";
    authenticated :boolean;

    constructor() {
        this.authenticated = false;
    }

    authenticate(username :string, password :string) :boolean {
        window.fetch('/api/users/login', {
            headers: {
                'Authorization': 'Basic ' + btoa(username + ':' + password)
            }
        })
            .then(function (response :any) {
                if (response.status == 200) {
                    this._username = username;
                    this._password = password;
                    this.authenticated = true;
                    console.log(this.authenticated);
                } else {
                    console.log('Authentication failure');
                    return false;
                }
            });
        return true;
    }

    getAuthHeader() :string {
        if (!this.authenticated) {
            console.log('Error: not authenticated.')
        } else {
            return btoa(this._username + ':' + this._password);
        }
    }
}
