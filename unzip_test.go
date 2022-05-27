package check_passport

import (
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

func Test_Unzip(t *testing.T) {
	err := testDB.Unzip(filepath.Join(testDst, archiveName), filepath.Join(testDst, testDirName))
	require.Nil(t, err)
}
