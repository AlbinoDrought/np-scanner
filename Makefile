all: ui dist/np-scanner

.PHONY: dist/np-scanner
dist/np-scanner:
	rm -rf internal/web/packaged && mkdir -p internal/web/packaged
	cp -ar ui/dist/. internal/web/packaged/.
	git archive HEAD -o internal/web/packaged/source.tar.gz
	go get
	go generate ./...
	go test ./...
	mkdir -p dist && rm -f dist/np-scanner && go build -o dist/np-scanner

.PHONY: ui
ui:
	cd ui && npm install && npm run build
