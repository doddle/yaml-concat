# Usage

This tool will:
- walk a `--path` of files looking for `*.yaml` and `*.yml` files
- on default hidden files (eg `.hidden.yaml` will be ignored (#TODO: allow searching for them also, or a custom exclude list?)
- for each yaml file found it will stream each individual document to stdout

Notes:
- This handles individual **documents** inside a file (`---` seperated docs in a files) consistently
- `stderr` is never used, everything is `stdout`
- comments and quoting are preserved
- some additional notes about what its doing can be outputted, these are always produced as safe  `# ` comments

`yaml-concat --exclude somefile.yaml --path ./source/files/`


# Sources

- some code was taken from: https://github.com/wangkuiyi/yamlfmt
- a modified go-yaml.v3 is used to keep sequence indentation to a minimum (https://github.com/go-yaml/yaml/issues/661 + https://github.com/starkers/yaml/commit/63f3856906e9106804ce495f3077d99340cdf9d9)


# Reasoning:

1. I need to safely stream files from a folder into stdout (and exclude certain ones)
2. python's vanilla pyyaml is a nightmare (doesn't care about quotes and goes on to even wreck types and then also their values..) (EG: `foo: '00123'` can become: `foo: 123` ..danger Will Robinson!)


# example:

input file: "sample.yaml"

raw content looks like:
```
apiVersion: acme/v1
kind: Deployment
spec:
  replicas: 3
  template:
    spec:
      containers:
      # much myimage, so wow
      - image: myimage:foo-bar
        name: someapp
        ports:
          - containerPort: 8080
            name: http
        env:
        - name: SOME_VAR
          value: '00123123'
```

output:

```
# source: ./sample.yaml[0]
apiVersion: acme/v1
kind: Deployment
spec:
  replicas: 3
  template:
    spec:
      containers:
      # much myimage, so wow
      - image: myimage:foo-bar
        name: someapp
        ports:
        - containerPort: 8080
          name: http
        env:
        - name: SOME_VAR
          value: '00123123'
---
# source: ./sample.yaml[1]
anotherDoc: true
handled: "just fine"
```
