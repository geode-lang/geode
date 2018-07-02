// This test framework is taken from
// https://github.com/ark-lang/ark/tree/master/tests
// It was just too good not to use.

package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/nickwanninger/geode/pkg/util/color"
)

// TestJob -
type TestJob struct {
	Name, Sourcefile               string
	CompilerArgs, RunArgs          []string
	CompilerError, RunError        int
	Input                          string
	CompilerOutput, ExpectedOutput string
}

type testResult struct {
	TestJob        TestJob
	CompilerError  int
	RunError       int
	CompilerOutput string
	ExpectedOutput string
}

func parseTestJob(filename string) (TestJob, error) {
	var job TestJob

	dataBytes, ferr := ioutil.ReadFile(filename)
	if ferr != nil {
		panic(ferr)
	}

	rawLines := strings.Split(string(dataBytes), "\n")
	configLines := make([]string, 0, len(rawLines))

	for _, line := range rawLines {
		if len(line) > 0 && line[0] == '#' {
			configLines = append(configLines, strings.Trim(line, "# "))
		}
	}

	config := strings.Join(configLines, "\n")

	if _, err := toml.Decode(config, &job); err != nil {
		return TestJob{}, err
	}
	return job, nil
}

// RunTests runs all the tests in some directory
func RunTests(testDirectory string) int {
	var dirs []string
	files := make(map[string][]string)

	// Find all toml files in test directory
	filepath.Walk(testDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relpath, err := filepath.Rel(testDirectory, path)
		if err != nil {
			return err
		}

		dir, file := filepath.Split(relpath)
		if info.IsDir() {
			if file == "." {
				file = ""
			} else {
				file += "/"
			}

			dirs = append(dirs, file)
		} else if strings.HasSuffix(file, ".g") {
			files[dir] = append(files[dir], file)
		}
		return nil
	})

	// Sort directories
	sort.Strings(dirs)

	var jobs []TestJob
	for _, dir := range dirs {
		// Sort files
		sort.Strings(files[dir])

		for _, file := range files[dir] {
			path := filepath.Join(testDirectory, dir, file)

			// Parse job file
			job, err := parseTestJob(path)
			if err != nil {
				fmt.Printf("%s\n", err.Error())
				return 1
			}
			job.Sourcefile = path

			jobs = append(jobs, job)
		}
	}

	// Do jobs

	results := make(chan testResult)

	go func() {
		outBuf := new(bytes.Buffer)
		for _, job := range jobs {
			outpath := fmt.Sprintf("%s_test", job.Sourcefile)

			// Compile the test program
			buildArgs := []string{"build"}
			buildArgs = append(buildArgs, job.CompilerArgs...)
			buildArgs = append(buildArgs, []string{"-o", outpath, job.Sourcefile}...)

			outBuf.Reset()

			var err error
			res := testResult{TestJob: job}

			res.CompilerError, err = runCommand(outBuf, "", "geode", buildArgs)
			if err != nil {
				fmt.Printf("Error while building test:\n%s\n", err.Error())
				os.Exit(1)
			}
			res.CompilerOutput = outBuf.String()

			if res.CompilerError != 0 {
				results <- res
				res.RunError = -1
				continue
			}

			// Run the test program
			outBuf.Reset()

			res.RunError, err = runCommand(outBuf, job.Input, fmt.Sprintf("./%s", outpath), job.RunArgs)
			if err != nil {
				fmt.Printf("Error while running test:\n%s\n", err.Error())
				os.Exit(1)
			}
			res.ExpectedOutput = outBuf.String()

			// Remove test executable
			if err := os.Remove(outpath); err != nil {
				fmt.Printf("Error while removing test executable:\n%s\n", err.Error())
				os.Exit(1)
			}

			results <- res
		}
		close(results)
	}()

	// Check results
	numSucceses := 0
	numTests := 0

	for res := range results {
		failure := false

		fmt.Printf("%s\n", res.TestJob.Name)

		// Check build errors

		if (res.TestJob.CompilerError == -1 && res.CompilerError != 0) || (res.CompilerError == res.TestJob.CompilerError) {
		} else {
			fmt.Printf("  CompilerError:\n    ")
			msg := color.Red("✗")
			fmt.Printf("%s. Expected %d\n    ", msg, res.TestJob.CompilerError)
			fmt.Printf("Got %d\n", res.CompilerError)
			failure = true
		}

		// Check run errors
		if res.CompilerOutput == res.TestJob.CompilerOutput {
		} else {
			fmt.Printf("  CompilerOutput:\n    ")
			msg := color.Red("✗")
			fmt.Printf("%s. Expected %s\n    ", msg, res.TestJob.CompilerOutput)
			fmt.Printf("Got %s\n", res.CompilerOutput)
			failure = true
		}

		// Check run errors
		if res.RunError == res.TestJob.RunError {
		} else {
			fmt.Printf("  RunError:\n    ")
			msg := color.Red("✗")
			fmt.Printf("%s. Expected %d\n    ", msg, res.TestJob.RunError)
			fmt.Printf("Got %d\n", res.RunError)
			failure = true
		}

		// Check run errors
		if res.ExpectedOutput == res.TestJob.ExpectedOutput {
		} else {
			fmt.Printf("  ExpectedOutput:\n    ")
			msg := color.Red("✗")
			fmt.Printf("%s. Expected %s\n    ", msg, res.TestJob.ExpectedOutput)
			fmt.Printf("Got %s\n", res.ExpectedOutput)
			failure = true
		}

		// Output result
		if !failure {
			numSucceses++
			fmt.Printf("%sTest Passed%s\n", color.TEXT_GREEN, color.TEXT_RESET)
		} else {
			fmt.Printf("%sTest Failed%s\n", color.TEXT_RED, color.TEXT_RESET)
		}
		numTests++

		fmt.Printf("\n")
	}

	fmt.Printf("Total: %d / %d tests ran succesfully\n\n", numSucceses, numTests)
	if numSucceses < numTests {
		return 1
	}
	return 0
}

func runCommand(out io.Writer, input string, cmd string, args []string) (int, error) {
	// Run the test program
	command := exec.Command(cmd, args...)
	command.Stdin = strings.NewReader(input)

	// Output handling
	ow := out
	command.Stdout, command.Stderr = ow, ow

	// Disable coloring for matching compiler output
	command.Env = append(os.Environ(), "COLOR=0")

	// Start the test
	if err := command.Start(); err != nil {
		return -1, err
	}

	// Check the exit status
	if err := command.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				return status.ExitStatus(), nil
			}
		} else {
			return -1, err
		}
	}

	return 0, nil
}
