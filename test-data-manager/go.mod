module github.com/mtlynch/whatgotdone/test-data-manager

go 1.16

replace github.com/mtlynch/whatgotdone/backend => ../backend

require (
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/mattn/go-sqlite3 v1.14.9
	github.com/mtlynch/whatgotdone/backend v0.0.0-00010101000000-000000000000
	gopkg.in/yaml.v2 v2.4.0
)
