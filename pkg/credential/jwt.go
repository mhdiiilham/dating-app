package credential

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/mhdiiilham/dating-app/entity"
	log "github.com/sirupsen/logrus"
)

// TokenClaims struct holds the data required for jwt payload.
type TokenClaims struct {
	jwt.StandardClaims
	ID    string `json:"id"`
	Email string `json:"email"`
}

// JwtGenerator struct holds required dependecies for Jwt Generator.
type JwtGenerator struct {
	applicationName string
	tokenDuration   time.Duration
	signingMethod   *jwt.SigningMethodHMAC
	signatureKey    string
}

// NewJwtGenerator function return new instance of Jwt Generator.
func NewJwtGenerator(
	applicationName string,
	tokenDuration time.Duration,
	signatureKey string,
) *JwtGenerator {
	return &JwtGenerator{
		applicationName: applicationName,
		tokenDuration:   tokenDuration,
		signingMethod:   jwt.SigningMethodHS256,
		signatureKey:    signatureKey,
	}
}

// CreateAccessToken function use for signed payload (User ID and User Email) to JWT Signed Token.
func (g JwtGenerator) CreateAccessToken(userID, email string) (accessToken string, err error) {
	claims := TokenClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    g.applicationName,
			ExpiresAt: time.Now().Add(g.tokenDuration).Unix(),
		},
		ID:    userID,
		Email: email,
	}

	token := jwt.NewWithClaims(g.signingMethod, claims)
	signedToken, err := token.SignedString([]byte(g.signatureKey))
	if err != nil {
		log.Warnf("[JwtGenerator.CreateAccessToken] error returned from token.SignedString: %v", err)
		return "", entity.ErrInvalidAccessToken
	}

	return signedToken, nil
}

// ParseToken function used to parsed signed token string.
func (g JwtGenerator) ParseToken(accessToken string) (*TokenClaims, error) {
	token, err := jwt.Parse(accessToken, func(t *jwt.Token) (interface{}, error) {
		if method, ok := t.Method.(*jwt.SigningMethodHMAC); !ok || method != g.signingMethod {
			log.Warn("[JwtGenerator.ParseToken] Invalid signing method")
			return nil, entity.ErrInvalidAccessToken
		}

		return []byte(g.signatureKey), nil
	})

	if err != nil {
		log.Warnf("[JwtGenerator.ParseToken] error returned from jwt.Parse: %v", err)
		return nil, entity.ErrInternalServerError
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		log.Warn("[JwtGenerator.ParseToken] failed parse claims to TokenClaims")
		return nil, entity.ErrInvalidAccessToken
	}

	email, _ := claims["email"].(string)
	id, _ := claims["id"].(string)

	return &TokenClaims{
		Email: email,
		ID:    id,
	}, nil
}
