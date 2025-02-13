package output

import (
	"github.com/xiaoyuer11223344/ksubdomain-fix/pkg/gologger"
	"github.com/xiaoyuer11223344/ksubdomain-fix/pkg/runner/result"
	"strings"
)

type ScreenOutputNoWidth struct {
}

func NewScreenOutputNoWidth() (*ScreenOutputNoWidth, error) {
	return &ScreenOutputNoWidth{}, nil
}
func (s *ScreenOutputNoWidth) WriteDomainResult(domain result.Result) error {
	var msg string
	var domains []string = []string{domain.Subdomain}
	for _, item := range domain.Answers {
		domains = append(domains, item)
	}
	msg = strings.Join(domains, " => ")
	gologger.Infof("%s\n", msg)
	return nil
}
func (s *ScreenOutputNoWidth) Close() {

}
