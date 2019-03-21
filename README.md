# gometaimports

A simple tool implementing [meta tags](https://golang.org/cmd/go/#hdr-Remote_import_paths)
for Golang packages remote imports.

## Usage

Prepare configuration file:

```bash
cp config.yaml.dist config.yaml

# edit config.yaml for your needs
```

Ensure modules are present in `vendor` directory:

```bash
go mod download
go mod vendor
```

### Local

To run locally, build and run the application:

```bash
go build
./gometaimports -config ./config.yaml
```

### Docker

Build docker image:

```bash
docker image build -t hypnoglow/gometaimports:dirty .
```

Run docker container:

```bash
docker container run -v $(pwd)/config.yaml:/etc/gometaimports.yaml -p 8080:8080 \
  hypnoglow/gometaimports:dirty gometaimports -config /etc/gometaimports.yaml
```

## Licence

[MIT](LICENCE).
