package check_passport

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
)

var (
	rePassport = regexp.MustCompile(`(?si)([0-9]{4}),([0-9]{6})`)
	reSeries   = regexp.MustCompile(`(?si)^[0-9]{4}$`)
)

func sortNumbers(dir string) (err error) {
	var (
		files   []os.DirEntry
		file    *os.File
		numbers []uint32
	)
	if files, err = os.ReadDir(dir); err != nil {
		return
	}

	for _, f := range files {
		if f.IsDir() || reSeries.MatchString(f.Name()) {
			continue
		}

		fp := filepath.Join(dir, f.Name())
		if file, err = os.OpenFile(fp, os.O_RDWR, 0); err != nil {
			err = fmt.Errorf("Open file %s %v ", f.Name(), err)
			return
		}

		fi, _ := f.Info()
		numbers = make([]uint32, fi.Size()/4)
		buf := make([]byte, 4)
		bufReader := bufio.NewReader(file)
		i := 0
		for {
			if _, err = bufReader.Read(buf); err == io.EOF {
				break
			}

			if err != nil {
				err = fmt.Errorf("buf.Read %s %v", f.Name(), err)
				_ = file.Close()
				return
			}

			var num uint32
			if num, err = bytesToNumber(buf); err != nil {
				err = fmt.Errorf("binary.Read %s %v", f.Name(), err)
				_ = file.Close()
				return
			}

			numbers[i] = num
			i++
		}

		sort.Slice(numbers, func(i, j int) bool {
			return numbers[i] < numbers[j]
		})

		var b []byte
		buf = make([]byte, fi.Size())
		for j, n := range numbers {
			if b, err = numberToByte(n); err != nil {
				err = fmt.Errorf("numberToByte %s %d %v ", file.Name(), n, err)
				return
			}

			for jj, bb := range b {
				buf[j*4+jj] = bb
			}
		}

		_, _ = file.Seek(0, 0)
		if _, err = file.Write(buf); err != nil {
			_ = file.Close()
			err = fmt.Errorf("Write %s %v ", file.Name(), err)
			return
		}

		err = file.Sync()
		_ = file.Close()
		if err != nil {
			err = fmt.Errorf("Sync %s %v ", file.Name(), err)
			return
		}
	}

	return
}

func numberToByte(num uint32) (b []byte, err error) {
	var bBuf bytes.Buffer
	if err = binary.Write(&bBuf, binary.BigEndian, num); err == nil {
		b = bBuf.Bytes()
	}
	return
}

func bytesToNumber(b []byte) (num uint32, err error) {
	err = binary.Read(bytes.NewReader(b), binary.BigEndian, &num)
	return
}
