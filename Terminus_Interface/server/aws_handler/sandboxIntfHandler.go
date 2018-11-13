package TERMINUS

import (
	"net/http"
	"io/ioutil"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func FormTreeFromDockerCompose(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Info("filename=", vars["filename"])
	log.Info("mainservicename=", vars["mainservicename"])

	urlExp := "http://sandbox_microservices:8083/formTreeFromDockerCompose/"+vars["filename"] +"/"+vars["mainservicename"]
	log.Info(urlExp)
	responseExp, err := http.Get(urlExp)
	if err != nil {
		log.Error(err)
		http.Error(w, "sandbox_microservices connection error", http.StatusBadRequest)
		return
	} else {
		defer responseExp.Body.Close()
		contentsExp, err := ioutil.ReadAll(responseExp.Body)
		if err != nil {
			log.Error(err)
			http.Error(w, "sandbox_microservices read error", http.StatusBadRequest)
			return
		}
		w.Write(contentsExp)
	}
}
func FormTreeFromDockerComposeMongo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Info("serviceName=", vars["serviceName"])
	log.Info("original=", vars["original"])

	urlExp := "http://sandbox_microservices:8083/formTreeFromDockerComposeMongo/"+vars["serviceName"] +"/"+vars["original"]
	log.Info(urlExp)
	responseExp, err := http.Get(urlExp)
	if err != nil {
		log.Error(err)
		http.Error(w, "sandbox_microservices connection error", http.StatusBadRequest)
		return
	} else {
		defer responseExp.Body.Close()
		contentsExp, err := ioutil.ReadAll(responseExp.Body)
		if err != nil {
			log.Error(err)
			http.Error(w, "sandbox_microservices read error", http.StatusBadRequest)
			return
		}
		w.Write(contentsExp)
	}
}
func GetModifiedDockerComposeYaml(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Info("serviceName=", vars["serviceName"])
	log.Info("original=", vars["original"])
	urlExp := "http://sandbox_microservices:8083/getModifiedDockerComposeYaml/"+vars["serviceName"] +"/"+vars["original"]
	log.Info(urlExp)
	responseExp, err := http.Get(urlExp)
	if err != nil {
		log.Error(err)
		http.Error(w, "sandbox_microservices connection error", http.StatusBadRequest)
		return
	} else {
		defer responseExp.Body.Close()
		contentsExp, err := ioutil.ReadAll(responseExp.Body)
		if err != nil {
			log.Error(err)
			http.Error(w, "sandbox_microservices read error", http.StatusBadRequest)
			return
		}
		w.Write(contentsExp)
	}
}

func GetModifiedServiceEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Info("serviceName=", vars["serviceName"])
	log.Info("original=", vars["original"])
	urlExp := "http://sandbox_microservices:8083/getModifiedServiceEndpoint/"+vars["serviceName"] +"/"+vars["original"]
	log.Info(urlExp)
	responseExp, err := http.Get(urlExp)
	if err != nil {
		log.Error(err)
		http.Error(w, "sandbox_microservices connection error", http.StatusBadRequest)
		return
	} else {
		defer responseExp.Body.Close()
		contentsExp, err := ioutil.ReadAll(responseExp.Body)
		if err != nil {
			log.Error(err)
			http.Error(w, "sandbox_microservices read error", http.StatusBadRequest)
			return
		}
		w.Write(contentsExp)
	}
}
func GetAvailableApps(w http.ResponseWriter, r *http.Request) {
	urlExp := "http://sandbox_microservices:8083/getAvailableApps"
	log.Info(urlExp)
	responseExp, err := http.Get(urlExp)
	if err != nil {
		log.Error(err)
		http.Error(w, "sandbox_microservices connection error", http.StatusBadRequest)
		return
	} else {
		defer responseExp.Body.Close()
		contentsExp, err := ioutil.ReadAll(responseExp.Body)
		if err != nil {
			log.Error(err)
			http.Error(w, "sandbox_microservices read error", http.StatusBadRequest)
			return
		}
		w.Write(contentsExp)
	}
}

func ParseYamlToYaml(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Info("filename=", vars["filename"])

	urlExp := "http://sandbox_microservices:8083/parseYamlToYaml/"+vars["filename"]
	log.Info(urlExp)
	responseExp, err := http.Get(urlExp)
	if err != nil {
		log.Error(err)
		http.Error(w, "sandbox_microservices connection error", http.StatusBadRequest)
		return
	} else {
		defer responseExp.Body.Close()
		contentsExp, err := ioutil.ReadAll(responseExp.Body)
		if err != nil {
			log.Error(err)
			http.Error(w, "sandbox_microservices read error", http.StatusBadRequest)
			return
		}
		w.Write(contentsExp)
	}
}
func ParseYamlToJson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Info("filename=", vars["filename"])

	urlExp := "http://sandbox_microservices:8083/parseyamltojson/"+vars["filename"]
	log.Info(urlExp)
	responseExp, err := http.Get(urlExp)
	if err != nil {
		log.Error(err)
		http.Error(w, "sandbox_microservices connection error", http.StatusBadRequest)
		return
	} else {
		defer responseExp.Body.Close()
		contentsExp, err := ioutil.ReadAll(responseExp.Body)
		if err != nil {
			log.Error(err)
			http.Error(w, "sandbox_microservices read error", http.StatusBadRequest)
			return
		}
		w.Write(contentsExp)
	}
}
func DeployserviceToRecordAPIAndResponse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Info("filename=", vars["filename"])
	log.Info("endpointapi=", vars["endpointapi"])
	log.Info("mainservicename=", vars["mainservicename"])

	urlExp := "http://sandbox_microservices:8083/deployserviceToRecordAPIAndResponse/"+vars["filename"]+"/"+vars["endpointapi"]+"/"+vars["mainservicename"]
	log.Info(urlExp)
	responseExp, err := http.Get(urlExp)
	if err != nil {
		log.Error(err)
		http.Error(w, "sandbox_microservices connection error", http.StatusBadRequest)
		return
	} else {
		defer responseExp.Body.Close()
		contentsExp, err := ioutil.ReadAll(responseExp.Body)
		if err != nil {
			log.Error(err)
			http.Error(w, "sandbox_microservices read error", http.StatusBadRequest)
			return
		}
		w.Write(contentsExp)
	}
}

func GetAllDockerComposeFileNames (w http.ResponseWriter, r *http.Request)  {

	urlExp := "http://sandbox_microservices:8083/getAllDockerComposeFileNames/"
	log.Info(urlExp)
	responseExp, err := http.Get(urlExp)
	if err != nil {
		log.Error(err)
		http.Error(w, "sandbox_microservices connection error", http.StatusBadRequest)
		return
	} else {
		defer responseExp.Body.Close()
		contentsExp, err := ioutil.ReadAll(responseExp.Body)
		if err != nil {
			log.Error(err)
			http.Error(w, "sandbox_microservices read error", http.StatusBadRequest)
			return
		}
		w.Write(contentsExp)
	}
}
func GetAllDockerComposeServiceNames (w http.ResponseWriter, r *http.Request)  {

	vars := mux.Vars(r)
	log.Info("filename=", vars["filename"])

	urlExp := "http://sandbox_microservices:8083/getServicesFromDockerCompose/"+vars["filename"]
	log.Info(urlExp)
	responseExp, err := http.Get(urlExp)
	if err != nil {
		log.Error(err)
		http.Error(w, "sandbox_microservices connection error", http.StatusBadRequest)
		return
	} else {
		defer responseExp.Body.Close()
		contentsExp, err := ioutil.ReadAll(responseExp.Body)
		if err != nil {
			log.Error(err)
			http.Error(w, "sandbox_microservices read error", http.StatusBadRequest)
			return
		}
		w.Write(contentsExp)
	}
}

func GetAllTreesFromDockerComposeMongo(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	log.Info("filename=", vars["filename"])
	log.Info("servicename=", vars["servicename"])

	urlExp := "http://sandbox_microservices:8083/getAllTreesFromDockerComposeMongo/"+vars["filename"]+"/"+vars["servicename"]
	log.Info(urlExp)
	responseExp, err := http.Get(urlExp)
	if err != nil {
		log.Error(err)
		http.Error(w, "sandbox_microservices connection error", http.StatusBadRequest)
		return
	} else {
		defer responseExp.Body.Close()
		contentsExp, err := ioutil.ReadAll(responseExp.Body)
		if err != nil {
			log.Error(err)
			http.Error(w, "sandbox_microservices read error", http.StatusBadRequest)
			return
		}
		w.Write(contentsExp)
	}
}

func GetAllTestsInformationSandbox(w http.ResponseWriter, r *http.Request)  {

	urlExp := "http://sandbox_microservices:8083/getAllTestsInformation"
	log.Info(urlExp)
	responseExp, err := http.Get(urlExp)
	if err != nil {
		log.Error(err)
		http.Error(w, "sandbox_microservices connection error", http.StatusBadRequest)
		return
	} else {
		defer responseExp.Body.Close()
		contentsExp, err := ioutil.ReadAll(responseExp.Body)
		if err != nil {
			log.Error(err)
			http.Error(w, "sandbox_microservices read error", http.StatusBadRequest)
			return
		}
		w.Write(contentsExp)
	}
}
