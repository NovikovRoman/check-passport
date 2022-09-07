package check_passport

import (
	"context"
	"os"
	"path/filepath"
	"strconv"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Update(t *testing.T) {
	var (
		err          error
		ok           bool
		numPassports int
	)

	_ = os.RemoveAll(testDst)

	err = emulOldDB(testDst, 10000)
	require.Nil(t, err)

	ctx := context.Background()

	testDB := NewDB(testDst, nil)

	ok, numPassports, err = testDB.Update(ctx)
	require.Nil(t, err)
	assert.True(t, ok)
	assert.Greater(t, numPassports, 0)

	symlink := filepath.Join(testDst, symlinkCurrent)

	var dir string
	dir, err = os.Readlink(symlink)
	require.Nil(t, err)
	require.NotEqual(t, dir, dirName+strconv.Itoa(10000))

	_ = os.RemoveAll(testDst)
}

func Test_UpdateAvailable(t *testing.T) {
	var (
		ok           bool
		numPassports int
		err          error
	)

	_ = os.RemoveAll(testDst)
	ctx := context.Background()

	testDB := NewDB(testDst, nil)
	ok, numPassports, err = testDB.UpdateAvailable(ctx, testDst)
	require.Nil(t, err)
	require.Greater(t, numPassports, 0)
	require.True(t, ok)

	err = emulOldDB(testDst, 10000)
	require.Nil(t, err)

	ok, numPassports, err = testDB.UpdateAvailable(ctx, testDst)
	require.Nil(t, err)
	require.Greater(t, numPassports, 0)
	require.True(t, ok)

	_ = os.RemoveAll(testDst)
	err = emulOldDB(testDst, numPassports)
	require.Nil(t, err)

	ok, _, err = testDB.UpdateAvailable(ctx, testDst)
	require.Nil(t, err)
	assert.False(t, ok)

	_ = os.RemoveAll(testDst)
}

func Test_downloadFile(t *testing.T) {
	_ = os.RemoveAll(testDst)
	ctx := context.Background()

	testDB := NewDB(testDst, nil)
	err := os.MkdirAll(testDst, dirPermission)
	require.Nil(t, err)

	err = testDB.downloadFile(ctx, filepath.Join(testDst, archiveName), archiveUrl)
	require.Nil(t, err)

	_ = os.RemoveAll(testDst)
}
