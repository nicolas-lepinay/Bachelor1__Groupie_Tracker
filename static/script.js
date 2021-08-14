// 𝑺𝒄𝒓𝒊𝒑𝒕 𝒈é𝒓𝒂𝒏𝒕 𝒍𝒂 𝒃𝒂𝒓𝒓𝒆 𝒅𝒆 𝒓𝒆𝒄𝒉𝒆𝒓𝒄𝒉𝒆 𝒅𝒚𝒏𝒂𝒎𝒊𝒒𝒖𝒆 𝒆𝒕 𝒍𝒂 𝓜𝓪𝓹 𝓑𝓸𝔁 :

const searchInput = document.getElementById("sch"); // Barre de recherche générale
const searchInput2 = document.getElementById("sch2"); // Petite barre de saisie pour valeur 1 ("Entre {v1} et {v2}")
const searchInput3 = document.getElementById("sch3"); // Petite barre de saisie pour valeur 2 ("Entre {v1} et {v2}")

let searchInputsBis = document.getElementsByClassName("search"); // Les 2 petites barres de saisies + les paragraphes "Entre" et "et" qui les entourent.
const selector = document.getElementById("sType"); // Sélecteur de champs (Artists, Members, First Album, etc.)
let artists = document.getElementsByClassName("artist"); // NodeList contenant toutes les div 'Artists'
let miniArtists = document.getElementsByClassName("mini-artist"); // NodeList contenant toutes les <a> 'mini-artist' dans le menu Hamburger

let searchTerm = "";
let value1 = 0;
let value2 = 2020;

let invisibleNames = document.getElementsByClassName("invisibleName");
let invisibleMembers = document.getElementsByClassName("invisibleMembers");
let invisibleLocations = document.getElementsByClassName("invisibleLocations");
let invisibleNumbers = document.getElementsByClassName("invisibleNumber");
let invisibleAlbums = document.getElementsByClassName("invisibleAlbum");
let invisibleCreations = document.getElementsByClassName("invisibleCreation");

let onMainPage = true;

// 𝑪𝒉𝒂𝒏𝒈𝒆𝒎𝒆𝒏𝒕 𝒅𝒖 𝒄𝒉𝒂𝒎𝒑 𝒅𝒆 𝒓𝒆𝒄𝒉𝒆𝒓𝒄𝒉𝒆 (𝘮𝘦𝘯𝘶 𝘥é𝘳𝘰𝘶𝘭𝘢𝘯𝘵) :
selector.addEventListener("change", (e) => {
    field = e.target.value;

    switch (field) {
        case "member":
            searchInput.style.display = "inline"; // La barre de recherche générale est visible
            [...searchInputsBis].forEach(element => element.style.display = "inline"); // Les 2 petites barres sont aussi visibles
            break;
        case "creationDate":
        case "firstAlbum":
            searchInput.style.display = "none";
            [...searchInputsBis].forEach(element => element.style.display = "inline");
            break;
        default: // "Artist" ou "Location"
            searchInput.style.display = "inline";
            [...searchInputsBis].forEach(element => element.style.display = "none");
            break;
    }
    checkNoResult();

});

// 𝑺𝒂𝒊𝒔𝒊𝒆 𝒅𝒂𝒏𝒔 𝒍𝒂 𝒃𝒂𝒓𝒓𝒆 𝒅𝒆 𝒓𝒆𝒄𝒉𝒆𝒓𝒄𝒉𝒆 (𝘪𝘯𝘱𝘶𝘵) :
searchInput.addEventListener("input", (e) => {
    searchTerm = e.target.value;
    switch (selector.value) {
        case "artist":
            searchNames(invisibleNames);
            break;
        case "location":
            searchNames(invisibleLocations);
            break;
        case "member":
            searchNumbers(invisibleMembers);
            break;
    }
    checkNoResult();
});

function searchNames(nodeList) {
    if (artists[0] != null) {
        onMainPage = true;
    } else {
        onMainPage = false;
    }

    // Pour i allant de 0 jusqu'au nombre d'artistes (51) :
    for (i = 0; i < miniArtists.length; i++) {
        // Si le contenu-texte du paragraphe 'invisibleXYZ' contient le terme recherché :
        if (nodeList[i].textContent.toLowerCase().includes(searchTerm.toLowerCase())) {
            //... alors la div 'artist' est visible, sinon elle est invisible :
            if (onMainPage) { artists[i].style.display = "block"; }
            miniArtists[i].style.display = "block";
        } else {
            if (onMainPage) { artists[i].style.display = "none"; }
            miniArtists[i].style.display = "none";
        }
    }
}

function searchNumbers(nodeList) {

    searchInput2.addEventListener("input", (e) => {
        value1 = e.target.value;
    });

    searchInput3.addEventListener("input", (e) => {
        value2 = e.target.value;
    });

    if (artists[0] != null) {
        onMainPage = true;
    } else {
        onMainPage = false;
    }

    for (i = 0; i < miniArtists.length; i++) {
        if (nodeList[i].textContent.toLowerCase().includes(searchTerm.toLowerCase()) && parseInt(invisibleNumbers[i].textContent) >= value1 && parseInt(invisibleNumbers[i].textContent) <= value2) {
            if (onMainPage) { artists[i].style.display = "block"; }
            miniArtists[i].style.display = "block";
        } else {
            if (onMainPage) { artists[i].style.display = "none"; }
            miniArtists[i].style.display = "none";
        }
    }
}


// 𝑺𝒂𝒊𝒔𝒊𝒆 𝒅𝒂𝒏𝒔 𝒖𝒏𝒆 𝒅𝒆𝒔 𝟐 𝒑𝒆𝒕𝒊𝒕𝒆𝒔 𝒃𝒂𝒓𝒓𝒆𝒔 𝒅𝒆 𝒔𝒂𝒊𝒔𝒊𝒆 (𝘷𝘢𝘭𝘦𝘶𝘳1 𝘦𝘵 𝘷𝘢𝘭𝘦𝘶𝘳2) :
searchInput2.addEventListener("input", (e) => {
    value1 = e.target.value;

    switch (selector.value) {
        case "firstAlbum":
            searchAlbum(invisibleAlbums);
            break;
        case "creationDate":
            searchCreation(invisibleCreations);
            break;
        case "member":
            searchNumbers(invisibleMembers);
    }
    checkNoResult();
});

searchInput3.addEventListener("input", (e) => {
    value2 = e.target.value;

    switch (selector.value) {
        case "firstAlbum":
            searchAlbum(invisibleAlbums);
            break;
        case "creationDate":
            searchCreation(invisibleCreations);
            break;
        case "member":
            searchNumbers(invisibleMembers);
    }
    checkNoResult();
});

function searchAlbum(nodeList) {
    if (artists[0] != null) {
        onMainPage = true;
    } else {
        onMainPage = false;
    }

    for (i = 0; i < miniArtists.length; i++) {
        if (parseInt(nodeList[i].textContent.slice(6)) >= value1 && parseInt(nodeList[i].textContent.slice(6)) <= value2) {
            if (onMainPage) { artists[i].style.display = "block"; }
            miniArtists[i].style.display = "block";
        } else {
            if (onMainPage) { artists[i].style.display = "none"; }
            miniArtists[i].style.display = "none";
        }
    }
}

function searchCreation(nodeList) {
    if (artists[0] != null) {
        onMainPage = true;
    } else {
        onMainPage = false;
    }

    for (i = 0; i < miniArtists.length; i++) {
        if (parseInt(nodeList[i].textContent) >= value1 && parseInt(nodeList[i].textContent) <= value2) {
            if (onMainPage) { artists[i].style.display = "block"; }
            miniArtists[i].style.display = "block";
        } else {
            if (onMainPage) { artists[i].style.display = "none"; }
            miniArtists[i].style.display = "none";
        }
    }
}

function checkNoResult() {
    if (artists[0] != null) {
        onMainPage = true;
    } else {
        onMainPage = false;
    }

    // S'il n'y a aucun résultat :
    let k = 0;
    for (i = 0; i < miniArtists.length; i++) {
        if (miniArtists[i].style.display == "block") {
            k++;
            break;
        }
    }
    if (k == 0) {
        if (onMainPage) { document.getElementById("no-result").style.display = "block"; }
        document.getElementById("mini-no-result").style.display = "block";
    } else {
        if (onMainPage) { document.getElementById("no-result").style.display = "none"; }
        document.getElementById("mini-no-result").style.display = "none";
    }
}


// 𝑭𝒆𝒕𝒄𝒉 𝒅𝒆 𝒍'𝑨𝑷𝑰 𝒅𝒆 𝑾𝒊𝒌𝒊𝒑é𝒅𝒊𝒂 : 
let data;

async function fetchData(name) {
    // Formattage des noms d'artistes non-reconnus immédiatement par Wikipédia :
    switch (name) {
        case "SOJA":
            name = "Soldiers+of+Jah+Army";
            break;
        case "Genesis":
        case "NWA":
        case "Muse":
            name += "+(groupe)";
            break;
        case "R3HAB":
            name = "Fadil+El+Ghoul";
            break;
    }
    // URL devant être suivi d'un terme de recherche (name) :
    data = await fetch("https://fr.wikipedia.org/w/api.php?format=json&action=query&prop=extracts&exintro&explaintext&redirects=1&origin=*&titles=" + name) // Ne pas oublier 'origin=*' dans l'URL pour éviter une erreur 'NO-CORS'
        .then(response => response.json())
}

// 𝑬𝒏𝒗𝒐𝒊 𝒅𝒆 𝒍𝒂 𝒅𝒆𝒔𝒄𝒓𝒊𝒑𝒕𝒊𝒐𝒏 𝑾𝒊𝒌𝒊𝒑é𝒅𝒊𝒂 𝒅𝒆 𝒍'𝒂𝒓𝒕𝒊𝒔𝒕𝒆 𝒅𝒂𝒏𝒔 𝒍𝒂 𝒅𝒊𝒗 "𝒎𝒂𝒊𝒏" 𝒅𝒆 𝒍𝒂 𝒑𝒂𝒈𝒆 𝒅𝒆𝒕𝒂𝒊𝒍.𝒉𝒕𝒎𝒍 : 
async function printHTML() {
    let main = document.getElementById("main2");
    if (main != null) {
        let name = document.getElementById("band").textContent;
        await fetchData(name);

        for (let id in data.query.pages) {
            let i = 900; // Je garde, au minimum, les 900 premiers caractères de la description Wikipédia.
            while (data.query.pages[id].extract[i] != "." && i < 1500) { // Tant que le caractère d'indice i n'est pas une point (fin d'une phrase), i++.
                i++;
            }
            // Je récupère la description Wikipédia (data.query.pages[id].extract),
            // je n'en garde que les i premiers caractères (.substring(0, i + 1)), 
            // et je supprime toute information indésirable entre crochets (phonétique, référence Wiki de type [1], etc.)
            main.textContent = data.query.pages[id].extract.substring(0, i + 1).replace(/\[(.*?)\]/g, '');
        }
    }
}

printHTML();

function modifyWikiURL() {
    if (document.getElementById("band") != null) {
        let name = document.getElementById("band").textContent;
        let iframeWiki = document.getElementById("iframeWiki");
        switch (name) {
            case "SOJA":
                iframeWiki.src = "https://fr.wikipedia.org/wiki/" + "Soldiers of Jah Army" + "?printable=yes";
                break;
            case "Genesis":
            case "NWA":
            case "Muse":
                iframeWiki.src = "https://fr.wikipedia.org/wiki/" + name + " (groupe)" + "?printable=yes";
                break;
            case "R3HAB":
                iframeWiki.src = "https://fr.wikipedia.org/wiki/" + "Fadil El Ghoul" + "?printable=yes";
                break;
        }
    }
}

modifyWikiURL();

// 𝑭𝒐𝒓𝒎𝒂𝒕𝒕𝒂𝒈𝒆 𝒅𝒆𝒔 𝒅𝒂𝒕𝒆𝒔 𝒅𝒆 𝒄𝒐𝒏𝒄𝒆𝒓𝒕 :
function formatCities() {
    if (artists[0] == null) {
        onMainPage = false;
    }

    if (!onMainPage) {
        let cities = document.getElementsByClassName("city");

        for (let i = 0; i < cities.length; i++) {
            cities[i].innerHTML = cities[i].textContent.replaceAll("-", ", "); // Remplacement des tirets par des virgules
            cities[i].innerHTML = cities[i].textContent.replaceAll("_", " "); // Remplacement des tirets bas par des espaces
            cities[i].innerHTML = cities[i].textContent.replaceAll("uk ", "UK "); // Capitalisation de "UK" et "USA"
            cities[i].innerHTML = cities[i].textContent.replaceAll("usa ", "USA ");
            cities[i].innerHTML = titleCase(cities[i].textContent); // La 1ère lettre de chaque mot est mise en majuscule
        };
    }
};

function formatDates() {
    if (artists[0] == null) {
        onMainPage = false;
    }

    if (!onMainPage) {
        if (document.getElementsByClassName("concert-dates") != null) {
            let dates = document.getElementsByClassName("concert-dates");

            for (let i = 0; i < dates.length; i++) {
                dates[i].innerHTML = dates[i].innerHTML.replaceAll("[", "");
                dates[i].innerHTML = dates[i].innerHTML.replaceAll("]", "");
                dates[i].innerHTML = dates[i].innerHTML.replaceAll("-", ".");
            };

            for (let i = 0; i < dates.length; i++) {
                dates[i].innerHTML = dates[i].innerHTML.replaceAll(" ", "<br>");
            };

            let firstAlbum = document.getElementById("first-album");
            firstAlbum.innerHTML = firstAlbum.innerHTML.replaceAll("-", ".");
        }
    }
};

formatCities();
formatDates();

function modifyOpacity(str) {
    switch (str) {
        case 'up':
            document.getElementById("map-banner").style.height = "0px";
            document.getElementById('map-subcontainer').style.opacity = 1;
            break;
        case 'down':
            document.getElementById("map-banner").style.height = "71px";
            document.getElementById('map-subcontainer').style.opacity = 0.0;
            break;
    }
}

function expand(element) {
    // Fonction que fait s'ouvrir / se fermer un iFrame selon qu'il est déjà ouvert ou fermé (toggle). Fait aussi se fermer le 2nd iFrame si on souhaite ouvrir le 1er (pour ne pas surcharger l'écran).
    // Fait aussi varier l'opacité de la map (car on ne peut pas faire varier sa taille)

    let iframeGoogle = document.getElementById("iframeGoogle");
    let iframeWiki = document.getElementById("iframeWiki");
    let banner = document.getElementById("map-banner");
    let map = document.getElementById("map-subcontainer");
    let wrapper = document.querySelector(".wrapper");

    switch (element) {
        case "wiki":
            if (iframeWiki.style.height == "") { // Par défaut, height == "" pour JavaScript car aucune taille n'est spécifié dans l'HTML. En revanche, elle est bien spécifiée dans le CSS ! (height: 95px)
                iframeWiki.style.height = "600px"; // J'ouvre Wiki
                iframeGoogle.style.height = ""; // Je ferme Google
                map.style.opacity = ""; // Je ferme la Map
                banner.style.height = "";
                wrapper.scrollBy({ // La div 'wrapper' scroll jusqu'à la div 'wiki' (moins une marge de 150)
                    top: document.querySelector("#wiki").getBoundingClientRect().top - 150,
                    left: 0,
                    behavior: 'smooth'
                });
            } else if (iframeWiki.style.height = !"") {
                iframeWiki.style.height = "";
            };
            break;
        case "google":
            if (iframeGoogle.style.height == "") {
                iframeGoogle.style.height = "600px"; // J'ouvre Google
                iframeWiki.style.height = ""; // Je ferme Wiki
                map.style.opacity = ""; // Je ferme la Map
                banner.style.height = "";
                wrapper.scrollBy({ // La div 'wrapper' scroll jusqu'à la div 'wiki' (moins une marge de 150)
                    top: document.querySelector("#wiki").getBoundingClientRect().top - 150,
                    left: 0,
                    behavior: 'smooth'
                });
            } else if (iframeGoogle.style.height != "") {
                iframeGoogle.style.height = "";
            };
            break;
        case "map":
            if (map.style.opacity == "" && banner.style.height == "") {
                map.style.opacity = 1; // J'ouvre la map
                banner.style.height = "0px";
                iframeWiki.style.height = ""; // Je ferme Wiki
                iframeGoogle.style.height = ""; // Je ferme Google
                wrapper.scrollBy({ // La div 'wrapper' scroll jusqu'à la div 'wiki' (moins une marge de 150)
                    top: document.querySelector("#wiki").getBoundingClientRect().top,
                    left: 0,
                    behavior: 'smooth'
                });
            } else {
                map.style.opacity = "";
                banner.style.height = "";
            };
            break;
    }
};


function titleCase(str) {
    str = str.split(' ');

    for (var i = 0; i < str.length; i++) {
        str[i] = str[i].charAt(0).toUpperCase() + str[i].slice(1);
    }

    return str.join(' ');
}


// -------------------- // 𝑴𝒂𝒑 ⭐ 𝑩𝒐𝒙 \\ --------------------
var map = L.mapbox.map('map')
L.mapbox.accessToken = 'pk.eyJ1IjoidGVuZWJyb3MiLCJhIjoiY2tub2QybHZmMHhyejJ3bGl2eDZvOHl3aCJ9.NQSuw0AwXS1dI04rHVKAAw';

function newLocation(value) {
    if (value != "0") {
        load(value).then(x => {
            map.setView([x.features[0].center[1], x.features[0].center[0]], 9)
        })
    }
}
async function load(value) {
    return await fetch('https://api.mapbox.com/geocoding/v5/mapbox.places/' + value + '.json?access_token=pk.eyJ1IjoidGVuZWJyb3MiLCJhIjoiY2tub2QybHZmMHhyejJ3bGl2eDZvOHl3aCJ9.NQSuw0AwXS1dI04rHVKAAw').then(res => res.json());
}
// Les petites pastilles \\
var x = document.getElementById("slct");
var optionValues = [];
for (let i = 1; i < x.length; i++) {
    optionValues.push(x.options[i].text);
}
var geojson = [];
optionValues.forEach(async element => await load(element).then(x => { geojson.push(x.features[0]); if (optionValues.length == geojson.length) { secondAsync() } }))

async function secondAsync() {
    map
        .setView([geojson[0].geometry.coordinates[1], geojson[0].geometry.coordinates[0]], 2)
        .addLayer(L.mapbox.styleLayer('mapbox://styles/mapbox/streets-v11'));
    geojson.forEach(function(marker) {
        L.mapbox.featureLayer({
            type: 'Feature',
            geometry: {
                type: 'Point',
                coordinates: marker.geometry.coordinates
            },
            properties: {
                title: '',
                description: '',
                'marker-size': 'large',
                'marker-color': '#BE9A6B',
                'marker-symbol': 'music'
            }
        }).on('click', function(e) {
            map.setView(e.latlng, 11);
        }).addTo(map);
    });
}
// -------------------------- \\ ⭐ // --------------------------