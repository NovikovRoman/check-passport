package check_passport

import (
	"testing"
)

func TestDB_IsValid(t *testing.T) {
	type args struct {
		series string
		number string
	}
	tests := []struct {
		name    string
		args    args
		wantOk  bool
		wantErr bool
	}{
		{
			name: "неверный регион 0000 000000",
			args: args{
				series: "0000",
				number: "000000",
			},
			wantOk:  false,
			wantErr: false,
		},
		{
			name: "неверный номер 0100 000000",
			args: args{
				series: "0100",
				number: "000000",
			},
			wantOk:  false,
			wantErr: false,
		},
		{
			name: "в БД нет данного номера 0100 000101",
			args: args{
				series: "0100",
				number: "000101",
			},
			wantOk:  true,
			wantErr: false,
		},
		{
			name: "есть в БД 0502 875433",
			args: args{
				series: "0502",
				number: "875433",
			},
			wantOk:  false,
			wantErr: false,
		},
	}

	db := NewDB(testDst, nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOk, err := db.IsValid(tt.args.series, tt.args.number)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsValid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOk != tt.wantOk {
				t.Errorf("IsValid() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
