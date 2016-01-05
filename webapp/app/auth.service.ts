export class AuthService {
    private _username :string = "";
    private _password :string = "";
    authenticated :boolean;

    constructor() {
        this.authenticated = false;
    }

    // authenticate returns a Promise, whose result will contain a boolean
    // value; true if authentication was successful and false otherwise
    authenticate(username :string, password :string) :Promise<boolean> {
        var promise = new Promise(function (resolve, reject) {
            window.fetch('/api/users/login', {
                headers: {
                    'Authorization': 'Basic ' + btoa(username + ':' + password)
                }
            }).then((response :any) => {
                if (response.status == 200) {
                    this._username = username;
                    this._password = password;
                    this.authenticated = true;
                    resolve(true);
                } else {
                    console.log('Authentication failure');
                    resolve(false);
                }
            })
        });
        return promise;
    }

    getAuthHeader() :string {
        if (!this.authenticated) {
            console.log('Error: not authenticated.')
        } else {
            return btoa(this._username + ':' + this._password);
        }
    }
}
