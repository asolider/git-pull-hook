package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var projects map[string]string

func init() {
	/*
		bytes, err := ioutil.ReadFile("./config.json")

		if err != nil {
			log.Fatalf("read config err: %s \n", err)
		}

		err = json.Unmarshal(bytes, &projects)
		if err != nil {
			log.Fatalf("format json config err: %s \n", err)
		}
	*/

	initLog()

	initConfig()

}

func initConfig() {
	config := viper.New()
	config.SetConfigFile("./config.yaml")
	err := config.ReadInConfig()
	if err != nil {
		log.Fatalf("read config e: %s", err)
	}

	config.Unmarshal(&projects)

	config.OnConfigChange(func(e fsnotify.Event) {
		config.Unmarshal(&projects)
		log.Printf("reload config: %s", projects)
	})
	config.WatchConfig()

	log.Printf(" config : %s \n", projects)
}

func initLog() {
	logFile, err := os.OpenFile("./run.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile)
	log.SetPrefix("[run-git-pull] ")
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
	return
}

func main() {
	http.HandleFunc("/git_pull", func(w http.ResponseWriter, r *http.Request) {
		project := r.URL.Query().Get("project")
		log.Printf("run git pull: %s \n", project)

		projectDir, exist := projects[project]

		if project == "" || !exist {
			log.Printf("project not config: %s \n", project)
			w.Write([]byte("project not config"))
			return
		}

		cmdStr := fmt.Sprintf("cd %s && git pull", projectDir)
		log.Printf("exec: %s \n", cmdStr)
		cmd := exec.Command("sh", "-c", cmdStr)
		out, err := cmd.CombinedOutput()

		if err != nil {
			log.Printf("git pull err: %s \n", err)
			log.Printf("\n %s \n", out)
		} else {
			log.Printf("\n %s \n", out)
		}
		w.Write(out)
	})

	http.ListenAndServe(":8081", nil)
}
