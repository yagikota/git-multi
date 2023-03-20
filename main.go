package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"sync"
	"sync/atomic"
)

var GITHUB_DOMAIN = "github.com"

func repoHTTPSURLFromSSHURL(repoURL string) (string, error) {
	repoURLParts := strings.Split(repoURL, ":")
	repoURLParts = strings.Split(repoURLParts[1], ".")
	url, err := url.JoinPath(GITHUB_DOMAIN, repoURLParts[0])
	if err != nil {
		return "", err
	}
	return url, nil
}

type Response struct {
	RepoURL        string
	ClonedReposNum int
	Error          error
}

func gitClone(repoSSHURL string, ch chan<- Response, clonedReposNum *int32) {
	// fmt.Println("git clone")
	var res Response
	repoHTTPSURL, err := repoHTTPSURLFromSSHURL(repoSSHURL)
	res.RepoURL = repoHTTPSURL
	if err != nil {
		res.Error = err
		ch <- res
		return
	}
	cmd := exec.Command("git", "clone", repoSSHURL)
	cmd.Stdout = os.Stdout
	if stderr, err := cmd.CombinedOutput(); err != nil {
		res.Error = errors.New(string(stderr))
	} else {
		atomic.AddInt32(clonedReposNum, 1)
	}
	res.ClonedReposNum = int(*clonedReposNum)
	ch <- res
	fmt.Println(repoHTTPSURL)
}

func main() {
	repoSSHURLs := []string{
		"git@github.com:yagikota/tapple_clone.git",
		"git@github.com:yagikota/wildcard-domain-scanning.git",
		"git@github.com:gin-gonic/gin.git",
		"git@github.com:labstack/echo.git",
		"git@github.com:golang/go.git",
		// "git@github.com:tensorflow/tensorflow.git",

	}
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
	// https://ludwig125.hatenablog.com/entry/2019/09/28/043127
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
	fmt.Printf("==> (%d/%d) repositories are successfully cloned.\n", clonedReposNum, totalReposNum)
	if buffer.Len() > 0 {
		fmt.Println("following repositories are not cloned:")
		fmt.Println(buffer.String())
	} else {
		fmt.Printf("All %d repositories are successfully cloned", totalReposNum)
	}
}
