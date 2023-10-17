package validator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateRegex(t *testing.T) {
	regex := `^010([0-9]{4})([0-9]{4})$`

	require.True(t, validateRegex(regex, "01011112222"))
	require.False(t, validateRegex(regex, "010-1111-2222"))
}

func TestIsSupportedProductSize(t *testing.T) {
	require.True(t, IsSupportedProductSize("small"))
	require.False(t, IsSupportedProductSize("medium"))
}
