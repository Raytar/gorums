package internalgorums

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/build"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"go/types"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/tools/go/loader"
)

const (
	devPkgPath = "./cmd/protoc-gen-gorums/dev"
)

// GenerateBundleFile generates a file with static definitions for Gorums.
func GenerateBundleFile(dst string) {
	pkgIdent, code := bundle(devPkgPath)
	src := "// Code generated by protoc-gen-gorums. DO NOT EDIT.\n" +
		"// Source files can be found in: " + devPkgPath + "\n\n" +
		"package internalgorums\n\n" +
		generatePkgMap(pkgIdent) +
		"var staticCode = `" + string(code) + "`\n"

	staticContent, err := format.Source([]byte(src))
	if err != nil {
		log.Fatalf("formatting failed: %v", err)
	}
	currentContent, err := ioutil.ReadFile(dst)
	if err != nil {
		log.Fatal(err)
	}
	if diff := cmp.Diff(currentContent, staticContent); diff != "" {
		fmt.Fprintf(os.Stderr, "change detected (-current +new):\n%s", diff)
		fmt.Fprintf(os.Stderr, "\nReview changes above; to revert use:\n")
		fmt.Fprintf(os.Stderr, "mv %s.bak %s\n", dst, dst)
	}
	err = ioutil.WriteFile(dst, []byte(staticContent), 0666)
	if err != nil {
		log.Fatal(err)
	}
}

func generatePkgMap(pkgs map[string]string) string {
	var keys []string
	for k := range pkgs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var buf bytes.Buffer
	buf.WriteString(`
	// pkgIdentMap maps from package name to one of the package's identifiers.
	// These identifiers are used by the Gorums protoc plugin to generate
	// appropriate import statements.
	`)
	buf.WriteString("var pkgIdentMap = map[string]string{\n")
	for _, imp := range keys {
		buf.WriteString("\t" + `"`)
		buf.WriteString(imp)
		buf.WriteString(`": "`)
		buf.WriteString(pkgs[imp])
		buf.WriteString(`",` + "\n")
	}
	buf.WriteString("}\n\n")
	return buf.String()
}

// findIdentifiers examines the given package to find all imported packages,
// and one used identifier in that imported package. These identifiers are
// used by the Gorums protoc plugin to generate appropriate import statements.
func findIdentifiers(fset *token.FileSet, info *loader.PackageInfo) map[string]string {
	packagePath := info.Pkg.Path()
	pkgIdent := make(map[string]string)
	for id, obj := range info.Uses {
		pos := fset.Position(id.Pos())
		if strings.Contains(pos.Filename, "zorums") {
			// ignore identifiers in zorums generated files
			continue
		}
		if pkg := obj.Pkg(); pkg != nil && pkg.Path() != packagePath {
			switch obj := obj.(type) {
			case *types.Const:
				// only need to store one identifier for each imported package
				pkgIdent[pkg.Path()] = obj.Name()

			case *types.TypeName:
				// only need to store one identifier for each imported package
				pkgIdent[pkg.Path()] = obj.Name()

			case *types.Var:
				// only need to store one identifier for each imported package
				pkgIdent[pkg.Path()] = obj.Name()

			case *types.Func:
				if typ := obj.Type(); typ != nil {
					if recv := typ.(*types.Signature).Recv(); recv != nil {
						// ignore functions on non-package types
						continue
					}
				}
				// only need to store one identifier for each imported package
				pkgIdent[pkg.Path()] = obj.Name()
			}
		}
	}
	return pkgIdent
}

var ctxt = &build.Default

// bundle returns a slice with the code for the given src package without imports.
// The returned map contains packages to be imported along with one identifier
// using the relevant import path.
// Loosly based on x/tools/cmd/bundle
func bundle(src string) (map[string]string, []byte) {
	conf := loader.Config{ParserMode: parser.ParseComments, Build: ctxt}
	conf.Import(src)
	lprog, err := conf.Load()
	if err != nil {
		log.Fatalf("failed to load Go package: %v", err)
	}
	// Because there was a single Import call and Load succeeded,
	// InitialPackages is guaranteed to hold the sole requested package.
	info := lprog.InitialPackages()[0]

	var out bytes.Buffer
	printFiles(&out, lprog.Fset, info)
	return findIdentifiers(lprog.Fset, info), out.Bytes()
}

func printFiles(out *bytes.Buffer, fset *token.FileSet, info *loader.PackageInfo) {
	for _, f := range info.Files {
		// filter files in dev package that shouldn't be bundled in template_static.go
		fileName := fset.File(f.Pos()).Name()
		if ignore(fileName) {
			continue
		}

		last := f.Package
		if len(f.Imports) > 0 {
			imp := f.Imports[len(f.Imports)-1]
			last = imp.End()
			if imp.Comment != nil {
				if e := imp.Comment.End(); e > last {
					last = e
				}
			}
		}

		// Pretty-print package-level declarations.
		// but no package or import declarations.
		var buf bytes.Buffer
		for _, decl := range f.Decls {
			if decl, ok := decl.(*ast.GenDecl); ok && decl.Tok == token.IMPORT {
				continue
			}
			beg, end := sourceRange(decl)
			printComments(out, f.Comments, last, beg)

			buf.Reset()
			format.Node(&buf, fset, &printer.CommentedNode{Node: decl, Comments: f.Comments})
			out.Write(buf.Bytes())
			last = printSameLineComment(out, f.Comments, fset, end)
			out.WriteString("\n\n")
		}
		printLastComments(out, f.Comments, last)
	}
}

// ignore files in dev folder with suffixes that shouldn't be bundled.
func ignore(file string) bool {
	for _, suffix := range []string{".proto", ".pb.go", "_test.go"} {
		if strings.HasSuffix(file, suffix) {
			return true
		}
	}
	return false
}

// sourceRange returns the [beg, end) interval of source code
// belonging to decl (incl. associated comments).
func sourceRange(decl ast.Decl) (beg, end token.Pos) {
	beg = decl.Pos()
	end = decl.End()

	var doc, com *ast.CommentGroup

	switch d := decl.(type) {
	case *ast.GenDecl:
		doc = d.Doc
		if len(d.Specs) > 0 {
			switch spec := d.Specs[len(d.Specs)-1].(type) {
			case *ast.ValueSpec:
				com = spec.Comment
			case *ast.TypeSpec:
				com = spec.Comment
			}
		}
	case *ast.FuncDecl:
		doc = d.Doc
	}

	if doc != nil {
		beg = doc.Pos()
	}
	if com != nil && com.End() > end {
		end = com.End()
	}

	return beg, end
}

func printPackageComments(out *bytes.Buffer, files []*ast.File) {
	// Concatenate package comments from all files...
	for _, f := range files {
		if doc := f.Doc.Text(); strings.TrimSpace(doc) != "" {
			for _, line := range strings.Split(doc, "\n") {
				fmt.Fprintf(out, "// %s\n", line)
			}
		}
	}
	// ...but don't let them become the actual package comment.
	fmt.Fprintln(out)
}

func printComments(out *bytes.Buffer, comments []*ast.CommentGroup, pos, end token.Pos) {
	for _, cg := range comments {
		if pos <= cg.Pos() && cg.Pos() < end {
			for _, c := range cg.List {
				fmt.Fprintln(out, c.Text)
			}
			fmt.Fprintln(out)
		}
	}
}

const infinity = 1 << 30

func printLastComments(out *bytes.Buffer, comments []*ast.CommentGroup, pos token.Pos) {
	printComments(out, comments, pos, infinity)
}

func printSameLineComment(out *bytes.Buffer, comments []*ast.CommentGroup, fset *token.FileSet, pos token.Pos) token.Pos {
	tf := fset.File(pos)
	for _, cg := range comments {
		if pos <= cg.Pos() && tf.Line(cg.Pos()) == tf.Line(pos) {
			for _, c := range cg.List {
				fmt.Fprintln(out, c.Text)
			}
			return cg.End()
		}
	}
	return pos
}
