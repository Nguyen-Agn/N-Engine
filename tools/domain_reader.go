package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

func main() {
	// 1. Lấy tham số là các tên type cần tìm
	flag.Parse()
	typesToFind := flag.Args()

	if len(typesToFind) == 0 {
		fmt.Println("Usage: go run tools/domain_reader.go <TypeName1> <TypeName2> ...")
		os.Exit(1)
	}

	// 2. Đường dẫn tới folder domain
	domainPath := "./domain"

	// 3. Parser toàn bộ folder domain
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, domainPath, nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("Error parsing directory: %v\n", err)
		os.Exit(1)
	}

	foundCount := 0
	
	// 4. Duyệt qua các package (thường chỉ có 1 là domain)
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			ast.Inspect(file, func(n ast.Node) bool {
				// Tìm định nghĩa Type (Struct hoặc Interface)
				ts, ok := n.(*ast.TypeSpec)
				if !ok {
					return true
				}

				// Kiểm tra xem name có trong danh sách cần tìm không
				name := ts.Name.Name
				if isMatch(name, typesToFind) {
					fmt.Printf("\n--- [ %s ] ---\n", name)
					printNode(fset, n)
					foundCount++
				}
				return true
			})
		}
	}

	if foundCount == 0 {
		fmt.Printf("No types matching [%s] found in %s\n", strings.Join(typesToFind, ", "), domainPath)
	}
}

func isMatch(name string, targets []string) bool {
	for _, t := range targets {
		if strings.EqualFold(name, t) {
			return true
		}
	}
	return false
}

func printNode(fset *token.FileSet, node ast.Node) {
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, node); err != nil {
		fmt.Printf("Error formatting node: %v\n", err)
		return
	}
	fmt.Println(buf.String())
}
