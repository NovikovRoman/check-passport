package check_passport

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func Test_Update(t *testing.T) {
	ok, numPassports, err := testDB.Update(nil)
	require.Nil(t, err)
	assert.True(t, ok)
	assert.Greater(t, numPassports, 0)

	err = os.Remove(filepath.Join(testDst, infoFile))
	require.Nil(t, err)
}

func Test_UpdateAvailable(t *testing.T) {
	setupForTest()

	ok, numPassports, err := testDB.UpdateAvailable(nil, testDst)
	require.Nil(t, err)
	require.Greater(t, numPassports, 0)
	require.True(t, ok)

	fp := filepath.Join(testDst, infoFile)
	_ = ioutil.WriteFile(fp, []byte("10000"), os.ModePerm)
	ok, _, err = testDB.UpdateAvailable(nil, testDst)
	assert.Nil(t, err)
	assert.True(t, ok)

	_ = ioutil.WriteFile(fp, []byte(strconv.Itoa(numPassports)), os.ModePerm)
	ok, _, err = testDB.UpdateAvailable(nil, testDst)
	assert.Nil(t, err)
	assert.False(t, ok)

	err = os.Remove(fp)
	require.Nil(t, err)
}

func Test_downloadFile(t *testing.T) {
	setupForTest()

	err := os.MkdirAll(testDst, dirPermission)
	require.Nil(t, err)

	err = testDB.downloadFile(nil, filepath.Join(testDst, archiveName), archiveUrl)
	require.Nil(t, err)
}
