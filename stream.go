package dbstream

import (
	"database/sql"
	"fmt"
)

type Scannable interface {
	Scan(rows *sql.Rows) error
}

type RowScannerFn[StreamType any] func(eachRow *sql.Rows, into *StreamType) error

// Stream allows iterating over a collection of typed structures
type Stream[StreamType Scannable] interface {
	Next() bool
	Current() (*StreamType, error)
	Close() error
}

func RowStreamFrom[StreamType any](rows *sql.Rows, rowScanner RowScannerFn[StreamType]) *RowStream[StreamType] {
	var obj StreamType
	return &RowStream[StreamType]{rows: rows, scanRow: rowScanner, current: obj}
}

type RowStream[StreamType any] struct {
	current   StreamType
	rows      *sql.Rows
	needsScan bool
	hasNext   bool
	scanRow   RowScannerFn[StreamType]
}

func (rs *RowStream[StreamType]) Current() (*StreamType, error) {
	if !rs.hasNext {
		return nil, fmt.Errorf("stream has no more elements or is closed")
	}
	if rs.needsScan {
		rs.scanRow(rs.rows, &rs.current)
	}
	return &rs.current, nil
}

func (rs *RowStream[StreamType]) Next() bool {
	rs.hasNext = rs.rows.Next()
	rs.needsScan = rs.hasNext
	return rs.hasNext
}

func (rs *RowStream[StreamType]) Close() error {
	rs.hasNext = false
	rs.needsScan = false
	return rs.rows.Close()
}
