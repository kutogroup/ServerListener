package models

import (
    m "github.com/kutogroup/kuto.api/models"
)

const (
    TableConns          = "conns"
    ColumnConnsID       = "id"
    ColumnConnsCreateAt = "create_at"
    ColumnConnsUpdateAt = "update_at"
    ColumnConnsDeleted  = "deleted"
    ColumnConnsServerID = "server_id"
    ColumnConnsConns    = "conns"
)

type Conns struct {
    ID       int64  `db:"id"        json:"id"        comment:""`
    CreateAt m.Time `db:"create_at" json:"create_at" comment:""`
    UpdateAt m.Time `db:"update_at" json:"update_at" comment:""`
    Deleted  int64  `db:"deleted"   json:"deleted"   comment:""`
    ServerID int64  `db:"server_id" json:"server_id" comment:""`
    Conns    int64  `db:"conns"     json:"conns"     comment:""`
}

func init() {
    m.ModelTables = append(m.ModelTables, Conns{})
}