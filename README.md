[![Go Reference](https://pkg.go.dev/badge/github.com/nofeaturesonlybugs/fscopy.svg)](https://pkg.go.dev/github.com/nofeaturesonlybugs/fscopy)
[![Go Report Card](https://goreportcard.com/badge/github.com/nofeaturesonlybugs/fscopy)](https://goreportcard.com/report/github.com/nofeaturesonlybugs/fscopy)
[![Build Status](https://travis-ci.com/nofeaturesonlybugs/fscopy.svg?branch=master)](https://travis-ci.com/nofeaturesonlybugs/fscopy)
[![codecov](https://codecov.io/gh/nofeaturesonlybugs/fscopy/branch/master/graph/badge.svg)](https://codecov.io/gh/nofeaturesonlybugs/fscopy)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Package `fscopy` implements file system copy functions not present in standard library.

```go
fscopy.Dir( "/home/nofeatures/dst-dir", "/home/nofeatures/source-dir" )
fscopy.File( "/home/nofeatures/dst-file", "/home/nofeatures/source-file" )
fscopy.Into( "/home/nofeatures/dst-dir", "/tmp/file1", "/tmp/file2", "/tmp/some-dir" )
```

## Circular References and symlinks

The types of environments where I intend to use this package are somewhat manicured and highly unlikely to experience issues related to symlinks or circular references. As such this package has no logic or consideration for symlinks. However `fscopy.Dir()` does have a low-effort mechanism to track `src-to-dst` copies that have already occurred and will not duplicate its efforts.

## Test Coverage

I don't expect test coverage for this package will ever achieve 100%. The majority of statements that are not covered are `return errors.Go(err)` (essentially `return err`) where `os.Stat()` or `os.File.Close()` or `io.Copy()` have returned a non-nil error. While I could put in the efforts to mock such operations I don't feel the extra complexity would be worth the trade off in code coverage.
