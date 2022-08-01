package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

//go:embed templates/main.go
var mainGo []byte

func main() {
	app := kingpin.New("gostudy", "Assists in creating code for go language learning")
	add(app)
	rm(app)
	kingpin.MustParse(app.Parse(os.Args[1:]))
}

func add(app *kingpin.Application) {
	cmd := app.Command("add", "Add new study.")
	studyName := cmd.Arg("studyName", "Name for study.").Required().String()
	open := cmd.Flag("open", "Open new file automatically with Visual Studio Code.").Short('o').Bool()

	cmd.Action(func(c *kingpin.ParseContext) error {
		addCommand(*studyName, *open)
		return nil
	})
}

func addCommand(studyName string, open bool) {
	// check: cmd dir exits
	const dir = "./cmd"
	if f, err := os.Stat(dir); os.IsNotExist(err) || !f.IsDir() {
		fmt.Printf("%s directory is not found.\n", dir)
		return
	}

	newDir := filepath.Join(dir, studyName)

	if err := os.Mkdir(newDir, 0777); err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	fmt.Printf("Created %s directory.\n", newDir)

	newFile := filepath.Join(newDir, "main.go")
	f, err := os.Create(newFile)
	if err != nil {
		fmt.Printf("error: %v", err)
	}

	defer f.Close()
	if _, err := f.Write(mainGo); err != nil {
		fmt.Printf("error: %v\n", err)
	}

	fmt.Printf("Created %s file\n", newFile)

	if open {
		openWithVSCode(newFile)
	}
}

func openWithVSCode(file string) {
	if _, err := exec.LookPath("code"); err != nil {
		return
	}

	if err := exec.Command("code", file).Run(); err != nil {
		fmt.Printf("error: %v\n", err)
	}
}

func rm(app *kingpin.Application) {
	cmd := app.Command("rm", "Remove the study")
	studyName := cmd.Arg("studyName", "Name for study.").Required().String()
	force := cmd.Flag("force", "Open new file automatically with Visual Studio Code").Short('f').Bool()

	cmd.Action(func(c *kingpin.ParseContext) error {
		rmCommand(*studyName, *force)
		return nil
	})
}

func rmCommand(studyName string, force bool) {
	if studyName == "" {
		fmt.Println("error: studyName is empty.")
		return
	}

	var dir = "./cmd/" + studyName

	if f, err := os.Stat(dir); os.IsNotExist(err) || !f.IsDir() {
		fmt.Printf("%s directory is not found.\n", dir)
		return
	}

	confirm := true
	if !force {
		fmt.Printf("Remove %s directory?  [y]es/[n]o ", dir)
		var yesNo string
		fmt.Scan(&yesNo)
		confirm = strings.EqualFold(yesNo, "y") || strings.EqualFold(yesNo, "yes")
	}

	if !confirm {
		return
	}

	if err := os.RemoveAll(dir); err != nil {
		fmt.Printf("error: %v\n", err)
	}

	fmt.Printf("Removed %s directory.\n", dir)
}
