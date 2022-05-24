package conv

import "testing"

func TestHex2Bin(t *testing.T) {
	type args struct {
		hex string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "hex2bin normal", args: args{hex: "2d5555555"}, want: "1011010101010101010101010101010101", wantErr: false},
		{name: "hex2bin error", args: args{hex: "ewor1132"}, want: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Hex2Bin(tt.args.hex)
			if (err != nil) != tt.wantErr {
				t.Errorf("Hex2Bin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Hex2Bin() = %v, want %v", got, tt.want)
			}
		})
	}
}
