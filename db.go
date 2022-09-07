package check_passport

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

const (
	dirName        = "db"
	archiveName    = "arch.bz2"
	dirPermission  = 0775
	filePermission = 0664
	infoFile       = "lastupdate"
	fileExt        = ".bin"
	symlinkCurrent = "current"
)

var (
	reAllowableSeries = regexp.MustCompile(`(?si)^\d{4}$`)
	reAllowableNumber = regexp.MustCompile(`(?si)^\d{6}$`)
)

type DB struct {
	client *http.Client
	dst    string
}

func NewDB(dst string, transport *http.Transport) (db *DB) {
	db = &DB{
		client: http.DefaultClient,
		dst:    dst,
	}

	if transport != nil {
		db.client.Transport = transport
	}
	return
}

//Exists существовует ли БД
func (db *DB) Exists() bool {
	symlink := filepath.Join(db.dst, symlinkCurrent)
	_, err := os.Stat(symlink)
	return err == nil
}

//AllowableSeries допустимая серия паспорта
func (db *DB) AllowableSeries(series string) (ok bool) {
	if ok = reAllowableSeries.MatchString(series); !ok {
		return
	}
	_, ok = OkatoRegions[series[0:2]]
	return
}

//AllowableNumber допустимый номер паспорта
func (db *DB) AllowableNumber(number string) (ok bool) {
	if ok = reAllowableNumber.MatchString(number); !ok {
		return
	}

	n, _ := strconv.ParseInt(number, 10, 64)
	return n > 100
}

func (db *DB) IsValid(series, number string) (ok bool, err error) {
	if !db.AllowableSeries(series) || !db.AllowableNumber(number) {
		return
	}

	n, _ := strconv.ParseInt(number, 10, 64)
	numberCheck := uint32(n)

	symlink := filepath.Join(db.dst, symlinkCurrent)
	if _, err = os.Stat(symlink); os.IsNotExist(err) {
		err = errors.New("DB not found. ")
		return
	}

	var (
		fi   os.FileInfo
		file *os.File
	)
	fp := filepath.Join(symlink, series+fileExt)
	if fi, err = os.Stat(fp); os.IsNotExist(err) {
		err = nil
		ok = true
		return
	}

	if file, err = os.Open(fp); err != nil {
		return
	}

	defer func() {
		_ = file.Close()
	}()

	start := int64(0)
	// Количество номеров. Номер занимает 4 байта.
	finish := fi.Size() / 4

	// поиск номера делением списка
	for {
		if start > finish || finish < start {
			ok = true
			return
		}

		middle := start + (finish-start)/2

		if _, err = file.Seek(middle*4, 0); err != nil {
			return
		}

		b := make([]byte, 4)
		if _, err = file.Read(b); err != nil {
			return
		}

		var num uint32
		if num, err = bytesToNumber(b); err != nil {
			return
		}

		if num == numberCheck { // найден в БД
			return
		}

		if num > numberCheck {
			finish = middle - 1

		} else {
			start = middle + 1
		}
	}
}

func (db *DB) IsInvalid(series, number string) (ok bool, err error) {
	if ok, err = db.IsValid(series, number); err == nil {
		ok = !ok
	}
	return
}
