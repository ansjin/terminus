package TERMINUS

import ("net/http"
	"text/template"
	log "github.com/sirupsen/logrus"
)
func Index(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("").ParseFiles("templates/index.html", "templates/base.html")
	// check your err
	if(err!=nil){
		log.Info("err")
	}else{
		err = tmpl.ExecuteTemplate(w, "base", "")
	}
}

func ApproachVMBootingTime(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("").ParseFiles("templates/approach_vm_booting_time.html", "templates/base.html")
	// check your err
	if(err!=nil){
		log.Info("err")
	}else{
		err = tmpl.ExecuteTemplate(w, "base", "")
	}
}
func LogsPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("").ParseFiles("templates/logs.html", "templates/base.html")
	// check your err
	if(err!=nil){
		log.Info("err")
	}else{
		err = tmpl.ExecuteTemplate(w, "base", "")
	}
}
func CompareLimits(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("").ParseFiles("templates/compare_limits.html", "templates/base.html")
	// check your err
	if(err!=nil){
		log.Info("err")
	}else{
		err = tmpl.ExecuteTemplate(w, "base", "")
	}
}
func InterInstancePerfCompare(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("").ParseFiles("templates/inter_instance_perf.html", "templates/base.html")
	// check your err
	if(err!=nil){
		log.Info("err")
	}else{
		err = tmpl.ExecuteTemplate(w, "base", "")
	}
}
func InterReplicasPerfCompare(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("").ParseFiles("templates/inter_replicas_perf.html", "templates/base.html")
	// check your err
	if(err!=nil){
		log.Info("err")
	}else{
		err = tmpl.ExecuteTemplate(w, "base", "")
	}
}
func MSCAnalyzed(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("").ParseFiles("templates/msc_analyzed.html", "templates/base.html")
	// check your err
	if(err!=nil){
		log.Info("err")
	}else{
		err = tmpl.ExecuteTemplate(w, "base", "")
	}
}
func MSCAnalyzedAll(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("").ParseFiles("templates/msc_analyzedAll.html", "templates/base.html")
	// check your err
	if(err!=nil){
		log.Info("err")
	}else{
		err = tmpl.ExecuteTemplate(w, "base", "")
	}
}
func MSCPredictionAll(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("").ParseFiles("templates/msc_prediction_all.html", "templates/base.html")
	// check your err
	if(err!=nil){
		log.Info("err")
	}else{
		err = tmpl.ExecuteTemplate(w, "base", "")
	}
}
func MSCAllTests(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("").ParseFiles("templates/msc_tests.html", "templates/base.html")
	// check your err
	if(err!=nil){
		log.Info("err")
	}else{
		err = tmpl.ExecuteTemplate(w, "base", "")
	}
}
func PodBootingTimeResults(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("").ParseFiles("templates/pod_boot_time.html", "templates/base.html")
	// check your err
	if(err!=nil){
		log.Info("err")
	}else{
		err = tmpl.ExecuteTemplate(w, "base", "")
	}
}
func Conduct_pod_booting_timetest(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("").ParseFiles("templates/conduct_pod_booting_timetest.html", "templates/base.html")
	// check your err
	if(err!=nil){
		log.Info("err")
	}else{
		err = tmpl.ExecuteTemplate(w, "base", "")
	}
}
func VMBootingTimeResults(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("").ParseFiles("templates/view_prev_results.html", "templates/base.html")
	// check your err
	if(err!=nil){
		log.Info("err")
	}else{
		err = tmpl.ExecuteTemplate(w, "base", "")
	}
}

func ConductTestMSCBruteForce(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("").ParseFiles("templates/msc_brute_force_test.html", "templates/base.html")
	// check your err
	if(err!=nil){
		log.Info("err")
	}else{
		err = tmpl.ExecuteTemplate(w, "base", "")
	}
}
func Sandboxing(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("").ParseFiles("templates/perform_sand_box.html", "templates/base.html")
	// check your err
	if(err!=nil){
		log.Info("err")
	}else{
		err = tmpl.ExecuteTemplate(w, "base", "")
	}
}
func SandboxingTree(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("").ParseFiles("templates/sand_box.html", "templates/base.html")
	// check your err
	if(err!=nil){
		log.Info("err")
	}else{
		err = tmpl.ExecuteTemplate(w, "base", "")
	}
}
func SandboxedYaml(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("").ParseFiles("templates/get_modified_yaml.html", "templates/base.html")
	// check your err
	if(err!=nil){
		log.Info("err")
	}else{
		err = tmpl.ExecuteTemplate(w, "base", "")
	}
}
func SandboxingTests(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("").ParseFiles("templates/sandbox_tests.html", "templates/base.html")
	// check your err
	if(err!=nil){
		log.Info("err")
	}else{
		err = tmpl.ExecuteTemplate(w, "base", "")
	}
}
func Swaggeryaml(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("").ParseFiles("docs/swagger/swagger.yaml")
	// check your err
	if(err!=nil){
		log.Info("err")
	}else{
		err = tmpl.ExecuteTemplate(w, "base", "")
	}
}
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("").ParseFiles("templates/404.html", "templates/base.html")
	// check your err
	if(err!=nil){
		log.Info("err")
	}else{
		err = tmpl.ExecuteTemplate(w, "base", "")
	}
}
