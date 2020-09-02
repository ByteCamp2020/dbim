package logic

import (
	"bdim/src/internal/logic/conf"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type dbc struct {
	conn *sql.DB
}

func NewDbC(c *conf.MySql) *dbc {
	DbC := &dbc{
		conn: connect(c),
	}
	DbC.CreateTable("message")
	return DbC
}

func connect(c *conf.MySql) *sql.DB {
	address := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.Username, c.Password, c.Hostname, c.Port, c.Database)
	db, err := sql.Open("mysql", address)
	if err != nil {
	}
	return db
}

func (d *dbc) Close() error {
	return d.conn.Close()
}

func (d *dbc) CreateTable (table string) error {
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

type Message struct {
	Id int   `json:"id"`
	Uid string `json:"uid"`
	Roomid int `json:"room_id"`
	Msg string `json:"message"`
	Timestamp int `json:"time_stamp"`
	Visible int `json:"visible"`
}

func (d *dbc) GetMessage(uid string, roomid string, timestamp string) ([]Message, error) {
	query, args := GetArg(uid, roomid, timestamp)
	var rows *sql.Rows
	var err error
	if args == nil {
		rows, err = d.conn.Query(query)
	} else {
		rows, err = d.conn.Query(query, args...)
	}
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return nil, err
	}
	defer rows.Close()

	res := []Message{}

	for rows.Next() {
		var u Message

		err := rows.Scan(&u.Id, &u.Uid, &u.Roomid, &u.Msg, &u.Timestamp, &u.Visible)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return nil, err
		}
		res = append(res, u)
	}
	return res, nil
}

func GetArg(uid string, roomid string, timestamp string) (string, []interface{}) {
	args := make([]interface{}, 0)
	//query := "select id, uid, roomid, msg, timestamp, visible from message"
	query := "select * from message "
	if uid == "" && roomid == "" && timestamp == "" {
		return query, nil
	}
	query = query + " where "
	if uid != "" {
		query += "uid = ?"
		args = append(args, uid)
	}
	if roomid != "" {
		if len(args) > 0 {
			query += " AND roomid = ?"
		} else {
			query += "roomid = ?"
		}
		args = append(args, roomid)
	}
	if timestamp != "" {
		if len(args) > 0 {
			query += " AND timestamp = ?"
		} else {
			query += "timestamp = ?"
		}
		args = append(args, timestamp)
	}
	return query, args
}
