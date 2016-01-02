csvschema
=========

csvschema reads 100% of a csv file and generates a postgres schema based on
the contents of the columns. The script expects a single header line as well
as usage of a `,` as the CSV delimiter.

## Usage

```
csvschema ./a_file.csv
```

A Postgres `CREATE` command will be written to STDOUT

## Installation

Go must be installed.

```
go install github.com/bloomapi/csvschema
```

This will checkout csvschema and install it. If `$GOPATH/bin` is in your PATH,
it will be available for execution after installation from the console.
