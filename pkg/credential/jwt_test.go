package credential_test

import (
	"testing"
	"time"

	"github.com/mhdiiilham/dating-app/pkg/credential"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccessToken(t *testing.T) {
	testCases := []struct {
		condition   string
		expectedErr error
		ID          string
		Email       string
	}{
		{
			condition:   "access token created",
			expectedErr: nil,
			ID:          "1",
			Email:       "hi@muhammadilham.xyz",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.condition, func(t *testing.T) {
			assert := assert.New(t)
			jwtGeneratorClient := credential.NewJwtGenerator("dating-app", 2*time.Hour, "the secret of kalimdor")
			accessToken, actualErr := jwtGeneratorClient.CreateAccessToken("1", "hi@muhammadilham.xyz")
			assert.Equal(tc.expectedErr, actualErr)

			claims, err := jwtGeneratorClient.ParseToken(accessToken)
			assert.NoError(err)
			assert.Equal(tc.ID, claims.ID)
			assert.Equal(tc.Email, claims.Email)
		})
	}
}
