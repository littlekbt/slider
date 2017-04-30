package slider

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"google.golang.org/api/iterator"
	"google.golang.org/appengine"

	"cloud.google.com/go/storage"
)

func init() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/slide", slideHandler)
}

type MdFiles struct {
	Files []string
}

func handler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/index.html")
	files := MdFiles{Files: getMdFiles(r)}
	t.Execute(w, files)
}

func getMdFiles(r *http.Request) []string {
	ctx := appengine.NewContext(r)
	client, err := storage.NewClient(ctx)

	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	it := client.Bucket("slider-store").Objects(ctx, nil)
	files := make([]string, 0)
	for {
		objAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			// TODO: Handle error.
		}
		r := regexp.MustCompile(`\.md`)
		if r.MatchString(objAttrs.Name) {
			files = append(files, objAttrs.Name)
		}
	}

	return files
}

type MdFile struct {
	Title string
	Body  string
}

func slideHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/slide.html")
	m := r.URL.Query()
	mdFile := getMdFile(r, m.Get("object"))
	t.Execute(w, mdFile)
}

func getMdFile(r *http.Request, object string) MdFile {
	ctx := appengine.NewContext(r)
	client, err := storage.NewClient(ctx)

	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	oh := client.Bucket("slider-store").Object(object)
	re, _ := oh.NewReader(ctx)
	b, _ := ioutil.ReadAll(re)
	re.Close()

	return MdFile{Title: object, Body: string(b[:])}
}
