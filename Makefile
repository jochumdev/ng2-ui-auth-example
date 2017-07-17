SHELL := /bin/bash

.PHONY: default
default:
	cd $(GOPATH)/src/github.com/pcdummy/ng2-ui-auth-example/server/ && go get -t -v -d ./... && go install -v $(DEBUG) ./...
	if [ ! -e "secrets.ini" ]; then cp secrets.ini.tmpl secrets.ini; fi
	if [ ! -d "json" ]; then mkdir json; fi
	if [ ! -d "keys" ]; then mkdir keys; fi
	if [ ! -e "keys/lxdwebd.rsa" ]; then openssl genrsa -out keys/lxdwebd.rsa 1024; fi
	if [ ! -e "keys/lxdwebd.rsa.pub" ]; then openssl rsa -in keys/lxdwebd.rsa -pubout > keys/lxdwebd.rsa.pub; fi
	@echo "ng2uiauthexampled built successfully"
