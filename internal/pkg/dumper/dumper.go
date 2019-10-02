package dumper

import (
	"bufio"
	"bytes"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"go-burp/internal/pkg/request"
	"log"
	"net/http"
	"net/http/httputil"
)

type Dumper struct {
	dbObj *sql.DB
	Log   *log.Logger
}

func NewDumper(log *log.Logger) *Dumper {
	d := &Dumper{
		Log: log,
	}
	dbObj, err := sql.Open("sqlite3", "requestbin.sqlite")
	if err != nil {
		panic(err)
	}

	err = dbObj.Ping()
	if err != nil {
		panic(err)
	}

	_, err = dbObj.Exec(`
CREATE TABLE IF NOT EXISTS requests(
id INTEGER PRIMARY KEY AUTOINCREMENT,
timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
host TEXT,
request_body TEXT
);	
`)
	if err != nil {
		panic(err)
	}
	_, err = dbObj.Exec(`
CREATE TABLE IF NOT EXISTS responses(
id INTEGER PRIMARY KEY AUTOINCREMENT,
timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
host TEXT,
response_body TEXT,
request_id INTEGER,
FOREIGN KEY(request_id) REFERENCES requests(id)
);	
`)
	if err != nil {
		panic(err)
	}

	d.dbObj = dbObj
	d.Log.Print("Dumper succesfully inited")
	return d
}

func (d *Dumper) DumpRequest(req *http.Request, toDataBase bool) (msg request.Message, err error) {
	dump, err := httputil.DumpRequest(req, true)
	msg = request.Message{}
	if err != nil {
		return msg, err
	}

	if toDataBase {
		_, _ = d.dbObj.Exec(`INSERT INTO requests (host, request_body) VALUES (?, ?)`, req.Host, string(dump))
		row := d.dbObj.QueryRow(`SELECT id, request_body FROM requests WHERE id = (SELECT last_insert_rowid());`)
		err := row.Scan(&msg.Id, &msg.Dump)
		if err != nil {
			panic(err)
		}
	}

	return msg, nil
}

func (d *Dumper) GetRequest(id int64) (req *http.Request) {
	row := d.dbObj.QueryRow("SELECT request_body FROM requests WHERE id=?", id)
	var body string
	err := row.Scan(&body)
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(bytes.NewReader([]byte(body)))
	req, err = http.ReadRequest(r)
	if err != nil {
		panic(err)
	}
	req.RequestURI = ""
	return req
}
