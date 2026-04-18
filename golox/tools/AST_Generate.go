package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

)

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// goFieldList converts "Expr left, Token operator, Expr right"
// into "left Expr, operator Token, right Expr" (Go param order)
func goFieldList(fieldList string) string {
	fields := strings.Split(fieldList, ", ")
	result := make([]string, len(fields))
	for i, field := range fields {
		parts := strings.SplitN(field, " ", 2)
		typeName, name := parts[0], parts[1]
		result[i] = name + " " + javaTypeToGo(typeName)
	}
	return strings.Join(result, ", ")
}

func javaTypeToGo(t string) string {
	switch t {
	case "Object":
		return "any"
	case "String":
		return "string"
	case "boolean":
		return "bool"
	case "double":
		return "float64"
	case "int":
		return "int"
	default:
		return t // Expr, Token, etc. pass through as-is
	}
}

func defineType(w *bufio.Writer, structName, fieldList string) {
	// Struct definition
	fmt.Fprintf(w, "type %s struct {\n", structName)

	fields := strings.Split(fieldList, ", ")
	for _, field := range fields {
		parts := strings.SplitN(field, " ", 2)
		typeName := parts[0]
		name := parts[1]
		// Capitalize field name for Go export, map Java types to Go types
		fmt.Fprintf(w, "\t%s %s\n", capitalize(name), javaTypeToGo(typeName))
	}

	fmt.Fprintln(w, "}")
	fmt.Fprintln(w)

	// Constructor function
	fmt.Fprintf(w, "func New%s(%s) *%s {\n", structName, goFieldList(fieldList), structName)
	fmt.Fprintf(w, "\treturn &%s{\n", structName)
	for _, field := range fields {
		parts := strings.SplitN(field, " ", 2)
		name := parts[1]
		fmt.Fprintf(w, "\t\t%s: %s,\n", capitalize(name), name)
	}
	fmt.Fprintln(w, "\t}")
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w)
}

func defineAst(outputDir, baseName string, types []string) error {
	path := filepath.Join(outputDir, baseName+".go")
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	fmt.Fprintln(w, "package lox")
	fmt.Fprintln(w)
	fmt.Fprintf(w, "type %s interface{}\n", baseName)
	fmt.Fprintln(w)

	for _, t := range types {
		parts := strings.SplitN(t, ":", 2)
		structName := strings.TrimSpace(parts[0])
		fields := strings.TrimSpace(parts[1])
		defineType(w, structName, fields)
	}

	return w.Flush()
}

func main() {
	if len(os.Args) != 1 {
		fmt.Println("Usage: generate_ast <output directory")
		os.Exit(64)
	}

	outputDir := os.Args[0]

	defineAst(outputDir, "Expr", []string{})
}
