package check_passport

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

const (
	// http://проверки.гувм.мвд.рф/info-service.htm?sid=2000
	website    = "http://xn--b1afk4ade4e.xn--b1ab2a0a.xn--b1aew.xn--p1ai/info-service.htm?sid=2000"
	archiveUrl = "https://проверки.гувм.мвд.рф/upload/expired-passports/list_of_expired_passports.csv.bz2"
)

var (
	reNumberPassports = regexp.MustCompile(`(?si)<h3[^>]*>.+?Скорректировано\s+сведений\s+в.+?(\d+)`)
)

func (db *DB) Update(ctx context.Context) (ok bool, numPassports int, err error) {
	// есть ли обновление
	if ok, numPassports, err = db.UpdateAvailable(ctx, db.dst); err != nil || !ok {
		return
	}

	if err = os.MkdirAll(db.dst, dirPermission); err != nil {
		return
	}

	src := filepath.Join(db.dst, archiveName)
	if err = db.downloadFile(ctx, src, archiveUrl); err != nil {
		return
	}

	newDir := dirName + "_" + strconv.Itoa(numPassports)
	dstUpdate := filepath.Join(db.dst, newDir)
	if err = db.Unzip(src, dstUpdate); err != nil {
		return
	}

	if err = sortNumbers(dstUpdate); err != nil {
		return
	}

	_ = os.Remove(src)

	err = ioutil.WriteFile(
		filepath.Join(db.dst, infoFile), []byte(strconv.Itoa(numPassports)), os.ModePerm)
	if err != nil {
		return
	}

	// обновить ссылку на директорию
	symlink := filepath.Join(db.dst, symlinkCurrent)
	oldDir := ""
	if oldDir, err = os.Readlink(symlink); err == nil {
		if err = os.Remove(symlink); err != nil {
			return
		}
		_ = os.RemoveAll(filepath.Join(db.dst, oldDir))
	}

	err = os.Symlink(newDir, filepath.Join(db.dst, symlinkCurrent))
	return
}

//UpdateAvailable есть ли обновления и сколько паспортов в БД
func (db *DB) UpdateAvailable(ctx context.Context, dst string) (ok bool, numPassports int, err error) {
	var (
		req  *http.Request
		resp *http.Response
	)

	if ctx == nil {
		req, _ = http.NewRequest(http.MethodGet, website, nil)
	} else {
		req, _ = http.NewRequestWithContext(ctx, http.MethodGet, website, nil)
	}

	if resp, err = db.client.Do(req); err != nil {
		return
	}
	if resp == nil {
		err = errors.New("Response is nil. ")
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	b, _ := ioutil.ReadAll(resp.Body)
	m := reNumberPassports.FindSubmatch(b)
	if len(m) != 2 {
		err = errors.New("Number passports not found. ")
		return
	}

	var np int64
	if np, err = strconv.ParseInt(string(m[1]), 10, 64); err != nil {
		return
	}
	numPassports = int(np)

	fp := filepath.Join(dst, infoFile)
	if _, err = os.Stat(fp); os.IsNotExist(err) {
		err = nil
		ok = true
		return
	}

	b, _ = ioutil.ReadFile(fp)
	lastNumPassports, _ := strconv.ParseInt(string(b), 10, 64)

	ok = numPassports != int(lastNumPassports)
	return
}

func (db *DB) downloadFile(ctx context.Context, filepath string, url string) (err error) {
	var (
		out  *os.File
		resp *http.Response
		req  *http.Request
	)

	if out, err = os.Create(filepath); err != nil {
		return
	}

	defer func() {
		_ = out.Close()
		if err != nil {
			err = NewErrArchiveDownload(err.Error())
		}
	}()

	if ctx == nil {
		req, _ = http.NewRequest(http.MethodGet, url, nil)
	} else {
		req, _ = http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	}

	if resp, err = db.client.Do(req); err != nil {
		return
	}

	if resp == nil {
		err = errors.New("Response is nil. ")
		return
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", resp.Status)
		return
	}

	_, err = io.Copy(out, resp.Body)
	return
}
