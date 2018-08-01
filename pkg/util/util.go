package util

import (
	"fmt"
	"io/ioutil"
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

	tmpcmd := command + " " + strings.Join(args, " ")
	maxLen := 500
	if len(tmpcmd) > maxLen {
		tmpcmd = tmpcmd[:maxLen-3] + "..."
	}
	title := fmt.Sprintf("Command Execution (%s)", tmpcmd)
	log.Timed(title, func() {
		cmd := exec.Command(command, args...)
		out, err = cmd.CombinedOutput()
	})

	if err != nil {
		return out, err
	}
	return out, err
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
		pth, err := ioutil.TempDir(path.Join(HomeDir(), ".geode_cache"), "")
		if err != nil {
			log.Fatal("Unable to get temp directory\n")
		}

		tmpdir = pth
	}

	return tmpdir

}

// PurgeCache -
func PurgeCache() {
	cacheDir := GetCacheDir()

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
