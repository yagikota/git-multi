package multiclone

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsSSHURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{
			name: "valid ssh url",
			url:  "git@github.com:username/repo.git",
			want: true,
		},
		{
			name: "invalid ssh url(prefix miss match)",
			url:  "https://github.com/username/repo.git",
			want: false,
		},
		{
			name: "invalid ssh url(suffix miss match)",
			url:  "git@github.com:username/repo",
			want: false,
		},
		{
			name: "invalid url(prefix and suffix miss match)",
			url:  "invalid-url",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isSSHURL(tt.url)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRepoHTTPSURLFromSSHURL(t *testing.T) {
	tests := []struct {
		name    string
		repoURL string
		want    string
		wantErr error
	}{
		{
			name:    "valid ssh url",
			repoURL: "git@github.com:username/repo.git",
			want:    "https://github.com/username/repo",
			wantErr: nil,
		},
		{
			name:    "invalid ssh url",
			repoURL: "invalid-url",
			want:    "",
			wantErr: errors.New("invalid SSH URL"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repoHTTPSURLFromSSHURL(tt.repoURL)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// func Test_gitClone(t *testing.T) {
// 	type args struct {
// 		repoSSHURL     string
// 		ch             chan<- Response
// 		clonedReposNum *int32
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			gitClone(tt.args.repoSSHURL, tt.args.ch, tt.args.clonedReposNum)
// 		})
// 	}
// }

// func TestMultiClone(t *testing.T) {
// 	type args struct {
// 		repoSSHURLs []string
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			MultiClone(tt.args.repoSSHURLs)
// 		})
// 	}
// }
