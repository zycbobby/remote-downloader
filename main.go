package main
import (
	"github.com/astaxie/beego"
	"github.com/zycbobby/auto-downloader/controllers"
	"os"
	"errors"
)

func init() {
	downloadDir := beego.AppConfig.String("download_storage")

	err := _createDownloadDir(downloadDir)
	if err != nil {
		defer func() {
			msg := recover().(string)
			beego.BeeLogger.Error(msg)
			os.Exit(-1)
		}()
		panic(err.Error())
	}
}

func _createDownloadDir(path string) error {
	if fileInfo, err := os.Stat(path); err == nil && !fileInfo.IsDir() {
		return errors.New("file exists : " + path)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, os.ModeDir); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/tasks/?:taskId", &controllers.TaskController{})

	// read conf/app.conf automatically
	beego.Run();
}
