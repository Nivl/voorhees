# Voorhees

Voorhees is a program that parses the depency tree to find dependencies that
might no longer be maintained.

```
❯ go list -json -m -u all | voorhees
+-----------------------------------+----------------------------+
|              MODULE               |        LAST UPDATE         |
+-----------------------------------+----------------------------+
| github.com/olekukonko/tablewriter | 13 months ago (2019/12/05) |
| github.com/pkg/errors             | 12 months ago (2020/01/14) |
| github.com/spf13/pflag            | 16 months ago (2019/09/18) |
+-----------------------------------+----------------------------+
```

## usage

```
go list -json -m -u all | voorhees [flags]

| Flag          | Description                                                               |
| ------------- | ------------------------------------------------------------------------- |
| --config, -c  | Path to the optional config file (default: ./.voorhees.yml)               |
| --help, -h    | Display the help options                                                  |
| --ignore, -i  | Coma separated list of packages to ignore                                 |
| --limit -l    | Number of months after which a dep is considered unmaintained (default 6) |
| --version, -v | Display the version number                                                |
```

## Installation

### Github Action

See [voorhees-github-action](https://github.com/Nivl/voorhees-github-action).

```
voorhees:
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v2
    - uses: actions/setup-go@v2

    - name: Generate go.list
      run: go list -json -m all > go.list

    - name: Run Voorhees
      uses: Nivl/voorhees-github-action@v1
```

### macOS

```
brew install Nivl/homebrew-tap/voorhees
```

### Linux and Windows

Download a binary from the [Release](https://github.com/Nivl/voorhees/releases) page.

### Docker

```
go list -json -m all | docker run --rm -i ghcr.io/nivl/voorhees:latest
```

### Latest master

```
go get -u github.com/Nivl/voorhees
```


## Configration

You can configure each package separatetly using `.voorhees.yml`.

Example:

```
version: 1
default:
  limit: 6 months
rules:
  github.com/olekukonko/tablewriter: 52 weeks
  github.com/pkg/errors: 10 months
  github.com/spf13/pflag: skip
```

- `version`: version of the config file (default is and always will be 1 for the current major version of Voorhees).
- `default`:
  - `limit`: Will only report the packages if it hasn't been updated in the last _N_ weeks or month.
- `rules`: contains a list of key/values. The keys represents the packages, the value, the rule you want to apply. The value can have any of the following format:
  - `ignore`: Ignore the package.
  - `skip`: Alias for `ignore`.
  - `N weeks`: Will only report the package if it hasn't been updated in _N_ weeks.
  - `N week`: Alias for `N weeks`
  - `N month`: Will only report the package if it hasn't been updated in _N_ months.
  - `N months`: Alias for `N months`.
