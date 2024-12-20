package shared

import (
	"fmt"
)

func GenerateError(error_message error) error {
	return fmt.Errorf("%s", fmt.Sprintf("An error ocurred!: %s", error_message))
}
