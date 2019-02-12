package generate

import (
	j "github.com/dave/jennifer/jen"
	tbu "github.com/pindamonhangaba/tabua"
	"github.com/pindamonhangaba/tabua/reverse"
)

// Generator generates types for a table
type Generator struct {
	PackagePath string
}

// Run generates a jenifer.File and returns the package name
func (g *Generator) Run(t reverse.Table) (*j.File, string) {
	return buildTable(t, g.PackagePath), packageFilename(t.Name)
}

func buildTable(t reverse.Table, pkgPath string) *j.File {
	pkgName := packageFilename(t.Name)
	tableName := camel(t.Name)
	file := j.NewFile(pkgName)

	file.HeaderComment("This file is generated - do not edit.")
	file.Line()

	// table struct
	fields := []j.Code{}
	for _, c := range t.Columns {
		colName := columnName(t.Name, c.Name)
		fields = append(fields, j.Id(colName).Id(colName).Tag(map[string]string{"db": c.Name, "json": camelLower(c.Name)}))
	}
	file.Type().Id(tableName).Struct(fields...)
	file.Line()

	// implement tabua.Namer
	file.Comment("Name implements the tabua.Namer interface.")
	file.Func().Params(
		j.Id("t").Id(tableName),
	).Id("Name").Params().String().Block(
		j.Return(j.Lit(t.Name)),
	)

	// implement tabua.Table
	tableColumns := []j.Code{}
	for _, c := range t.Columns {
		colName := columnName(t.Name, c.Name)
		tableColumns = append(tableColumns, j.Id("t").Dot(colName))
	}
	file.Comment("Columns implements the tabua.Table interface.")
	file.Func().Params(
		j.Id("t").Id(tableName),
	).Id("Columns").Params().Index().Qual("github.com/pindamonhangaba/tabua", "Column").Block(
		j.Return(j.Index().Qual("github.com/pindamonhangaba/tabua", "Column").Values(tableColumns...)),
	)
	file.Line()

	tableConstraints := []j.Code{}
	for _, c := range t.Constraints {
		colName := cstname(c.Name)
		tableConstraints = append(tableConstraints, j.Id(colName))
	}
	file.Comment("Constraints implements the tabua.Table interface.")
	file.Func().Params(
		j.Id("t").Id(tableName),
	).Id("Constraints").Params().Index().Qual("github.com/pindamonhangaba/tabua", "Constraint").Block(
		j.Return(j.Index().Qual("github.com/pindamonhangaba/tabua", "Constraint").Values(tableConstraints...)),
	)
	file.Line()

	// constraints types
	// implement tabua.Constrainer
	for _, c := range t.Constraints {
		n := cstname(c.Name)

		file.Commentf("%s is a constraintforthe table \"%s\", a %s", n, tableName, c.Type)
		file.Const().Id(n).Op("=").Lit(c.Name)

		file.Comment("Type implements tbu.Constrainer")
		file.Func().Params(
			j.Id("c").Id(n),
		).Id("Type").Params().Qual("github.com/pindamonhangaba/tabua", "ConstraintType").Block(
			j.Return(j.Qual("github.com/pindamonhangaba/tabua", "ConstraintType").Parens(j.Lit(c.Type))),
		)
		file.Comment("Definition implements tbu.Constrainer")
		file.Func().Params(
			j.Id("c").Id(n),
		).Id("Definition").Params().String().Block(
			j.Return(j.Lit(c.Definition)),
		)

		switch tbu.ConstraintType(c.Type) {
		case tbu.ConstraintUnique:
			cols := []j.Code{}
			for _, c := range c.ColumnsLocal {
				colName := colName(c.Table, c.Column)
				cols = append(cols, j.Id(colName))
			}
			file.Comment("Uniques implements tbu.UniqueConstrainer")
			file.Func().Params(
				j.Id("c").Id(n),
			).Id("Uniques").Params().Index().Qual("github.com/pindamonhangaba/tabua", "Column").Block(
				j.Return(j.Index().Qual("github.com/pindamonhangaba/tabua", "Column").Values(cols...)),
			)
			break
		case tbu.ConstraintCheck:
			cols := []j.Code{}
			for _, c := range c.ColumnsLocal {
				colName := colName(c.Table, c.Column)
				cols = append(cols, j.Id(colName))
			}
			file.Comment("Columns implements tbu.CheckConstrainer")
			file.Func().Params(
				j.Id("c").Id(n),
			).Id("Columns").Params().Index().Qual("github.com/pindamonhangaba/tabua", "Column").Block(
				j.Return(j.Index().Qual("github.com/pindamonhangaba/tabua", "Column").Values(cols...)),
			)
			break
		case tbu.ConstraintPK:
			cols := []j.Code{}
			for _, c := range c.ColumnsLocal {
				colName := colName(c.Table, c.Column)
				cols = append(cols, j.Id(colName))
			}
			file.Comment("Keys implements tbu.PKConstrainer")
			file.Func().Params(
				j.Id("c").Id(n),
			).Id("Keys").Params().Index().Qual("github.com/pindamonhangaba/tabua", "Column").Block(
				j.Return(j.Index().Qual("github.com/pindamonhangaba/tabua", "Column").Values(cols...)),
			)
			break
		case tbu.ConstraintFK:
			cols := []j.Code{}
			colsf := []j.Code{}
			for _, c := range c.ColumnsLocal {
				colName := colName(c.Table, c.Column)
				cols = append(cols, j.Id(colName))
			}
			for _, c := range c.ColumnsForeign {
				colName := colName(c.Table, c.Column)
				colsf = append(colsf, j.Qual(pkgPath+packageFilename(c.Table), colName))
			}
			file.Comment("Key implements tbu.FKConstrainer")
			file.Func().Params(
				j.Id("c").Id(n),
			).Id("Key").Params().Qual("github.com/pindamonhangaba/tabua", "FK").Block(
				j.Return(j.Qual("github.com/pindamonhangaba/tabua", "FK").Values(j.Dict{
					j.Id("From"): j.Index().Qual("github.com/pindamonhangaba/tabua", "Column").Values(cols...),
					j.Id("To"):   j.Index().Qual("github.com/pindamonhangaba/tabua", "Column").Values(colsf...),
				})),
			)
			break
		}

	}
	file.Line()

	// implement tabua.Column
	for _, c := range t.Columns {
		colName := columnName(t.Name, c.Name)
		ctype := reType(c, c.NonNull)
		file.Commentf("%s is the column type for the table \"%s\", a %s", colName, tableName, ctype.String())
		if len(ctype.Name()) == 0 || len(ctype.PkgPath()) == 0 {
			file.Type().Id(colName).Id(ctype.String())
		} else {
			file.Type().Id(colName).Qual(ctype.PkgPath(), ctype.Name())
		}

		file.Comment("Name implements the tabua.Namer interface.")
		file.Func().Params(
			j.Id("c").Id(colName),
		).Id("Name").Params().String().Block(
			j.Return(j.Lit(c.Name)),
		)

		file.Comment("SQLType implements the tabua.Column interface.")
		file.Func().Params(
			j.Id("c").Id(colName),
		).Id("SQLType").Params().String().Block(
			j.Return(j.Lit(c.UDTName + "," + c.DataType)),
		)

		file.Comment("NonNull implements the tabua.Column interface.")
		file.Func().Params(
			j.Id("c").Id(colName),
		).Id("NonNull").Params().Bool().Block(
			j.Return(j.Lit(c.NonNull)),
		)

		file.Comment("Table implements the tabua.Column interface.")
		file.Func().Params(
			j.Id("c").Id(colName),
		).Id("Table").Params().Qual("github.com/pindamonhangaba/tabua", "Table").Block(
			j.Return(j.Id(tableName).Block()),
		)
	}
	return file
}
