package TERMINUS

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
		"Init",
		strings.ToUpper("Post"),
		"/init",
		InitUserConfig,
	},

	Route{
		"ListAllInstances",
		strings.ToUpper("Get"),
		"/instances",
		ListAllInstances,
	},

	Route{
		"InitKubeCluster",
		strings.ToUpper("Get"),
		"/initKubeCluster",
		InitKubeCluster,
	},

	Route{
		"StartKubeClusterKops",
		strings.ToUpper("Get"),
		"/startKubeCluster",
		StartKubeClusterKops,
	},
	Route{
		"UpdateKubeCluster",
		strings.ToUpper("Get"),
		"/updateKubeCluster",
		UpdateKubeCluster,
	},

	Route{
		"ValidateKubeCluster",
		strings.ToUpper("Get"),
		"/validateKubeCluster",
		ValidateKubeCluster,
	},

	Route{
		"CreateDashboardKops",
		strings.ToUpper("Get"),
		"/createDashboardKops",
		CreateDashboardKops,
	},

	Route{
		"DeleteDashboardKops",
		strings.ToUpper("Get"),
		"/deleteDashboardKops",
		DeleteDashboardKops,
	},

	Route{
		"CreateHeapsterKops",
		strings.ToUpper("Get"),
		"/createHeapsterKops",
		CreateHeapsterKops,
	},

	Route{
		"DeleteHeapsterKops",
		strings.ToUpper("Get"),
		"/deleteHeapsterKops",
		DeleteHeapsterKops,
	},

	Route{
		"createMetricsConfigurationKops",
		strings.ToUpper("Get"),
		"/createMetricsConfigurationKops",
		CreateMetricsConfigurationKops,
	},

	Route{
		"DeleteMetricsConfigurationKops",
		strings.ToUpper("Get"),
		"/deleteMetricsConfigurationKops",
		DeleteMetricsConfigurationKops,
	},

	Route{
		"CreateApplicationKops",
		strings.ToUpper("Get"),
		"/createApplicationKops",
		CreateApplicationKops,
	},

	Route{
		"DeleteApplicationKops",
		strings.ToUpper("Get"),
		"/deleteApplicationKops",
		DeleteApplicationKops,
	},

	Route{
		"GetClusterInfoKops",
		strings.ToUpper("Get"),
		"/getClusterInfoKops",
		GetClusterInfoKops,
	},

	Route{
		"GetPodsRunningKops",
		strings.ToUpper("Get"),
		"/getPodsRunningKops",
		GetPodsRunningKops,
	},

	Route{
		"GetServicesRunningKops",
		strings.ToUpper("Get"),
		"/getServicesRunningKops",
		GetServicesRunningKops,
	},

	Route{
		"DeployAndStartLoadTesting",
		strings.ToUpper("Get"),
		"/deployAndStartLoadTesting",
		DeployAndStartLoadTesting,
	},
	Route{
		"GetExternalServiceIp",
		strings.ToUpper("Get"),
		"/getExternalServiceIp/{appname}",
		GetExternalServiceIp,
	},
	Route{
		"GetPasswordDashboard",
		strings.ToUpper("Get"),
		"/getPasswordDashboard",
		GetPasswordDashboard,
	},

	Route{
		"GetTokenDashboard",
		strings.ToUpper("Get"),
		"/getTokenDashboard",
		GetTokenDashboard,
	},

	Route{
		"DeleteKubeCluster",
		strings.ToUpper("Get"),
		"/deleteKubeCluster",
		DeleteKubeCluster,
	},

	Route{
		"RunLoadGen",
		strings.ToUpper("Get"),
		"/runLoadGen",
		RunLoadGen,
	},
	Route{
		"GetMetrics",
		strings.ToUpper("Get"),
		"/getMetrics",
		GetMetricsAPI,
	},
}
