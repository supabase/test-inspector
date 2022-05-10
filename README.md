# test-inspector

Check your test results against the reference run and compare coverage for multiple versions of your product.

This project provides a less restrictive and more flexible approach then BDD to manage tests for different versions of the same project.

For example you have an `Application` which exposes some `API` and few `libraries` to use this `API` in multiple languages. You can think of one of these as the standard or reference and treat it's test suite as such.

This repo will help you to compare test coverage and feature implementation between all versions of your project and the reference test suite. Also is provides some utilities to help you maintain all versions' test suites in the same state.

## Usage

CLI app to inspect your test suite against reference test run and upload results to test-inspector.

Available Commands:

- `completion` Generate the autocompletion script for the specified shell
- `help` Help about any command
- `inspect` inspect test results comparing to the reference run for your project
- `print` print reference test results for your project
- `upload` upload latest results to test-inspector

Flags:

- `--config` string config file (default is $HOME/.test-inspector.yaml)
- `-h`, `--help` help for test-inspector
- `-H`, `--host` url for test-inspector backend (default "https://gryakvuryfsrgjohzhbq.supabase.co")
- `-w`, `--password` test-inspector user password
- `-f`, `--resultsPath` path to the directory with allure results (default "./allure-results")
- `-t`, `--type` report type (possible values: allure, junit) (default "allure")
- `-u`, `--user` test-inspector user email
- `-v`, `--versionID` version ID in test-inspector (required)

## Web UI

Small web UI to look at some comparison charts.

Go to `./web` to learn more.

## Some details

This project is built with:

- backend: <https://supabase.com>
- cli: `go 1.18`, `cobra`, `go-junit`, `supabase-go`
- web-gui: `node`, `nuxt`, `supabase-js`, deployed to `fly`
