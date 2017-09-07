# Babelfish Go Client [![Build Status](https://travis-ci.org/bblfsh/client-go.svg?branch=master)](https://travis-ci.org/bblfsh/client-go) [![codecov](https://codecov.io/gh/bblfsh/client-go/branch/master/graph/badge.svg)](https://codecov.io/gh/bblfsh/client-go)

Babelfish Go client library provides functionality to both
connect to the Babelfish server to parse code
(obtaining an [UAST](https://doc.bblf.sh/uast/specification.html) as a result)
and to analyse UASTs with the functionality provided by [libuast](https://github.com/bblfsh/libuast).

## Installation

```
$ make
```

## Usage

Users of this library will want use it as basis for their own code analysis,
but a standalone app is also provided to easily give a try to its features:

```
$ cli -e localhost:9432 -f sample.py -q "/compilation_unit//identifier"

```
