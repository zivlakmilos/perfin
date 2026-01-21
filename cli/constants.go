package cli

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

const (
	primaryColor   = (lipgloss.Color("#FF06B7"))
	secondaryColor = (lipgloss.Color("#767676"))
	errorColor     = (lipgloss.Color("#FF0000"))
	successColor   = (lipgloss.Color("#00FF00"))
)

func getUrl(base, endpoint string) string {
	return fmt.Sprintf("%s%s", base, endpoint)
}
