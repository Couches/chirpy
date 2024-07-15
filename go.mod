module github.com/Couches/chirpy

go 1.22.4

replace github.com/Couches/chirpy-database v0.0.0 => ./internal/chirpy-database

replace github.com/Couches/auth v0.0.0 => ./internal/auth

require (
	github.com/Couches/chirpy-database v0.0.0
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/joho/godotenv v1.5.1
	golang.org/x/crypto v0.25.0
)
