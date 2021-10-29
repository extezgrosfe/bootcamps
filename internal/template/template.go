package template

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/extezgrosfe/bootcamps/pkg/color"
)

type TemplateManager interface {
	PlaceTemplateInRepo() error
	ReplaceImportsInRepo() error
}

type templateManager struct {
	name     string
	username string
	path     string
}

func NewTemplateManager(name string, username string) TemplateManager {
	return &templateManager{
		name:     name,
		username: username,
		path:     "./" + name,
	}
}

func (tm *templateManager) PlaceTemplateInRepo() error {
	color.Print("white", "Looking for template folder...")
	if err := findTemplateFolder(); err != nil {
		color.Print("red", err.Error())
		return err
	}

	color.Print("white", "Template folder found!")

	// copy template folder content into repo folder
	cmd := fmt.Sprintf("cp -r ./template/* %s", tm.path)
	err := exec.Command("bash", "-c", cmd).Run()
	if err != nil {
		color.Print("red", fmt.Sprintf("Couldn't copy template folder: %s", err.Error()))
		return err
	}

	color.Print("white", "Template folder copied!")
	return nil
}

func (tm *templateManager) ReplaceImportsInRepo() error {
	err := filepath.Walk(tm.path, tm.visit)
	if err != nil {
		color.Print("red", fmt.Sprintf("Couldn't replace imports: %s", err.Error()))
		return err
	}

	return nil
}

// findTemplateFolder finds the template folder in the current directory
func findTemplateFolder() error {
	// check if a "template" folder exists
	// if not, create one

	if _, err := os.Stat("./template"); os.IsNotExist(err) {
		// clone template folder from github repo https://github.com/extezgrosfe/bootcamp-template.git
		color.Print("white", "Template folder not found, cloning template folder from github...")
		cmd := exec.Command("git", "clone", "https://github.com/extezgrosfe/bootcamp-template.git", "template")
		err := cmd.Run()
		if err != nil {
			color.Print("red", fmt.Sprintf("Couldn't clone template folder: %s", err.Error()))
			return err
		}
	}

	return nil
}

func (tm *templateManager) visit(path string, fi os.FileInfo, err error) error {

	if err != nil {
		return err
	}

	if fi.IsDir() {
		return nil //
	}

	matched, _ := filepath.Match("*go*", fi.Name())

	if matched {
		read, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		//fmt.Println(string(read))

		newContents := strings.Replace(string(read), "usuario/repositorio", tm.username+"/"+tm.name, -1)

		err = os.WriteFile(path, []byte(newContents), 0644)
		if err != nil {
			return err
		}

	}

	return nil
}
