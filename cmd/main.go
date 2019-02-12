package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pindamonhangaba/tabua/generate"
	"github.com/pindamonhangaba/tabua/reverse"
)

var (
	dbStringFlag = flag.String("db", "user=postgres password=postgres dbname=postgres sslmode=disable", "database connection string")
	pathFlag     = flag.String("p", "./", "path to save models")
	packageFlag  = flag.String("pkg", "generated/models", "package path")
	filterFlag   = flag.String("f", "", "filter tables to reverse")
	schemaFlag   = flag.String("sch", "", "database schema, default is 'public'")
)

func main() {
	flag.Parse()

	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	filter := reverse.Filter{Schema: "public"}
	if len(*filterFlag) > 0 {
		filter.Tables = strings.Split(*filterFlag, ",")
	}
	if len(*schemaFlag) > 0 {
		filter.Schema = *schemaFlag
	}

	db, err := sqlx.Connect("postgres", *dbStringFlag)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	r, _ := reverse.New(db, reverse.SQLFromPsql)
	tables, err := r.Run(filter)
	if err != nil {
		panic(err)
	}
	gen := generate.Generator{PackagePath: *packageFlag}

	buf := &bytes.Buffer{}
	for _, t := range tables {
		f, pkg := gen.Run(t)
		path := *pathFlag + pkg
		filename := path + "/" + pkg + ".go"
		err = os.MkdirAll(path, os.ModeDir)
		if err != nil {
			panic(err)
		}
		err = f.Render(buf)
		if err != nil {
			panic(err)
		}
		if err := ioutil.WriteFile(filename, buf.Bytes(), 0644); err != nil {
			panic(err)
		}
		buf.Reset()
	}
	log.Println("finished", pwd)
}
