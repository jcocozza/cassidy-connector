package api

import "testing"

func Test_validateKeys(t *testing.T) {
	type args struct {
		keys []StreamType
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test1", args{keys: []StreamType{"asdf"}}, true},
		{"test2", args{keys: []StreamType{"distance"}}, false},
		{"test3", args{keys: []StreamType{"distance", "asdf"}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateKeys(tt.args.keys); (err != nil) != tt.wantErr {
				t.Errorf("validateKeys() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
