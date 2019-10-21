package astrix

import (
	"bytes"
	"go/format"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
)

// GetFiles returns a list of files from the pkg of your choice.
// To be sure, just use an absolute path which can be retrieved from runtime:
//  package main
//
//  import (
//    "fmt"
//    "path/filepath"
//    "runtime"
//
//    "github.com/andersfylling/ggi"
//   )
//
//   var (
//     _, b, _, _ = runtime.Caller(0)
//     basepath   = filepath.Dir(b)
//     genpath    = "/generate/testing" // diff path; from root pkg to this files pkg
//   )
//
//   func main() {
//	   path := basepath[:len(basepath)-len(genpath)]
//	   files, err := astrix.FindFiles(path, regexp.MustCompile(".*\_gen\.go"))
//	   if err != nil {
//       panic(err)
//     }
//
//     // TODO: make sure these prints the .go files in your desired directory
//     for _, f := range files {
//       fmt.Println(f.Name())
//     }
//   }
func FindFiles(path string, except *regexp.Regexp) (files []os.FileInfo, err error) {
	f, err := os.Open(path)
	if err != nil {
		return files, err
	}
	fileInfo, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		f := file.Name()
		if file.IsDir() || !strings.HasSuffix(f, ".go") || (except != nil && except.MatchString(f)) {
			continue
		}

		files = append(files, file)
	}

	return files, nil
}

type Data struct {
	PackageName string
}

// ToCamelCase takes typical CONST names such as T_AS and converts them to TAs.
// TEST_EIN_TRES => TestEinTres
func ToCamelCase(s string) string {
	b := []byte(strings.ToLower(s))
	for i := range b {
		if b[i] == '_' && i < len(b)-1 {
			b[i+1] ^= 0x20
		}
	}
	s = strings.Replace(string(b), "-", "", -1)
	return s
}

func makeFile(d *Data, tplFile, target string) {
	fMap := template.FuncMap{
		"ToUpper":      strings.ToUpper,
		"ToLower":      strings.ToLower,
		"Capitalize":   func(s string) string { return strings.ToUpper(s[0:1]) + s[1:] },
		"Decapitalize": func(s string) string { return strings.ToLower(s[0:1]) + s[1:] },
		"ToCamelCase":  ToCamelCase,
	}

	// Open & parse our template
	tpl := template.Must(template.New(path.Base(tplFile)).Funcs(fMap).ParseFiles(tplFile))

	// Execute the template, inserting all the event information
	var b bytes.Buffer
	if err := tpl.Execute(&b, d); err != nil {
		panic(err)
	}

	// Format it according to gofmt standards
	formatted, err := format.Source(b.Bytes())
	if err != nil {
		panic(err)
	}

	// And write it.
	if err = ioutil.WriteFile(target, formatted, 0644); err != nil {
		panic(err)
	}
}
