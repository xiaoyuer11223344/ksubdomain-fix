set CGO_LDFLAGS = "-Wl,-static -L/usr/lib/x86_64-linux-gnu/libpcap.a -lpcap -Wl,-Bdynamic"
set GOOS = "linux"
set GOARCH = "amd64"
go build -o ./ksubdomain ./cmd/ksubdomain/