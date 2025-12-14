package prompt

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
)

type Inputs struct {
	URL          string
	Method       string
	ContentType  string
	Headers      []string
	OutputChoice string
	Body         string
}

func Run(current Inputs) (Inputs, error) {
	urlValue := defaultString(current.URL, "https://")
	methodChoice := pickDefault([]string{"GET", "POST", "PUT", "PATCH", "DELETE"}, strings.TrimSpace(current.Method), "GET")
	contentType := defaultString(strings.TrimSpace(current.ContentType), "application/json")
	headersText := strings.Join(current.Headers, "\n")
	bodyText := current.Body
	outputChoice := pickDefault([]string{
		"Body only (default)",
		"Status + response headers + body (-i)",
		"Verbose (-v)",
	}, strings.TrimSpace(current.OutputChoice), "Body only (default)")

	group := huh.NewGroup(
		buildURLInput(&urlValue),
		buildMethodSelect(&methodChoice),
		buildContentTypeInput(&contentType),
		buildHeadersField(&headersText),
		buildBodyField(&bodyText),
		buildOutputSelect(&outputChoice),
	)

	km := huh.NewDefaultKeyMap()
	km.Text.NewLine = key.NewBinding(key.WithKeys("shift+enter", "ctrl+j"), key.WithHelp("shift+enter / ctrl+j", "new line"))

	form := huh.NewForm(group).WithKeyMap(km)
	if err := runForm(form); err != nil {
		return Inputs{}, err
	}

	return Inputs{
		URL:          strings.TrimSpace(urlValue),
		Method:       strings.TrimSpace(methodChoice),
		ContentType:  strings.TrimSpace(contentType),
		Headers:      splitAndClean(headersText),
		Body:         strings.TrimSpace(bodyText),
		OutputChoice: strings.TrimSpace(outputChoice),
	}, nil
}

func ConfirmYes(message string) (bool, error) {
	var ok bool
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(message).
				Value(&ok),
		),
	)
	return ok, runForm(form)
}

func ConfirmWithEdit(copySupported bool) (string, error) {
	options := []string{
		"Run",
		"Edit inputs",
	}
	if copySupported {
		options = append(options, "Copy command & exit")
	}
	options = append(options, "Cancel")
	choice := options[0]
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Proceed?").
				Options(opts(options)...).
				Value(&choice),
		),
	)
	if err := runForm(form); err != nil {
		return "", err
	}
	return choice, nil
}

func required(v string) error {
	if strings.TrimSpace(v) == "" {
		return fmt.Errorf("required")
	}
	return nil
}

func opts(options []string) []huh.Option[string] {
	optObjs := make([]huh.Option[string], 0, len(options))
	for _, opt := range options {
		optObjs = append(optObjs, huh.NewOption(opt, opt))
	}
	return optObjs
}

func splitAndClean(text string) []string {
	var out []string
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		out = append(out, line)
	}
	return out
}

func pickDefault(options []string, current, fallback string) string {
	cur := strings.TrimSpace(current)
	for _, opt := range options {
		if cur == opt {
			return cur
		}
	}
	return fallback
}

func defaultString(val, fallback string) string {
	if strings.TrimSpace(val) == "" {
		return fallback
	}
	return val
}

func runForm(f *huh.Form) error {
	return f.WithTheme(huh.ThemeCharm()).Run()
}

func buildURLInput(value *string) *huh.Input {
	return huh.NewInput().
		Title("URL *").
		Prompt("").
		Value(value).
		Validate(required)
}

func buildMethodSelect(value *string) *huh.Select[string] {
	return huh.NewSelect[string]().
		Title("HTTP method *").
		Options(opts([]string{"GET", "POST", "PUT", "PATCH", "DELETE"})...).
		Value(value)
}

func buildContentTypeInput(value *string) *huh.Input {
	return huh.NewInput().
		Title("Content-Type").
		Prompt("").
		Value(value)
}

func buildHeadersField(value *string) *huh.Text {
	return huh.NewText().
		Title("Headers").
		Placeholder("Key: value").
		Lines(3).
		ExternalEditor(false).
		Value(value)
}

func buildBodyField(value *string) *huh.Text {
	return huh.NewText().
		Title("Body").
		Lines(3).
		ExternalEditor(false).
		Value(value)
}

func buildOutputSelect(value *string) *huh.Select[string] {
	return huh.NewSelect[string]().
		Title("Output *").
		Options(opts([]string{
			"Body only (default)",
			"Status + response headers + body (-i)",
			"Verbose (-v)",
		})...).
		Value(value)
}

func CopyToClipboard(text string) error {
	return clipboard.WriteAll(text)
}

func CopySupported() bool {
	return !clipboard.Unsupported
}
