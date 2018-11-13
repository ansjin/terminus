package sandbox_microservices

import (
	"net/http"
	"strings"
	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	// Declare the static file directory and point it to the directory we just made
	staticFileDirectory := http.Dir("./assets/")
	// Declare the handler, that routes requests to their respective filename.
	// The fileserver is wrapped in the `stripPrefix` method, because we want to
	// remove the "/assets/" prefix when looking for files.
	// For example, if we type "/assets/index.html" in our browser, the file server
	// will look for only "index.html" inside the directory declared above.
	// If we did not strip the prefix, the file server would look for "./assets/assets/index.html", and yield an error
	staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDirectory))
	// The "PathPrefix" method acts as a matcher, and matches all routes starting
	// with "/assets/", instead of the absolute route itself
	router.PathPrefix("/assets/").Handler(staticFileHandler).Methods("GET")

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hello World!")
	http.Redirect(w, r, "assets/", http.StatusSeeOther)
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},

	Route{
		"ParseYamlToJson",
		strings.ToUpper("GET"),
		"/parseyamltojson/{filename}",
		ParseYamlToJson,
	},
	Route{
		"GetServicesFromDockerCompose",
		strings.ToUpper("GET"),
		"/getServicesFromDockerCompose/{filename}",
		GetServicesFromDockerCompose,
	},

	Route{
		"ParseYamlToYaml",
		strings.ToUpper("GET"),
		"/parseYamlToYaml/{filename}",
		ParseYamlToYaml,
	},
	Route{
		"DeployserviceToRecordAPIAndResponse",
		strings.ToUpper("GET"),
		"/deployserviceToRecordAPIAndResponse/{filename}/{endpointapi}/{mainservicename}",
		DeployserviceToRecordAPIAndResponse,
	},
	Route{
		"FormTreeFromDockerCompose",
		strings.ToUpper("GET"),
		"/formTreeFromDockerCompose/{filename}/{mainservicename}",
		FormTreeFromDockerCompose,
	},
	Route{
		"Test",
		strings.ToUpper("GET"),
		"/test/{filename}/{mainservicename}",
		Test12,
	},
	Route{
		"FormTreeFromDockerComposeMongo",
		strings.ToUpper("GET"),
		"/formTreeFromDockerComposeMongo/{serviceName}/{original}/",
		FormTreeFromDockerComposeMongo,
	},
	Route{
		"GetAllTreesFromDockerComposeMongo",
		strings.ToUpper("GET"),
		"/getAllTreesFromDockerComposeMongo/{filename}/{serviceName}/",
		GetAllTreesFromDockerComposeMongo,
	},
	Route{
		"GetModifiedDockerComposeYaml",
		strings.ToUpper("GET"),
		"/getModifiedDockerComposeYaml/{serviceName}/{original}/",
		GetModifiedDockerComposeYaml,
	},

	Route{
		"GetAvailableApps",
		strings.ToUpper("GET"),
		"/getAvailableApps",
		GetAvailableApps,
	},


	Route{
		"getModifiedServiceEndpoint",
		strings.ToUpper("GET"),
		"/getModifiedServiceEndpoint/{serviceName}/{original}/",
		GetModifiedServiceEndpoint,
	},
	Route{
		"GetAllDockerComposeFileNames",
		strings.ToUpper("GET"),
		"/getAllDockerComposeFileNames/",
		GetAllDockerComposeFileNames,
	},


	Route{
		"GetAllTestsInformation",
		strings.ToUpper("GET"),
		"/getAllTestsInformation/",
		GetAllTestsInformation,
	},

}
