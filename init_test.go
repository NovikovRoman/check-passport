package check_passport

const (
	testDst     = "temptest"
	testDirName = "db_test"
)

var testDB *DB

func init() {
	setupForTest()
}

func setupForTest() {
	if testDB == nil {
		testDB = NewDB(testDst, nil)
	}
}
