package app

import (
	"errors"
	"fmt"
	"strings"

	"curlb/internal/curlcmd"
	"curlb/internal/prompt"
	"curlb/internal/runner"
	"github.com/charmbracelet/huh"
)

func Run(extraCurlOpts []string) error {
	inputs := prompt.Inputs{
		Method:       "GET",
		ContentType:  "application/json",
		OutputChoice: "Body only (default)",
	}
	var err error
	if inputs, err = prompt.Run(inputs); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			fmt.Println("Canceled.")
			return nil
		}
		return err
	}

	for {
		method := strings.ToUpper(strings.TrimSpace(inputs.Method))
		headers := append([]string{}, inputs.Headers...)
		if strings.TrimSpace(inputs.ContentType) != "" {
			headers = append([]string{fmt.Sprintf("Content-Type: %s", strings.TrimSpace(inputs.ContentType))}, headers...)
		}
		outputArgs := curlcmd.OutputArgsFromChoice(inputs.OutputChoice)
		bodyArg := curlcmd.NormalizeBody(inputs.Body)

		args := []string{}
		if method != "" {
			args = append(args, "-X", method)
		}

		for _, h := range headers {
			args = append(args, "-H", h)
		}

		if bodyArg != "" {
			args = append(args, "--data", bodyArg)
		}

		args = append(args, outputArgs...)

		// User-specified options go last to allow overrides.
		args = append(args, extraCurlOpts...)

		if inputs.URL != "" {
			args = append(args, inputs.URL)
		}

		previewOneLine := curlcmd.CommandPreview("curl", args)
		fmt.Println(previewOneLine)

		choice, err := prompt.ConfirmWithEdit(prompt.CopySupported())
		if err != nil {
			if errors.Is(err, huh.ErrUserAborted) {
				fmt.Println("Canceled.")
				return nil
			}
			return err
		}
		switch choice {
		case "Run":
			return runner.RunCurl(args)
		case "Edit inputs":
			inputs, err = prompt.Run(inputs)
			if err != nil {
				if errors.Is(err, huh.ErrUserAborted) {
					fmt.Println("Canceled.")
					return nil
				}
				return err
			}
			continue
		case "Copy command & exit":
			if err := prompt.CopyToClipboard(previewOneLine); err != nil {
				fmt.Println("Clipboard unavailable; here is the command:")
				fmt.Println(previewOneLine)
				return nil
			}
			fmt.Println("Copied command.")
			return nil
		default:
			fmt.Println("Cancelled.")
			return nil
		}
	}
}
