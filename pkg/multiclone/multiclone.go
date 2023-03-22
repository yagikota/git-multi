package multiclone

import (
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

func gitClone(repoSSHURL string, ch chan<- *Result, clonedReposNum *int32) {
	result := NewResult()
	repoHTTPSURL, err := repoHTTPSURLFromSSHURL(repoSSHURL)
	result.RepoURL = repoHTTPSURL
	if err != nil {
		result.Error = err
		ch <- result
		return
	}
	cmd := exec.Command("git", "clone", repoSSHURL)
	if stderr, err := cmd.CombinedOutput(); err != nil {
		result.Error = errors.New(string(stderr))
	} else {
		atomic.AddInt32(clonedReposNum, 1)
	}
	result.ClonedReposNum = int(*clonedReposNum)
	ch <- result
}

func MultiClone(repoSSHURLs []string, maxGoroutine int) {
	handler := NewMultiCloneHandler(repoSSHURLs, maxGoroutine)
	fmt.Printf("==> Cloning %d repositories:\n", handler.TotalReposNum())
	handler.PrintRepoSSHURLs()

	ch := make(chan *Result, maxGoroutine)
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

	for result := range ch {
		if err := result.Error; err != nil {
			// store in buffer and display at the end
			fmt.Fprintf(handler.Buffer(), result.RepoURL+"\n")
		} else {
			// on progress result
			result.PrintOnProgressResult(handler.TotalReposNum())
		}
	}

	// final result
	fmt.Println(strings.Repeat("=", 100))
	handler.printFinalResult(int(clonedReposNum))
}
