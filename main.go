package main

import (
	"encoding/json"
	"fmt"
	textTemplate "html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// β­ πΌπΉπ³ πππ π¨π·π° :
var artistsURL string = "https://groupietrackers.herokuapp.com/api/artists"
var locationsURL string = "https://groupietrackers.herokuapp.com/api/locations"
var datesURL string = "https://groupietrackers.herokuapp.com/api/dates"
var relationURL string = "https://groupietrackers.herokuapp.com/api/relation"

// β­ π«Γ©πππππππππ πππ πΊπππππππππ :
type API struct {
	Artists   Artists
	Locations Locations
	Dates     Dates
	Relations Relations
	ID        int
}

type Artists []struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    []string `json:"locations"`
}

type Locations struct {
	Index []struct {
		ID        int      `json:"id"`
		Locations []string `json:"locations"`
	} `json:"index"`
}

type Dates struct {
	Index []struct {
		ID    int      `json:"id"`
		Dates []string `json:"dates"`
	} `json:"index"`
}

type Relations struct {
	Index []struct {
		ID            int                 `json:"id"`
		DateLocations map[string][]string `json:"datesLocations"`
	} `json:"index"`
}

var All API

// β­β­β­ π­πππππππ ππͺπ²π· β­β­β­ :
func main() {
	// πͺππππππ πππ ππππππππ ππ πππππππ 'ππππππ' πππ ππ πππππππ :
	fs := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// π­ππππ πππ π π¨π·π° :
	parseJSON(artistsURL, &All.Artists)
	parseJSON(locationsURL, &All.Locations)
	parseJSON(datesURL, &All.Dates)
	parseJSON(relationURL, &All.Relations)

	// π¨ππππ πππ π³ππππππππ ππππ ππ ππππππ π¨ππππππ :
	for index := range All.Artists {
		All.Artists[index].Locations = All.Locations.Index[index].Locations
	}

	// π­πππππππππ πππ ππππ ππ ππππππππ ππ ππππππ "π½ππππ, π·πππ" :
	for i := range All.Artists { // Pour chaque Γ©lΓ©ment d'indice i dans All.Artists...
		for j := range All.Artists[i].Locations { // ...et pour chaque Γ©lement d'indice j dans All.Artists[i].Locations...
			All.Artists[i].Locations[j] = strings.ReplaceAll(All.Artists[i].Locations[j], "_", " ")  // Dans chaque Γ©lement, je remplace les "_" par des espaces...
			All.Artists[i].Locations[j] = strings.ReplaceAll(All.Artists[i].Locations[j], "-", ", ") // ...et les "-" par des virgules.
			All.Artists[i].Locations[j] = strings.ReplaceAll(All.Artists[i].Locations[j], ", uk", ", UK")
			All.Artists[i].Locations[j] = strings.ReplaceAll(All.Artists[i].Locations[j], ", usa", ", USA")
			All.Artists[i].Locations[j] = strings.Title(All.Artists[i].Locations[j]) // Je mets la 1Γ¨re lettre de chaque mot en majuscule (pour faire beau)
		}
	}

	// π―Γ©πππππππππ :
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/artist", detailHandler)

	// π³ππππππππ ππ πππππππ :
	fmt.Println("Listening server at port 8000.")
	http.ListenAndServe(":8000", nil)
}

// β­ π­πππππππ ππ πΓ©πππΓ©ππππππ / ππππππππ ππ π±πΊπΆπ΅ :
func parseJSON(myURL string, ptr interface{}) { // πππ est un pointeur

	// Je vais chercher l'API de l'URL, et stocke le rΓ©sultat dans πππ :
	res, err := http.Get(myURL)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	// Je lis πππ, et stocke le rΓ©sultat dans ππππ sous forme de tableau de bytes :
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	// Unmarshal dΓ©chiffre ππππ (qui est chiffrΓ© en JSON), et stocke le rΓ©sultat dans la variable dont l'adresse est stockΓ©e dans le pointeur πππ.
	json.Unmarshal(body, &ptr)
}

// β­ π­πππππππ πππππ―ππππππ ππππ ππ πππππππ­πππ (π¦πΉΓ©π€πΆπ΅π¦ π­π¦ π΅π¦π?π±π­π’π΅π¦ π’πππ£.ππ©π’π‘) :
func mainHandler(w http.ResponseWriter, r *http.Request) {

	// GESTION DU STATUT '404' :
	if r.URL.Path != "/" {
		http.Error(w, "404 PAGE NOT FOUND", http.StatusNotFound)
		return
	}

	// GESTION DES REQUEST METHODS :
	switch r.Method {

	// π MΓ©thode 'GET' β Lorsqu'on arrive sur la page main.html pour la 1Γ¨re fois :
	case "GET":
		tmpl, err := textTemplate.ParseFiles("./static/common.html", "./static/main.html", "./static/noresult.html")
		if err != nil {
			http.Error(w, "500 INTERNAL SERVER ERROR", http.StatusInternalServerError)
			log.Fatal(err)
		}
		tmpl.ExecuteTemplate(w, "common", All) // Envoyer tous les artistes

	// π MΓ©thode 'POST' β Lorsqu'on appuie sur le bouton 'Valider' pour effectuer une recherche dans la barre :
	case "POST":
		tmpl, err := textTemplate.ParseFiles("./static/common.html", "./static/main.html", "./static/noresult.html")
		if err != nil {
			http.Error(w, "500 INTERNAL SERVER ERROR", http.StatusInternalServerError)
			log.Fatal(err)
		}

		result := filter(w, r)
		if result.ID != -1 {
			tmpl.ExecuteTemplate(w, "common", result) // Envoyer les artistes filtrΓ©s
		} else {
			tmpl.ExecuteTemplate(w, "no-result", "Aucun rΓ©sultat...") // Envoyer la page "No Result"
		}
	}
}

// β­ π­πππππππ πππππππ―ππππππ ππππ ππ πππππππ­πππ (π¦πΉΓ©π€πΆπ΅π¦ π­π¦ π΅π¦π?π±π­π’π΅π¦ πππ©πππ‘.ππ©π’π‘) :
func detailHandler(w http.ResponseWriter, r *http.Request) {

	// π MΓ©thode 'GET' uniquement :
	All.ID = Atoi(r.URL.Query().Get("id")) - 1
	tmpl, err := textTemplate.ParseFiles("./static/common.html", "./static/detail.html")
	if err != nil {
		http.Error(w, "500 INTERNAL SERVER ERROR", http.StatusInternalServerError)
		log.Fatal(err)
	}
	tmpl.ExecuteTemplate(w, "common", All)
}

// β­ π­πππππππ ππ ππππππππ πππ ππππππππ ππππ ππ ππππππππ πππππ―ππππππ :
func filter(w http.ResponseWriter, r *http.Request) API {
	var ToSend API
	var ArtistsToSend Artists

	var Inputs struct {
		searchType string
		toSearch   string
		value1     string
		value2     string
	}

	// Cas oΓΉ il n'y aura aucun rΓ©sultat Γ  afficher :
	ToSend.ID = -1

	// Lecture et rΓ©cupΓ©ration de la requΓͺte :
	body, _ := ioutil.ReadAll(r.Body)        // Par exemple, body = [searchType=member&toSearch=Freddie&between=1&and=7]. Donc body contient toutes les valeurs des paramΓ¨tres existants dans le template HTML.
	query, _ := url.ParseQuery(string(body)) // Par exemple, query = map[searchType: [member], toSearch: [Freddie], between: [1], and: [7] ]. ParseQuery() analyse body (castΓ© en string) et crΓ©e une map en fonction des caractΓ¨res '&' et '='.

	for key, value := range query {
		switch key {
		case "searchType":
			Inputs.searchType = value[0]
		case "toSearch":
			Inputs.toSearch = value[0]
		case "between":
			Inputs.value1 = value[0]
		case "and":
			Inputs.value2 = value[0]
		}
	}

	// Filtrage d'aprΓ¨s l'input de l'utilisateur :
	for i := 0; i < len(All.Artists); i++ {
		if !checkArtists(Inputs.searchType, Inputs.toSearch, i) {
			continue
		}
		if !checkLocations(Inputs.searchType, Inputs.toSearch, i) {
			continue
		}
		if !checkMembers(Inputs.searchType, Inputs.toSearch, Inputs.value1, Inputs.value2, i) {
			continue
		}
		if !checkFirstAlbum(Inputs.searchType, Inputs.value1, Inputs.value2, i) {
			continue
		}
		if !checkCreationDate(Inputs.searchType, Inputs.value1, Inputs.value2, i) {
			continue
		}
		ArtistsToSend = append(ArtistsToSend, All.Artists[i])
		ToSend.ID = i // S'il n'y a aucun artiste Γ  ajouter, ToSend.ID restera Γ©gal Γ  -1
	}

	ToSend.Artists = ArtistsToSend
	return ToSend
}

func checkArtists(searchType string, toSearch string, index int) bool {

	// Si l'utilisateur n'a pas choisi le champ 'Artist', cette fonction doit return true :
	if searchType != "artist" {
		return true
	}

	// On met en minucule le nom de l'artiste[i] et le terme saisi toSearch :
	artistName := strings.ToLower(All.Artists[index].Name)
	toSearch = strings.ToLower(toSearch)

	// Si le nom de l'artiste contient le terme recherchΓ©, on return true. Sinon, false :
	if strings.Contains(artistName, toSearch) {
		return true
	}
	return false
}

func checkLocations(searchType string, toSearch string, index int) bool {

	// Si l'utilisateur n'a pas choisi le champ 'Artist', cette fonction doit return true :
	if searchType != "location" {
		return true
	}

	// On met en minucule le terme recherchΓ©, et on convertit l'array de strings contenant les locations en une seule string 'locations' (grΓ’ce Γ  strings.Join), qu'on met aussi en minuscule :
	locations := strings.ToLower(strings.Join(All.Artists[index].Locations, " / ")) // Chaque nom de location est sΓ©parΓ© par un " / "
	toSearch = strings.ToLower(toSearch)

	// Si les noms des locations contiennent le terme recherchΓ©, on return true. Sinon, false :
	if strings.Contains(locations, toSearch) {
		return true
	}
	return false
}

func checkMembers(searchType string, toSearch string, inputValue1 string, inputValue2 string, index int) bool {

	// Si l'utilisateur n'a pas choisi le champ 'Member', cette fonction doit return true :
	if searchType != "member" {
		return true
	}

	// On met en minucule le terme recherchΓ©, et on convertit l'array de strings contenant les noms des membres en une seule string 'membersNames' (grΓ’ce Γ  strings.Join), qu'on met aussi en minuscule :
	toSearch = strings.ToLower(toSearch)
	membersNames := strings.ToLower(strings.Join(All.Artists[index].Members, " "))
	numberOfMembers := len(All.Artists[index].Members) // Nombre de membres dans le groupe

	// Si la value1 ou la value2 n'a pas Γ©tΓ© renseignΓ©e par l'utilisateur, on les remplace chacune par une version min. / max. par dΓ©faut :
	if inputValue1 == "" {
		inputValue1 = "0"
	}
	if inputValue2 == "" {
		inputValue2 = "1000"
	}

	value1, err := strconv.Atoi(inputValue1)
	if err != nil {
		return false
	}
	value2, err := strconv.Atoi(inputValue2)
	if err != nil {
		return false
	}
	// Si les noms des membres contiennent le terme recherchΓ©, et que le nombre de membres est compris entre value1 et value2, on return true :
	if strings.Contains(membersNames, toSearch) && value1 <= numberOfMembers && value2 >= numberOfMembers {
		return true
	}
	return false
}

func checkFirstAlbum(searchType string, inputValue1 string, inputValue2 string, index int) bool {

	// Si l'utilisateur n'a pas choisi le champ 'First Album', cette fonction doit return true :
	if searchType != "firstAlbum" {
		return true
	}

	// Si la value1 ou la value2 n'a pas Γ©tΓ© renseignΓ©e par l'utilisateur, on les remplace chacune par une version min. / max. par dΓ©faut :
	if inputValue1 == "" {
		inputValue1 = "0"
	}
	if inputValue2 == "" {
		inputValue2 = "3000"
	}

	value1, err := strconv.Atoi(inputValue1)
	if err != nil {
		return false
	}
	value2, err := strconv.Atoi(inputValue2)
	if err != nil {
		return false
	}

	// On rΓ©cupΓ¨re la date :
	fullDate := All.Artists[index].FirstAlbum             // Date complΓ¨te au format DD-MM-YYY
	year, err := strconv.Atoi(fullDate[len(fullDate)-4:]) // AnnΓ©e uniquement, convertie en int.
	if err != nil {
		return false
	}

	// Si l'annΓ©e est comprise entre value1 et value2, on return true :
	if value1 <= year && value2 >= year {
		return true
	}
	return false
}

func checkCreationDate(searchType string, inputValue1 string, inputValue2 string, index int) bool {

	// Si l'utilisateur n'a pas choisi le champ 'First Album', cette fonction doit return true :
	if searchType != "creationDate" {
		return true
	}

	// Si la value1 ou la value2 n'a pas Γ©tΓ© renseignΓ©e par l'utilisateur, on les remplace chacune par une version min. / max. par dΓ©faut :
	if inputValue1 == "" {
		inputValue1 = "0"
	}
	if inputValue2 == "" {
		inputValue2 = "3000"
	}

	value1, err := strconv.Atoi(inputValue1)
	if err != nil {
		return false
	}
	value2, err := strconv.Atoi(inputValue2)
	if err != nil {
		return false
	}

	// On rΓ©cupΓ¨re la date :
	creationDate := All.Artists[index].CreationDate // AnnΓ©e uniquement, dΓ©jΓ  de type int.

	// Si l'annΓ©e est comprise entre value1 et value2, on return true :
	if value1 <= creationDate && value2 >= creationDate {
		return true
	}
	return false
}

func Atoi(s string) int {
	runes := []rune(s)
	i := 0
	multi := 1
	for l := len(runes) - 1; l > -1; l-- {
		if 47 < s[l] && s[l] < 58 || s[l] == 43 || s[l] == 45 {
			if 47 < s[l] && s[l] < 58 {
				i = i + int(s[l]-'0')*multi
				multi = multi * 10
			} else {
				if s[l] == 43 {
					if l == len(runes)-1 || len(runes) > 1 && s[l+1] == 43 || len(runes) > 1 && s[l+1] == 45 || len(runes) == 1 {
						i = 0
						return i
					} else {
						i = i + 0
					}
				}
				if s[l] == 45 {
					if l == len(runes)-1 || len(runes) > 1 && s[l+1] == 45 || len(runes) > 1 && s[l+1] == 43 || len(runes) == 1 {
						i = 0
						return i
					} else {
						i = -i
					}
				}
			}
		} else {
			i = 0
			l = -1
		}
	}
	return i
}
