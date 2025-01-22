package masking

import "testing"

func TestLogMaskingConfig(t *testing.T) {
	type args struct {
		a string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "normal", args: args{a: `"password":"123456"`}, want: `"password":"123456"`},
		{name: "normal masking", args: args{a: `"password":"123456",`}, want: `"password":"","`},
		{name: "normal masking addrs", args: args{a: `"Addrs":["password"],`}, want: `"Addrs":[""],`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LogMaskingConfig(tt.args.a); got != tt.want {
				t.Errorf("LogMaskingConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
