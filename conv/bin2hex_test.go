package conv

import "testing"

func TestBin2Hex(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "bin2hex normal", args: args{str: "1011010101010101010101010101010101"}, want: "2d5555555", wantErr: false},
		{name: "bin2hex error", args: args{str: "123245"}, want: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Bin2Hex(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("Bin2Hex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Bin2Hex() = %v, want %v", got, tt.want)
			}
		})
	}
}
