package main

import (
	"log"
	"net/http"
	"os"
	"regexp"
	"text/template"
)

type Page struct {
	Title string
	Body  []byte
}

// Template Caching, Error Handling, Validation, Function Literals and Closures

// Catching the error condition in each handler introduces a lot of repeated code.
// What if we could wrap each of the handlers in a function that does this validation
// and error checking? Go's function literals provide a powerful means of abstracting
// functionality that can help us here.

// ParseFiles, MustCompile panics
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

// func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
// 	// t, _ := template.ParseFiles(tmpl + ".html")
// 	// t.Execute(w, p)
// 	t, err := template.ParseFiles(tmpl + ".html")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	err = t.Execute(w, p)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }

// Template Caching
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// As geting the title part is repeating in each func, so write separate wrapper func
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	// func viewHandler(w http.ResponseWriter, r *http.Request) {

	// title := r.URL.Path[len("/view/"):]
	// p, _ := loadPage(title)
	// fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)

	// title := r.URL.Path[len("/view/"):]
	// title, err := getTitle(w, r)
	// if err != nil {
	// 	return
	// }

	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	// t, _ := template.ParseFiles("view.html")
	// t.Execute(w, p)
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	// func editHandler(w http.ResponseWriter, r *http.Request) {
	// title := r.URL.Path[len("/edit/"):]
	// p, err := loadPage(title)
	// if err != nil {
	// 	p = &Page{Title: title}
	// }
	// fmt.Fprintf(w, "<h1>Editing %s</h1>"+
	// 	"<form action=\"/save/%s\" method=\"POST\">"+
	// 	"<textarea name=\"body\">%s</textarea><br>"+
	// 	"<input type=\"submit\" value=\"Save\">"+
	// 	"</form>",
	// 	p.Title, p.Title, p.Body)

	// title := r.URL.Path[len("/edit/"):]
	// title, err := getTitle(w, r)
	// if err != nil {
	// 	return
	// }

	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	// t, _ := template.ParseFiles("edit.html")
	// t.Execute(w, p)
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	// func saveHandler(w http.ResponseWriter, r *http.Request) {
	// title := r.URL.Path[len("/save/"):]

	// title, err := getTitle(w, r)
	// if err != nil {
	// 	return
	// }

	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

// func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
// 	m := validPath.FindStringSubmatch(r.URL.Path)
// 	if m == nil {
// 		http.NotFound(w, r)
// 		return "", errors.New("invalid Page Title")
// 	}
// 	return m[2], nil // The title is the second subexpression.
// }

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func main() {
	// p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
	// p1.save()
	// p2, _ := loadPage("TestPage")
	// fmt.Println(string(p2.Body))

	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
