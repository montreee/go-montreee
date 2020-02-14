package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func checkJavaVersion(javaCommand string) error {
	cmd := exec.Command(javaCommand, "--version")
	return cmd.Run()
}

func main() {

	posibleJavaHomeDirEnvs := []string{"MONTREEE_JVM", "JVM", "JDK", "JAVA_HOME", "JRE_HOME"}

	javaCommand := "java"

	for _, env := range posibleJavaHomeDirEnvs {

		dir, isSet := os.LookupEnv(env)
		if !isSet {
			continue
		}

		javaCmdToCkeck := filepath.FromSlash(dir + "/bin/java")

		err := checkJavaVersion(javaCmdToCkeck)
		if err != nil {
			continue
		}

		javaCommand = javaCmdToCkeck

		break
	}

	err := checkJavaVersion(javaCommand)
	if err != nil {
		fmt.Println(err)
		return
	}

	args := append([]string{"-cp", "lib/*", "CLI"}, os.Args[1:]...)
	startJvmCommand := exec.Command(javaCommand, args...)
	startJvmCommand.Stdout = os.Stdout
	startJvmCommand.Stderr = os.Stderr
	_ = startJvmCommand.Run()
}
