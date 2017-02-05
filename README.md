Bench [![Go Report Card](https://goreportcard.com/badge/github.com/lnsp/bench)](https://goreportcard.com/report/github.com/lnsp/bench) [![Build Status](https://travis-ci.org/lnsp/bench.svg?&branch=develop)](https://travis-ci.org/lnsp/bench)
========

Bench is a file patch system using HTTP(S) and SHA-1 hashing.

## Installation
Run `go get github.com/lnsp/bench/cmd/bench` to install the `bench` binary.

## Usage guide
**Important:** Your working folder is automatically made the target folder for patch generation and fetching. To change this behavior, use the `--target "my-target"` flag to chose your target independently from your working directory.

### Fetching a patch
If you first getting started you may want to fetch a patch from a source (either a web server or a local folder).
This can be done using the `bench fetch` command.

- If you patch your target for the first time, you have to include a `--source "patch-source"` flag to specify your patch source. This can either be a folder on your local file system like `/etc/foo/bar` or a web address like `https://example.com/foo/bar/`.
- If you already have patched once or more, the stakes are high that your existing patch file has a source included. You can try it without specifying a special source, but if it fails, try to repeat the step above.

### Generating a patch
If you want to generate a patch, you can just enter `patch generate`. If you want to make it easier for your user to patch from your source, include the `--source "my-source"` flag to simplify their patching process (but ensure that the source can be mapped to your data).

## License
Copyright 2016 Lennart Espe. All rights reserved.

Use of this source code is governed by a MIT-style license that can be found in the LICENSE.md file.
