package main

import (
	"os"
	"os/exec"
)

func main() {
	args := append([]string{"-cp", "lib/*", "CLI"}, os.Args[1:]...)
	startJvmCommand := exec.Command("java", args...)
	startJvmCommand.Stdout = os.Stdout
	startJvmCommand.Stderr = os.Stderr
	_ = startJvmCommand.Run()
}
