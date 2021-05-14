package gomodule

import (
	"bytes"
	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
	"strings"
	"testing"
)

func TestSimpleTestFactory(t *testing.T) {
	ctx := blueprint.NewContext()

	ctx.MockFileSystem(map[string][]byte{
		"Blueprints": []byte(`
			tested_binary {
				name: "test-out",
				srcs: ["test-src.go"],
				testSrcs: ["test-src_test.go"],
				pkg: ".",
				vendorFirst: false
			}
		`),
		"test-src.go":      nil,
		"test-src_test.go": nil,
	})

	ctx.RegisterModuleType("tested_binary", SimpleTestFactory)

	cfg := bood.NewConfig()

	_, errs := ctx.ParseBlueprintsFiles(".", cfg)
	if len(errs) != 0 {
		t.Fatalf("Syntax errors in the test blueprint file: %s", errs)
	}

	_, errs = ctx.PrepareBuildActions(cfg)
	if len(errs) != 0 {
		t.Errorf("Unexpected errors while preparing build actions: %s", errs)
	}
	buffer := new(bytes.Buffer)
	if err := ctx.WriteBuildFile(buffer); err != nil {
		t.Errorf("Error writing ninja file: %s", err)
	} else {
		text := buffer.String()
		t.Logf("Gennerated ninja build file:\n%s", text)
		if !strings.Contains(text, "out/bin/test-out: g.gomodule.binaryBuild | test-src.go") {
			t.Errorf("Generated ninja file does not have build of the test module")
		}
		if !strings.Contains(text, "out/test/test-out: g.gomodule.binaryTest | test-src_test.go test-src.go") {
			t.Errorf("Generated ninja file does not have build of the test module")
		}
		if strings.Contains(text, "build vendor: g.gomodule.vendor | go.mod") {
			t.Errorf("Generated ninja file does not have vendor build rule")
		}
	}
}

func TestSimpleTestOptionalFactory(t *testing.T) {
	ctx := blueprint.NewContext()

	ctx.MockFileSystem(map[string][]byte{
		"Blueprints": []byte(`
			tested_binary {
				name: "test-out",
				srcs: [],
				optional: true,
				testSrcs: ["test-src_test.go"],
				pkg: ".",
				vendorFirst: false
			}
		`),
		"test-src.go":      nil,
		"test-src_test.go": nil,
	})

	ctx.RegisterModuleType("tested_binary", SimpleTestFactory)

	cfg := bood.NewConfig()

	_, errs := ctx.ParseBlueprintsFiles(".", cfg)
	if len(errs) != 0 {
		t.Fatalf("Syntax errors in the test blueprint file: %s", errs)
	}

	_, errs = ctx.PrepareBuildActions(cfg)
	if len(errs) != 0 {
		t.Errorf("Unexpected errors while preparing build actions: %s", errs)
	}
	buffer := new(bytes.Buffer)
	if err := ctx.WriteBuildFile(buffer); err != nil {
		t.Errorf("Error writing ninja file: %s", err)
	} else {
		text := buffer.String()
		t.Logf("Gennerated ninja build file:\n%s", text)
		if strings.Contains(text, "out/bin/test-out: g.gomodule.binaryBuild | test-src.go") {
			t.Errorf("Generated ninja file does not have build of the test module")
		}
		if !strings.Contains(text, "out/test/test-out: g.gomodule.binaryTest | test-src_test.go") {
			t.Errorf("Generated ninja file does not have build of the test module")
		}
		if strings.Contains(text, "build vendor: g.gomodule.vendor | go.mod") {
			t.Errorf("Generated ninja file does not have vendor build rule")
		}
		if strings.Contains(text, "default out/test/test-out") {
			t.Errorf("Generated ninja file has default target")
		}
	}
}
