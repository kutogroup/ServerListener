package models

import (
    m "kuto/models"
)

const (
    TableServer               = "server"
    ColumnServerID            = "id"
    ColumnServerCreateAt      = "create_at"
    ColumnServerUpdateAt      = "update_at"
    ColumnServerDeleted       = "deleted"
    ColumnServerHost          = "host"
    ColumnServerPort          = "port"
    ColumnServerTitle         = "title"
    ColumnServerChartColor    = "chart_color"
    ColumnServerUsername      = "username"
    ColumnServerReceiverStart = "receiver_start"
    ColumnServerTransmitStart = "transmit_start"
)

type Server struct {
    ID            int64  `db:"id"             json:"id"             comment:""`
    CreateAt      m.Time `db:"create_at"      json:"create_at"      comment:""`
    UpdateAt      m.Time `db:"update_at"      json:"update_at"      comment:""`
    Deleted       int64  `db:"deleted"        json:"deleted"        comment:""`
    Host          string `db:"host"           json:"host"           comment:""`
    Port          int64  `db:"port"           json:"port"           comment:""`
    Title         string `db:"title"          json:"title"          comment:""`
    ChartColor    string `db:"chart_color"    json:"chart_color"    comment:""`
    Username      string `db:"username"       json:"username"       comment:""`
    ReceiverStart string `db:"receiver_start" json:"receiver_start" comment:""`
    TransmitStart string `db:"transmit_start" json:"transmit_start" comment:""`
}

func init() {
    m.ModelTables = append(m.ModelTables, Server{})
}