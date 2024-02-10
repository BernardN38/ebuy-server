package databaseinit

import (
	"database/sql"
)

func InitializeProductTypesTable(db *sql.DB) error {
	// Execute SQL statements to insert common product types

	sqlStatements := `
		INSERT INTO product_types (type_name) VALUES ('Electronics') ON CONFLICT (type_name) DO NOTHING;
		INSERT INTO product_types (type_name) VALUES ('Clothing') ON CONFLICT (type_name) DO NOTHING;
		INSERT INTO product_types (type_name) VALUES ('Books') ON CONFLICT (type_name) DO NOTHING;
		INSERT INTO product_types (type_name) VALUES ('Home and Kitchen') ON CONFLICT (type_name) DO NOTHING;
		INSERT INTO product_types (type_name) VALUES ('Toys and Games') ON CONFLICT (type_name) DO NOTHING;
		INSERT INTO product_types (type_name) VALUES ('Sports and Outdoors') ON CONFLICT (type_name) DO NOTHING;
		INSERT INTO product_types (type_name) VALUES ('Beauty and Personal Care') ON CONFLICT (type_name) DO NOTHING;
		INSERT INTO product_types (type_name) VALUES ('Automotive') ON CONFLICT (type_name) DO NOTHING;
		INSERT INTO product_types (type_name) VALUES ('Health and Household') ON CONFLICT (type_name) DO NOTHING;
		INSERT INTO product_types (type_name) VALUES ('Tools and Home Improvement') ON CONFLICT (type_name) DO NOTHING;
		INSERT INTO product_types (type_name) VALUES ('Grocery') ON CONFLICT (type_name) DO NOTHING;
		INSERT INTO product_types (type_name) VALUES ('Movies and TV') ON CONFLICT (type_name) DO NOTHING;
		INSERT INTO product_types (type_name) VALUES ('Music') ON CONFLICT (type_name) DO NOTHING;
		INSERT INTO product_types (type_name) VALUES ('Pet Supplies') ON CONFLICT (type_name) DO NOTHING;
	`

	_, err := db.Exec(sqlStatements)
	if err != nil {
		return err
	}
	return nil
}
