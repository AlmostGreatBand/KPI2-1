package gomodule

import (
	"fmt"
	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
	"path"
)

var (
	pctx = blueprint.NewPackageContext("github.com/AlmostGreatBand/KPI2-1/gomodule")

	goBuild = pctx.StaticRule("binaryBuild", blueprint.RuleParams{
		Command:     "cd $workDir && go build -o $outPath $pkg",
		Description: "build go command $pkg",
	}, "workDir", "outPath", "pkg")

	goVendor = pctx.StaticRule("vendor", blueprint.RuleParams{
		Command:     "cd $workDir && go mod vendor",
		Description: "vendor dependencies of $name",
	}, "workDir", "name")

	goTest = pctx.StaticRule("binaryTest", blueprint.RuleParams{
		Command:     "cd $workDir && go test -v ${pkg} > ${outPath}",
		Description: "test go command $pkg",
	}, "workDir", "outPath", "pkg")
)

type testedBinaryModule struct {
	blueprint.SimpleName

	properties struct {
		Pkg string
		TestPkg string
		Srcs []string
		TestSrcs []string
		SrcsExclude []string
		TestSrcsExclude []string
		VendorFirst bool

		Deps []string
	}
}

func (tb *testedBinaryModule) DynamicDependencies(blueprint.DynamicDependerModuleContext) []string {
	return tb.properties.Deps
}

func (tb *testedBinaryModule) GenerateBuildActions(ctx blueprint.ModuleContext) {
	name := ctx.ModuleName()
	config := bood.ExtractConfig(ctx)
	config.Debug.Printf("Adding build and test actions for go binary module '%s'", name)

	buildPath := path.Join(config.BaseOutputDir, "bin", name)
	testPath := path.Join(config.BaseOutputDir, "test", name)

	inputErrors := false

	var buildInputs []string
	buildExclude := append(tb.properties.SrcsExclude, tb.properties.TestSrcs...)
	for _, src := range tb.properties.Srcs {
		if matches, err := ctx.GlobWithDeps(src, buildExclude); err == nil {
			buildInputs = append(buildInputs, matches...)
		} else {
			ctx.PropertyErrorf("srcs", "Cannot resolve build-files that match pattern %s", src)
			inputErrors = true
		}
	}

	var testInputs []string
	for _, src := range tb.properties.TestSrcs {
		if matches, err := ctx.GlobWithDeps(src, tb.properties.TestSrcsExclude); err == nil {
			testInputs = append(testInputs, matches...)
		} else {
			ctx.PropertyErrorf("testSrcs", "Cannot resolve test-files that match pattern %s", src)
			inputErrors = true
		}
	}

	testInputs = append(testInputs, buildInputs...)

	if inputErrors {
		return
	}

	if tb.properties.VendorFirst {
		vendorDirPath := path.Join(ctx.ModuleDir(), "vendor")
		ctx.Build(pctx, blueprint.BuildParams{
			Description: fmt.Sprintf("Vendor dependencies of %s", name),
			Rule:        goVendor,
			Outputs:     []string{vendorDirPath},
			Implicits:   []string{path.Join(ctx.ModuleDir(), "go.mod")},
			Optional:    true,
			Args: map[string]string{
				"workDir": ctx.ModuleDir(),
				"name":    name,
			},
		})
		buildInputs = append(buildInputs, vendorDirPath)
	}

	ctx.Build(pctx, blueprint.BuildParams{
		Description: fmt.Sprintf("Build %s as Go binary", name),
		Rule:        goBuild,
		Outputs:     []string{buildPath},
		Implicits:   buildInputs,
		Args: map[string]string{
			"outPath":	  buildPath,
			"workDir":    ctx.ModuleDir(),
			"pkg":        tb.properties.Pkg,
		},
	})

	ctx.Build(pctx, blueprint.BuildParams{
		Description: fmt.Sprintf("Test %s as Go binary", name),
		Rule:        goTest,
		Outputs:     []string{testPath},
		Implicits:   testInputs,
		Args: map[string]string{
			"outPath": 	  testPath,
			"workDir":    ctx.ModuleDir(),
			"pkg":        tb.properties.TestPkg,
		},
	})
}

func SimpleTestFactory() (blueprint.Module, []interface{}) {
	mType := &testedBinaryModule{}
	return mType, []interface{}{&mType.SimpleName.Properties, &mType.properties}
}
