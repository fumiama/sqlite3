// Copyright 2017 The Sqlite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"go/scanner"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"log"

	"github.com/cznic/cc"
	"github.com/cznic/ccgo"
	"github.com/cznic/ccir"
	"github.com/cznic/internal/buffer"
	"github.com/cznic/strutil"
	"github.com/cznic/xc"
)

var (
	cpp      = flag.Bool("cpp", false, "")
	dict     = xc.Dict
	errLimit = flag.Int("errlimit", 10, "")
	filter   = flag.String("re", "", "")
	ndebug   = flag.Bool("ndebug", false, "")
	noexec   = flag.Bool("noexec", false, "")
	oLog     = flag.Bool("log", false, "")
	trace    = flag.Bool("trc", false, "")
	yydebug  = flag.Int("yydebug", 0, "")
)

const (
	prologue = `/*

%s
*/

// Code generated by ccgo DO NOT EDIT.

package bin

import (
	"fmt"
	"math"
	"os"
	"path"
	"runtime"
	"unsafe"

	"github.com/cznic/crt"
	"github.com/edsrzf/mmap-go"

)

const minAlloc = 2<<5

var (
	inf = math.Inf(1)
)

func ftrace(s string, args ...interface{}) {
	_, fn, fl, _ := runtime.Caller(1)
	fmt.Fprintf(os.Stderr, "# %%s:%%d: %%v\n", path.Base(fn), fl, fmt.Sprintf(s, args...))
	os.Stderr.Sync()
}

func Init(heapSize, heapReserve int) int {
	heap, err := mmap.MapRegion(nil, heapSize+heapReserve, mmap.RDWR, mmap.ANON, 0)
	if err != nil {
		panic(err)
	}

	tls := crt.NewTLS()
	crt.RegisterHeap(unsafe.Pointer(&heap[0]), int64(heapSize+heapReserve))
	crt.X__register_stdfiles(tls, Xstdin, Xstdout, Xstderr)
	return int(Xinit(tls, int32(heapSize)))
}
`
)

func findRepo(s string) string {
	s = filepath.FromSlash(s)
	for _, v := range strings.Split(strutil.Gopath(), string(os.PathListSeparator)) {
		p := filepath.Join(v, "src", s)
		fi, err := os.Lstat(p)
		if err != nil {
			continue
		}

		if fi.IsDir() {
			wd, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}

			if p, err = filepath.Rel(wd, p); err != nil {
				log.Fatal(err)
			}

			return p
		}
	}
	return ""
}

func errStr(err error) string {
	switch x := err.(type) {
	case scanner.ErrorList:
		if len(x) != 1 {
			x.RemoveMultiples()
		}
		var b bytes.Buffer
		for i, v := range x {
			if i != 0 {
				b.WriteByte('\n')
			}
			b.WriteString(v.Error())
			if i == 9 {
				fmt.Fprintf(&b, "\n\t... and %v more errors", len(x)-10)
				break
			}
		}
		return b.String()
	default:
		return err.Error()
	}
}

func build(predef string, tus [][]string, opts ...cc.Opt) ([]*cc.TranslationUnit, []byte) {
	ndbg := ""
	if *ndebug {
		ndbg = "#define NDEBUG 1"
	}
	var build []*cc.TranslationUnit
	tus = append(tus, []string{ccir.CRT0Path})
	for _, src := range tus {
		model, err := ccir.NewModel()
		if err != nil {
			log.Fatal(err)
		}

		ast, err := cc.Parse(
			fmt.Sprintf(`
%s
#define _CCGO 1
#define __arch__ %s
#define __os__ %s
#include <builtin.h>
%s
`, ndbg, runtime.GOARCH, runtime.GOOS, predef),
			src,
			model,
			append([]cc.Opt{
				cc.AllowCompatibleTypedefRedefinitions(),
				cc.EnableEmptyStructs(),
				cc.EnableImplicitFuncDef(),
				cc.EnableNonConstStaticInitExpressions(),
				cc.ErrLimit(*errLimit),
				cc.SysIncludePaths([]string{ccir.LibcIncludePath}),
			}, opts...)...,
		)
		if err != nil {
			log.Fatal(errStr(err))
		}

		build = append(build, ast)
	}

	var out buffer.Bytes
	if err := ccgo.New(build, &out); err != nil {
		log.Fatal(err)
	}

	return build, out.Bytes()
}

func macros(buf io.Writer, ast *cc.TranslationUnit) {
	fmt.Fprintf(buf, `const (
`)
	var a []string
	for k, v := range ast.Macros {
		if v.Value != nil && v.Type.Kind() != cc.Bool {
			switch fn := v.DefTok.Position().Filename; {
			case
				fn == "builtin.h",
				fn == "<predefine>",
				strings.HasPrefix(fn, "predefined_"):
				// ignore
			default:
				a = append(a, string(dict.S(k)))
			}
		}
	}
	sort.Strings(a)
	for _, v := range a {
		m := ast.Macros[dict.SID(v)]
		if m.Value == nil {
			log.Fatal("TODO")
		}

		switch t := m.Type; t.Kind() {
		case
			cc.Int, cc.UInt, cc.Long, cc.ULong, cc.LongLong, cc.ULongLong,
			cc.Float, cc.LongDouble, cc.Bool:
			fmt.Fprintf(buf, "X%s = %v\n", v, m.Value)
		case cc.Ptr:
			switch t := t.Element(); t.Kind() {
			case cc.Char:
				fmt.Fprintf(buf, "X%s = %q\n", v, dict.S(int(m.Value.(cc.StringLitID))))
			default:
				log.Fatalf("%v", t.Kind())
			}
		default:
			log.Fatalf("%v", t.Kind())
		}
	}

	a = a[:0]
	for _, v := range ast.Declarations.Identifiers {
		switch x := v.Node.(type) {
		case *cc.DirectDeclarator:
			d := x.TopDeclarator()
			id, _ := d.Identifier()
			if x.EnumVal == nil {
				break
			}

			a = append(a, string(dict.S(id)))
		default:
			log.Fatalf("%T", x)
		}
	}
	sort.Strings(a)
	for _, v := range a {
		dd := ast.Declarations.Identifiers[dict.SID(v)].Node.(*cc.DirectDeclarator)
		fmt.Fprintf(buf, "X%s = %v\n", v, dd.EnumVal)
	}
	fmt.Fprintf(buf, ")\n")
}

func main() {
	const repo = "sqlite.org/sqlite-amalgamation-3180000/"

	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
	flag.Parse()
	pth := findRepo(repo)
	if pth == "" {
		log.Fatalf("repository not found: %v", repo)
		return
	}

	asta, src := build(
		`
		#define HAVE_USLEEP 1
		#define SQLITE_DEBUG 1
		#define SQLITE_ENABLE_API_ARMOR 1
		#define SQLITE_ENABLE_MEMSYS5 1
		#define SQLITE_USE_URI 1
		`,
		[][]string{
			{"main.c"},
			{filepath.Join(pth, "sqlite3.c")},
		},
		cc.EnableAnonymousStructFields(),
		cc.IncludePaths([]string{pth}),
	)

	var b bytes.Buffer
	lic, err := ioutil.ReadFile("SQLITE-LICENSE")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(&b, prologue, lic)
	macros(&b, asta[0])
	b.Write(src)
	b2, err := format.Source(b.Bytes())
	if err != nil {
		b2 = b.Bytes()
	}
	if err := os.MkdirAll("internal/bin", 0775); err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(fmt.Sprintf("internal/bin/bin_%s_%s.go", runtime.GOOS, runtime.GOARCH), b2, 0664); err != nil {
		log.Fatal(err)
	}
}
