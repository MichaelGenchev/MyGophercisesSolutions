package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

/*

I need to create a Graph from structs Nodes where each Node Is a Story with Title, Story, Options[] where options is another node.

From parsing the json I need to make the connections from the story arc IDS that are in the jsons "a" : {}, "a" is the key for this story.

Then once I have these I can display them with HTMLs where I will have each HTML page with title as the title of the story, body with the text and maybe links to the next stories.

If we don't want the INTRO to be hardcoded in the json we can tweak by adding,
"start_arc" "..." And our problem on decoding reads start_arc and uses that when / is requested.
*/

type AdventureHandler struct {
	adventure Adventure
	tmpl      *template.Template
}

func (h AdventureHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	arc := strings.TrimPrefix(r.URL.Path, "/story/")
	if arc == "" || arc == "/" {
		arc = h.adventure.Start
	}

	page, ok := h.adventure.Arcs[arc]
	if !ok {
		http.NotFound(w, r)
		return
	}

	err := h.tmpl.Execute(w, page)
	if err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

func main() {
	// Read the json file/

	jsonFile, err := os.ReadFile("gopher.json")
	if err != nil {
		log.Fatal(err)
	}

	adventure, err := NewAdventure(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	tmpl := template.Must(template.ParseFiles("index.html"))

	h := AdventureHandler{
		tmpl:      tmpl,
		adventure: *adventure,
	}

	log.Fatal(http.ListenAndServe(":8080", h))
}
