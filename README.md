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

## How to use in your process

1. Take a look at the existing project at <https://test-inspector.fly.dev/projects/1>

   - You can select the reference version (supabase-js currently) and check out what tests are currently passing for it (<https://test-inspector.fly.dev/projects/1/versions/2>)

2. Download the latest version of the CLI app from <https://github.com/supabase/test-inspector/suites/8289986777/artifacts/363373298>

   - When you download it on MacOS, you may need to allow it to run. You can do it by right-clicking on the file with Ctrl and selecting `Open` from the context menu.

3. Run the CLI app with `./test-inspector -v 2 print` command to see the detailed test results for the reference version of the library (you can find the version ID on the test-inspector website). You will also see steps for each test case so you can make your tests as close as possible to the reference version.

4. To upload your test results to test-inspector you need your test runner to generate testrun report in `allure` (<https://docs.qameta.io/allure/>) or `junit` (<https://www.ibm.com/docs/en/developer-for-zos/14.1?topic=formats-junit-xml-format>) format. Check out how to do that for your programming language.

5. You can compare your test results with the reference version by running `./test-inspector -v 5 inspect -t junit -f ./path/to/junit.xml` or `./test-inspector -v 5 inspect -f ./path/to/allure-results` command. You will see the comparison chart for each test case and the overall coverage.

6. To upload results to test-inspector you need to register at <https://test-inspector.fly.dev/login>. Find the version ID for a library you are testing (for example Python is #3). And run the following command `./test-inspector -u username@example.com -w $INSPECTOR_PASSWORD -v $VERSION_ID -f ./allure-results upload -l $LAUNCH_NAME`. Or the same command but with a path to `junit` report with `-t junit` option.
