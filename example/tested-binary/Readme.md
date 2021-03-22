# Installation

```go get - u github.com/AlmostGreatBand/KPI2-1/build/cmd/bood```

this command installs **bood** to your GOPATH, so you should run it either with an absolute path or add **bood** as ENV variable

## Examples

changes only in test-build files
```
zhenyka@MacBook-Pro-Dream tested-binary % ../../build/out/bin/bood
INFO 2021/03/22 23:55:23 Ninja build file is generated at out/build.ninja
INFO 2021/03/22 23:55:23 Starting the build now
[1/1] Test tested_binary as Go binary
```

wrong test(or tests)
```
INFO 2021/03/22 23:53:56 Ninja build file is generated at out/build.ninja
INFO 2021/03/22 23:53:56 Starting the build now
[2/2] Test tested_binary as Go binary
FAILED: out/test/tested_binary 
cd . && go test -v ./... > out/test/tested_binary
ninja: build stopped: subcommand failed.
INFO 2021/03/22 23:53:57 Error invoking ninja build. See logs above.

```

no changes
```
bood
INFO 2021/03/22 23:51:24 Ninja build file is generated at out/build.ninja
INFO 2021/03/22 23:51:24 Starting the build now
ninja: no work to do.
```