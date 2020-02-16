package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	envMontreeeDir        = "MONTREEE"
	defaultMontreeeDir    = "."
	envMontreeeWorkingDir = "MONTREEE_WORKING_DIR"
	montreeeClasspath     = "/lib/*"
	montreeeMainClass     = "CLI"
)

var (
	javaHomeDirEnvs = []string{
		"MONTREEE_JAVA",
		"MONTREEE_JAVA_HOME",
		"MONTREEE_JRE",
		"MONTREEE_JRE_HOME",
		"MONTREEE_JVM",
		"MONTREEE_JVM_HOME",
		"MONTREEE_JDK",
		"MONTREEE_JDK_HOME",
		"JAVA",
		"JAVA_HOME",
		"JRE",
		"JRE_HOME",
		"JVE",
		"JVM_HOME",
		"JDE",
		"JDK_HOME",
	}
	javaExecuteables = []string{
		"/bin/java",
		"/jre/bin/java",
	}
)

var ErrJavaResolveFromEnvironment = errors.New("failed to resolve Java command from environment")

func main() {
	javaCmd := resolveJavaCmd()

	err := exec.Command(javaCmd, "--version").Run() // check if the binary to this command exists
	if err != nil {
		printJavaConfigNotFound()

		return
	}

	montreeeDir := defaultMontreeeDir

	if dir, ok := os.LookupEnv(envMontreeeDir); ok {
		montreeeDir = dir
	}

	// todo decide on using a environment variable to add extra vm args
	args := append([]string{"-cp", montreeeDir + montreeeClasspath, montreeeMainClass}, os.Args[1:]...)
	montreeeCmd := exec.Command(javaCmd, args...)

	if dir, ok := os.LookupEnv(envMontreeeWorkingDir); ok {
		montreeeCmd.Dir = dir
	}

	montreeeCmd.Stdout = os.Stdout
	montreeeCmd.Stderr = os.Stderr

	_ = montreeeCmd.Run()
}

func resolveJavaCmd() string {
	javaCmd, err := resolveJavaCmdFromEnvironment()
	if err != nil {
		javaCmd = "java"
	}
	return javaCmd
}

func resolveJavaCmdFromEnvironment() (string, error) {
	for _, env := range javaHomeDirEnvs {
		dir, isSet := os.LookupEnv(env)
		if !isSet {
			continue
		}

		for _, binDir := range javaExecuteables {
			checkCmd := filepath.FromSlash(dir + binDir)

			err := exec.Command(checkCmd, "--version").Run()
			if err != nil {
				continue
			}

			return checkCmd, nil
		}
	}

	return "", ErrJavaResolveFromEnvironment
}

func printJavaConfigNotFound() {
	fmt.Println("failed to find a valid java installation")
	fmt.Println("try to add java to the path or set one of the following environment variables:")
	fmt.Println()
	fmt.Println(strings.Join(javaHomeDirEnvs, "\n"))
}
