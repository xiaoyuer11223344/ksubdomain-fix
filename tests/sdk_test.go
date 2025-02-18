package tests

import (
	"context"
	"github.com/xiaoyuer11223344/ksubdomain-fix/pkg/gologger"
	"github.com/xiaoyuer11223344/ksubdomain-fix/pkg/options"
	"github.com/xiaoyuer11223344/ksubdomain-fix/pkg/runner"
	"github.com/xiaoyuer11223344/ksubdomain-fix/pkg/runner/outputter"
	"github.com/xiaoyuer11223344/ksubdomain-fix/pkg/runner/outputter/output"
	"github.com/xiaoyuer11223344/ksubdomain-fix/pkg/runner/processbar"
	"testing"
)

func TestSDK(t *testing.T) {
	process := processbar.ScreenProcess{}
	screenPrinter, _ := output.NewScreenOutput(false)

	domains := []string{"www.hacking8.com", "x.hacking8.com"}
	domainChanel := make(chan string)
	go func() {
		for _, d := range domains {
			domainChanel <- d
		}
		close(domainChanel)
	}()
	opt := &options.Options{
		Rate:        options.Band2Rate("1m"),
		Domain:      domainChanel,
		DomainTotal: 2,
		Resolvers:   options.GetResolvers(""),
		Silent:      false,
		TimeOut:     10,
		Retry:       3,
		Method:      runner.VerifyType,
		DnsType:     "a",
		Writer: []outputter.Output{
			screenPrinter,
		},
		ProcessBar: &process,
		EtherInfo:  options.GetDeviceConfig(),
	}
	opt.Check()
	r, err := runner.New(opt)
	if err != nil {
		gologger.Fatalf(err.Error())
	}
	ctx := context.Background()
	r.RunEnumeration(ctx)
	r.Close()
}
