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

var conf ServerConf

var projectMap map[string]string = make(map[string]string)

type Project struct {
	ProjectName string `mapstructure:"name"`
	ProjectPath string `mapstructure:"path"`
}

type ServerConf struct {
	ServerPort int `mapstructure:"server_port"`
	Projects   []Project
}

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

	config.Unmarshal(&conf)

	config.OnConfigChange(func(e fsnotify.Event) {
		config.Unmarshal(&conf)
		log.Printf("reload config: %v", conf)
	})
	config.WatchConfig()

	log.Printf(" config : %v \n", conf)
}

func initLog() {
	logFile, err := os.OpenFile("./run.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile)
	log.SetPrefix("[run-git-pull] ")
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
}

func getProjectPath(projectName string) (projectDir string, exist bool) {
	if projectName == "" || len(conf.Projects) == 0 {
		return "", false
	}

	for _, p := range conf.Projects {
		if p.ProjectName == projectName {
			return p.ProjectPath, true
		}
	}

	return "", false
}

func main() {
	http.HandleFunc("/git_pull", func(w http.ResponseWriter, r *http.Request) {
		project := r.URL.Query().Get("project")
		log.Printf("run git pull: %s \n", project)

		projectDir, exist := getProjectPath(project)

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
			log.Printf("%s \n", out)
		} else {
			log.Printf("\n %s \n", out)
		}
		w.Write(out)
	})
	fmt.Printf("listen on 0.0.0.0:%d \n", conf.ServerPort)
	http.ListenAndServe(fmt.Sprintf(":%d", conf.ServerPort), nil)
}
