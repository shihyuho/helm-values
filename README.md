[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/shihyuho/helm-values/blob/master/LICENSE)

# Helm Values Plugin

[Helm](https://github.com/helm/helm) doesn't support specify *value.yaml* when packaging chart archive. Therefore this plugin helps developers merge one or more YAML files of values for easily packaging different environments Helm Charts.

[![asciicast](https://asciinema.org/a/195144.png)](https://asciinema.org/a/195144)


## Install

Fetch the latest binary release of helm-values and install it:
 
```sh
$ helm plugin install https://github.com/shihyuho/helm-values
```

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
  -f, --values valueFiles      specify values in a YAML file (can specify multiple) (default [])
```

## Example

The structure is like:

```sh
.
├── mychart
│   ├── .helmignore
│   ├── Chart.yaml
│   ├── charts
│   ├── templates
│   └── values.yaml
└── myenv
    ├── dev.yaml
    ├── sit.yaml
    └── uat.yaml
```

The script for package different environments chart archive:

```sh
# Merge sit values.yaml
$ helm values mychart --values myenv/sit.yaml --output-dir mychart

# Package
$ helm package mychart

# Restore values.yaml
$ mv mychart/values.yaml.bak mychart/values.yaml
```
