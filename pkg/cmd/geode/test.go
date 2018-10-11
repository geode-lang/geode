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
	"path"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"
	"unicode"

	"github.com/BurntSushi/toml"
	"github.com/geode-lang/geode/pkg/arg"
	"github.com/geode-lang/geode/pkg/util"
	"github.com/geode-lang/geode/pkg/util/color"
	"github.com/geode-lang/geode/pkg/util/log"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// TestJob -
type TestJob struct {
	Name, sourcefile                         string
	CompilerArgs, RunArgs                    []string
	compilerError, RunStatus, CompilerStatus int
	Input                                    string
	compilerOutput, RunOutput                string
}

type testResult struct {
	TestJob        TestJob
	compilerError  int
	RunStatus      int
	CompilerStatus int
	compilerOutput string
	RunOutput      string
	timetaken      time.Duration
}

func parseTestJob(filename string) (TestJob, error) {
	var job TestJob

	configPath := path.Join(path.Dir(filename), "test.toml")

	if _, err := toml.DecodeFile(configPath, &job); err != nil {
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
			job.sourcefile = path

			jobs = append(jobs, job)
		}
	}

	// Do jobs

	results := make(chan testResult, len(jobs))

	testRunCount := 0

	util.RunCommand("geode", "clean")

	go func() {
		for _, job := range jobs {

			start := time.Now()
			outBuf := new(bytes.Buffer)
			outpath := fmt.Sprintf("%s_test", job.sourcefile)

			// Compile the test program
			buildArgs := []string{"build"}
			buildArgs = append(buildArgs, job.CompilerArgs...)
			buildArgs = append(buildArgs, "-o", outpath, job.sourcefile)

			outBuf.Reset()

			var err error
			res := testResult{TestJob: job}

			res.compilerError, err = runCommand(outBuf, "", "geode", buildArgs)
			if err != nil {
				fmt.Printf("Error while building test:\n%s\n", err.Error())
				os.Exit(1)
			}
			res.compilerOutput = outBuf.String()

			if res.compilerError != 0 {
				results <- res
				res.RunStatus = -1
				return
			}

			// Run the test program
			outBuf.Reset()

			res.RunStatus, err = runCommand(outBuf, job.Input, fmt.Sprintf("./%s", outpath), job.RunArgs)
			if err != nil {
				fmt.Printf("Error while running test:\n%s\n", err.Error())
				os.Exit(1)
			}

			res.RunOutput = outBuf.String()

			// Remove test executable
			if err := os.Remove(outpath); err != nil {
				fmt.Printf("Error while removing test executable:\n%s\n", err.Error())
				os.Exit(1)
			}

			t := time.Now()

			elapsed := t.Sub(start)

			res.timetaken = elapsed
			results <- res
			testRunCount++

			if testRunCount == len(jobs) {
				close(results)
			}
		}
	}()

	// Check results
	numSucceses := 0
	numTests := 0

	index := 0

	for res := range results {
		index++
		failure := false

		errBuf := &bytes.Buffer{}

		// Check build errors

		if (res.TestJob.compilerError == -1 && res.compilerError != res.TestJob.CompilerStatus) || (res.compilerError == res.TestJob.compilerError) {
		} else {
			fmt.Fprintf(errBuf, "CompilerStatus:\n")
			fmt.Fprintf(errBuf, "Expected: %d\n", res.TestJob.compilerError)
			fmt.Fprintf(errBuf, "Got:      %d\n", res.TestJob.CompilerStatus)
			failure = true
		}

		// Check run errors
		if res.RunStatus == res.TestJob.RunStatus {
		} else {
			fmt.Fprintf(errBuf, "RunStatus:\n")
			fmt.Fprintf(errBuf, "Expected: %d\n", res.TestJob.RunStatus)
			fmt.Fprintf(errBuf, "Got:      %d\n", res.RunStatus)
			failure = true
		}

		// Check run errors
		if res.RunOutput == res.TestJob.RunOutput {
		} else {

			dmp := diffmatchpatch.New()
			expected := res.TestJob.RunOutput
			got := res.RunOutput
			diffs := dmp.DiffMain(expected, got, false)

			fmt.Fprintf(errBuf, "RunOutput:\n")
			fmt.Fprintf(errBuf, "Expected: %q\n", expected)
			fmt.Fprintf(errBuf, "Got:      %q\n", got)

			fmt.Fprintf(errBuf, "diff:\n%s\n", dmp.DiffPrettyText(diffs))
			failure = true
		}

		ok := fmt.Sprintf("%sOKAY%s", color.TEXT_GREEN, color.TEXT_RESET)

		// Output result
		if !failure {
			fmt.Printf("(%d)\t%s %s\n", index, ok, res.TestJob.Name)
			numSucceses++
			// fmt.Printf("%sTest Passed%s\n", color.TEXT_GREEN, color.TEXT_RESET)
		} else {
			failed := fmt.Sprintf("%sFAIL%s", color.TEXT_RED, color.TEXT_RESET)
			fmt.Printf("(%d)\t%s %s\n", index, failed, res.TestJob.Name)
			fmt.Printf("%s\n", errBuf.String())
		}
		numTests++

	}

	fmt.Printf("->\t%d/%d (%.0f%%) tests ran succesfully\n\n", numSucceses, numTests, float64(numSucceses)/float64(numTests)*100)
	if numSucceses < numTests {
		os.Exit(1)
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

var testTemplate = `# {{NAME}}
is main

func main int {
	# Test code here.
	return 0;
}
`

func isValidPathRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

func cleanTestName(name string) string {
	src := []rune(name)
	runes := make([]rune, 0, len(src))

	for _, r := range src {
		if isValidPathRune(r) {
			runes = append(runes, r)
		} else if r == ' ' {
			runes = append(runes, '-')
		}
	}
	return string(runes)
}

// CreateTestCMD uese the os args to create a test
func CreateTestCMD() {
	name := cleanTestName(*arg.NewTestName)
	dirPath := fmt.Sprintf("./tests/%s", name)
	sourcePath := path.Join(dirPath, fmt.Sprintf("%s.g", name))

	configPath := path.Join(dirPath, "test.toml")

	stats, _ := os.Stat(dirPath)
	if stats != nil && stats.IsDir() {
		log.Fatal("Test %q already exists\n", name)
	}
	os.MkdirAll(dirPath, os.ModePerm)

	fileContent := strings.Replace(testTemplate, "{{NAME}}", *arg.NewTestName, -1)

	ioutil.WriteFile(sourcePath, []byte(fileContent), os.ModePerm)

	// Write the config

	job := TestJob{}
	job.Name = *arg.NewTestName

	buff := &bytes.Buffer{}
	toml.NewEncoder(buff).Encode(job)

	ioutil.WriteFile(configPath, buff.Bytes(), os.ModePerm)

	fmt.Printf("Test %q created in %q\n", name, sourcePath)
}
