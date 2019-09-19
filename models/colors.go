package models

import (
    m "github.com/kutogroup/kuto.api/models"
)

const (
    TableColors          = "colors"
    ColumnColorsID       = "id"
    ColumnColorsCreateAt = "create_at"
    ColumnColorsUpdateAt = "update_at"
    ColumnColorsDeleted  = "deleted"
    ColumnColorsColor    = "color"
)

type Colors struct {
    ID       int64  `db:"id"        json:"id"        comment:""`
    CreateAt m.Time `db:"create_at" json:"create_at" comment:""`
    UpdateAt m.Time `db:"update_at" json:"update_at" comment:""`
    Deleted  int64  `db:"deleted"   json:"deleted"   comment:""`
    Color    string `db:"color"     json:"color"     comment:""`
}

func init() {
    m.ModelTables = append(m.ModelTables, Colors{})
}