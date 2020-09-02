package dao

import (
	"bdim/src/internal/dbworker/conf"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type Dao struct {
	conn *sql.DB
}

func New(c *conf.MySql) *Dao {
	dao := &Dao{
		conn: connect(c),
	}
	dao.CreateTable("message")
	return dao
}

func connect(c *conf.MySql) *sql.DB {
	address := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.Username, c.Password, c.Hostname, c.Port, c.Database)
	db, err := sql.Open("mysql", address)
	if err != nil {
	}
	return db
}

func (d *Dao) Close() error {
	return d.conn.Close()
}

func (d *Dao) CreateTable (table string) error {
	sql := `
    CREATE TABLE IF NOT EXISTS ` + table + ` (
    id INT(11) NOT NULL AUTO_INCREMENT,
    uid VARCHAR(255) NOT NULL,
    roomid INT(11) NOT NULL,
    msg VARCHAR(255) NOT NULL,
    timestamp INT(11) NOT NULL,
    visible TINYINT(1) NOT NULL,
    PRIMARY KEY (id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 AUTO_INCREMENT=1;`

	fmt.Println("\n" + sql + "\n")
	smt, _ := d.conn.Prepare(sql)
	_, err := smt.Exec()
	if (err != nil) {
		return err
	}
	return nil
}

func (d *Dao) AddMessage(uid string, roomid int32, msg string, timestamp int32, visible bool) {
	var vis int8
	if visible == true {
		vis = 1
	} else {
		vis = 0
	}
	execStr := fmt.Sprintf(`INSERT INTO message
		(uid, roomid, msg, timestamp, visible)
		VALUES ('%s', %v, '%s', %v, %v)`,
		uid, roomid, msg, timestamp, vis)
	ret, err := d.conn.Exec(execStr)
	if err != nil {

	}
	if LastInsertId, err := ret.LastInsertId(); nil == err {
		fmt.Println("Dao.mySql:LastInsertId:", LastInsertId)
	}
	if RowsAffected, err := ret.RowsAffected(); nil == err {
		fmt.Println("Dao.mySql:RowsAffected:", RowsAffected)
	}
}
