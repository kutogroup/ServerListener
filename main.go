package main

import (
	m "ServerListener/models"
	"encoding/json"
	"fmt"
	"kuto/pkg"
	"kuto/utils"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

var ticks = 0
var db = pkg.NewDatabase("sl", "localhost:3306", "root", "root")
var logger = pkg.NewLogger(os.Stdout, true)

func main() {
	var servers []m.Server
	err := db.Select(&servers, "id>0")
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
					logger.I("start get server, username=%s, host=%s", s.Username, s.Host)
					sri, err := strconv.ParseInt(s.ReceiverStart, 10, 64)
					if err != nil {
						logger.E("get server receive err=%s", err)
						continue
					}

					sti, err := strconv.ParseInt(s.TransmitStart, 10, 64)
					if err != nil {
						logger.E("get server transmit err=%s", err)
						continue
					}

					r := strings.TrimRight(utils.CommandGetResult("./receive", s.Username, s.Host), "\n")
					t := strings.TrimRight(utils.CommandGetResult("./transmit", s.Username, s.Host), "\n")

					speed := &m.Speed{}
					ri, err := strconv.ParseInt(r, 10, 64)
					if err != nil {
						logger.E("get receive err=%s", err)
						continue
					}

					ti, err := strconv.ParseInt(t, 10, 64)
					if err != nil {
						logger.E("get transmit err=%s", err)
						continue
					}

					speed.ServerID = s.ID
					speed.Receive = strconv.FormatInt(ri-sri, 10)
					speed.Transmit = strconv.FormatInt(ti-sti, 10)
					logger.I("r=%s, t=%s", speed.Receive, speed.Transmit)
					err = db.Insert(speed)
					if err != nil {
						logger.E("insert db failed, err=%s", err)
						continue
					}

					c := utils.CommandGetResult("./conn", s.Username, s.Host, strconv.FormatInt(s.Port, 10))
					c = strings.Trim(c, " ")
					c = strings.Trim(c, "\n")
					logger.I("c=%s", c)
					conn := &m.Conns{}
					num, err := strconv.ParseInt(c, 10, 64)
					if err != nil {
						logger.E("parse conn to int failed, err=%s", err)
						continue
					}
					conn.Conns = num
					conn.ServerID = s.ID
					db.Insert(conn)
				}
			}
		}
	}()

	r := httprouter.New()
	r.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Write([]byte("Not available"))
	})
	r.GET("/speed", Speed)
	r.GET("/conns", Conns)
	r.ServeFiles("/html/*filepath", http.Dir("html/"))
	log.Fatal(http.ListenAndServe(":9090", r))
}

func initServer(servers []m.Server) {
	for _, s := range servers {
		r := strings.TrimRight(utils.CommandGetResult("./receive", s.Username, s.Host), "\n")
		t := strings.TrimRight(utils.CommandGetResult("./transmit", s.Username, s.Host), "\n")
		logger.I("host=%s, r=%s, t=%s", s.Host, r, t)

		ri, err := strconv.ParseInt(r, 10, 64)
		if err != nil {
			logger.E("init server receive err=%s", err)
			continue
		}

		ti, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			logger.E("init server transmit err=%s", err)
			continue
		}

		s.ReceiverStart = strconv.FormatInt(ri, 10)
		s.TransmitStart = strconv.FormatInt(ti, 10)
		err = db.Update(&s, m.ColumnServerReceiverStart, m.ColumnServerTransmitStart)
		if err != nil {
			logger.E("update db failed, err=%s", err)
		}
	}
}

//Speed 流量监听
func Speed(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var speeds []m.Speed

	y, mo, d := time.Now().Date()

	r.ParseForm()
	sid := r.Form.Get("sid")
	cs := r.Form.Get("startDate")
	ce := r.Form.Get("endDate")

	if len(cs) != 10 {
		cs = fmt.Sprintf("%04d-%02d-%02d", y, mo, d)
	}

	if len(ce) != 10 {
		ce = fmt.Sprintf("%04d-%02d-%02d 23:59:59", y, mo, d)
	} else {
		ce = ce + " 23:59:59"
	}
	fmt.Println("cs=" + cs + ", ce=" + ce)

	err := db.Select(&speeds,
		fmt.Sprintf("%s=%s AND %s>='%s' AND %s<='%s'",
			m.ColumnSpeedServerID,
			sid,
			m.ColumnServerCreateAt, cs,
			m.ColumnServerCreateAt, ce))

	if err != nil {
		logger.E("select today failed, err=%s", err)
		w.WriteHeader(401)
		return
	}

	b, err := json.Marshal(speeds)
	if err != nil {
		w.WriteHeader(402)
		return
	}

	w.Write(b)
}

//Conns 获取服务器连接数
func Conns(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var conns []m.Conns

	y, mo, d := time.Now().Date()

	r.ParseForm()
	sid := r.Form.Get("sid")
	cs := r.Form.Get("startDate")
	ce := r.Form.Get("endDate")

	if len(cs) != 10 {
		cs = fmt.Sprintf("%04d-%02d-%02d", y, mo, d)
	}

	if len(ce) != 10 {
		ce = fmt.Sprintf("%04d-%02d-%02d 23:59:59", y, mo, d)
	} else {
		ce = ce + " 23:59:59"
	}
	fmt.Println("cs=" + cs + ", ce=" + ce)

	err := db.Select(&conns,
		fmt.Sprintf("%s=%s AND %s>='%s' AND %s<='%s'",
			m.ColumnSpeedServerID,
			sid,
			m.ColumnServerCreateAt, cs,
			m.ColumnServerCreateAt, ce))

	if err != nil {
		logger.E("select today failed, err=%s", err)
		w.WriteHeader(401)
		return
	}

	b, err := json.Marshal(conns)
	if err != nil {
		w.WriteHeader(402)
		return
	}

	w.Write(b)
}