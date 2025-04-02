module github.com/kien-tn/chirpy

go 1.23.3

require (
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/kien-tn/chirpy/internal/auth v0.0.0-20250401190131-450811b775ff
	github.com/lib/pq v1.10.9
)

require (
	github.com/golang-jwt/jwt/v5 v5.2.2 // indirect
	golang.org/x/crypto v0.36.0 // indirect
)

replace github.com/kien-tn/chirpy/internal/auth => ./internal/auth
