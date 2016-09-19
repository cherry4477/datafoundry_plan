package initialize

import (
	"encoding/json"
	"fmt"
	"github.com/asiainfoLDP/datafoundry_plan/ds"
	"github.com/asiainfoLDP/datahub_commons/mq"
	"github.com/astaxie/beego/logs"
	"github.com/miekg/dns"
	"net"
	"os"
	"sync/atomic"
	"time"
	"unsafe"
)

const (
	Platform_Local      = "local"
	Platform_DaoCloud   = "daocloud"
	Platform_DaoCloudUT = "daocloud_ut"
	Platform_DataOS     = "dataos"

	SENDER     = "datafoundry_plan"
	CHANNELLEN = 3000
)

var (
	log   *logs.BeeLogger
	theMQ unsafe.Pointer

	Platform = Platform_DaoCloud

	Debug = false
)

func init() {
	//initLog()

	logs.SetAlermSendingCallback(sendAlarm)
}

func InitLog() *logs.BeeLogger {

	log = logs.NewLogger(CHANNELLEN)

	err := log.SetLogger("console", "")
	if err != nil {
		log.Error("set logger err:", err)
		return nil
	}

	log.EnableFuncCallDepth(true)

	if Debug == false {
		fmt.Println("mode is info...")
		log.SetLevel(logs.LevelInfo)
	} else {
		fmt.Println("mode is debug...")
		log.SetLevel(logs.LevelDebug)
	}

	return log
}

func InitMQ() {
RETRY:

	//kafkas := net.JoinHostPort(KafkaAddrPort())
	ip, port := "10.1.235.98", "9092"
	kafkas := fmt.Sprintf("%s:%s", ip, port)
	log.Info("connectMQ, kafkas =", kafkas)

	messageQueue, err := mq.NewMQ([]string{kafkas}) // ex. {"192.168.1.1:9092", "192.168.1.2:9092"}
	if err != nil {
		log.Error("connectMQ error:", err.Error())
		time.Sleep(10 * time.Second)
		goto RETRY
	}

	q := &ds.MQ{MessageQueue: messageQueue}

	atomic.StorePointer(&theMQ, unsafe.Pointer(q))

	log.Info("MQ inited successfully.")
}

func sendAlarm(msg string) {
	// todo: need a buffered channel to store the alarms?
	q := getMQ()
	if q == nil {
		return
	}

	event := ds.AlarmEvent{Sender: SENDER, Content: msg, Send_time: time.Now()}

	b, err := json.Marshal(&event)
	if err != nil {
		log.Error("Marshal err:", err)
		return
	}

	_, _, err = q.MessageQueue.SendSyncMessage("to_alarm.json", []byte(""), b)
	if err != nil {
		log.Error("sendAlarm (to_alarm.json) error: ", err)
		return
	}
	return
}

func getMQ() *ds.MQ {
	return (*ds.MQ)(atomic.LoadPointer(&theMQ))
}

func KafkaAddrPort() (string, string) {
	switch Platform {
	case Platform_DaoCloud:
		entryList := dnsExchange(os.Getenv("kafka_service_name"))

		for _, v := range entryList {
			if v.Port == "9092" {
				return v.Ip, v.Port
			}
		}
	case Platform_DataOS:
		return os.Getenv(os.Getenv("ENV_NAME_KAFKA_ADDR")), os.Getenv(os.Getenv("ENV_NAME_KAFKA_PORT"))
	case Platform_DaoCloudUT:
		fallthrough
	case Platform_Local:
		return os.Getenv("MQ_KAFKA_ADDR"), os.Getenv("MQ_KAFKA_PORT")
	}

	return "", ""
}

func dnsExchange(srvName string) []*ds.DnsEntry {
	fiilSrvName := fmt.Sprintf("%s.service.consul", srvName)
	agentAddr := net.JoinHostPort(consulAddrPort())
	log.Debug("DNS query %s @ %s", fiilSrvName, agentAddr)

	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(fiilSrvName), dns.TypeSRV)
	m.RecursionDesired = true

	c := &dns.Client{Net: "tcp"}
	r, _, err := c.Exchange(m, agentAddr)
	if err != nil {
		log.Error("dns  exchange err:", err)
		return nil
	}
	if r.Rcode != dns.RcodeSuccess {
		log.Warn("dns query err:", r.Rcode)
		return nil
	}

	/*
		entries := make([]*dnsEntry, 0, 16)
		for i := len(r.Answer) - 1; i >= 0; i-- {
			answer := r.Answer[i]
			log.DefaultLogger().Debugf("r.Answer[%d]=%s", i, answer.String())

			srv, ok := answer.(*dns.SRV)
			if ok {
				m.SetQuestion(dns.Fqdn(srv.Target), dns.TypeA)
				r1, _, err := c.Exchange(m, agentAddr)
				if err != nil {
					log.DefaultLogger().Warningf("dns query error: %s", err.Error())
					continue
				}

				for j := len(r1.Answer) - 1; j >= 0; j-- {
					answer1 := r1.Answer[j]
					log.DefaultLogger().Debugf("r1.Answer[%d]=%v", i, answer1)

					a, ok := answer1.(*dns.A)
					if ok {
						a.A is the node ip instead of service ip
						entries = append(entries,  &dnsEntry{ip: a.A.String(), port: fmt.Sprintf("%d", srv.Port)})
					}
				}
			}
		}

		return entries
	*/

	if len(r.Extra) != len(r.Answer) {
		e := fmt.Sprintf("len(r.Extra)(%d) != len(r.Answer)(%d)", len(r.Extra), len(r.Answer))
		log.Warn(e)
		return nil
	}

	num := len(r.Extra)
	entries := make([]*ds.DnsEntry, num)
	index := 0
	for i := 0; i < num; i++ {
		a, oka := r.Extra[i].(*dns.A)
		s, oks := r.Answer[i].(*dns.SRV)
		if oka && oks {
			entries[index] = &ds.DnsEntry{Ip: a.A.String(), Port: fmt.Sprintf("%d", s.Port)}
			index++
		}
	}

	return entries[:index]
}

func consulAddrPort() (string, string) {
	return os.Getenv("CONSUL_SERVER"), os.Getenv("CONSUL_DNS_PORT")
}
