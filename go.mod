module github.com/Couches/chirpy

go 1.22.4

replace github.com/Couches/chirpy-database v0.0.0 => ./internal/chirpy-database

require github.com/Couches/chirpy-database v0.0.0

require golang.org/x/crypto v0.25.0 // indirect
