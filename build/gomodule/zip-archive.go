package gomodule

import (
	"fmt"
	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
	"path"
	"strings"
)

var (
	goZip = pctx.StaticRule("goZip", blueprint.RuleParams{
		Command:     "cd $workDir && zip $outPath $inputs",
		Description: "zip files $inputs",
	}, "workDir", "inputs", "outPath")
)

type zipModule struct {
	blueprint.SimpleName

	properties struct {
		Srcs []string
		Deps []string
	}
}

func (z *zipModule) DynamicDependencies(blueprint.DynamicDependerModuleContext) []string {
	return z.properties.Deps
}

func (z *zipModule) GenerateBuildActions(ctx blueprint.ModuleContext) {
	name := ctx.ModuleName()
	config := bood.ExtractConfig(ctx)
	config.Debug.Printf("Adding zip actions for '%s'", name)

	zipPath := path.Join(config.BaseOutputDir, "archives", name)

	inputErrors := false

	var zipInputs []string
	for _, src := range z.properties.Srcs {
		if matches, err := ctx.GlobWithDeps(src, nil); err == nil {
			zipInputs = append(zipInputs, matches...)
		} else {
			ctx.PropertyErrorf("zip", "Cannot resolve files that match pattern %s", src)
			inputErrors = true
		}
	}

	if inputErrors {
		return
	}

	ctx.Build(pctx, blueprint.BuildParams{
		Description: fmt.Sprintf("Zip %s", name),
		Rule:        goZip,
		Outputs:     []string{zipPath},
		Implicits:   zipInputs,
		Args: map[string]string{
			"outPath": zipPath,
			"workDir": ctx.ModuleDir(),
			"inputs":  strings.Join(zipInputs, " "),
		},
	})
}

func SimpleZipFactory() (blueprint.Module, []interface{}) {
	mType := &zipModule{}
	return mType, []interface{}{&mType.SimpleName.Properties, &mType.properties}
}
