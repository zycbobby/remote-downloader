package controllers
import (
	"github.com/astaxie/beego"
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
	. "github.com/zycbobby/auto-downloader/models"
	"os/exec"
	"log"
	"github.com/astaxie/beego/context"
)

var tasks = make(map[string]Task)

type TaskController struct {
	beego.Controller
}

// this function will be called every time a new request come (to initialize the request and responseWriter)
func (self *TaskController) Init(ctx *context.Context, controllerName, actionName string, app interface{}) {
	self.Controller.Init(ctx, controllerName, actionName, app)
}

func (self *TaskController) Get() {
	w := self.Ctx.ResponseWriter
	r := self.Ctx.Request

	_taskId := self.Ctx.Input.Param(":taskId")
	if "" != _taskId {
		self.Ctx.WriteString("get task : " + _taskId)
	} else {
		taskUrl := r.URL.Query().Get("taskUrl")
		switch taskUrl {

		// taskUrl not present in the query
		case "":
			if js, err := json.Marshal(tasks); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.Write(js)
			}
		default:
			if t, present := tasks[taskUrl]; present {
				if js, err := json.Marshal(t); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				} else {
					w.Header().Set("Content-Type", "application/json")
					w.Write(js)
				}
			} else {
				http.Error(w, "Task not register", http.StatusNotFound)
			}

		}
	}
}

func (self *TaskController) Post() {
	w := self.Ctx.ResponseWriter
	r := self.Ctx.Request
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
		tasks[t.Url] = t


		if js, err := json.Marshal(t); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)

			// start to handle
			go func(task *Task) {
				t.Status = STATUS_DOWNLOADING
				cmd := exec.Command("aria2c", t.Url)
				err := cmd.Start()
				if err != nil {
					log.Fatal(err)
				}
				err = cmd.Wait()
				// use beego log instead
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