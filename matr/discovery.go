package matr

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/matr-builder/matr/parser"
)

const (
	defaultMatrFile    = "Matrfile.go"
	defaultCacheFolder = ".matr"
)

var (
	matrFilePath string
	helpFlag     bool
	versionFlag  bool
)

// Run is the primary entrypoint to matrs cli tool.
// This is where the matrfile path is resolved, compiled and executed
func Run() {
	// TODO: clean up this shit show
	flag.StringVar(&matrFilePath, "matrfile", "./", "path to Matrfile")
	flag.BoolVar(&helpFlag, "h", false, "Display usage info")
	flag.BoolVar(&versionFlag, "v", false, "Display version")
	flag.Parse()
	if versionFlag {
		fmt.Println(Version)
		return
	}

	args := flag.Args()

	if helpFlag {
		args = append([]string{"-h"}, args...)
	}

	cmds, err := parseMatrfile(matrFilePath)
	if err != nil {
		flag.Usage()
		if helpFlag && flag.Arg(0) == "" {
			fmt.Print("\nTargets:\n  No Matrfile.go or Matrfile found\n")
			return
		}

		fmt.Print("\n  " + err.Error() + "\n")
		return
	}

	matrCachePath, err := build(matrFilePath, cmds)
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		return
	}

	if err := run(matrCachePath, args...); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		return
	}
}

func parseMatrfile(path string) ([]parser.Command, error) {
	var err error
	var cmds []parser.Command

	matrFilePath, err = filepath.Abs(matrFilePath)
	if err != nil {
		return cmds, err
	}

	matrFilePath, err = getMatrfilePath(matrFilePath)
	if err != nil {
		return cmds, err
	}

	cmds, err = parser.Parse(matrFilePath)
	if err != nil {
		return cmds, err
	}

	return cmds, nil
}

func run(matrCachePath string, args ...string) error {
	c := exec.Command(filepath.Join(matrCachePath, "matr"), args...)
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	return c.Run()
}

func build(matrCachePath string, cmds []parser.Command) (string, error) {
	var b bytes.Buffer

	matrPath, matrFile := filepath.Split(matrFilePath)
	matrCachePath = filepath.Join(matrPath, ".matr")

	if dir, err := os.Stat(matrCachePath); err != nil || !dir.IsDir() {
		if err := os.Mkdir(matrCachePath, 0777); err != nil {
			return "", err
		}
	}

	f, err := os.OpenFile(filepath.Join(matrCachePath, "main.go"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if !symlinkValid(matrCachePath) {
		os.Remove(filepath.Join(matrCachePath, defaultMatrFile))
		if err := os.Symlink(filepath.Join(matrPath, matrFile), filepath.Join(matrCachePath, defaultMatrFile)); err != nil {
			if os.IsExist(err) {
				return "", err
			}
		}
	}

	if err := generate(cmds, &b); err != nil {
		return "", err
	}

	io.Copy(f, &b)
	// TODO: check if we need to rebuild
	cmd := exec.Command("go", "build", "-tags", "matr", "-o", filepath.Join(matrCachePath, "matr"),
		filepath.Join(matrCachePath, "Matrfile.go"),
		filepath.Join(matrCachePath, "main.go"),
	)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return matrCachePath, cmd.Run()
}

func getMatrfilePath(matrFilePath string) (string, error) {
	matrFilePath, err := filepath.Abs(matrFilePath)
	if err != nil {
		return "", err
	}

	fp, err := os.Stat(matrFilePath)
	if err != nil {
		return "", errors.New("unable to find Matrfile: " + matrFilePath)
	}

	if !fp.IsDir() {
		return matrFilePath, nil
	}

	matrFilePath = filepath.Join(matrFilePath, "Matrfile")

	if _, err = os.Stat(matrFilePath + ".go"); err == nil {
		return matrFilePath + ".go", nil
	}

	if _, err := os.Stat(matrFilePath); err == nil {
		return matrFilePath, nil
	}

	return "", errors.New("unable to find Matrfile")
}

func symlinkValid(path string) bool {
	pth, err := os.Readlink(filepath.Join(path, "Matrfile.go"))
	if err != nil {
		return false
	}

	if _, err := os.Stat(pth); err != nil {
		return false
	}
	return true
}
