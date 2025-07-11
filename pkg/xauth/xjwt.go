package xauth

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xsecurity"
	"github.com/golang-jwt/jwt/v5"
)

type JWTManager interface {
	GenerateToken(claims jwt.MapClaims, validityDuration time.Duration) (string, error)
	ValidateToken(token string) (jwt.MapClaims, error)
	GetClaimBy(token, claimKey string) (any, error)
	GetClaims(token string) (jwt.MapClaims, error)
}

// Predefined errors for JWT operations
var (
	ErrInvalidSigningMethod = errors.New("invalid signing method")
	ErrInvalidToken         = errors.New("invalid token")
	ErrClaimNotFound        = errors.New("claim not found")
	ErrTokenExpired         = errors.New("token has expired")
	ErrTokenNotYetValid     = errors.New("token is not yet valid")
)

// _JWTManager handles JWT generation and validation
type _JWTManager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	issuer     string
}

// NewJWTManager initializes a new _JWTManager
func NewJWTManager(privateKeyPath, issuer string) (*_JWTManager, error) {
	privateKey, err := xsecurity.LoadPrivateKey(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load jwt private key: %w", err)
	}

	return &_JWTManager{
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
		issuer:     issuer,
	}, nil
}

// GenerateToken generates a new JWT token
func (j *_JWTManager) GenerateToken(claims jwt.MapClaims, validityDuration time.Duration) (string, error) {
	claims["iss"] = j.issuer
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(validityDuration).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(j.privateKey)
}

// ValidateToken validates a JWT token and returns its claims
func (j *_JWTManager) ValidateToken(token string) (jwt.MapClaims, error) {
	_token, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, ErrInvalidSigningMethod
		}
		return j.publicKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, ErrTokenNotYetValid
		}
		return nil, ErrInvalidToken
	}

	if claims, ok := _token.Claims.(jwt.MapClaims); ok && _token.Valid {
		return claims, nil
	}
	return nil, ErrInvalidToken
}

// ExtractClaim extracts a specific claim from the token
func (j *_JWTManager) GetClaimBy(token, claimKey string) (any, error) {
	claims, err := j.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	claimValue, ok := claims[claimKey]
	if !ok {
		return nil, ErrClaimNotFound
	}

	return claimValue, nil
}

// ExtractClaim extracts a specific claim from the token
func (j *_JWTManager) GetClaims(token string) (jwt.MapClaims, error) {
	claims, err := j.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	return claims, nil
}
