package sequence

import (
	"database/sql"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// 建立MySQL链接 执行REPLACE INSERT语句
// REPLACE INTO sequence (stub) value ('a');
// SELECT LAST_INSERT_ID;

const sqlReplaceIntoStub = `REPLACE INTO sequence (stub) VALUES ('a')`

type Mysql struct {
	conn sqlx.SqlConn
}

func NewMySQL(dsn string) Sequence {
	return &Mysql{
		conn: sqlx.NewMysql(dsn),
	}
}

// Next 取下一个号
func (m *Mysql) Next() (seq uint64, err error) {
	// prepare 做准备 -> 预编译
	var stmt sqlx.StmtSession
	stmt, err = m.conn.Prepare(sqlReplaceIntoStub)
	if err != nil {
		logx.Errorw("m.conn.Prepare failed", logx.LogField{Key: "err", Value: err.Error()})
		return 0, err
	}
	defer stmt.Close()

	// 执行
	var rest sql.Result
	rest, err = stmt.Exec()
	if err != nil {
		logx.Errorw("stmt.Exec failed", logx.LogField{Key: "err", Value: err.Error()})
		return 0, err
	}

	// 获取刚插入的ID作为被转为短链的初始数据
	var lid int64
	lid, err = rest.LastInsertId()
	if err != nil {
		logx.Errorw("rest.LastInsertId failed", logx.LogField{Key: "err", Value: err.Error()})
		return 0, err
	}
	return uint64(lid), nil
}
