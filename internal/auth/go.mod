module github.com/Couches/auth

go 1.22.5

replace github.com/Couches/chirpy-database v0.0.0 => ../chirpy-database

require (
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	golang.org/x/crypto v0.25.0 // indirect
)

require github.com/Couches/chirpy-database v0.0.0
