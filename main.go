package main
import (
	"github.com/astaxie/beego"
	"github.com/zycbobby/auto-downloader/controllers"
)

func main() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/tasks/?:taskId", &controllers.TaskController{})

	// read conf/app.conf automatically
	beego.Run();
}
