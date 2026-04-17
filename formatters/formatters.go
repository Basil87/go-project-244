package formatters

import "code/diff"

// GetFormatter returns the formatter function for the given format name.
// Recognised values: "plain", "json". Any other value returns FormatStylish.
func GetFormatter(format string) func([]diff.DiffNode) string {
	switch format {
	case "plain":
		return FormatPlain
	case "json":
		return FormatJSON
	default:
		return FormatStylish
	}
}
