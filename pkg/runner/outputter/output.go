package outputter

import (
	"github.com/xiaoyuer11223344/ksubdomain-fix/pkg/runner/result"
)

type Output interface {
	WriteDomainResult(domain result.Result) error
	Close()
}
