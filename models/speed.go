package models

import (
	m "github.com/kutogroup/kuto.api/models"
)

const (
	TableSpeed          = "speed"
	ColumnSpeedID       = "id"
	ColumnSpeedCreateAt = "create_at"
	ColumnSpeedUpdateAt = "update_at"
	ColumnSpeedDeleted  = "deleted"
	ColumnSpeedServerID = "server_id"
	ColumnSpeedReceive  = "receive"
	ColumnSpeedTransmit = "transmit"
)

type Speed struct {
	ID       int64  `db:"id"        json:"id"        comment:""`
	CreateAt m.Time `db:"create_at" json:"create_at" comment:""`
	UpdateAt m.Time `db:"update_at" json:"update_at" comment:""`
	Deleted  int64  `db:"deleted"   json:"deleted"   comment:""`
	ServerID int64  `db:"server_id" json:"server_id" comment:""`
	Receive  string `db:"receive"   json:"receive"   comment:""`
	Transmit string `db:"transmit"  json:"transmit"  comment:""`
}

func init() {
	m.ModelTables = append(m.ModelTables, Speed{})
}
