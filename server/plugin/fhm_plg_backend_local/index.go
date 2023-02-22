package fhm_plg_backend_local

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	. "github.com/mickael-kerjean/filestash/server/common"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	Backend.Register("fhm local", FhmLocal{})
}

type FhmLocal struct {
	exclude_extensions []string
}

func (this FhmLocal) Init(params map[string]string, app *App) (IBackend, error) {
	if err := bcrypt.CompareHashAndPassword(
		[]byte(Config.Get("auth.admin").String()),
		[]byte(params["password"]),
	); err != nil {
		return nil, ErrAuthenticationFailed
	}

	var exclude_extensions []string

	for _, extension := range strings.Split(params["exclude_extensions"], ",") {
		exclude_extensions = append(exclude_extensions, strings.TrimSpace(extension))
	}

	return FhmLocal{
		exclude_extensions: exclude_extensions,
	}, nil
}

func (this FhmLocal) LoginForm() Form {
	return Form{
		Elmnts: []FormElement{
			{
				Name:  "type",
				Type:  "hidden",
				Value: "fhm local",
			},
			{
				Name:        "password",
				Type:        "password",
				Placeholder: "Admin Password",
			},
			{
				Name:        "path",
				Type:        "text",
				Placeholder: "Path",
			},
			{
				Name:        "exclude_extensions",
				Type:        "text",
				Default:     "jtmb",
				Description: "Comma separated list of file extensions, without the leading '.', that should be omitted from directory listings",
			},
		},
	}
}

func (this FhmLocal) Home() (string, error) {
	return os.UserHomeDir()
}

func (this FhmLocal) Ls(path string) ([]os.FileInfo, error) {
	// open directory
	f, err := SafeOsOpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}

	// read directory
	all_results, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}

	var results []os.FileInfo

	// iterate through results
	for i := 0; i < len(all_results); i++ {
		result := all_results[i]

		if !result.IsDir() {
			// exclude files with extensions configured in "exclude_extensions"
			extension := strings.TrimPrefix(filepath.Ext(result.Name()), ".")

			if ContainsStr(this.exclude_extensions, extension) {
				continue
			}
		}

		results = append(results, result)
	}

	return results, nil
}

func (this FhmLocal) Cat(path string) (io.ReadCloser, error) {
	return SafeOsOpenFile(path, os.O_RDONLY, os.ModePerm)
}

func (this FhmLocal) Mkdir(path string) error {
	return SafeOsMkdir(path, 0755)
}

func (this FhmLocal) Rm(path string) error {
	return SafeOsRemoveAll(path)
}

func (this FhmLocal) Mv(from, to string) error {
	return SafeOsRename(from, to)
}

func (this FhmLocal) Save(path string, content io.Reader) error {
	f, err := SafeOsOpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, content)
	return err
}

func (this FhmLocal) Touch(path string) error {
	f, err := SafeOsOpenFile(path, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	if _, err = f.Write([]byte("")); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

func ContainsStr(list []string, search string) bool {
	for i := 0; i < len(list); i++ {
		if list[i] == search {
			return true
		}
	}

	return false
}
