package models

import (
	m "github.com/kutogroup/kuto.api/models"
)

const (
	TableConnections          = "connections"
	ColumnConnectionsID       = "id"
	ColumnConnectionsCreateAt = "create_at"
	ColumnConnectionsUpdateAt = "update_at"
	ColumnConnectionsDeleted  = "deleted"
	ColumnConnectionsTCP      = "tcp"
	ColumnConnectionsUDP      = "udp"
	ColumnConnectionsServerID = "server_id"
)

type Connections struct {
	ID       int64  `db:"id"        json:"id"        comment:""`
	CreateAt m.Time `db:"create_at" json:"create_at" comment:""`
	UpdateAt m.Time `db:"update_at" json:"update_at" comment:""`
	Deleted  int64  `db:"deleted"   json:"deleted"   comment:""`
	TCP      int64  `db:"tcp"       json:"tcp"       comment:""`
	UDP      int64  `db:"udp"       json:"udp"       comment:""`
	ServerID int64  `db:"server_id" json:"server_id" comment:""`
}

func init() {
	m.ModelTables = append(m.ModelTables, Connections{})
}
