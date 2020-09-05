package htmlviewer

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Images struct {
	Files []string
}

var images Images
var imagespath string
var templates = template.Must(template.New("main").Funcs(template.FuncMap{
	"EndDiv": func(fname string) bool {
		if strings.Contains(fname, "year") {
			return true
		} else {
			return false
		}
	},
	"BeginDiv": func(fname string) bool {
		if strings.Contains(fname, "day") {
			return true
		} else {
			return false
		}
	},
}).ParseFiles("/usr/share/rrdreader/templates/images.html"))

func RunHTMLServer(address string, output string) {
	imagespath = output
	http.HandleFunc("/", RootHandler)
	fs := http.FileServer(http.Dir(output))
	http.Handle("/images/", http.StripPrefix("/images/", fs))
	http.ListenAndServe(address, nil)
}

func renderTemplate(w http.ResponseWriter, tmpl string, img *Images) {
	templates.Funcs(template.FuncMap{"hideIt": func(fname string) bool {
		if strings.Contains(fname, "year") || strings.Contains(fname, "month") {
			return true
		} else {
			return false
		}
	}})
	err := templates.ExecuteTemplate(w, tmpl+".html", img)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	images.Files = ImageSearcher(imagespath)
	renderTemplate(w, "images", &images)
}

func ImageSearcher(dirname string) (result []string) {

	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		log.Fatalf("ImageSearcher error:%s", err)
	}
	dirBase := filepath.Base(dirname)
	err := filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if filepath.Ext(path) == ".png" {
				fileBase := filepath.Base(path)
				uriFile := filepath.Join(dirBase, fileBase)
				result = append(result, uriFile)
			}
		}
		return nil
	})

	if len(result) == 0 {
		log.Fatalf("No files in directory")
	}

	if err != nil {
		log.Fatalf("ImageSearcher error: %s", err)
	}
	// Changes *-week.png and *-month.png in
	// list to view in the template
	for i := 0; i < len(result); i++ {
		if strings.Contains(result[i], "-week") {
			var temp string
			temp = result[i]
			result[i] = result[i-1]
			result[i-1] = temp
		}
	}
	return result
}
