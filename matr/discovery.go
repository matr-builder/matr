package matr

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"unicode"
	"unicode/utf8"

	"github.com/matr-builder/matr/parser"
)

const (
	defaultMatrFile    = "Matrfile.go"
	defaultCacheFolder = ".matr"
)

var (
	matrFilePath string
	helpFlag     bool
)

// Run is the primary entrypoint to matrs cli tool.
// This is where the matrfile path is resolved, compiled and executed
func Run() {
	// TODO: clean up this shit show
	flag.StringVar(&matrFilePath, "matrfile", "./", "path to Matrfile")
	flag.BoolVar(&helpFlag, "help", false, "Display usage info")

	cmds, err := parseMatrfile(matrFilePath)
	if err != nil {
		log.Fatal(err)
	}

	flag.Usage = usage(cmds)
	flag.Parse()

	cmd := flag.Arg(0)
	validCmd := false

	if cmd == "" {
		cmd = "default"
	}

	for _, c := range cmds {
		if strings.ToLower(c.Name) != cmd {
			continue
		}
		validCmd = true
		break
	}

	if helpFlag || !validCmd {
		if cmd != "default" {
			os.Stderr.WriteString("No handler found for target: " + cmd + "\n\n")
		}
		flag.Usage()
		return
	}

	matrCachePath, err := build(matrFilePath, cmds)
	if err != nil {
		log.Fatal(err)
	}

	if err := run(matrCachePath, flag.Args()...); err != nil {
		log.Fatal(err)
	}
}

func usage(cmds []parser.Command) func() {
	return func() {
		if cmd := flag.Arg(0); cmd != "" {
			for _, c := range cmds {
				if lowerFirst(c.Name) == cmd {
					fmt.Println("matr " + cmd + " :\n")
					fmt.Println(c.Doc)
					fmt.Print("\n")
					return
				}
			}
			os.Stderr.WriteString("No handler found for target: " + cmd + "\n\n")
		}

		fmt.Println("\nUsage: matr <opts> [target] args...")

		fmt.Println("\nOptions:")
		flag.PrintDefaults()

		fmt.Println("\nTargets:")
		tw := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		for _, cmd := range cmds {
			if !cmd.IsExported || cmd.Name == "Default" {
				continue
			}
			fmt.Fprintf(tw, "	%s\t%s\n", lowerFirst(cmd.Name), cmd.Summary)
		}
		tw.Flush()
		fmt.Println(" ")
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

func generate(cmds []parser.Command, w io.Writer) error {
	// Create a new template and parse the letter into it.
	t := template.Must(template.New("letter").Funcs(template.FuncMap{
		"cmdname": func(name string) string {
			s := strings.Replace(name, "_", ":", -1)
			r, n := utf8.DecodeRuneInString(s)
			return string(unicode.ToLower(r)) + s[n:]
		},
	}).Parse(defaultTemplate))
	return t.Execute(w, cmds)
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
		return "", errors.New("Error: unable to find Matrfile: " + matrFilePath)
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

	return "", errors.New("Error: unable to find Matrfile")
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

func lowerFirst(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[n:]
}
