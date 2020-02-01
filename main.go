package main

import (
	m "ServerListener/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kutogroup/kuto.api/pkg"
	"github.com/kutogroup/kuto.api/utils"

	"github.com/julienschmidt/httprouter"
)

var ticks = 0
var servers []m.Server
var db = pkg.NewDatabase("sl", "localhost:3306", "root", "root")
var email = pkg.NewEmail("sun.zg@outlook.com", "Szg20130515", "smtp-mail.outlook.com", 587, false)
var logger *pkg.WahaLogger
var replaceIPTable = make(map[int64]bool)
var emailTable = make(map[int64]bool)

func main() {
	logFile, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}

	logger = pkg.NewLogger(logFile, true)
	err = db.Select(&servers, "id>0 ORDER BY title")
	if err != nil {
		panic(err)
	}

	if len(os.Args) >= 3 && os.Args[1] == "replace_ip" {
		fmt.Println("start replace_ip")
		fmt.Println(utils.CommandGetResult("./aws/aws_replace_ip", os.Args[2], "-R"))
		return
	}

	go func() {
		for {
			today := time.Now().Format("2006-01-02")
			for _, s := range servers {
				var conns []m.Connections
				err := db.Select(&conns, fmt.Sprintf("%s=%d and create_at>'%s' ORDER BY create_at DESC limit 0, 2", m.ColumnConnectionsID, s.ID, today))

				if err == nil {
					if len(conns) == 2 {
						if (conns[0].TCP < 10 && conns[0].UDP < 10) &&
							(conns[1].TCP < 10 && conns[1].UDP < 10) {

							if _, ok := emailTable[s.ID]; !ok {
								emailTable[s.ID] = true
								email.Send("kutogroup@outlook.com", "Server blocked", fmt.Sprintf("title=%s, host=%s", s.Title, s.Host))
							}

							if s.Auto == 0 {
								continue
							}

							if _, ok := replaceIPTable[s.ID]; ok {
								logger.I("need to replace id, but replace yet, s=%s, ip=%s", s.Title, s.Host)
							} else {
								logger.I("need to replace id, s=%s, ip=%s", s.Title, s.Host)
								replaceIPTable[s.ID] = true
							}

							utils.CommandGetResult("./aws/aws_replace_ip", s.Host, "-R")
						} else {
							delete(replaceIPTable, s.ID)
							delete(emailTable, s.ID)
						}
					}
				} else {
					logger.I("get conns failed, err=%s", err)
				}
			}

			time.Sleep(time.Minute)
		}
	}()

	go func() {
		for {
			if ticks%1 == 0 {
				//每隔一分钟刷新服务器
				var ns []m.Server
				err := db.Select(&ns, "id>0 ORDER BY title")
				if err == nil {
					servers = ns
				}
			}

			if ticks%10 == 0 {
				for _, s := range servers {
					go func(s m.Server) {
						resp, err := http.Get(fmt.Sprintf("http://%s:31676/conn", s.Host))
						if err != nil {
							logger.E("get server conn failed")
							return
						}
						body, err := ioutil.ReadAll(resp.Body)
						if err != nil {
							logger.E("read conn body failed")
							return
						}

						result := string(body)
						subs := strings.Split(result, ",")
						if len(subs) == 2 {
							tcp, err := strconv.ParseInt(subs[0], 10, 64)
							if err != nil {
								logger.E("convert tcp failed")
								return
							}
							udp, err := strconv.ParseInt(subs[1], 10, 64)
							if err != nil {
								logger.E("convert udp failed")
								return
							}

							conn := &m.Connections{}
							conn.ServerID = s.ID
							conn.TCP = tcp
							conn.UDP = udp
							db.Insert(conn)
						}
					}(s)
				}
			}

			ticks = ticks + 1

			time.Sleep(time.Minute)
		}
	}()

	r := httprouter.New()
	r.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Write([]byte("Not available"))
	})
	r.GET("/conns/", Conns)
	r.GET("/single/:id", Single)
	r.ServeFiles("/html/*filepath", http.Dir("html/"))
	log.Fatal(http.ListenAndServe(":9090", r))
}

//Speed 流量监听
func Speed(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	y, mo, d := time.Now().Date()

	r.ParseForm()

	var country = ps.ByName("country")
	result := make([]interface{}, 0)
	for _, server := range servers {
		if len(country) > 0 && server.Country != "ALL" && !strings.Contains(server.Country, country) {
			continue
		}

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

		var speeds []m.Speed

		err := db.Select(&speeds,
			fmt.Sprintf("%s=%d AND %s>='%s' AND %s<='%s' ORDER BY %s DESC LIMIT 0,1",
				m.ColumnSpeedServerID,
				server.ID,
				m.ColumnServerCreateAt, cs,
				m.ColumnServerCreateAt, ce,
				m.ColumnServerID))

		if err != nil {
			logger.E("select today failed, err=%s", err)
			w.WriteHeader(401)
			return
		}

		result = append(result, map[string]interface{}{
			m.ColumnServerTitle:      server.Title,
			m.ColumnServerChartColor: server.ChartColor,
			"speeds":                 speeds,
		})
	}

	b, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(402)
		return
	}

	w.Write(b)
}

//Conns 获取服务器连接数
func Conns(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	y, mo, d := time.Now().Date()

	r.ParseForm()

	result := make([]interface{}, 0)
	var country = ps.ByName("country")
	for _, server := range servers {
		if len(country) > 0 && server.Country != "ALL" && !strings.Contains(server.Country, country) {
			continue
		}

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

		var conns []m.Connections

		err := db.Select(&conns,
			fmt.Sprintf("%s=%d AND %s>='%s' AND %s<='%s' ORDER BY %s DESC LIMIT 0,1",
				m.ColumnConnectionsServerID,
				server.ID,
				m.ColumnConnectionsCreateAt, cs,
				m.ColumnConnectionsCreateAt, ce,
				m.ColumnServerID))

		if err != nil {
			logger.E("select today failed, err=%s", err)
			w.WriteHeader(401)
			return
		}

		result = append(result, map[string]interface{}{
			m.ColumnServerTitle:         server.Title,
			m.ColumnServerChartColor:    server.ChartColor,
			m.ColumnConnectionsServerID: server.ID,
			"conns":                     conns,
		})
	}

	b, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(402)
		return
	}

	w.Write(b)
}

//Single 获取服务器连接数
func Single(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	y, mo, d := time.Now().Date()

	r.ParseForm()

	var idStr = ps.ByName("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.E("no id specific, %s", err)
		w.WriteHeader(402)
		return
	}

	var server m.Server
	for _, s := range servers {
		if s.ID == id {
			server = s
			break
		}
	}

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

	var conns []m.Connections
	err = db.Select(&conns,
		fmt.Sprintf("%s=%d AND %s>='%s' AND %s<='%s'",
			m.ColumnConnectionsServerID,
			server.ID,
			m.ColumnConnectionsCreateAt, cs,
			m.ColumnConnectionsCreateAt, ce))

	if err != nil {
		logger.E("select today failed, err=%s", err)
		w.WriteHeader(401)
		return
	}

	result := map[string]interface{}{
		m.ColumnServerTitle:         server.Title,
		m.ColumnServerChartColor:    server.ChartColor,
		m.ColumnConnectionsServerID: server.ID,
		"conns":                     conns,
	}

	b, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(402)
		return
	}

	w.Write(b)
}
