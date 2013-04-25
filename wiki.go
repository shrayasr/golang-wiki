package main
import (
	"fmt"
	"errors"
	"regexp"
	"io/ioutil"
	"net/http"
	"html/template"
)

type Page struct {
	Title	string
	Body	[]byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	return &Page{Title: title, Body: body}, nil
}

/* ----- GLOBAL ----- */
const lenPath = len("/view/")
var templates = template.Must(template.ParseFiles("edit.html","view.html"))
var titleValidator = regexp.MustCompile("^[a-zA-Z0-9]+$")
/* ------------------ */

func getTitle(w http.ResponseWriter, r *http.Request) (title string, err error) {

	title = r.URL.Path[lenPath:]

	if !titleValidator.MatchString(title) {
		http.NotFound(w,r)
		err = errors.new("Invalid Page Title")
	}

	return
}

func renderTemplate(w http.ResponseWriter, templateFile string, p *Page) {

	err := templates.ExecuteTemplate(w, templateFile+".html", p)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func viewHandler(w http.ResponseWriter, r *http.Request) {

	title,err := getTitle(w,r)

	if err != nil {
		return
	}

	fmt.Printf("Rendering %s\n",title)

	p, err := loadPage(title)

	if err != nil {
		http.Redirect(w,r, "/edit/"+title,http.StatusFound)
		return
	}

	renderTemplate(w,"view",p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {

	title,err := getTitle(w,r)

	if err != nil {
		return
	}

	fmt.Printf("Editing %s\n",title)

	p, err := loadPage(title)

	if err != nil {
		p = &Page{Title: title}
	}

	renderTemplate(w,"edit",p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {

	title,err := getTitle(w,r)

	if err != nil {
		return
	}

	fmt.Printf("Saving %s\n",title)

	body := r.FormValue("body")

	p := &Page{Title:title,Body: []byte(body)}
	err := p.save()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w,r,"/view/"+title,http.StatusFound)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w,"Wiki server running")
}

func main() {

	http.HandleFunc("/",indexHandler)
	http.HandleFunc("/view/",viewHandler)
	http.HandleFunc("/edit/",editHandler)
	http.HandleFunc("/save/",saveHandler)

	fmt.Println("Listening on 8080");

	http.ListenAndServe(":8080",nil)

}
