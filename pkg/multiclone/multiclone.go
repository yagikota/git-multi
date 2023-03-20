package multiclone

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"os/exec"
	"strings"
	"sync"
	"sync/atomic"
)

const (
	gitHubDomain    = "https://github.com"
	gitHubSSHPrefix = "git@github.com:"
	gitHubSSHSuffix = ".git"
)

type Response struct {
	RepoURL        string
	ClonedReposNum int
	Error          error
}

// https://docs.github.com/en/get-started/getting-started-with-git/about-remote-repositories#about-remote-repositories
func isSSHURL(repoURL string) bool {
	return strings.HasPrefix(repoURL, gitHubSSHPrefix) && strings.HasSuffix(repoURL, gitHubSSHSuffix)
}

func repoHTTPSURLFromSSHURL(repoURL string) (string, error) {
	if !isSSHURL(repoURL) {
		return "", errors.New("invalid SSH URL")
	}
	// git@github.com:username/repo.git -> username/repo
	trimmed := strings.TrimPrefix(repoURL, gitHubSSHPrefix)
	trimmed = strings.TrimSuffix(trimmed, gitHubSSHSuffix)
	url, err := url.JoinPath(gitHubDomain, trimmed)
	if err != nil {
		return "", err
	}
	return url, nil
}

func gitClone(repoSSHURL string, ch chan<- Response, clonedReposNum *int32) {
	var res Response
	repoHTTPSURL, err := repoHTTPSURLFromSSHURL(repoSSHURL)
	res.RepoURL = repoHTTPSURL
	if err != nil {
		res.Error = err
		ch <- res
		return
	}
	cmd := exec.Command("git", "clone", repoSSHURL)
	if stderr, err := cmd.CombinedOutput(); err != nil {
		res.Error = errors.New(string(stderr))
	} else {
		atomic.AddInt32(clonedReposNum, 1)
	}
	res.ClonedReposNum = int(*clonedReposNum)
	ch <- res
}

func MultiClone(repoSSHURLs []string) {
	totalReposNum := int32(len(repoSSHURLs))
	repoHTTPSURLs := make([]string, totalReposNum)
	fmt.Printf("==> Cloning %d repositories:\n", totalReposNum)
	for i, repoURL := range repoSSHURLs {
		url, err := repoHTTPSURLFromSSHURL(repoURL)
		if err != nil {
			continue
		}
		repoHTTPSURLs[i] = url
		fmt.Printf("%s ...\n", url)
	}

	maxGoroutines := 10
	ch := make(chan Response, maxGoroutines)
	wg := new(sync.WaitGroup)
	var clonedReposNum int32
	for _, repoSSHURL := range repoSSHURLs {
		wg.Add(1)
		go func(repoSSHURL string) {
			defer wg.Done()
			gitClone(repoSSHURL, ch, &clonedReposNum)
		}(repoSSHURL)
	}

	// wait until all the above goroutines are finished in another goroutine
	// close the channel when done
	go func() {
		defer close(ch)
		wg.Wait()
	}()

	var buffer bytes.Buffer
	for res := range ch {
		if err := res.Error; err != nil {
			// store in buffer and display at the end
			fmt.Fprintf(&buffer, res.RepoURL+"\n")
		} else {
			fmt.Printf("%s (%d/%d)\n", res.RepoURL, res.ClonedReposNum, totalReposNum)
		}
	}

	fmt.Println(strings.Repeat("=", 100))
	if clonedReposNum == totalReposNum {
		fmt.Printf("==> All repositories have successfully cloned (%d/%d)", clonedReposNum, totalReposNum)
		return
	}
	if clonedReposNum > 1 {
		fmt.Printf("==> some repositories have successfully cloned (%d/%d)\n", clonedReposNum, totalReposNum)
		return
	}
	if clonedReposNum == 1 {
		fmt.Printf("==> one repository has successfully cloned (%d/%d)\n", clonedReposNum, totalReposNum)
		return
	}
	fmt.Printf("==> (%d/%d) repositories are successfully cloned.\n", clonedReposNum, totalReposNum)
	if buffer.Len() > 0 {
		fmt.Println("following repositories are not cloned:")
		fmt.Println(buffer.String())
	} else {
		fmt.Printf("All %d repositories are successfully cloned", totalReposNum)
	}
}
