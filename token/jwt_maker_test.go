package token

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"
	"github.com/thehaung/simplebank/util/randutil"
	"testing"
	"time"
)

func TestJwtMaker(t *testing.T) {
	maker, err := NewJwtMaker(randutil.StringWithQuantity(32))
	require.NoError(t, err)

	userName := randutil.Owner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(userName, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, userName, payload.Username)

	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredJwtToken(t *testing.T) {
	maker, err := NewJwtMaker(randutil.StringWithQuantity(32))
	require.NoError(t, err)

	token, payload, err := maker.CreateToken(randutil.Owner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidJwtTokenAlgNone(t *testing.T) {
	payload, err := NewPayload(randutil.Owner(), time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJwtMaker(randutil.StringWithQuantity(32))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}

func TestInvalidSecretKeySize(t *testing.T) {
	payload, err := NewPayload(randutil.Owner(), time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	maker, err := NewJwtMaker(randutil.StringWithQuantity(30))
	require.Error(t, err)
	require.Nil(t, maker)
}
