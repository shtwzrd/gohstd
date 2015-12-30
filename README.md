#gohstd

A web service for hosting your command history remotely, to be queried by all your machines via the `gohst` command-line client

## Building from source

1. Install Go for your platform
2. Set the GOPATH environment variable
3. Fetch gohstd
Install it via `go get` like so:
    `go get github.com/warreq/gohstd`
_or_
   `git clone https://github.com/warreq/gohstd`
   `cd gohstd` 
   `go build`
4. (_Optional_) Build the web application
   `npm install`
   `npm run build`

## Installation

1. Drop the gohstd binary somewhere on your PATH
2. Set the environment variable `DATABASE_URL` to a postgres database connection string (format _postgres://username:password@host:port/db-name_) 
2. (Optional) Set the environment variable `PORT` to an available port. The default is 8080 when PORT is not provided. 
3. Execute gohstd! 
