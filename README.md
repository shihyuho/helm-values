[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/shihyuho/helm-values/blob/master/LICENSE)

# Helm Values Plugin

This is a Helm plugin to help developers merge one or more YAML files of values

[![asciicast](https://asciinema.org/a/195144.png)](https://asciinema.org/a/195144)

## Usage

```sh
$ helm values [flags] CHART
```

### Flags

```sh
Flags:
      --backup-suffix string   suffix append to values.yaml if values.yaml already exist in output-dir (default ".bak")
  -h, --help                   help for helm
  -o, --output-dir string      writes the merged values to files in output-dir instead of stdout
  -f, --values valueFiles      specify one or more YAML files of values (default [])
```

## Install

Fetch the latest binary release of helm values and install it:
 
```sh
$ helm plugin install https://github.com/shihyuho/helm-values
```
