module auth

go 1.20

require (
	github.com/4el0ve4ek/restaraunt-api/library v0.0.0-00010101000000-000000000000
	github.com/go-chi/chi/v5 v5.0.8
	github.com/golang-jwt/jwt/v5 v5.0.0
	github.com/pkg/errors v0.9.1
	golang.org/x/crypto v0.9.0
)

require (
	github.com/lib/pq v1.10.9 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/4el0ve4ek/restaraunt-api/library => ../library