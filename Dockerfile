FROM golang:1.12-alpine3.9 AS build

ARG VERSION="unknown"

WORKDIR /opt/gometaimports/

COPY . .

RUN CGO_ENABLED=0 go build -mod=vendor -ldflags "-X main.version=${VERSION} -X main.templateDir=/usr/share/gometaimports/templates"

FROM alpine:3.9

COPY --from=build /opt/gometaimports/gometaimports /usr/local/bin/gometaimports
COPY --from=build /opt/gometaimports/templates /usr/share/gometaimports/templates

CMD ["gometaimports"]
