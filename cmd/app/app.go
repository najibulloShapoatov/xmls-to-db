package app

import (
	"photouploader/pkg/config"
	"photouploader/pkg/db"
	"photouploader/pkg/log"
	"photouploader/pkg/reader"
	"photouploader/pkg/uploader"
	"strings"
	"time"
)

const (
	logFileName = "./log/uploader_%s.log"
)

//Run ..
func Run() error {

	//Initialize config
	{
		config.Init("./default.conf")
	}

	//Initalize logger
	{
		logFileName := strings.Replace(logFileName, "%s", time.Now().Format("20060102_150405"), 1)
		log.Init(logFileName, config.GetLogLevel())
	}

	//Load database configs
	{
		dbconf := config.LoadDBConfigs("DB")
		db.Init(&dbconf)
	}

	fileLocation := config.GetFileLocation()
	fileArchive := config.GetFileArchive()
	fileRootPath := config.GetFileRootPath()

	files, err := reader.GetFiles(fileLocation)
	if err != nil {
		log.Error("Get Files returned error =>", err)
		panic(err)
	}

	return uploader.Upload(files, fileLocation, fileArchive, fileRootPath)
}
