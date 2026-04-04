package formatters

import "code/diff"

func GetFormatter(format string) func([]diff.DiffNode) string {
	switch format {
	case "plain":
		return FormatPlain
	default:
		return FormatStylish
	}
}
