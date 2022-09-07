package check_passport

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Unzip(t *testing.T) {
	_ = os.RemoveAll(testDst)

	testDB := NewDB(testDst, nil)

	err := os.MkdirAll(testDst, dirPermission)
	require.Nil(t, err)

	ctx := context.Background()
	err = testDB.downloadFile(ctx, filepath.Join(testDst, archiveName), archiveUrl)
	require.Nil(t, err)

	err = testDB.Unzip(filepath.Join(testDst, archiveName), testDst)
	require.Nil(t, err)

	_ = os.RemoveAll(testDst)
}
