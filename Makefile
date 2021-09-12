all: dist/np-scanner ui

.PHONY: dist/np-scanner
dist/np-scanner:
	rm -rf internal/web/packaged && mkdir -p internal/web/packaged
	cp -ar ui/dist/. internal/web/packaged/.
	GO111MODULE=off go get github.com/GeertJohan/go.rice/rice
	go get
	go generate ./...
	go test ./...
	mkdir -p dist && rm -f dist/np-scanner && go build -o dist/np-scanner

.PHONY: ui
ui:
	cd ui && npm install && npm run build
