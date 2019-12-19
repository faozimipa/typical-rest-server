package typrails_test

import (
	"errors"
	"regexp"
	"testing"

	"github.com/typical-go/typical-rest-server/pkg/typrails"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestFetcher(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	testcases := []struct {
		data map[string]string
		err  error
		*typrails.Entity
	}{
		{
			data: map[string]string{
				"column1": "type1",
				"column2": "type2",
			},
			err: errors.New("\"id\" with underlying data type \"int4\" is missing; \"updated_at\" with underlying data type \"timestamp\" is missing; \"created_at\" with underlying data type \"timestamp\" is missing"),
		},
		{
			data: map[string]string{
				"id":         "int4",
				"name":       "varchar",
				"created_at": "timestamp",
				"updated_at": "timestamp",
			},
			Entity: &typrails.Entity{
				Name:           "book",
				Type:           "Book",
				Table:          "books",
				Cache:          "BOOKS",
				ProjectPackage: "some-package",
				Fields: []typrails.Field{
					{Name: "ID", Type: "int64", Udt: "int4", Column: "id", StructTag: "`json:\"id\"`"},
					{Name: "Name", Type: "string", Udt: "varchar", Column: "name", StructTag: "`json:\"name\"`"},
					{Name: "CreatedAt", Type: "time.Time", Udt: "timestamp", Column: "created_at", StructTag: "`json:\"created_at\"`"},
					{Name: "UpdatedAt", Type: "time.Time", Udt: "timestamp", Column: "updated_at", StructTag: "`json:\"updated_at\"`"},
				},
				Forms: []typrails.Field{
					{Name: "Name", Type: "string", Udt: "varchar", Column: "name", StructTag: "`json:\"name\"`"},
				},
			},
		},
		{
			data: map[string]string{
				"id":         "int",
				"created_at": "timestamp",
				"updated_at": "timestamp",
			},
			err: errors.New("\"id\" with underlying data type \"int4\" is missing"),
		},
	}
	fetcher := typrails.Fetcher{DB: db}
	query := regexp.QuoteMeta("SELECT column_name, udt_name FROM information_schema.COLUMNS WHERE table_name = ?")
	for _, tt := range testcases {
		rows := sqlmock.NewRows([]string{"column_name", "data_type"})
		for key, value := range tt.data {
			rows.AddRow(key, value)
		}
		mock.ExpectQuery(query).WithArgs("books").WillReturnRows(rows)
		entity, err := fetcher.Fetch("some-package", "books")
		require.EqualValues(t, tt.err, err)
		require.EqualValues(t, tt.Entity, entity)
	}
}

func TestEntityName(t *testing.T) {
	testcases := []struct {
		TableName  string
		EntityName string
	}{
		{TableName: "book", EntityName: "book"},
		{TableName: "books", EntityName: "book"},
		{TableName: "aliases", EntityName: "alias"},
	}
	for _, tt := range testcases {
		require.Equal(t, tt.EntityName, typrails.EntityName(tt.TableName))
	}
}