# Build SPA
FROM node:16-alpine3.12 as SPA

COPY ./ui /ui
COPY Makefile /
WORKDIR /

RUN apk add --no-cache make && make ui

# Build binary
FROM golang:1.17.1-alpine as builder

RUN apk add --no-cache make git

COPY . $GOPATH/src/go.albinodrought.com/neptunes-pride
WORKDIR $GOPATH/src/go.albinodrought.com/neptunes-pride

## Embed SPA
COPY --from=SPA /ui/dist $GOPATH/src/go.albinodrought.com/neptunes-pride/ui/dist

ENV CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

RUN make dist/np-scanner && mv dist/np-scanner /np-scanner

# Lightweight Runtime Env
FROM gcr.io/distroless/base-debian10
COPY --from=builder /np-scanner /np-scanner
CMD ["/np-scanner", "serve"]
