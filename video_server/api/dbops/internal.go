package dbops

import (
	"database/sql"
	"log"
	"strconv"
	"sync"
	"video_server/api/defs"
)

func InserSession(sid string, ttl int64, username string) error {
	ttlstr := strconv.FormatInt(ttl, 10)
	stmtIns, err := dbConn.Prepare("INSERT INTO sessions (session_id, TTL, username) VALUES(?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = stmtIns.Exec(sid, ttlstr, username)
	stmtIns.Close()
	if err != nil {
		return err
	}
	return nil
}

func RetrieveSession(sid string) (*defs.SimpleSession, error) {
	ss := &defs.SimpleSession{}
	stmtOut, err := dbConn.Prepare("SELECT TTL, username FROM sessions WHERE session_id = ?")
	if err != nil {
		return nil, err
	}

	var ttl string
	var username string
	err = stmtOut.QueryRow(sid).Scan(&ttl, &username)
	stmtOut.Close()
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if ttlint, err := strconv.ParseInt(ttl, 10, 64); err == nil {
		ss.TTL = ttlint
		ss.Username = username
	} else {
		return nil, err
	}
	return ss, nil
}

func RetrieveAllSessions() (*sync.Map, error) {
	m := &sync.Map{}
	stmtOut, err := dbConn.Prepare("SELECT * FROM sessions")
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}

	rows, err := stmtOut.Query()
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}

	for rows.Next() {
		var id string
		var ttlstr string
		var username string

		if er := rows.Scan(&id, &ttlstr, &username); er != nil {
			log.Printf("retrive sessions error：%s", er)
			break
		}

		if ttl, err1 := strconv.ParseInt(ttlstr, 10, 64); err1 == nil {
			ss := &defs.SimpleSession{
				Username: username,
				TTL:      ttl,
			}
			m.Store(id, ss)
			log.Printf("session id：%s, ttl：%d", id, ss.TTL)
		}

	}
	return m, nil
}

func DeleteSession(sid string) error {
	stmtOut, err := dbConn.Prepare("DELETE FROM sessions WHERE session_id = ?")
	if err != nil {
		log.Printf("%s", err)
		return err
	}

	if _, err := stmtOut.Query(sid); err != nil {
		return err
	}

	return nil
}
