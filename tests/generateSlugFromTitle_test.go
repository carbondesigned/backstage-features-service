package tests

import (
	"testing"

	"github.com/carbondesigned/backstage-features-service/utils"
)

func TestGenerateSlugFromTitle(t *testing.T) {
	type args struct {
		title string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test if the function returns the correct slug",
			args: args{
				title: "This is a test",
			},
			want: "this-is-a-test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.GenerateSlugFromTitle(tt.args.title); got != tt.want {
				t.Errorf("GenerateSlugFromTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}
