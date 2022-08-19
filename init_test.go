package check_passport

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

const (
	testDst = "temptest"
	testData = "testdata"
)

func emulOldDB(dst string, numPassports int) (err error) {
	oldDB := strconv.Itoa(numPassports)
	oldDir := dirName + "_" + oldDB
	if err = os.MkdirAll(filepath.Join(dst, oldDir), dirPermission); err != nil {
		return
	}

	// lastupdate
	err = ioutil.WriteFile(filepath.Join(testDst, infoFile), []byte(oldDB), fs.ModePerm)
	if err != nil {
		return
	}

	symlink := filepath.Join(dst, symlinkCurrent)
	err = os.Symlink(oldDir, symlink)
	return
}
