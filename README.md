# Voorhees

[![Go Report Card](https://goreportcard.com/badge/github.com/nivl/voorhees)](https://goreportcard.com/report/github.com/nivl/voorhees)

Voorhees is a program that parses the depency tree to find dependency that
might no longer be maintained.

## Install

`go get -u github.com/Nivl/voorhees`

## usage

`voorhees [flags]`

| Flag         | Description                                                           |
| ------------ | --------------------------------------------------------------------- |
| --time -t    | specify time after wich we assume a dep might no longer be maintained |
| --ignore, -i | coma separated list of packages to ignore                             |
| --indirect   | check indirect modules                                                |
