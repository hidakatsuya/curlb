package curlcmd

import (
	"bytes"
	"encoding/json"
	"strings"
)

func CommandPreview(cmd string, args []string) string {
	escaped := make([]string, 0, len(args)+1)
	escaped = append(escaped, ShellEscape(cmd))
	for _, arg := range args {
		escaped = append(escaped, ShellEscape(arg))
	}
	return strings.Join(escaped, " ")
}

func ShellEscape(s string) string {
	if s == "" {
		return "''"
	}
	if strings.IndexFunc(s, func(r rune) bool {
		return r == ' ' || r == '"' || r == '\'' || r == '\\' || r == '$' || r == '`' || r == '|' || r == '&' || r == ';' || r == '<' || r == '>' || r == '(' || r == ')' || r == '{' || r == '}'
	}) == -1 {
		return s
	}
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

func OutputArgsFromChoice(choice string) []string {
	switch choice {
	case "Status + response headers + body (-i)":
		return []string{"-i"}
	case "Verbose (-v)":
		return []string{"-v"}
	default:
		return nil
	}
}

func NormalizeBody(body string) string {
	body = strings.TrimSpace(body)
	if body == "" {
		return ""
	}

	var buf bytes.Buffer
	if err := json.Compact(&buf, []byte(body)); err == nil {
		return buf.String()
	}

	parts := strings.Fields(body)
	return strings.Join(parts, " ")
}
