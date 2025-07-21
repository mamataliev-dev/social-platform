// Package security provides abstractions and implementations for password hashing
// and verification. It supports Single Responsibility and Dependency Inversion principles.
package security

// Hasher defines methods for hashing and verifying passwords. It enables
// Dependency Inversion and Interface Segregation for password security.
type Hasher interface {
	HashPassword(password string) (string, error)
	VerifyPassword(hash, password string) error
}
