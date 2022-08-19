package check_passport

import (
	"bytes"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_sortNumbers(t *testing.T) {
	var (
		buf     []byte
		bufSort []byte
		err     error
	)

	_ = os.RemoveAll(testDst)
	err = os.MkdirAll(testDst, dirPermission)
	require.Nil(t, err)

	filename := filepath.Join(testDst, "sort")

	data := []int{
		124320,
		987652,
		231432,
		753443,
	}

	buf, err = getBytes(data)
	require.Nil(t, err)

	err = ioutil.WriteFile(filename, buf, fs.ModePerm)
	require.Nil(t, err)

	err = sortNumbers(testDst)
	require.Nil(t, err)

	sort.Sort(sort.IntSlice(data))
	bufSort, err = getBytes(data)
	require.Nil(t, err)

	buf, err = ioutil.ReadFile(filename)
	require.Nil(t, err)

	require.Equal(t, bytes.Compare(buf, bufSort), 0)

	_ = os.RemoveAll(testDst)
}

func getBytes(data []int) (buf []byte, err error) {
	buf = []byte{}
	for _, num := range data {
		var b []byte
		if b, err = numberToByte(uint32(num)); err != nil {
			return
		}
		buf = append(buf, b...)
	}

	return
}
