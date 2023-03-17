package hashutil

import (
	"github.com/stretchr/testify/require"
	"github.com/thehaung/simplebank/util/randutil"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestPassword(t *testing.T) {
	password := randutil.StringWithQuantity(6)
	hashedPassword1, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1)

	err = CheckPassword(password, hashedPassword1)
	require.NoError(t, err)

	wrongPassword := randutil.StringWithQuantity(6)
	err = CheckPassword(wrongPassword, hashedPassword1)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashedPassword2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2)
	require.NotEqual(t, hashedPassword1, hashedPassword2)
}

func BenchmarkPassword(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = CheckPassword("secret", "$2a$10$I00w.S7ELh37J1UpRsGYruqVgAM5jGsQUxSfrNmmqflYytXf4CVY2")
	}
}
