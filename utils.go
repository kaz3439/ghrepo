package main

import (
  "fmt"
  "strconv"
  "bufio"
  "os"
)

func PromptPersonalGithubToken() string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("your personal github token: ")
	scanner.Scan()
	githubToken := scanner.Text()
	return githubToken
}

func GetRepositryField(name string, field interface{}, prompt bool) interface{} {
	if prompt == false {
		return field
	}

	var defaultPrompt string
	switch field.(type) {
	case string:
		defaultPrompt = field.(string)
	case bool:
		defaultPrompt = strconv.FormatBool(field.(bool))
	case int:
		if field.(int) != 0 {
			defaultPrompt = strconv.FormatInt(int64(field.(int)), 10)
		} else {
			defaultPrompt = ""
		}
	default:
		defaultPrompt = ""
	}

	fmt.Printf("%s [%s] : ", name, defaultPrompt)
	scan := bufio.NewScanner(os.Stdin)
	scan.Scan()
	text := scan.Text()
	switch field.(type) {
	case bool:
		if res, err := strconv.ParseBool(text); err != nil {
			field = res
		}
	case int:
		if res, err := strconv.ParseInt(text, 10, 0); err != nil {
			field = int(res)
		}
	default:
		field = text
	}
	return field
}
