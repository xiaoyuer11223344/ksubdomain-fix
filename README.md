
# 特点

1. 无状态子域名爆破，基于pcap高速发包，通过gopacket实现
2. 支持重试机制
3. 泛解析
4. 基于域名规律自动生成FUZZ爆破点
5. 提供SDK接口，支持knative调用

# 用法

golang + libpcap库

```bash
   enum, e    枚举域名
   verify, v  验证模式
   test       测试本地网卡的最大发送速度
   help, h    Shows a list of commands or help for one command
```

# 模式

## 验证模式

提供完整的域名列表，ksubdomain负责快速获取结果

```bash
./ksubdomain verify -h

NAME:
   ksubdomain verify - 验证模式

USAGE:
   ksubdomain verify [command options] [arguments...]

OPTIONS:
   --filename value, -f value   验证域名文件路径
   --band value, -b value       宽带的下行速度，可以5M,5K,5G (default: "2m")
   --resolvers value, -r value  dns服务器文件路径，一行一个dns地址
   --output value, -o value     输出文件名
   --silent                     使用后屏幕将仅输出域名 (default: false)
   --retry value                重试次数,当为-1时将一直重试 (default: 3)
   --timeout value              超时时间 (default: 6)
   --stdin                      接受stdin输入 (default: false)
   --only-domain, --od          只打印域名，不显示ip (default: false)
   --not-print, --np            不打印域名结果 (default: false)
   --dns-type value             dns类型 1为a记录 2为ns记录 5为cname记录 16为txt (default: 1)
   --help, -h                   show help (default: false)
```

```
从文件读取 
./ksubdomain v -f dict.txt

从stdin读取
echo "www.hacking8.com"|./ksubdomain v --stdin

读取ns记录
echo "hacking8.com" | ./ksubdomain v --stdin --dns-type 2
```

## 枚举模式
只提供一级域名，指定域名字典或使用ksubdomain内置字典，枚举所有二级域名

```bash
./ksubdomain enum -h

NAME:
   ksubdomain enum - 枚举域名

USAGE:
   ksubdomain enum [command options] [arguments...]

OPTIONS:
   --band value, -b value          宽带的下行速度，可以5M,5K,5G (default: "2m")
   --resolvers value, -r value     dns服务器文件路径，一行一个dns地址
   --output value, -o value        输出文件名
   --silent                        使用后屏幕将仅输出域名 (default: false)
   --retry value                   重试次数,当为-1时将一直重试 (default: 3)
   --timeout value                 超时时间 (default: 6)
   --stdin                         接受stdin输入 (default: false)
   --only-domain, --od             只打印域名，不显示ip (default: false)
   --not-print, --np               不打印域名结果 (default: false)
   --dns-type value                dns类型 1为a记录 2为ns记录 5为cname记录 16为txt (default: 1)
   --domain value, -d value        爆破的域名
   --domainList value, --dl value  从文件中指定域名
   --filename value, -f value      字典路径
   --skip-wild                     跳过泛解析域名 (default: false)
   --level value, -l value         枚举几级域名，默认为2，二级域名 (default: 2)
   --level-dict value, --ld value  枚举多级域名的字典文件，当level大于2时候使用，不填则会默认
   --help, -h                      show help (default: false)
```

```
./ksubdomain e -d baidu.com

从stdin获取
echo "baidu.com"|./ksubdomain e --stdin
```
# SDK调用


```go
package main

import (
	"context"
	"github.com/xiaoyuer11223344/ksubdomain-fix/core/gologger"
	"github.com/xiaoyuer11223344/ksubdomain-fix/core/options"
	"github.com/xiaoyuer11223344/ksubdomain-fix/runner"
	"github.com/xiaoyuer11223344/ksubdomain-fix/runner/outputter"
	"github.com/xiaoyuer11223344/ksubdomain-fix/runner/outputter/output"
	"github.com/xiaoyuer11223344/ksubdomain-fix/runner/processbar"
	"strings"
)

func main() {
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
```

options的参数结构

```go
type Options struct {
	Rate        int64              // 每秒发包速率
	Domain      io.Reader          // 域名输入
	DomainTotal int                // 扫描域名总数
	Resolvers   []string           // dns resolvers
	Silent      bool               // 安静模式
	TimeOut     int                // 超时时间 单位(秒)
	Retry       int                // 最大重试次数
	Method      string             // verify模式 enum模式 test模式
	DnsType     string             // dns类型 a ns aaaa
	Writer      []outputter.Output // 输出结构
	ProcessBar  processbar.ProcessBar
	EtherInfo   *device.EtherTable // 网卡信息
}
```

1. ksubdomain底层接口只是一个dns验证器，如果要通过一级域名枚举，需要把全部的域名都放入`Domain`字段中，可以看enum参数是怎么写的 `cmd/ksubdomain/enum.go`
2. Write参数是一个outputter.Output接口，用途是如何处理DNS返回的接口，ksubdomain已经内置了三种接口在 `runner/outputter/output`中，主要作用是把数据存入内存、数据写入文件、数据打印到屏幕，可以自己实现这个接口，实现自定义的操作。
3. ProcessBar参数是一个processbar.ProcessBar接口，主要用途是将程序内`成功个数`、`发送个数`、`队列数`、`接收数`、`失败数`、`耗时`传递给用户，实现这个参数可以时时获取这些。
4. EtherInfo是*device.EtherTable类型，用来获取网卡的信息，一般用函数`options.GetDeviceConfig()`即可自动获取网卡配置。



## 特性和Tips

- 无状态爆破，有失败重发机制，速度极快
- 中文帮助，-h会看到中文帮助
- 两种模式，枚举模式和验证模式，枚举模式内置10w字典
- 将网络参数简化为了-b参数，输入你的网络下载速度如-b 5m，将会自动限制网卡发包速度。
- 可以使用./ksubdomain test来测试本地最大发包数
- 获取网卡改为了全自动并可以根据配置文件读取。
- 会有一个时时的进度条，依次显示成功/发送/队列/接收/失败/耗时 信息。
- 不同规模的数据，调整 --retry --timeout参数即可获得最优效果
- 当--retry为-1，将会一直重试直到所有成功。
- 支持爆破ns记录

## 与massdns、dnsx对比

使用100w字典，在4H5M的网络环境下测试

|          | ksubdomain                                                   | massdns                                                      | dnsx                                                         |
| -------- | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| 支持系统 | Windows/Linux/Darwin                                         | Windows/Linux/Darwin                                         | Windows/Linux/Darwin                                         |
| 功能 | 支持验证和枚举 | 只能验证 | 只能验证 |
| 发包方式 | pcap网卡发包                                                 | epoll,pcap,socket                                            | socket                                                       |
| 命令行 | time ./ksubdomain v -b 5m -f d2.txt -o ksubdomain.txt -r dns.txt --retry 3 --np | time ./massdns -r dns.txt -t AAAA -w massdns.txt d2.txt --root -o L | time ./dnsx -a -o dnsx.txt -r dns.txt -l d2.txt -retry 3 -t 5000 |
| 备注   | 加了--np 防止打印过多                                        |                                                              |                                                              |
| 结果   | 耗时:1m28.273s<br />成功个数:1397                            | 耗时:3m29.337s<br />成功个数:1396                            | 耗时:5m26.780s <br />成功个数:1396                           |

ksubdomain只需要1分半，速度远远比massdns、dnsx快~

## 参考

- 原ksubdomain https://github.com/knownsec/ksubdomain
- 从 Masscan, Zmap 源码分析到开发实践 <https://paper.seebug.org/1052/>
- ksubdomain 无状态域名爆破工具介绍 <https://paper.seebug.org/1325/>
- [ksubdomain与massdns的对比](https://mp.weixin.qq.com/s?__biz=MzU2NzcwNTY3Mg==&mid=2247484471&idx=1&sn=322d5db2d11363cd2392d7bd29c679f1&chksm=fc986d10cbefe406f4bda22f62a16f08c71f31c241024fc82ecbb8e41c9c7188cfbd71276b81&token=76024279&lang=zh_CN#rd) 
