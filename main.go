package main

import (
	m "ServerListener/models"
	"kuto/config"
	"kuto/pkg"
	"kuto/utils"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var ticks = 0

var h = NewHTTP(
	nil,
	pkg.NewDatabase("sl", "localhost:3306", "root", "root"),
	pkg.NewLogger(os.Stdout, true),
	pkg.NewCache(config.CacheHost, config.CachePoolSize, time.Minute))

func main() {
	var servers []m.Server
	err := h.DB.Select(&servers, "id>0")
	if err != nil {
		panic(err)
	}

	if servers[0].ReceiverStart == "0" {
		initServer(servers)
	}

	ticker := time.NewTicker(time.Minute)
	go func() {
		for range ticker.C {
			now := time.Now()
			if now.Day() == 1 && now.Hour() == 0 && now.Minute() == 1 {
				initServer(servers)
				continue
			}

			ticks = ticks + 1

			if ticks%5 == 0 {
				for _, s := range servers {
					h.Logger.I("start get server, username=%s, host=%s", s.Username, s.Host)
					sri, err := strconv.ParseInt(s.ReceiverStart, 10, 64)
					if err != nil {
						h.Logger.E("get server receive err=%s", err)
						continue
					}

					sti, err := strconv.ParseInt(s.TransmitStart, 10, 64)
					if err != nil {
						h.Logger.E("get server transmit err=%s", err)
						continue
					}

					r := strings.TrimRight(utils.CommandGetResult("./receive", s.Username, s.Host), "\n")
					t := strings.TrimRight(utils.CommandGetResult("./transmit", s.Username, s.Host), "\n")
					h.Logger.I("r=%s, t=%s", r, t)

					speed := &m.Speed{}
					ri, err := strconv.ParseInt(r, 10, 64)
					if err != nil {
						h.Logger.E("get receive err=%s", err)
						continue
					}

					ti, err := strconv.ParseInt(t, 10, 64)
					if err != nil {
						h.Logger.E("get transmit err=%s", err)
						continue
					}

					speed.ServerID = s.ID
					speed.Receive = strconv.FormatInt(ri-sri, 10)
					speed.Transmit = strconv.FormatInt(ti-sti, 10)
					err = h.DB.Insert(speed)
					if err != nil {
						h.Logger.E("insert db failed, err=%s", err)
						continue
					}

					c := utils.CommandGetResult("./conn", s.Username, s.Host, strconv.FormatInt(s.Port, 10))
					h.Logger.I("c=%s", c)
					conn := &m.Conns{}
					num, err := strconv.ParseInt(c, 10, 64)
					if err != nil {
						continue
					}
					conn.Conns = num
					conn.ServerID = s.ID
					h.DB.Insert(conn)
				}
			}
		}
	}()

	h.Serve(":80", time.Minute)
}

func initServer(servers []m.Server) {
	for _, s := range servers {
		r := strings.TrimRight(utils.CommandGetResult("./receive", s.Username, s.Host), "\n")
		t := strings.TrimRight(utils.CommandGetResult("./transmit", s.Username, s.Host), "\n")
		h.Logger.I("r=%s, t=%s", r, t)

		ri, err := strconv.ParseInt(r, 10, 64)
		if err != nil {
			h.Logger.E("init server receive err=%s", err)
			continue
		}

		ti, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			h.Logger.E("init server transmit err=%s", err)
			continue
		}

		s.ReceiverStart = strconv.FormatInt(ri, 10)
		s.TransmitStart = strconv.FormatInt(ti, 10)
		err = h.DB.Update(&s, m.ColumnServerReceiverStart, m.ColumnServerTransmitStart)
		if err != nil {
			h.Logger.E("update db failed, err=%s", err)
		}
	}
}

//Index 首页
func (c *HTTP) Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not available"))
}

//NetStatus 网络状态
func (c *HTTP) NetStatus(w http.ResponseWriter, r *http.Request) {
}
