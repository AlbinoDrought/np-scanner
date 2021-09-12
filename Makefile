.PHONY: dist/neptunes-scanner
dist/neptunes-scanner: # ui
	rm -rf internal/web/packaged && mkdir -p internal/web/packaged
	cp -ar ui/dist/. internal/web/packaged/.
	go get
	go generate ./...
	go test ./...
	mkdir -p dist && rm -f dist/np-scanner && go build -o dist/np-scanner

.PHONY: ui
ui:
	cd ui && npm install && npm run build
