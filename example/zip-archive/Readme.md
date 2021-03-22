# Installation

```go get - u github.com/AlmostGreatBand/KPI2-1/build/cmd/bood```

this command installs bood to your GOPATH, so you should run it either with an absolute path or add bood as ENV variable

## Examples

create zip-archive
```
bood
INFO 2021/03/23 00:17:56 Ninja build file is generated at out/build.ninja
INFO 2021/03/23 00:17:56 Starting the build now
[1/1] Zip my-archive
updating: file_to_zip.txt (stored 0%)
updating: folder_to_zip/file_to_zip-2.txt (stored 0%)
updating: file_to_zip-3.txt (stored 0%)
```

no changes
```
bood
INFO 2021/03/23 00:30:16 Ninja build file is generated at out/build.ninja
INFO 2021/03/23 00:30:16 Starting the build now
ninja: no work to do.
```