package check_passport

import (
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

func Test_sortNumbers(t *testing.T) {
	err := sortNumbers(filepath.Join(testDst, testDirName))
	require.Nil(t, err)
}
