package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

const (
	integrationDir = "integration"
	buildDir       = "test"
	instructions   = "stack[Widget]"
)

func main() {
	wd, err := os.Getwd()
	logFailureOnError("finding working directory", err)

	buildInIntegrationDir := path.Join(integrationDir, buildDir)
	_, err = os.Stat(buildInIntegrationDir)
	if logFailureOnError(fmt.Sprintf("getting statistics on %s", buildInIntegrationDir), err) {
		logFailureOnError("creating build directory", os.MkdirAll(buildInIntegrationDir, 0755))
	}

	build := exec.Command("go", "build", "-o", "integration/test/build", ".")
	logFailureOnError("executing integration build", build.Start())
	logFailureOnError("completing integration build", build.Wait())

	logFailureOnError(fmt.Sprintf("changing to test directory (%s)", buildInIntegrationDir),
		os.Chdir(wd))
	run := exec.Command(filepath.Join(buildInIntegrationDir)+"/build", "-pathTo", "integration/example/widgetStack.go", instructions)
	run.Dir = wd
	stdout, err := run.StdoutPipe()
	if err != nil {
		log.Fatalf("integration: stdout pipe unavailable: %v", err)
	}

	stderr, err := run.StderrPipe()
	if err != nil {
		log.Fatalf("integration: stderr pipe unavailable: %v", err)
	}

	logFailureOnError("executing integration test", run.Start())
	logFrom(stdout)
	logFrom(stderr)
	logFailureOnError("completing integration test", run.Wait())
}

func logFailureOnError(msg string, err error) bool {
	if err != nil {
		log.Printf("integration: failed %s: %v", msg, err)
		return true
	}
	return false
}

func logFrom(r io.ReadCloser) {
	b := make([]byte, 0, 1024)
	r.Read(b)
	out := "<empty>"
	if strings.Trim(string(b), " ") != "" {
		out = fmt.Sprintf("bytes(%d)[%v]: ", len(b), b) + string(b)
	}
	log.Println(out)
}
