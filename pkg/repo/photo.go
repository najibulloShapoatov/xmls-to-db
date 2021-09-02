package repo

import (
	"database/sql"
	"photouploader/pkg/db"
	"photouploader/pkg/log"

	//"photouploader/pkg/log"
	"time"

	"github.com/google/uuid"
)

//Photo  ...
type Photo struct {
	ID           int64
	ExternalID   string
	UID          uuid.UUID
	Index        int
	VehiclePlate string
	ServerIP     string
	Address      string
	VTime        time.Time
	OriginalPath string
	FilePath     string
	UploadedAt   time.Time
	FileRootPath string
}

/*
//CreateOrUpdate ...
func (p *Photo) CreateOrUpdate() {

	db := db.GetDB()
	var ph Photo

	err := db.Where("external_id = ? AND index = ?", p.ExternalID, p.Index).First(&ph).Error
	if err != nil {
		db.Create(&p)
		//log.Info("not found data=>", p.ExternalID, err, "creating data =>", p.ExternalID, p.Index, p.OriginalPath)

	} else {

		ph.ExternalID = p.ExternalID
		ph.UID = p.UID
		ph.Index = p.Index
		ph.VehiclePlate = p.VehiclePlate
		ph.ServerIP = p.ServerIP
		ph.Address = p.Address
		ph.VTime = p.VTime
		ph.OriginalPath = p.OriginalPath
		ph.FilePath = p.FilePath
		ph.UploadedAt = p.UploadedAt
		db.Save(&ph)
		p = &ph
		//log.Info("updating data =>", p.ExternalID, p.Index, p.OriginalPath)
	}
}
*/
//Save ....
func (p *Photo) Save() {

	query := `select id from photos where external_id = $1 AND index = $2`

	err := db.GetDB().QueryRow(query, p.ExternalID, p.Index).Scan(&p.ID)
	if err == sql.ErrNoRows {
		p.Insert()
	} else {
		p.Update()
	}
	log.Debug("func (p *Photo) Save() ->", err)

}

func (p *Photo) Insert() error {
	db := db.GetDB()

	query := `INSERT INTO "public"."photos"( "external_id", "uid", "index", "vehicle_plate", "server_ip", "address", "v_time", "original_path", "file_path", "uploaded_at", "file_root_path") 
	VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);`
	_, err := db.Exec(query,
		p.ExternalID,
		p.UID,
		p.Index,
		p.VehiclePlate,
		p.ServerIP,
		p.Address,
		p.VTime,
		p.OriginalPath,
		p.FilePath,
		p.UploadedAt,
		p.FileRootPath,
	)
	if err != nil {
		log.Error("func (p *Photo) Insert() error", err)
	}
	return err
}

func (p *Photo) Update() error {

	query := `UPDATE "public"."photos" SET "external_id" = $1, "uid" = $2, "index" = $3, 
	"vehicle_plate" = $4, "server_ip" = $5, "address" = $6, "v_time" = $7, 
	"original_path" = $8, "file_path" = $9, "uploaded_at" = $10, "uploaded_at" = $10
	 WHERE "id" = $11;
	`
	_, err := db.GetDB().Exec(query,
		p.ExternalID,
		p.UID,
		p.Index,
		p.VehiclePlate,
		p.ServerIP,
		p.Address,
		p.VTime,
		p.OriginalPath,
		p.FilePath,
		p.UploadedAt,
		p.ID,
	)
	if err != nil {
		log.Error("func (p *Photo) Update() error", err)
	}
	return err
}
