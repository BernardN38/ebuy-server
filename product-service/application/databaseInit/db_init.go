package databaseinit

import (
	"database/sql"
)

func InitializeProductTypesTable(db *sql.DB) error {
	// Execute SQL statements to insert common product types
	sqlStatements := `
		INSERT INTO product_types (type_name) VALUES ('Electronics');
		INSERT INTO product_types (type_name) VALUES ('Clothing');
		INSERT INTO product_types (type_name) VALUES ('Books');
		INSERT INTO product_types (type_name) VALUES ('Home and Kitchen');
		INSERT INTO product_types (type_name) VALUES ('Toys and Games');
		INSERT INTO product_types (type_name) VALUES ('Sports and Outdoors');
		INSERT INTO product_types (type_name) VALUES ('Beauty and Personal Care');
		INSERT INTO product_types (type_name) VALUES ('Automotive');
		INSERT INTO product_types (type_name) VALUES ('Health and Household');
		INSERT INTO product_types (type_name) VALUES ('Tools and Home Improvement');
		INSERT INTO product_types (type_name) VALUES ('Grocery');
		INSERT INTO product_types (type_name) VALUES ('Movies and TV');
		INSERT INTO product_types (type_name) VALUES ('Music');
		INSERT INTO product_types (type_name) VALUES ('Pet Supplies');
		INSERT INTO product_types (type_name) VALUES ('Office Products');
	`

	_, err := db.Exec(sqlStatements)
	if err != nil {
		return err
	}
	return nil
}
