package gomodule

import (
	"bytes"
	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
	"strings"
	"testing"
)

func TestSimpleZipFactory(t *testing.T) {
	ctx := blueprint.NewContext()

	ctx.MockFileSystem(map[string][]byte{
		"Blueprints": []byte(`
			zip_archive {
			  name: "zip-out",
			  srcs: ["zip-src.txt"],
			}
		`),
		"zip-src.txt": nil,
	})

	ctx.RegisterModuleType("zip_archive", SimpleZipFactory)

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
		if !strings.Contains(text, "build out/archives/zip-out.zip: g.gomodule.goZip | zip-src.txt") {
			t.Errorf("Generated ninja file does not have build of the zip module")
		}
	}
}
