# Helm Values Plugin

This is a Helm plugin to help developers merge one or more YAML files of values

[![asciicast](https://asciinema.org/a/StHw7ZiwrnJOoCr9ZoyRJpUKK.png)](https://asciinema.org/a/StHw7ZiwrnJOoCr9ZoyRJpUKK)

## Usage

```sh
$ helm values [flags] CHART
```

### Flags

```sh
  -h, --help                help for helm
  -o, --output string       write to a file, instead of stdout
  -f, --values valueFiles   specify one or more YAML files of values (default [])
  -v, --verbose             enable verbose output.
```

## Install

Fetch the latest binary release of helm values and install it:
 
```sh
$ helm plugin install https://github.com/shihyuho/helm-values
```
