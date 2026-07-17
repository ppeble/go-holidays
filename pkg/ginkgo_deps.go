//go:build tools

// Package holidays: this file is a dependency anchor for the Ginkgo/Gomega
// test framework (see epic go-holidays-hk6). It exists purely so
// `go mod tidy` keeps github.com/onsi/ginkgo/v2 and github.com/onsi/gomega
// pinned in go.mod/go.sum before any test file imports them. It is excluded
// from normal builds and tests by the "tools" build tag and never compiled
// as part of the package.
package holidays

import (
	_ "github.com/onsi/ginkgo/v2"
	_ "github.com/onsi/gomega"
)
