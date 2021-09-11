.PHONY: dist/neptunes-scanner
dist/neptunes-scanner:
	mkdir -p dist && go get && go test ./... && go build -o dist/np-scanner
