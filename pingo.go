package pingo

import (
	"log"
	"net"
	"strconv"
	"time"

	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/cloudwatch"
)

const (
	CloudWatchOkValue   = float64(2)
	CloudWatchFailValue = float64(1)
	CloudWatchNamespace = "pingo"
)

type Host struct {
	Hostname string
	Port     int
}

type Handler func(host Host, status bool, t time.Time) error

func (host Host) Address() string {
	return host.Hostname + ":" + strconv.Itoa(host.Port)
}

func LoggingHandler(host Host, status bool, t time.Time) error {
	text := "FAIL"
	if status == true {
		text = "OK"
	}
	log.Printf("PING %s %s", host.Address(), text)
	return nil
}

func NewCloudWatchHandler(region string) Handler {
	auth, err := aws.EnvAuth()
	if err != nil {
		log.Fatalln(err)
	}
	c, err := cloudwatch.NewCloudWatch(auth, aws.Regions[region].CloudWatchServicepoint)
	if err != nil {
		log.Fatalln(err)
	}
	return func(host Host, status bool, t time.Time) error {
		value := CloudWatchFailValue
		if status == true {
			value = CloudWatchOkValue
		}
		metric := cloudwatch.MetricDatum{
			MetricName: host.Address(),
			Value:      value,
		}
		_, err := c.PutMetricDataNamespace([]cloudwatch.MetricDatum{metric}, CloudWatchNamespace)
		return err
	}
}

func Ping(host Host, timeout time.Duration, handlers []Handler) {
	address := host.Address()
	now := time.Now()
	conn, err := net.DialTimeout("tcp", address, timeout)
	ok := err == nil
	if conn != nil {
		conn.Close()
	}
	for _, handler := range handlers {
		err := handler(host, ok, now)
		if err != nil {
			log.Println(err)
		}
	}
}
