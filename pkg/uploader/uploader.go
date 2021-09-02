package uploader

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"photouploader/pkg/log"
	"photouploader/pkg/repo"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

//RowData ...
type RowData struct {
	XMLName xml.Name `xml:"ROWDATA"`
	Rows    []Row    `xml:"ROW"`
}

// Row structure
type Row struct {
	XMLName           xml.Name `xml:"ROW"`
	EXTERNALID        string   `xml:"XH"`
	VEHICLEPLATE      string   `xml:"HPHM"`
	VIOLATIONTIME     string   `xml:"WFSJ"`
	VIOLATIONLOCATION string   `xml:"WFDZ"`
	INDEX             string   `xml:"ZJWJSX"`
	SERVERIP          string   `xml:"ZJWJIP"`
	FILEPATH          string   `xml:"ZJWJLJ"`
}

//Upload ...
func Upload(files []string, fileLocation, fileArchive string, fileRootPath string) error {

	log.Info("in upload folder _>", files, "files")
	//wg := sync.WaitGroup{}
	//wg.Add(len(files))

	for _, file := range files {

		//go func(file, fileLocation, fileArchive string) {
			//defer wg.Done()
			err := proceedFile(fileLocation + "//" + file, fileRootPath)
			if err != nil {
				log.Info(" proceedFile err ->", file)
			}
			err = fileMove(file, fileLocation, fileArchive)
			if err != nil {
				log.Info("err ->", file)
			}
		//}(file, fileLocation, fileArchive)

	}
	//wg.Wait()
	log.Info("uploaded files->", files)

	return nil
}

func fileMove(file, from, to string) error {

	year, month, day := time.Now().Date()

	var Path = to + fmt.Sprintf("//%v//%v//%v//", year, int(month), day)
	_, err := os.Stat(Path)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(Path, 0755)
		if errDir != nil {
			log.Error("Error not created folder ", err)
			return err
		}
	}
	from = from + "//" + file
	Path = Path + file
	err = os.Rename(from, Path)
	if err != nil {
		log.Error("Error not move file ", err)
		return err
	}

	return nil
}

func proceedFile(filename string, fileRootPath string) error {

	var data RowData

	// Open our xmlFile
	xmlFile, err := os.Open(filename)
	// if we os.Open returns an error then handle it
	if err != nil {
		return err
	}

	log.Info("Successfully file opened _> ", filename)

	// defer the closing of our xmlFile so that we can parse it later on
	defer xmlFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		return err
	}

	// we unmarshal our byteArray which contains our
	// xmlFiles content into 'users' which we defined above
	xml.Unmarshal(byteValue, &data)

	log.Info("readed _> ", len(data.Rows), " pices.")

	if err = SaveToDB(data.Rows, fileRootPath); err != nil {
		return err
	}
	log.Info("||||Uploaded!||||| _> ", filename)
	return nil

}

//SaveToDB ...
func SaveToDB(rows []Row, fileRootPath string) error {

	var err error

	//const size = 10_000
	const size = 10_00
	var parts int = len(rows) / size
	wg := sync.WaitGroup{}

	for i := 0; i <= parts; i++ {
		wg.Add(1)
		var l int = i * size
		var r int = (i + 1) * size
		if r > len(rows) {
			r = len(rows)
		}

		go func(data []Row, err error) {
			defer wg.Done()

			for _, row := range data {

				idx, err := strconv.Atoi(row.INDEX)
				if err != nil {
					log.Error("not parsed Index photo =>", row, err)
				}
				var filePath string
				year := strings.Split(row.FILEPATH, "/")[2][:4]
				if year == "WeiF" {
					filePath = strings.Replace(row.FILEPATH, "/capture", "", 1)
				} else {
					filePath = strings.Replace(row.FILEPATH, "capture", year, 1)
				}

				photo := &repo.Photo{
					ExternalID:   row.EXTERNALID,
					UID:          uuid.New(),
					Index:        idx,
					VehiclePlate: row.VEHICLEPLATE,
					ServerIP:     row.SERVERIP,
					Address:      row.VIOLATIONLOCATION,
					VTime:        convertToTime(row.VIOLATIONTIME, nil),
					OriginalPath: row.FILEPATH,
					FilePath:     filePath,
					UploadedAt:   time.Now(),
					FileRootPath: fileRootPath,
				}
				//photo.CreateOrUpdate()
				photo.Save()

			}
		}(rows[l:r], err)
	}
	wg.Wait()

	return nil
}

func convertToTime(s string, ls []string) time.Time {

	var lts = []string{"02.01.2006 15:04:05", "01.02.06", "01-02-06", "02.01.2006", "02.01.06", "2006-01-02", "06-01-02", "2006"}
	lts = append(lts, ls...)
	s = strings.TrimSpace(s)
	var bt time.Time
	var err error
	for _, l := range lts {
		bt, err = time.Parse(l, s)
		if err == nil {
			return bt
		}
	}

	return bt
}
