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

// ⭐ 𝑼𝑹𝑳 𝒅𝒆𝒔 𝑨𝑷𝑰 :
var artistsURL string = "https://groupietrackers.herokuapp.com/api/artists"
var locationsURL string = "https://groupietrackers.herokuapp.com/api/locations"
var datesURL string = "https://groupietrackers.herokuapp.com/api/dates"
var relationURL string = "https://groupietrackers.herokuapp.com/api/relation"

// ⭐ 𝑫é𝒄𝒍𝒂𝒓𝒂𝒕𝒊𝒐𝒏 𝒅𝒆𝒔 𝑺𝒕𝒓𝒖𝒄𝒕𝒖𝒓𝒆𝒔 :
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

// ⭐⭐⭐ 𝑭𝒐𝒏𝒄𝒕𝒊𝒐𝒏 𝓜𝓪𝓲𝓷 ⭐⭐⭐ :
func main() {
	// 𝑪𝒉𝒂𝒓𝒈𝒆𝒓 𝒍𝒆𝒔 𝒇𝒊𝒄𝒉𝒊𝒆𝒓𝒔 𝒅𝒖 𝒅𝒐𝒔𝒔𝒊𝒆𝒓 '𝒔𝒕𝒂𝒕𝒊𝒄' 𝒔𝒖𝒓 𝒍𝒆 𝒔𝒆𝒓𝒗𝒆𝒖𝒓 :
	fs := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// 𝑭𝒆𝒕𝒄𝒉 𝒅𝒆𝒔 𝟒 𝑨𝑷𝑰 :
	parseJSON(artistsURL, &All.Artists)
	parseJSON(locationsURL, &All.Locations)
	parseJSON(datesURL, &All.Dates)
	parseJSON(relationURL, &All.Relations)

	// 𝑨𝒋𝒐𝒖𝒕 𝒅𝒆𝒔 𝑳𝒐𝒄𝒂𝒕𝒊𝒐𝒏𝒔 𝒅𝒂𝒏𝒔 𝒍𝒂 𝒔𝒕𝒓𝒖𝒄𝒕 𝑨𝒓𝒕𝒊𝒔𝒕𝒔 :
	for index := range All.Artists {
		All.Artists[index].Locations = All.Locations.Index[index].Locations
	}

	// 𝑭𝒐𝒓𝒎𝒂𝒕𝒕𝒂𝒈𝒆 𝒅𝒆𝒔 𝒏𝒐𝒎𝒔 𝒅𝒆 𝒍𝒐𝒄𝒂𝒕𝒊𝒐𝒏 𝒂𝒖 𝒇𝒐𝒓𝒎𝒂𝒕 "𝑽𝒊𝒍𝒍𝒆, 𝑷𝒂𝒚𝒔" :
	for i := range All.Artists { // Pour chaque élément d'indice i dans All.Artists...
		for j := range All.Artists[i].Locations { // ...et pour chaque élement d'indice j dans All.Artists[i].Locations...
			All.Artists[i].Locations[j] = strings.ReplaceAll(All.Artists[i].Locations[j], "_", " ")  // Dans chaque élement, je remplace les "_" par des espaces...
			All.Artists[i].Locations[j] = strings.ReplaceAll(All.Artists[i].Locations[j], "-", ", ") // ...et les "-" par des virgules.
			All.Artists[i].Locations[j] = strings.ReplaceAll(All.Artists[i].Locations[j], ", uk", ", UK")
			All.Artists[i].Locations[j] = strings.ReplaceAll(All.Artists[i].Locations[j], ", usa", ", USA")
			All.Artists[i].Locations[j] = strings.Title(All.Artists[i].Locations[j]) // Je mets la 1ère lettre de chaque mot en majuscule (pour faire beau)
		}
	}

	// 𝑯é𝒃𝒆𝒓𝒈𝒆𝒎𝒆𝒏𝒕 :
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/artist", detailHandler)

	// 𝑳𝒂𝒏𝒄𝒆𝒎𝒆𝒏𝒕 𝒅𝒖 𝒔𝒆𝒓𝒗𝒆𝒖𝒓 :
	fmt.Println("Listening server at port 8000.")
	http.ListenAndServe(":8000", nil)
}

// ⭐ 𝑭𝒐𝒏𝒄𝒕𝒊𝒐𝒏 𝒅𝒆 𝒓é𝒄𝒖𝒑é𝒓𝒂𝒕𝒊𝒐𝒏 / 𝒔𝒕𝒐𝒄𝒌𝒂𝒈𝒆 𝒅𝒖 𝑱𝑺𝑶𝑵 :
func parseJSON(myURL string, ptr interface{}) { // 𝒑𝒕𝒓 est un pointeur

	// Je vais chercher l'API de l'URL, et stocke le résultat dans 𝒓𝒆𝒔 :
	res, err := http.Get(myURL)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	// Je lis 𝒓𝒆𝒔, et stocke le résultat dans 𝒃𝒐𝒅𝒚 sous forme de tableau de bytes :
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	// Unmarshal déchiffre 𝒃𝒐𝒅𝒚 (qui est chiffré en JSON), et stocke le résultat dans la variable dont l'adresse est stockée dans le pointeur 𝒑𝒕𝒓.
	json.Unmarshal(body, &ptr)
}

// ⭐ 𝑭𝒐𝒏𝒄𝒕𝒊𝒐𝒏 𝒎𝒂𝒊𝒏𝑯𝒂𝒏𝒅𝒍𝒆𝒓 𝒑𝒐𝒖𝒓 𝒍𝒆 𝒉𝒂𝒏𝒅𝒍𝒆𝑭𝒖𝒏𝒄 (𝘦𝘹é𝘤𝘶𝘵𝘦 𝘭𝘦 𝘵𝘦𝘮𝘱𝘭𝘢𝘵𝘦 𝙢𝙖𝙞𝙣.𝙝𝙩𝙢𝙡) :
func mainHandler(w http.ResponseWriter, r *http.Request) {

	// GESTION DU STATUT '404' :
	if r.URL.Path != "/" {
		http.Error(w, "404 PAGE NOT FOUND", http.StatusNotFound)
		return
	}

	// GESTION DES REQUEST METHODS :
	switch r.Method {

	// 🍔 Méthode 'GET' — Lorsqu'on arrive sur la page main.html pour la 1ère fois :
	case "GET":
		tmpl, err := textTemplate.ParseFiles("./static/common.html", "./static/main.html", "./static/noresult.html")
		if err != nil {
			http.Error(w, "500 INTERNAL SERVER ERROR", http.StatusInternalServerError)
			log.Fatal(err)
		}
		tmpl.ExecuteTemplate(w, "common", All) // Envoyer tous les artistes

	// 🍔 Méthode 'POST' — Lorsqu'on appuie sur le bouton 'Valider' pour effectuer une recherche dans la barre :
	case "POST":
		tmpl, err := textTemplate.ParseFiles("./static/common.html", "./static/main.html", "./static/noresult.html")
		if err != nil {
			http.Error(w, "500 INTERNAL SERVER ERROR", http.StatusInternalServerError)
			log.Fatal(err)
		}

		result := filter(w, r)
		if result.ID != -1 {
			tmpl.ExecuteTemplate(w, "common", result) // Envoyer les artistes filtrés
		} else {
			tmpl.ExecuteTemplate(w, "no-result", "Aucun résultat...") // Envoyer la page "No Result"
		}
	}
}

// ⭐ 𝑭𝒐𝒏𝒄𝒕𝒊𝒐𝒏 𝒅𝒆𝒕𝒂𝒊𝒍𝑯𝒂𝒏𝒅𝒍𝒆𝒓 𝒑𝒐𝒖𝒓 𝒍𝒆 𝒉𝒂𝒏𝒅𝒍𝒆𝑭𝒖𝒏𝒄 (𝘦𝘹é𝘤𝘶𝘵𝘦 𝘭𝘦 𝘵𝘦𝘮𝘱𝘭𝘢𝘵𝘦 𝙙𝙚𝙩𝙖𝙞𝙡.𝙝𝙩𝙢𝙡) :
func detailHandler(w http.ResponseWriter, r *http.Request) {

	// 🍔 Méthode 'GET' uniquement :
	All.ID = Atoi(r.URL.Query().Get("id")) - 1
	tmpl, err := textTemplate.ParseFiles("./static/common.html", "./static/detail.html")
	if err != nil {
		http.Error(w, "500 INTERNAL SERVER ERROR", http.StatusInternalServerError)
		log.Fatal(err)
	}
	tmpl.ExecuteTemplate(w, "common", All)
}

// ⭐ 𝑭𝒐𝒏𝒄𝒕𝒊𝒐𝒏 𝒅𝒆 𝒇𝒊𝒍𝒕𝒓𝒂𝒈𝒆 𝒅𝒆𝒔 𝒂𝒓𝒕𝒊𝒔𝒕𝒆𝒔 𝒑𝒐𝒖𝒓 𝒍𝒂 𝒇𝒐𝒏𝒄𝒕𝒊𝒐𝒏 𝒎𝒂𝒊𝒏𝑯𝒂𝒏𝒅𝒍𝒆𝒓 :
func filter(w http.ResponseWriter, r *http.Request) API {
	var ToSend API
	var ArtistsToSend Artists

	var Inputs struct {
		searchType string
		toSearch   string
		value1     string
		value2     string
	}

	// Cas où il n'y aura aucun résultat à afficher :
	ToSend.ID = -1

	// Lecture et récupération de la requête :
	body, _ := ioutil.ReadAll(r.Body)        // Par exemple, body = [searchType=member&toSearch=Freddie&between=1&and=7]. Donc body contient toutes les valeurs des paramètres existants dans le template HTML.
	query, _ := url.ParseQuery(string(body)) // Par exemple, query = map[searchType: [member], toSearch: [Freddie], between: [1], and: [7] ]. ParseQuery() analyse body (casté en string) et crée une map en fonction des caractères '&' et '='.

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

	// Filtrage d'après l'input de l'utilisateur :
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
		ToSend.ID = i // S'il n'y a aucun artiste à ajouter, ToSend.ID restera égal à -1
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

	// Si le nom de l'artiste contient le terme recherché, on return true. Sinon, false :
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

	// On met en minucule le terme recherché, et on convertit l'array de strings contenant les locations en une seule string 'locations' (grâce à strings.Join), qu'on met aussi en minuscule :
	locations := strings.ToLower(strings.Join(All.Artists[index].Locations, " / ")) // Chaque nom de location est séparé par un " / "
	toSearch = strings.ToLower(toSearch)

	// Si les noms des locations contiennent le terme recherché, on return true. Sinon, false :
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

	// On met en minucule le terme recherché, et on convertit l'array de strings contenant les noms des membres en une seule string 'membersNames' (grâce à strings.Join), qu'on met aussi en minuscule :
	toSearch = strings.ToLower(toSearch)
	membersNames := strings.ToLower(strings.Join(All.Artists[index].Members, " "))
	numberOfMembers := len(All.Artists[index].Members) // Nombre de membres dans le groupe

	// Si la value1 ou la value2 n'a pas été renseignée par l'utilisateur, on les remplace chacune par une version min. / max. par défaut :
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
	// Si les noms des membres contiennent le terme recherché, et que le nombre de membres est compris entre value1 et value2, on return true :
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

	// Si la value1 ou la value2 n'a pas été renseignée par l'utilisateur, on les remplace chacune par une version min. / max. par défaut :
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

	// On récupère la date :
	fullDate := All.Artists[index].FirstAlbum             // Date complète au format DD-MM-YYY
	year, err := strconv.Atoi(fullDate[len(fullDate)-4:]) // Année uniquement, convertie en int.
	if err != nil {
		return false
	}

	// Si l'année est comprise entre value1 et value2, on return true :
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

	// Si la value1 ou la value2 n'a pas été renseignée par l'utilisateur, on les remplace chacune par une version min. / max. par défaut :
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

	// On récupère la date :
	creationDate := All.Artists[index].CreationDate // Année uniquement, déjà de type int.

	// Si l'année est comprise entre value1 et value2, on return true :
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
