package links

import (
	"testing"

	mock_keys "github.com/abdybaevae/url-shortener/mocks/pkg/services/keys"
	"github.com/abdybaevae/url-shortener/pkg/errors"
	gomock "github.com/golang/mock/gomock"
)

const DefaultLink = "https://google.com/path/to/something"
const DefaultKey = "some-key"

func TestShorten(t *testing.T) {
	type deps struct {
		keyService *mock_keys.MockKeyService
	}
	tt := []struct {
		name    string
		want    string
		wantErr error
		args    string
		prepare func(d *deps)
	}{
		{
			name: "should shorten link",
			want: DefaultKey,
			args: DefaultLink,
			prepare: func(d *deps) {
				t.Log("prepare called")
				d.keyService.EXPECT().Get().Return(DefaultKey, nil)
			},
		},
		{
			name:    "should fail on empty link",
			wantErr: errors.InvalidLink,
			args:    "",
		},
		{
			name:    "should fail on invalid link",
			wantErr: errors.InvalidLink,
			args:    "brokenlink",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			d := deps{
				keyService: mock_keys.NewMockKeyService(ctrl),
			}
			if tc.prepare != nil {
				tc.prepare(&d)
			}
			linkService := NewLinkService(d.keyService)
			got, gotErr := linkService.Shorten(tc.args)
			if tc.wantErr != nil {
				if tc.wantErr != gotErr {
					t.Errorf("want %v got %v", tc.wantErr, gotErr)
				}
				return

			}
			if got != tc.want {
				t.Errorf("want %v got %v", tc.want, got)
			}
		})

	}
}
