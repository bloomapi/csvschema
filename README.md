csvschema
=========
[![Build Status](https://travis-ci.org/bloomapi/csvschema.svg?branch=master)](https://travis-ci.org/bloomapi/csvschema)

csvschema reads 100% of a csv file and generates a postgres schema based on
the contents of the columns. The script expects a single header line as well
as usage of a `,` as the CSV delimiter.

csvschema is fast. Scanning 100% of a 4.7 million row CSV file of 5.3GB takes 6
minutes on a 3ghz mid-2014 Macbook Pro.

## Usage

```
csvschema ./a_file.csv
```

A Postgres `CREATE` command will be written to STDOUT

## Installation

Go must be installed.

```
go get github.com/bloomapi/csvschema
```

This will checkout csvschema and install it. If `$GOPATH/bin` is in your PATH,
it will be available for execution after installation from the console.
