package main

import "github.com/sethvargo/go-githubactions"

const apiTokenInput = "api_token"

func main() {
	gha := githubactions.New()

	token := gha.GetInput(apiTokenInput)
	if token == "" {
		gha.Fatalf("%s required", apiTokenInput)
	}

	gha.Infof("Hello, world!")
	gha.Infof("%s=%s", apiTokenInput, token)
}
