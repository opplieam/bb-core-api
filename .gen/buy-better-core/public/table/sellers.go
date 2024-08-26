//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var Sellers = newSellersTable("public", "sellers", "")

type sellersTable struct {
	postgres.Table

	// Columns
	ID   postgres.ColumnInteger
	Name postgres.ColumnString
	URL  postgres.ColumnString

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type SellersTable struct {
	sellersTable

	EXCLUDED sellersTable
}

// AS creates new SellersTable with assigned alias
func (a SellersTable) AS(alias string) *SellersTable {
	return newSellersTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new SellersTable with assigned schema name
func (a SellersTable) FromSchema(schemaName string) *SellersTable {
	return newSellersTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new SellersTable with assigned table prefix
func (a SellersTable) WithPrefix(prefix string) *SellersTable {
	return newSellersTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new SellersTable with assigned table suffix
func (a SellersTable) WithSuffix(suffix string) *SellersTable {
	return newSellersTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newSellersTable(schemaName, tableName, alias string) *SellersTable {
	return &SellersTable{
		sellersTable: newSellersTableImpl(schemaName, tableName, alias),
		EXCLUDED:     newSellersTableImpl("", "excluded", ""),
	}
}

func newSellersTableImpl(schemaName, tableName, alias string) sellersTable {
	var (
		IDColumn       = postgres.IntegerColumn("id")
		NameColumn     = postgres.StringColumn("name")
		URLColumn      = postgres.StringColumn("url")
		allColumns     = postgres.ColumnList{IDColumn, NameColumn, URLColumn}
		mutableColumns = postgres.ColumnList{NameColumn, URLColumn}
	)

	return sellersTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:   IDColumn,
		Name: NameColumn,
		URL:  URLColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
