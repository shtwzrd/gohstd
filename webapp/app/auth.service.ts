export class AuthService {
    private username :string = "";
    private password :string = "";
    private authenticated :boolean;

    constructor() {
        this.authenticated = false;
    }

    // authenticate returns a Promise, whose result will contain a boolean
    // value; true if authentication was successful and false otherwise
    authenticate(username :string, password :string) :Promise<boolean> {
        var promise = new Promise((resolve, reject) => {
            window.fetch('/api/users/login', {
                headers: {
                    'Authorization': 'Basic ' + btoa(username + ':' + password)
                }
            }).then((response :any) => {
                if (response.status == 200) {
                    this.setCredentials(username, password);
                    resolve(true);
                } else {
                    console.log('Authentication failure');
                    resolve(false);
                }
            })
        });
        return promise;
    }

    setCredentials(username :string, password :string) {
        this.username = username;
        this.password = password;
        this.authenticated = true;
    }

    getAuthHeader() :string {
        if (!this.authenticated) {
            console.log('Error: not authenticated.')
        } else {
            return btoa(this.username + ':' + this.password);
        }
    }
}
