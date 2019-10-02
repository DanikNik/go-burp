package dumper

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
)

type Dumper struct {
	dbObj *sql.DB
	Log   *logrus.Logger
}

func NewDumper(log *logrus.Logger) *Dumper {
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

	d.dbObj = dbObj
	d.Log.Log(logrus.InfoLevel, "dumper inited successfully")
	return d
}

func (d *Dumper) DumpRequest(req *http.Request, toDataBase bool) (body string, err error) {
	dump, err := httputil.DumpRequest(req, true)
	if err != nil {
		return "", err
	}

	if toDataBase {
		_, _ = d.dbObj.Exec(`INSERT INTO requests (host, request_body) VALUES (?, ?)`, req.Host, string(dump))
	}

	return string(dump), nil
}
