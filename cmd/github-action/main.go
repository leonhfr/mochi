package main

import "github.com/sethvargo/go-githubactions"

func main() {
	gha := githubactions.New()
	gha.Infof("Hello, world!")
}
