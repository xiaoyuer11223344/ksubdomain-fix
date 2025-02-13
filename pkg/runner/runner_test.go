package runner

import (
	"context"
	"github.com/xiaoyuer11223344/ksubdomain-fix/pkg"
	"github.com/xiaoyuer11223344/ksubdomain-fix/pkg/gologger"
	"github.com/xiaoyuer11223344/ksubdomain-fix/pkg/options"
	"github.com/xiaoyuer11223344/ksubdomain-fix/pkg/runner/outputter"
	"github.com/xiaoyuer11223344/ksubdomain-fix/pkg/runner/outputter/output"
	processbar2 "github.com/xiaoyuer11223344/ksubdomain-fix/pkg/runner/processbar"
	"testing"
)

func TestRunner(t *testing.T) {
	process := processbar2.FakeScreenProcess{}
	screenPrinter, _ := output.NewScreenOutputNoWidth()
	domains := []string{"stu.baidu.com", "haokan.baidu.com"}
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
		TimeOut:     5,
		Retry:       1,
		Method:      VerifyType,
		DnsType:     "a",
		Writer: []outputter.Output{
			screenPrinter,
		},
		ProcessBar: &process,
		EtherInfo:  options.GetDeviceConfig(),
	}
	opt.Check()
	r, err := New(opt)
	if err != nil {
		gologger.Fatalf(err.Error())
	}
	ctx := context.Background()
	r.RunEnumeration(ctx)
	r.Close()
}

func TestRunnerEnum(t *testing.T) {
	process := processbar2.ScreenProcess{}
	screenPrinter, _ := output.NewScreenOutputNoWidth()
	domains := pkg.GetDefaultSubdomainData()
	domainChanel := make(chan string)
	go func() {
		for _, d := range domains {
			domainChanel <- d + ".baidu.com"
		}
		close(domainChanel)
	}()
	opt := &options.Options{
		Rate:        options.Band2Rate("1m"),
		Domain:      domainChanel,
		DomainTotal: len(domains),
		Resolvers:   options.GetResolvers(""),
		Silent:      false,
		TimeOut:     5,
		Retry:       1,
		Method:      VerifyType,
		DnsType:     "a",
		Writer: []outputter.Output{
			screenPrinter,
		},
		ProcessBar: &process,
		EtherInfo:  options.GetDeviceConfig(),
	}
	opt.Check()
	r, err := New(opt)
	if err != nil {
		gologger.Fatalf(err.Error())
	}
	ctx := context.Background()
	r.RunEnumeration(ctx)
	r.Close()
}

func TestManyRunner(t *testing.T) {
	for i := 0; i < 5; i++ {
		TestRunner(t)
	}
}
