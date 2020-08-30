package dao

import ("database/sql"
		"bdim/src/internal/dbworker/conf"
		"fmt"
)

type Dao struct {
	conn *sql.DB
}

func New(c *conf.MySql) *Dao {
	dao := &Dao {
		conn: connect(c),
	}
	return dao
}

func connect(c *conf.MySql) *sql.DB {
	address := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.Username, c.Password, c.Hostname, c.Port, c.Database)
	db, err := sql.Open("mysql", address)
	if err != nil {

	}
	return db
}

func (d *Dao) AddMessage(uid string, roomid int, msg string, timestamp int, visible bool) {
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
