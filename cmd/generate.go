/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

type TemplateData struct {
	PackageName string
	StructName  string
	DBType      string
	DBSnippet   string
}

const controllerTemplate = `package {{.PackageName}}

type {{.StructName}}Controller struct {

}

func New{{.StructName}}Controller() *{{.StructName}}Controller {
	return &{{.StructName}}Controller{}
}
`

const repositoryTemplate = `package {{.PackageName}}

import "database/sql"

type {{.StructName}}Repository struct {
	DB *sql.DB
}

func New{{.StructName}}Repository(db *sql.DB) *{{.StructName}}Repository {
	return &{{.StructName}}Repository{DB: db}
}

/*
Konfigurasi Driver {{.DBType}}:
{{.DBSnippet}}
*/
`

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		folderName := strings.ToLower(moduleName)
		chosenDB := strings.ToLower(dbType)

		structName := strings.Title(folderName)

		fmt.Printf("Mulai membuat boilerplate untuk modul: %s (Database: %s)...\n", folderName, chosenDB)

		err := os.MkdirAll(folderName, 0755)

		if err != nil {
			fmt.Printf("Gagal membuat folder: %v\n", err)
			return
		}

		var dbSnippet string
		if chosenDB == "mysql" {
			dbSnippet = `// Driver: github.com/go-sql-driver/mysql
// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?parseTime=true"
// db, err := sql.Open("mysql", dsn)`
		
		}else {
			chosenDB = "postgres"
			dbSnippet = `// Driver: github.com/lib/pq
// dsn := "host=localhost user=postgres password=secret dbname=mydb sslmode=disable"
// db, err := sql.Open("postgres", dsn)`
		}
		data := TemplateData{
			PackageName: folderName,
			StructName:  structName,
			DBType:      chosenDB,
			DBSnippet:   dbSnippet,
		}
		templates := map[string]string{
			"controller.go": controllerTemplate,
			"repository.go": repositoryTemplate,
		}


		for fileName, tmplStr := range templates {
			tmpl, err := template.New(fileName).Parse(tmplStr)
			if err != nil {
				fmt.Printf("Gagal memproses template %s: %v\n", fileName, err)
				continue
			}

			var buf bytes.Buffer
			err = tmpl.Execute(&buf, data)
			if err != nil {
				fmt.Printf("Gagal mengeksekusi template %s: %v\n", fileName, err)
				continue
			}

			filePath := filepath.Join(folderName, fileName)
			err = os.WriteFile(filePath, buf.Bytes(), 0644)
			if err != nil {
				fmt.Printf("Gagal membuat file %s: %v\n", fileName, err)
			} else {
				fmt.Printf("✔️ Berhasil membuat %s (via Template)\n", filePath)
			}
		}

		fmt.Println("Boilerplate dinamis berhasil di-generate!")
	},
}

var moduleName string
var dbType string

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVarP(&moduleName, "name", "n", "", "Nama modul atau fitur (wajib diisi)")
	generateCmd.MarkFlagRequired("name")
	generateCmd.Flags().StringVarP(&dbType, "db", "d", "postgres", "Jenis database (postgres/mysql)")
}
