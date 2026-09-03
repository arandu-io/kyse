module github.com/arandu-io/kyse

go 1.26

require (
	github.com/arandu-io/framework v0.43.0
	github.com/arandu-io/hesape v0.22.0
)

require (
	golang.org/x/crypto v0.55.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
)

retract (
	v0.15.1 // Requires Framework v0.41.0, which reads _csrf, while Hesape and Dialog emit _token.
	v0.15.0 // Requires Framework v0.41.0 while Hesape emits _token; Dialog also emits _csrf.
)
