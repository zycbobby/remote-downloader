package main
import (
	"net/http"
	"fmt"
	"log"
	"time"
	"encoding/json"
	"io/ioutil"
	"os/exec"
)


const (
	STATUS_READY = iota
	STATUS_DOWNLOADING
	STATUS_FAIL
	STATUS_FINISH
)

type STATUS int

var GLOBAL_TASK_ID_INCREMENTOR = 1

type Task struct {
	Id int
	Url string
	Status STATUS
	DownloadLink string // mapping to the file system
}

type DownloadHandler struct {
	tasks map[string]Task
}

func (self *DownloadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// define the routes here
	switch r.Method {
	case "GET":
		switch r.URL.Path {
		case "/tasks":
			if taskUrl := r.URL.Query().Get("taskUrl"); taskUrl != "" {
				if t, present := self.tasks[taskUrl]; present {
					if js, err := json.Marshal(t);err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					} else {
						w.Header().Set("Content-Type", "application/json")
						w.Write(js)
					}
				} else {
					http.Error(w, "Task not register", http.StatusNotFound)
				}
			} else {
				if js, err := json.Marshal(self.tasks);err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				} else {
					w.Header().Set("Content-Type", "application/json")
					w.Write(js)
				}
			}
		}
	case "POST":
		switch r.URL.Path {
		case "/tasks":
			fmt.Println("create tasks")
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				err = fmt.Errorf("read request body error:%v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			kv := make(map[string]interface{})
			if err = json.Unmarshal(b, &kv); err != nil {
				err = fmt.Errorf("Unmarshal body error:%v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if _url, ok := kv["url"]; ok {
				t := Task{
					Id : GLOBAL_TASK_ID_INCREMENTOR,
					Url : _url.(string),
					Status: STATUS_READY,
					DownloadLink: "",
				}
				GLOBAL_TASK_ID_INCREMENTOR++
				self.tasks[t.Url] = t


				if js, err := json.Marshal(t);err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				} else {
					w.Header().Set("Content-Type", "application/json")
					w.Write(js)

					// start to handle
					go func(task *Task) {
						t.Status = STATUS_DOWNLOADING
						cmd := exec.Command("aria2c", t.Url)
						err:=cmd.Start()
						if err != nil {
							log.Fatal(err)
						}
						err = cmd.Wait()
						log.Printf("Download %s finished with error: %v", task.Url, err)
						if nil == err {
							t.Status = STATUS_FINISH
						} else {
							t.Status = STATUS_FAIL
						}
					}(&t)
					return
				}
			} else {
				err = fmt.Errorf("url not present in the request")
				w.WriteHeader(http.StatusBadRequest)
				return
			}


		}

	}
}
func main() {
	s := &http.Server{
		Addr: ":5678",
		Handler: &DownloadHandler{
			tasks  : make(map[string]Task),
		},
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(s.ListenAndServe());
}
