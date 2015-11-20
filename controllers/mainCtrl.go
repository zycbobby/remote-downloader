package controllers
import (
	"github.com/astaxie/beego"
	"encoding/json"
	"net/http"
)


type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	this.Data["Website"] = "beego.me"
	this.Data["Email"] = "astaxie@gmail.com"
	if js, err := json.Marshal(this.Data); err == nil {
		this.Ctx.WriteString(string(js))
	} else {
		this.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		this.Ctx.WriteString(err.Error())
	}
}