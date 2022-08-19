package check_passport

import (
	"bufio"
	"compress/bzip2"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func (db *DB) Unzip(src, dst string) (err error) {
	var (
		f       *os.File
		buf     io.Reader
		reader  io.Reader
		scanner *bufio.Scanner
		b       []byte
	)

	if err = os.MkdirAll(dst, dirPermission); err != nil {
		return
	}

	if f, err = os.OpenFile(src, 0, 0); err != nil {
		return
	}
	defer func() {
		_ = f.Close()
	}()

	cache := NewDataCache()
	buf = bufio.NewReader(f)
	reader = bufio.NewReader(bzip2.NewReader(buf))
	scanner = bufio.NewScanner(reader)
	scanner.Scan() // пропуск заголовков

	for scanner.Scan() {
		m := rePassport.FindStringSubmatch(scanner.Text())
		if len(m) != 3 {
			continue
		}

		if _, ok := OkatoRegions[m[1][0:2]]; !ok {
			continue
		}

		fPath := filepath.Join(dst, m[1]+fileExt)
		if f, err = cache.OpenFile(fPath); err != nil {
			err = fmt.Errorf("OpenFile %s %v", fPath, err)
			_ = cache.CloseAllFile()
			return
		}

		num, _ := strconv.ParseInt(m[2], 10, 32)
		if b, err = numberToByte(uint32(num)); err != nil {
			err = fmt.Errorf("NumberToByte %s %s%s", fPath, m[1], m[2])
			_ = cache.CloseAllFile()
			continue
		}

		if err = cache.AddData(fPath, b); err != nil {
			err = fmt.Errorf("AddData %s %s%s", fPath, m[1], m[2])
			_ = cache.CloseAllFile()
			continue
		}
	}

	if errs := cache.CloseAllFile(); len(errs) > 0 {
		sErr := make([]string, len(errs))
		i := 0
		for name, e := range errs {
			sErr[i] = fmt.Sprintf("Close %s %v", name, e)
			i++
		}
		err = errors.New(strings.Join(sErr, "\n"))
	}

	return
}
