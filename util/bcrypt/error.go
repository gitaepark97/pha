package bcrypt

import "fmt"

func errHashPassword(err error) error {
	return fmt.Errorf("failed to hash password: %w", err)
}
