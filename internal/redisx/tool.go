package redisx

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"sort"
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
	rdb = Client()
	obj := rdb.DebugObject(ctx, queueName)

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

func DbSize(ctx context.Context) (int64, error) {
	rdb = Client()
	return rdb.DBSize(ctx).Result()
}

func ZCard(ctx context.Context, key string) int64 {
	rdb = Client()
	return rdb.ZCard(ctx, key).Val()
}

func HGetAll(ctx context.Context, key string) (map[string]string, error) {
	rdb = Client()
	return rdb.HGetAll(ctx, key).Result()
}

func HSet(ctx context.Context, key string, data map[string]any) error {
	rdb = Client()
	return rdb.HSet(ctx, key, data).Err()
}

func Del(ctx context.Context, key string) error {
	rdb = Client()
	return rdb.Del(ctx, key).Err()
}

func ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	rdb = Client()
	return rdb.ZScan(ctx, key, cursor, match, count).Result()
}

func XRevRange(ctx context.Context, stream, start, stop string) ([]redis.XMessage, error) {
	rdb = Client()
	return rdb.XRevRange(ctx, stream, start, stop).Result()
}

func ZRemRangeByScore(ctx context.Context, key, min, max string) error {
	rdb = Client()
	return rdb.ZRemRangeByScore(ctx, key, min, max).Err()
}

func XRangeN(ctx context.Context, stream string, start, stop string, count int64) ([]redis.XMessage, error) {
	rdb = Client()
	return rdb.XRangeN(ctx, stream, start, stop, count).Result()
}

func Monitor(ctx context.Context) string {
	rdb = Client()
	return rdb.Do(ctx, "MONITOR").String()
}

func Keys(ctx context.Context, key string) ([]string, error) {
	rdb = Client()
	cmd := rdb.Keys(ctx, key)
	queues, err := cmd.Result()
	if err != nil {
		return nil, err
	}
	return queues, nil
}
func Info(ctx context.Context) (map[string]string, error) {
	rdb = Client()
	infoStr, err := rdb.Info(ctx).Result()
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

func handleInfos(data string) map[string]any {
	if strings.Contains(data, "\r\n") {
		data = strings.ReplaceAll(data, "\r\n", "\n")
	}
	dataM := make(map[string]any, 0)
	arr := strings.Split(data, "\n")
	for _, m := range arr[1:] {
		if !strings.Contains(m, ":") {
			continue
		}
		s := strings.Split(m, ":")
		dataM[s[0]] = s[1]
	}
	return dataM
}

// Server
// redis command: info server
func Server(ctx context.Context) (map[string]any, error) {
	rdb = Client()
	server, err := rdb.Info(ctx, "Server").Result()
	if err != nil {
		return nil, err
	}
	return handleInfos(server), nil
}

// Clients
// redis command: info clients
func Clients(ctx context.Context) (map[string]any, error) {
	rdb = Client()
	clients, err := rdb.Info(ctx, "Clients").Result()
	if err != nil {
		return nil, err
	}
	m := handleInfos(clients)
	return m, nil
}

// Persistence
// redis command: info persistence
func Persistence(ctx context.Context) (map[string]any, error) {
	rdb = Client()
	persistence, err := rdb.Info(ctx, "Persistence").Result()
	if err != nil {
		return nil, err
	}
	return handleInfos(persistence), nil
}

// Memory
// redis command:info memory
func Memory(ctx context.Context) (map[string]any, error) {
	rdb = Client()
	memory, err := rdb.Info(ctx, "MEMORY").Result()
	if err != nil {
		return nil, err
	}
	m := handleInfos(memory)
	return m, nil
}

// Stats
// redis command:info stats
func Stats(ctx context.Context) (map[string]any, error) {
	rdb = Client()
	stats, err := rdb.Info(ctx, "Stats").Result()
	if err != nil {
		return nil, err
	}
	return handleInfos(stats), nil
}

// Keyspace
// redis command:info keyspace
func KeySpace(ctx context.Context) ([]map[string]any, error) {
	rdb = Client()
	spaces, err := rdb.Info(ctx, "Keyspace").Result()
	if err != nil {
		return nil, err
	}
	m := handleInfos(spaces)
	slic := make([]map[string]any, 0)
	for i, v := range m {
		nmap := make(map[string]any, 0)
		vv := strings.Split(cast.ToString(v), ",")
		for _, sv := range vv {
			lv := strings.Split(sv, "=")
			nmap[lv[0]] = lv[1]
		}
		nmap["dbname"] = i
		slic = append(slic, nmap)
	}
	return slic, nil
}

// Commands
// sort in reverse order based on the `usec_per_call` field
type Commands []map[string]any

func (t Commands) Len() int {
	return len(t)
}

func (t Commands) Less(i, j int) bool {
	return cast.ToFloat64(t[j]["usec_per_call"]) < cast.ToFloat64(t[i]["usec_per_call"])
}

func (t Commands) Swap(i, j int) {
	t[j], t[i] = t[i], t[j]
}

// CommandStats
// redis command: info Commandstats
func CommandStats(ctx context.Context) ([]map[string]any, error) {
	rdb = Client()
	command, err := rdb.Info(ctx, "Commandstats").Result()
	if err != nil {
		return nil, err
	}

	if strings.Contains(command, "\r\n") {
		command = strings.ReplaceAll(command, "\r\n", "\n")
	}
	commands := strings.Split(command, "\n")

	var commandMap Commands
	for _, m := range commands[1:] {
		if !strings.Contains(m, ":") {
			continue
		}
		s := strings.Split(m, ":")
		key := strings.ReplaceAll(s[0], "cmdstat_", "")
		val := s[1]
		vals := strings.Split(val, ",")
		nmap := make(map[string]any, 0)
		nmap["command"] = key
		for _, v := range vals {
			vv := strings.Split(v, "=")
			nmap[vv[0]] = vv[1]
		}
		commandMap = append(commandMap, nmap)
	}
	sort.Sort(commandMap)
	return commandMap, nil
}

func ClientList(ctx context.Context) ([]map[string]any, error) {
	rdb = Client()
	cmd := rdb.ClientList(ctx)
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
	rdb = Client()
	cmd := rdb.ZRange(ctx, match, page, pageSize)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	result, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	njson := json.Json

	length, err := rdb.ZLexCount(ctx, match, "-", "+").Result()
	if err != nil {
		return nil, err
	}
	d := make([]map[string]any, 0, pageSize)
	for _, v := range result {

		cmd := rdb.ZRank(ctx, match, v)
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

		ttl := time.Until(cast.ToTime(njson.Get(payloadByte, "ExpireTime").ToString())).Seconds()
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

	rdb = Client()
	// get queues
	cmd := rdb.Keys(ctx, QueueKey(BqConfig.Redis.Prefix))
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
