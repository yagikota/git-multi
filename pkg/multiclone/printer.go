package multiclone

import (
	"bytes"
	"fmt"
)

// type MultiClonePrinter interface {
// 	PrintOnProgressResult(totalReposNum int)
// }

type Result struct {
	RepoURL        string
	ClonedReposNum int
	Error          error
}

func NewResult() *Result {
	return new(Result)
}

func (r *Result) PrintOnProgressResult(totalReposNum int) {
	fmt.Printf("%s (%d/%d)\n", r.RepoURL, r.ClonedReposNum, totalReposNum)
}

type MultiCloneHandler struct {
	maxGoroutine  int
	repoSSHURLs   []string
	repoHTTPSURLs []string
	buffer        bytes.Buffer
}

func NewMultiCloneHandler(repoSSHURLs []string, maxGoroutine int) *MultiCloneHandler {
	repoHTTPSURLs := make([]string, 0, len(repoSSHURLs))
	for _, repoURL := range repoSSHURLs {
		url, err := repoHTTPSURLFromSSHURL(repoURL)
		if err != nil {
			continue
		}
		repoHTTPSURLs = append(repoHTTPSURLs, url)
	}
	return &MultiCloneHandler{
		maxGoroutine:  maxGoroutine,
		repoSSHURLs:   repoSSHURLs,
		repoHTTPSURLs: repoHTTPSURLs,
	}
}

func (h *MultiCloneHandler) TotalReposNum() int {
	return len(h.repoSSHURLs)
}

func (h *MultiCloneHandler) ClonedReposNum() int {
	return len(h.repoSSHURLs)
}

func (h *MultiCloneHandler) PrintRepoSSHURLs() {
	for _, repoURL := range h.repoSSHURLs {
		fmt.Printf("%s ...\n", repoURL)
	}
}

func (h *MultiCloneHandler) Buffer() *bytes.Buffer {
	return &h.buffer
}

func (h *MultiCloneHandler) printFinalResult(clonedReposNum int) {
	totalReposNum := h.TotalReposNum()
	fmt.Printf("==> (%d/%d) success\n", clonedReposNum, totalReposNum)
	if h.buffer.Len() > 0 {
		fmt.Println("following repositories are not cloned:")
		fmt.Println(h.buffer.String())
	} else {
		fmt.Printf("All %d repositories are successfully cloned", totalReposNum)
	}
}
