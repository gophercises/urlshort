# How to use this program

_To use this program, a copy of the [sqlite](https://sqlite.org/) database software may be required._

Build the program using `go build` to obtain a binary. Alternatively, `go run main.go` can be used.

Usage: ./binary [options]

Available options:

- -yaml "path to YAML file"
- -json "path to JSON file"
- -db "path to sqlite3 database"
- -help show this help screen

Examples of database sources that can be used are `urlmap.yml`, `urlmap.json`, and `url_import.sql`, provided in this directory.

`url_import.sql` is used to generate an sqlite3 persistent database which will be used by default if no option is passed to the program.
If you wish to use the `-db` database option, run the command `sqlite urls.db < url_import.sql` on first use to set up the database.