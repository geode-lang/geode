package util

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/geode-lang/geode/pkg/util/log"
	"github.com/shibukawa/configdir"
)

// RunCommand executes a command and returns stdout from it.
func RunCommand(command string, args ...string) ([]byte, error) {
	var out []byte
	var err error

	// fmt.Printf("%s %s\n", command, strings.Join(args, " "))

	log.Verbose("%s %s\n", command, strings.Join(args, " "))

	tmpcmd := command + " " + strings.Join(args, " ")
	maxLen := 500
	if len(tmpcmd) > maxLen {
		tmpcmd = tmpcmd[:maxLen-3] + "..."
	}
	title := fmt.Sprintf("Command Execution (%s)", tmpcmd)
	fullcommand := fmt.Sprintf("%s %s", command, strings.Join(args, " "))
	log.Timed(title, func() {
		cmd := exec.Command("bash", "-c", fullcommand)
		out, err = cmd.CombinedOutput()
	})

	if err != nil {
		return out, err
	}
	return out, err
}

// RunCommandStr is a wrapper around RunCommand that returns a string instead
func RunCommandStr(command string, args ...string) (string, error) {
	b, e := RunCommand(command, args...)
	return string(b), e
}

// BashCmd runs a command in a bash context
func BashCmd(command string) (string, error) {
	return RunCommandStr("bash", "-c", fmt.Sprintf("\"%s\"", command))
}

// StdLibDir returns the stdlib directory path
func StdLibDir() string {
	libpath := os.Getenv("GEODELIB")
	if libpath == "" {
		libpath = "/usr/local/lib/geodelib"
	}
	return libpath
}

// StdLibFile takes a path in the stdlib and
// joins it to the directory path
func StdLibFile(p string) string {
	return path.Join(StdLibDir(), p)
}

// HomeDir will return the home directory of the current user.
func HomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

// GetCacheDir get's the OS Specific cache directory
func GetCacheDir() string {
	configDirs := configdir.New("nick-wanninger", "geode-lang")
	cache := configDirs.QueryCacheFolder()
	return cache.Path
}

var tmpdir string

// GetTmp returns a temporary directory
func GetTmp() string {

	if tmpdir == "" {
		pth, err := ioutil.TempDir(path.Join(HomeDir(), ".geode/tmp"), "")
		if err != nil {
			log.Fatal("Unable to get temp directory\n")
		}
		os.MkdirAll(pth, os.ModePerm)

		tmpdir = pth
	}

	return tmpdir

}

// PurgeCache -
func PurgeCache() {
	cacheDir := GetCacheDir()

	os.MkdirAll(cacheDir, os.ModePerm)

	files, _ := ioutil.ReadDir(cacheDir)
	// if err != nil {
	// 	log.Fatal("Unable to search cache for files\n")
	// }

	now := time.Now()

	cacheInvalidationTimeout := 10 * time.Minute

	for _, f := range files {
		if now.Sub(f.ModTime()) > cacheInvalidationTimeout {
			os.Remove(path.Join(cacheDir, f.Name()))
		}
	}

}

// RandomHex returns a random hex string of n bytes in length
func RandomHex(n int) string {
	bytes := make([]byte, n)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// HashFile takes a path and hashes it efficiently into sha1
func HashFile(path string) string {
	var returnMD5String string
	file, err := os.Open(path)
	if err != nil {
		return returnMD5String
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String
	}
	hashInBytes := hash.Sum(nil)[:16]
	returnMD5String = hex.EncodeToString(hashInBytes)
	return returnMD5String

}

// QuickHash returns the first couple chars from a sha256
func QuickHash(in string, l int) string {
	if l > 64 {
		l = 64
	}
	return fmt.Sprintf("%x", sha256.Sum256([]byte(in)))[:l]
}

// EatError takes an error and if it is nil, ignores it.
// otherwise it is a fatal error
func EatError(e error) {
	if e != nil {
		log.Fatal("%s\n", e)
	}
}

// Jaro returns the jaro distance between two strings
func Jaro(a, b string) float64 {
	la := float64(len(a))
	lb := float64(len(b))

	// match range = max(len(a), len(b)) / 2 - 1
	matchRange := int(math.Floor(math.Max(la, lb)/2.0)) - 1
	matchRange = int(math.Max(0, float64(matchRange-1)))
	var matches, halfs float64
	transposed := make([]bool, len(b))

	for i := 0; i < len(a); i++ {
		start := int(math.Max(0, float64(i-matchRange)))
		end := int(math.Min(lb-1, float64(i+matchRange)))

		for j := start; j <= end; j++ {
			if transposed[j] {
				continue
			}

			if a[i] == b[j] {
				if i != j {
					halfs++
				}
				matches++
				transposed[j] = true
				break
			}
		}
	}

	if matches == 0 {
		return 0
	}

	transposes := math.Floor(float64(halfs / 2))

	return ((matches / la) + (matches / lb) + (matches-transposes)/matches) / 3.0
}
