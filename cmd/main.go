package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var montreeeDirEnv = "MONTREEE"

var defaultMontreeeDir = "."

var montreeeWorkingDirEnv = "MONTREEE_WORKING_DIR"

var montreeeClasspath = "/lib/*"

var montreeeMainClass = "CLI"

var defaultJavaCmd = "java"

var posibleJavaHomeDirEnvs = []string{
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
	"JRE",
	"JVM_HOME",
	"JRE",
	"JDK_HOME",
}

var posibleJavaExecuteables = []string{
	"/bin/java",
	"/jre/bin/java",
}

func main() {

	javaCmd := resolveJavaCmd()

	err := checkJavaVersion(javaCmd)
	if err != nil {
		printJavaConfigurationNotFound()
		return
	}

	montreeeDir := defaultMontreeeDir

	if dir, isSet := os.LookupEnv(montreeeDirEnv); isSet {
		montreeeDir = dir
	}

	//todo decide on using a environment variable to add extra vm args
	args := append([]string{"-cp", montreeeDir + montreeeClasspath, montreeeMainClass}, os.Args[1:]...)
	montreeeCmd := exec.Command(javaCmd, args...)

	if dir, isSet := os.LookupEnv(montreeeWorkingDirEnv); isSet {
		montreeeCmd.Dir = dir
	}

	montreeeCmd.Stdout = os.Stdout
	montreeeCmd.Stderr = os.Stderr
	_ = montreeeCmd.Run()
}

func resolveJavaCmd() string {
	javaCmd, err := resolveJavaCmdFromEnvironment()
	if err != nil {
		javaCmd = defaultJavaCmd
	}
	return javaCmd
}

func resolveJavaCmdFromEnvironment() (string, error) {
	for _, env := range posibleJavaHomeDirEnvs {

		dir, isSet := os.LookupEnv(env)
		if !isSet {
			continue
		}

		for _, binDir := range posibleJavaExecuteables {

			javaCmdToCkeck := filepath.FromSlash(dir + binDir)

			err := checkJavaVersion(javaCmdToCkeck)
			if err != nil {
				continue
			}

			return javaCmdToCkeck, nil
		}
	}
	return "", NewErrFailedToResolveJavaCommandFromEnvironment()
}

type ErrFailedToResolveJavaCommandFromEnvironment struct{}

func NewErrFailedToResolveJavaCommandFromEnvironment() *ErrFailedToResolveJavaCommandFromEnvironment {
	return &ErrFailedToResolveJavaCommandFromEnvironment{}
}

func (e ErrFailedToResolveJavaCommandFromEnvironment) Error() string {
	return "Failed to resolve Java command from environment."
}

func checkJavaVersion(javaCmd string) error {
	cmd := exec.Command(javaCmd, "--version")
	return cmd.Run()
}

func printJavaConfigurationNotFound() {
	fmt.Println("Failed to find a valid java installation.")
	fmt.Println("Try to add java to the path or set one of the following environment variables.")
	fmt.Println()
	for _, env := range posibleJavaHomeDirEnvs {
		fmt.Println(env)
	}
}
