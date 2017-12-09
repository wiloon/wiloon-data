package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/wiloon/wiloon-log/log"
)

type Database struct {
	conf Config
	conn *sql.DB
}

func NewDatabase(config Config) Database {
	db := Database{conf: config}
	db.Connect()
	return db
}
func (db *Database) Connect() {
	conf := db.conf
	conn, err := sql.Open("mysql", conf.Username+":"+conf.Password+"@tcp("+conf.Address+")/"+conf.DatabaseName+"?charset=utf8")
	if err != nil {
		log.Info("failed to connect %v, err:%v", conf.DatabaseName, err)
	}
	db.conn = conn
}

func (db *Database) Close() {
	db.conn.Close()
}

func (db *Database) Count(table string) int {
	stmt := "select count(*) as c from " + table
	rows, err := db.conn.Query(stmt)
	if err != nil {
		log.Info("failed to query, sql:%v, err:%v", stmt, err)
	}
	defer rows.Close()
	var c int
	for rows.Next() {
		rows.Scan(&c)
	}
	return c
}

func (db *Database) Find(stmt string, args ...interface{}) []map[string]interface{} {
	var result []map[string]interface{}
	log.Info("sql:%v, args:%v", stmt, args)
	rows, err := db.conn.Query(stmt, args...)
	if err != nil {
		log.Info("failed to query, sql:%v, err:%v", stmt, err)
	}
	defer rows.Close()

	cols, _ := rows.Columns()

	for rows.Next() {
		row := make(map[string]interface{})
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			log.Info("failed to scan rows.", err)
		}

		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			row[colName] = *val
		}
		result = append(result, row)
	}
	return result
}

func (db *Database) Save(stmt string, args ...interface{}) {

	result, err := db.conn.Exec(stmt, args...)
	if err != nil {
		log.Info("failed to query, sql:%v, err:%v", stmt, err)
	}
	log.Debug("result:", result)

}
