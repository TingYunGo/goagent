// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.
// +build linux
// +build amd64 arm64
// +build cgo

package kafkagoframe

import (
	"context"
	"fmt"
	"reflect"
	"runtime"

	"github.com/TingYunGo/goagent"
	"github.com/TingYunGo/goagent/libs/tystring"
	kafka "github.com/segmentio/kafka-go"
)

const (
	KafkaGoProducerIndexContext = tingyun3.StorageMQKafka + 4
)

var skipTokens = []string{
	"github.com/segmentio/kafka-go",
	"github.com/TingYunGo/goagent",
}

func readBrokers(r *kafka.Reader) string {
	brokers := ""
	for _, broker := range r.Config().Brokers {
		if len(brokers) > 0 {
			brokers = brokers + "," + broker
		} else {
			brokers = broker
		}
	}
	return brokers
}
func getInstance(message *kafka.Message) string {
	return fmt.Sprintf("%s:%d", message.Topic, message.Partition)
}
func getTraceID(message *kafka.Message) string {
	for _, head := range message.Headers {
		if head.Key == "X-Tingyun" {
			return string(head.Value[:])
		}
	}
	return ""
}
func getmethodNameByAddr(p uintptr) string {
	if r := runtime.FuncForPC(p); r == nil {
		return ""
	} else {
		return r.Name()
	}
}
func matchMethod(method, matcher string) bool {
	return tystring.SubString(method, 0, len(matcher)) == matcher
}

//go:noinline
func getCallerName() string {

	localList := [4]uintptr{0, 0, 0, 0}
	stackList := localList[:]
	count := runtime.Callers(2, stackList)
	for i := 0; i < count; i++ {
		method := getmethodNameByAddr(stackList[i] - 1)

		if matchMethod(method, skipTokens[0]) {
			continue
		}
		if matchMethod(method, skipTokens[1]) {
			continue
		}
		return method
	}
	return ""
}

func appendHeader(component *tingyun3.Component, msg *kafka.Message) {
	if msg == nil || component == nil {
		return
	}
	if !tingyun3.ReadServerConfigBool(tingyun3.ServerConfigBoolMQEnabled, false) {
		return
	}
	if trackID := component.CreateTrackID(); len(trackID) > 0 {
		msg.Headers = append(msg.Headers, kafka.Header{
			Key:   "X-Tingyun",
			Value: []byte(trackID),
		})
	}
}

//go:noinline
func writerWriteMessages(w *kafka.Writer, ctx context.Context, msgs ...kafka.Message) error {
	trampoline.arg1 = *trampoline.idpointer + trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapwriterWriteMessages(w *kafka.Writer, ctx context.Context, msgs ...kafka.Message) error {

	var component *tingyun3.Component = nil
	if len(msgs) > 0 {

		if action := tingyun3.GetAction(); action != nil {
			topic := msgs[0].Topic
			broker := w.Addr.String()

			if component = action.CreateMQComponent("Kafka", false, broker, topic); component != nil {
				component.SetMethod(getCallerName())
				appendHeader(component, &msgs[0])
			}
		}
	}

	err := writerWriteMessages(w, ctx, msgs...)

	if component != nil {
		if err != nil {
			component.SetException(err, "WriterWriteMessages", 2)
		}
		component.End(1)
	}
	return err
}

//go:noinline
func ReaderFetchMessage(r *kafka.Reader, ctx context.Context) (kafka.Message, error) {
	trampoline.arg1 = *trampoline.idpointer + trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return kafka.Message{}, nil
}

//go:noinline
func WrapReaderFetchMessage(r *kafka.Reader, ctx context.Context) (kafka.Message, error) {
	onConsumerActionEnd(false)
	m, err := ReaderFetchMessage(r, ctx)

	if err == nil && tingyun3.ReadServerConfigBool(tingyun3.ServerConfigBoolMQEnabled, false) {
		onConsumerActionStart(r, &m)
	}
	return m, err

}

type consumerLocal struct {
	callerName string
	action     *tingyun3.Action
}

func onConsumerActionStart(r *kafka.Reader, message *kafka.Message) {
	if action := tingyun3.GetAction(); action != nil {
		return
	}
	callerName := getCallerName()
	instname := getInstance(message)

	if action, _ := tingyun3.CreateAction(instname, callerName); action != nil {

		brokers := readBrokers(r)

		action.SetConsumer("Kafka", brokers, message.Topic)
		action.SetName("CLIENTIP", brokers)
		action.SetURL(fmt.Sprintf("kafka://%s/%s/%d", brokers, message.Topic, message.Partition))

		if trackID := getTraceID(message); len(trackID) > 0 {
			action.SetTrackID(trackID)
			action.SetBackEnabled(false)
		}

		tingyun3.LocalSet(KafkaGoProducerIndexContext, &consumerLocal{
			callerName: callerName,
			action:     action,
		})
		tingyun3.SetAction(action)
	}
}
func onConsumerActionEnd(lastEnd bool) {

	if c := tingyun3.LocalGet(KafkaGoProducerIndexContext); c != nil {

		local := c.(*consumerLocal)
		if callerName := getCallerName(); local.callerName == callerName || lastEnd {

			local.action.Finish()
			tingyun3.LocalDelete(KafkaGoProducerIndexContext)
			tingyun3.LocalClear()
		}
	}
}

func init() {

	tingyun3.Register(reflect.ValueOf(WrapReaderFetchMessage).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapwriterWriteMessages).Pointer())
}
