package service

import (
	"database/sql"
	"github.com/atongen/gosaic/model"
)

type GidxService struct {
	DB *sql.DB
}

func NewGidxService(db *sql.DB) *GidxService {
	return &GidxService{DB: db}
}

func (gidxService *GidxService) FindById(id int64) (*model.Gidx, error) {
	gidx := &model.Gidx{Id: id}
	rows, err := gidxService.DB.Query("select path, md5sum, width, height from gidx where id = ? limit 1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&gidx.Path, &gidx.Md5sum, &gidx.Width, &gidx.Height)
		if err != nil {
			return nil, err
		}
	}
	err = rows.Err()
	if err != nil {
		return gidx, err
	}
	return gidx, nil
}

func (gidxService *GidxService) FindGidxByMd5sum(md5sum string) (*model.Gidx, error) {
	gidx := &model.Gidx{Md5sum: md5sum}
	rows, err := gidxService.DB.Query("select id, path, width, height from gidx where md5sum = ? limit 1", md5sum)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&gidx.Id, &gidx.Path, &gidx.Width, &gidx.Height)
		if err != nil {
			return nil, err
		}
	}
	err = rows.Err()
	if err != nil {
		return gidx, err
	}
	return gidx, nil
}

func (gidxService *GidxService) ExistsByMd5sum(md5sum string) (bool, error) {
	var exists int
	rows, err := gidxService.DB.Query("select 1 from gidx where md5sum = ? limit 1", md5sum)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&exists)
		if err != nil {
			return false, err
		}
	}
	err = rows.Err()
	if err != nil {
		return false, err
	}
	if exists == 1 {
		return true, nil
	} else {
		return false, nil
	}

}

func (gidxService *GidxService) Create(gidx *model.Gidx) error {
	stmt, err := gidxService.DB.Prepare("INSERT INTO gidx(path, md5sum, width, height) VALUES(?, ?, ?, ?)")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(gidx.Path, gidx.Md5sum, gidx.Width, gidx.Height)
	if err != nil {
		return err
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		return err
	}
	gidx.Id = lastId
	return nil
}

func (gidxService *GidxService) Update(gidx *model.Gidx) error {
	stmt, err := gidxService.DB.Prepare("UPDATE gidx set path = ?, md5sum = ?, width = ?, height = ? where id = ? limit 1")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(gidx.Path, gidx.Md5sum, gidx.Width, gidx.Height, gidx.Id)
	if err != nil {
		return err
	}
	return nil
}

func (gidxService *GidxService) Delete(id int64) error {
	stmt, err := gidxService.DB.Prepare("delete from gidx where id = ? limit 1")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}
