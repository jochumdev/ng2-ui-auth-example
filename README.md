## ng2-ui-auth-example - OAuth with [ng2-ui-auth](https://github.com/ronzeidman/ng2-ui-auth) and a Go server

### Online Demo

- I run this at [ng2-satellizer.pc-dummy.net](https://ng2-satellizer.pc-dummy.net).

### Requirements to run your own

- Go 1.7+ (I use [gvm](https://github.com/moovweb/gvm) for that).
- Nodejs (I use [nvm](https://github.com/creationix/nvm) inside a pyenv venv).
- GNU make to build the go server.
- An account on at least one of Google, Facebook, Github to test it.

### Usage

1. $ go get github.com/pcdummy/ng2-ui-auth-example
2. $ cd $GOPATH/src/github.com/pcdummy/ng2-ui-auth-example
3. $ cp secrets.ini.tmpl to secrets.ini
4. Goto Google, Facebook and Github and create an OAuth app.
5. Paste the ClientId and the Secret into your secrets.ini.
6. Install node, npm and go 1.7+
7. Run "$ make" to build the go server
8. Build the angular client: $ cd static/; npm i; npm run dev && npm run prod
9. Run the server: $ cd ..; ng2uiauthexampled --config dev.ini serve

### TODO

- Move the Go server and its components to server/.

### Authors

- René Jochum @pcdummy

### License

MIT
