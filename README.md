# gometaimports

A simple tool implementing [meta tags](https://golang.org/cmd/go/#hdr-Remote_import_paths)
for Golang packages remote imports.

## Usage

Prepare configuration file:

```bash
cp config.yaml.dist config.yaml

# edit config.yaml for your needs
```

To run locally, build and run the application:

```bash
go build
./gometaimports -config ./config.yaml
```

## Licence

[MIT](LICENCE).
