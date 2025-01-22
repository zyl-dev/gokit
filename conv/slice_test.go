package conv

import (
	"reflect"
	"testing"
)

func Test_RemoveDuplicateElement(t *testing.T) {
	type args struct {
		list []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"name", args{list: []string{"111"}}, []string{"111"}},
		{"name", args{list: []string{"111", "1112"}}, []string{"111", "1112"}},
		{"name", args{list: []string{"111", "111"}}, []string{"111"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveDuplicateElement(tt.args.list); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveDuplicateElement() = %v, want %v", got, tt.want)
			}
		})
	}
}
