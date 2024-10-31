package redisx

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/retail-ai-inc/beanq/v3/helper/json"
	"github.com/retail-ai-inc/beanq/v3/helper/stringx"
	"github.com/spf13/cast"
)

type ObjectStruct struct {
	ValueAt          string
	Encoding         string
	RefCount         int
	SerizlizedLength int
	Lru              int
	LruSecondsIdle   int
}

func Object(ctx context.Context, queueName string) (objstr ObjectStruct) {
	client = Client()
	obj := client.DebugObject(ctx, queueName)

	str, _ := obj.Result()
	// Value at:0x7fc38fe77cc0 refcount:1 encoding:stream serializedlength:12 lru:7878503 lru_seconds_idle:3
	valueAt := "Value at"
	if strings.HasPrefix(str, valueAt) {
		str = strings.ReplaceAll(str, valueAt, "value_at")
	}

	strs := strings.Split(str, " ")

	for _, s := range strs {
		sarr := strings.Split(s, ":")
		if len(sarr) >= 2 {
			switch sarr[0] {
			case "value_at":
				objstr.ValueAt = sarr[1]
			case "refcount":
				objstr.RefCount = cast.ToInt(sarr[1])
			case "encoding":
				objstr.Encoding = sarr[1]
			case "serializedlength":
				objstr.SerizlizedLength = cast.ToInt(sarr[1])
			case "lru":
				objstr.Lru = cast.ToInt(sarr[1])
			case "lru_seconds_idle":
				objstr.LruSecondsIdle = cast.ToInt(sarr[1])
			}
		}
	}
	return
}
func Keys(ctx context.Context, key string) ([]string, error) {
	client = Client()
	cmd := client.Keys(ctx, key)
	fmt.Println(cmd.String())
	queues, err := cmd.Result()
	if err != nil {
		return nil, err
	}
	return queues, nil
}
func Info(ctx context.Context) (map[string]string, error) {
	client = Client()
	infoStr, err := client.Info(ctx).Result()
	if err != nil {
		return nil, err
	}
	info := make(map[string]string)
	lines := strings.Split(infoStr, "\r\n")
	for _, l := range lines {
		kv := strings.Split(l, ":")
		if len(kv) == 2 {
			info[kv[0]] = kv[1]
		}
	}
	return info, nil

}

func ClientList(ctx context.Context) ([]map[string]any, error) {
	client = Client()
	cmd := client.ClientList(ctx)
	if err := cmd.Err(); err != nil {
		return nil, err
	}
	data, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	arr := strings.Split(data, "\n")
	ldata := make(map[string]any, 0)
	rdata := make([]map[string]any, 0, 10)
	for _, v := range arr {
		nv := strings.Split(v, " ")
		for _, nvv := range nv {
			vals := strings.Split(nvv, "=")
			if vals[0] == "age" {
				if vals[1] == "0" {
					continue
				}
			}
			if len(vals) < 2 {
				continue
			}
			ldata[vals[0]] = vals[1]
			rdata = append(rdata, ldata)
		}
	}
	return rdata, nil
}
func ZRange(ctx context.Context, match string, page, pageSize int64) (map[string]any, error) {
	client = Client()
	cmd := client.ZRange(ctx, match, page, pageSize)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	result, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	njson := json.Json

	length, err := client.ZLexCount(ctx, match, "-", "+").Result()
	if err != nil {
		return nil, err
	}
	d := make([]map[string]any, 0, pageSize)
	for _, v := range result {

		cmd := client.ZRank(ctx, match, v)
		key, err := cmd.Result()
		if err != nil {
			continue
		}
		payloadByte := stringx.StringToByte(v)
		npayload := njson.Get(payloadByte, "Payload")
		addTime := njson.Get(payloadByte, "AddTime")
		runTime := njson.Get(payloadByte, "RunTime")
		group := njson.Get(payloadByte, "Group")

		queuestr := njson.Get(payloadByte, "Queue").ToString()
		queues := strings.Split(queuestr, ":")
		queue := queuestr
		if len(queues) >= 4 {
			queue = queues[2]
		}

		ttl := cast.ToTime(njson.Get(payloadByte, "ExpireTime").ToString()).Sub(time.Now()).Seconds()
		d = append(d, map[string]any{"key": key, "ttl": fmt.Sprintf("%.3f", ttl), "addTime": addTime, "runTime": runTime, "group": group, "queue": queue, "payload": npayload})

	}
	return map[string]any{"data": d, "total": length}, nil
}

type Msg struct {
	Id      string `json:"id"`
	Level   string
	Info    string
	Payload any `json:"payload"`

	AddTime     string    `json:"addTime"`
	ExpireTime  time.Time `json:"expireTime"`
	RunTime     string    `json:"runTime"`
	BeginTime   time.Time
	EndTime     time.Time
	ExecuteTime time.Time
	Topic       string `json:"topic"`
	Channel     string `json:"channel"`
	Consumer    string `json:"consumer"`
	Score       string
}

type Stream struct {
	Prefix   string `json:"prefix"`
	Channel  string `json:"channel"`
	Topic    string `json:"topic"`
	MoodType string `json:"moodType"`
	State    string `json:"state"`
	Size     int    `json:"size"`
	Idle     int    `json:"idle"`
}

func QueueInfo(ctx context.Context) (any, error) {

	client = Client()
	// get queues
	cmd := client.Keys(ctx, QueueKey(BqConfig.Redis.Prefix))
	queues, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	data := make(map[string][]Stream, 0)
	for _, queue := range queues {

		arr := strings.Split(queue, ":")
		if len(arr) < 4 {
			continue
		}

		obj := Object(ctx, queue)

		stream := Stream{
			Prefix:   arr[0],
			Channel:  arr[1],
			Topic:    arr[2],
			MoodType: arr[3],
			State:    "Run",
			Size:     obj.SerizlizedLength,
			Idle:     obj.LruSecondsIdle,
		}
		data[arr[1]] = append(data[arr[1]], stream)
	}

	return data, nil
}
func ScheduleQueueKey(prefix string) string {
	return strings.Join([]string{prefix, "*", "zset"}, ":")
}
func QueueKey(prefix string) string {
	return strings.Join([]string{prefix, "*", "stream"}, ":")
}
