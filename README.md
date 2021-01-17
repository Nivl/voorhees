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

## Installation

### Macos

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

## usage

```
go list -json -m -u all | voorhees [flags]
```

| Flag         | Description                                                               |
| ------------ | ------------------------------------------------------------------------- |
| --limit -l   | number of weeks after which a dep is considered unmaintained (default 26) |
| --ignore, -i | coma separated list of packages to ignore                                 |

```

```
