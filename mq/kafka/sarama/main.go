// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.
// +build linux
// +build amd64 arm64
// +build cgo

package saramaframe

import (
	"fmt"
	"reflect"
	"runtime"
	"sync"
	"time"
	"unsafe"

	"github.com/TingYunGo/goagent"
	"github.com/TingYunGo/goagent/libs/tystring"
	"github.com/TingYunGo/goagent/runtime"
	"github.com/IBM/sarama"
)

const (
	SaramaProducerIndexContext = tingyun3.StorageMQKafka + 0
	SaramaConsumerIndexContext = tingyun3.StorageMQKafka + 1
	SaramaConsumerIndexStore   = tingyun3.StorageMQKafka + 2
)

type asyncProducer struct {
	client    sarama.Client
	conf      *sarama.Config
	ownClient bool
}
type syncProducer struct {
	producer *asyncProducer
	wg       sync.WaitGroup
}

func appendHeader(component *tingyun3.Component, msg *sarama.ProducerMessage) {
	if msg == nil || component == nil {
		return
	}
	if !tingyun3.ReadServerConfigBool(tingyun3.ServerConfigBoolMQEnabled, false) {
		return
	}
	if trackID := component.CreateTrackID(); len(trackID) > 0 {
		msg.Headers = append(msg.Headers, sarama.RecordHeader{
			Key:   []byte("X-Tingyun"),
			Value: []byte(trackID),
		})
	}
}

//go:noinline
func syncProducerSendMessage(sp *syncProducer, msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	trampoline.arg1 = *trampoline.idpointer + trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return 0, 0, nil
}

//go:noinline
func WrapsyncProducerSendMessage(sp *syncProducer, msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	var producer *tingyun3.Component = nil
	action := tingyun3.GetAction()
	if action != nil {
		brokers := ""
		for _, broker := range sp.producer.client.Brokers() {
			if len(brokers) > 0 {
				brokers = brokers + "," + broker.Addr()
			} else {
				brokers = broker.Addr()
			}
		}
		if producer = action.CreateMQComponent("Kafka", false, brokers, msg.Topic); producer != nil {
			method := tingyun3.GetCallerName(2)
			producer.SetMethod(method)
			appendHeader(producer, msg)
		}
	}
	begin := time.Now()
	partition, offset, err = syncProducerSendMessage(sp, msg)
	if action != nil && producer != nil {

		producer.FixBegin(begin)
		if err != nil {
			producer.SetException(err, "syncProducer.SendMessage", 2)
		}
		producer.End(1)
	}
	return
}

//go:noinline
func syncProducerSendMessages(sp *syncProducer, msgs []*sarama.ProducerMessage) error {
	trampoline.arg1 = *trampoline.idpointer + trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapsyncProducerSendMessages(sp *syncProducer, msgs []*sarama.ProducerMessage) (err error) {
	var producer *tingyun3.Component = nil
	action := tingyun3.GetAction()
	if action != nil {
		brokers := ""
		for _, broker := range sp.producer.client.Brokers() {
			if len(brokers) > 0 {
				brokers = brokers + "," + broker.Addr()
			} else {
				brokers = broker.Addr()
			}
		}
		if len(msgs) > 0 {
			if producer = action.CreateMQComponent("Kafka", false, brokers, msgs[0].Topic); producer != nil {
				method := tingyun3.GetCallerName(2)
				producer.SetMethod(method)
				appendHeader(producer, msgs[0])
			}
		}
	}
	begin := time.Now()
	err = syncProducerSendMessages(sp, msgs)
	if action != nil && producer != nil {

		producer.FixBegin(begin)
		if err != nil {
			producer.SetException(err, "syncProducer.SendMessages", 2)
		}
		producer.End(1)
	}
	return
}

//go:noinline
func asyncProducerInput(p *asyncProducer) chan<- *sarama.ProducerMessage {
	trampoline.arg1 = *trampoline.idpointer + trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapasyncProducerInput(p *asyncProducer) chan<- *sarama.ProducerMessage {
	//github.com/IBM/sarama.(*asyncProducer).Input
	action := tingyun3.GetAction()
	if action != nil {
		brokers := ""
		for _, broker := range p.client.Brokers() {
			if len(brokers) > 0 {
				brokers = brokers + "," + broker.Addr()
			} else {
				brokers = broker.Addr()
			}
		}
		if len(brokers) > 0 {
			tingyun3.LocalSet(SaramaProducerIndexContext, brokers)
			action.OnEnd(func() {
				tingyun3.LocalDelete(SaramaProducerIndexContext)
			})
		}
	}
	return asyncProducerInput(p)
}

func chanSendHandler(channelMsgType string, ep unsafe.Pointer, callerpc uintptr) bool {

	if channelMsgType == "*sarama.ProducerMessage" {

		if action := tingyun3.GetAction(); action != nil {

			callerName := runtime.FuncForPC(callerpc).Name()

			if callerName == "github.com/IBM/sarama.(*syncProducer).SendMessage" ||
				callerName == "github.com/IBM/sarama.(*asyncProducer).finishTransaction" {
				return true
			}

			brokers := ""
			topic := ""
			var msg *sarama.ProducerMessage = nil

			if c := tingyun3.LocalGet(SaramaProducerIndexContext); c != nil {
				brokers = c.(string)
			}
			if ep != nil {
				type pointerTopointer struct {
					p unsafe.Pointer
				}
				if pp := (*pointerTopointer)(ep); pp.p != nil {
					msg = (*sarama.ProducerMessage)(pp.p)
					topic = msg.Topic
				}
			}

			producer := action.CreateMQComponent("Kafka", false, brokers, topic)
			producer.SetMethod(callerName)
			if producer != nil && msg != nil {
				appendHeader(producer, msg)
			}
			producer.End(1)
		}
		return true
	}
	return false
}

func getmethodNameByAddr(p uintptr) string {
	if r := runtime.FuncForPC(p); r == nil {
		return ""
	} else {
		return r.Name()
	}
}

type consumerLocal struct {
	from       string
	callerName string
	action     *tingyun3.Action
}

func onConsumerActionStart(inst string, callerName string, message *sarama.ConsumerMessage) {
	if action := tingyun3.GetAction(); action != nil {
		return
	}
	brokers := ""
	if c := tingyun3.LocalGet(SaramaConsumerIndexStore); c != nil {
		brokers = c.(string)
	}

	instname := fmt.Sprintf("%s:%s:%d", inst, message.Topic, message.Partition)
	if action, _ := tingyun3.CreateAction(instname, callerName); action != nil {
		action.SetConsumer("Kafka", brokers, message.Topic)
		action.SetName("CLIENTIP", brokers)
		action.SetURL(fmt.Sprintf("kafka://%s/%s/%d", brokers, message.Topic, message.Partition))
		for _, item := range message.Headers {
			if string(item.Key[:]) == "X-Tingyun" {
				action.SetTrackID(string(item.Value[:]))
				action.SetBackEnabled(false)
				break
			}
		}

		tingyun3.LocalSet(SaramaConsumerIndexContext, &consumerLocal{
			from:       inst,
			callerName: callerName,
			action:     action,
		})
		tingyun3.SetAction(action)
	}
}
func onConsumerActionEnd(inst, callerName string) {

	if c := tingyun3.LocalGet(SaramaConsumerIndexContext); c != nil {
		local := c.(*consumerLocal)
		if local.from != inst || local.callerName != callerName {
			return
		}
		local.action.Finish()
		tingyun3.LocalDelete(SaramaConsumerIndexContext)
		consumerBroker := tingyun3.LocalGet(SaramaConsumerIndexStore)
		tingyun3.LocalClear()
		if consumerBroker != nil {
			tingyun3.LocalSet(SaramaConsumerIndexStore, consumerBroker)
		}
	}
}

func chanRecvHandler(channelMsgType string, ep unsafe.Pointer, block bool) wrapruntime.ChanRecvHandler {
	if channelMsgType == "*sarama.ConsumerMessage" {
		if callerName := getCallerName(); callerName != "github.com/IBM/sarama.(*partitionConsumer).responseFeeder" {
			onConsumerActionEnd("chan", callerName)
			return &chanrecvHandler{
				ep: ep,
			}
		}
	}
	return nil
}

type chanrecvHandler struct {
	ep unsafe.Pointer
}

func (ch *chanrecvHandler) Ret(selected, received bool) {

	if !received {
		return
	}
	if !tingyun3.ReadServerConfigBool(tingyun3.ServerConfigBoolMQEnabled, false) {
		return
	}
	type pointerTopointer struct {
		p unsafe.Pointer
	}
	if pp := (*pointerTopointer)(ch.ep); pp.p != nil {
		msg := (*sarama.ConsumerMessage)(pp.p)
		callerName := getCallerName()
		onConsumerActionStart("chan", callerName, msg)
	}
}

func matchMethod(method, matcher string) bool {
	return tystring.SubString(method, 0, len(matcher)) == matcher
}

//go:noinline
func getCallerName() string {

	localList := [6]uintptr{0, 0, 0, 0}
	stackList := localList[:]
	count := runtime.Callers(3, stackList)
	for i := 0; i < count; i++ {
		method := getmethodNameByAddr(stackList[i] - 1)
		if matchMethod(method, "runtime.chanrecv") {
			continue
		}
		if matchMethod(method, "github.com/TingYunGo/goagent") {
			continue
		}
		return method
	}
	return ""
}
func selectHandler(channelMsgType string) wrapruntime.SelectHandler {
	if channelMsgType == "*sarama.ConsumerMessage" {
		callerName := getCallerName()
		if callerName == "github.com/IBM/sarama.(*partitionConsumer).responseFeeder" {
			return nil
		}
		onConsumerActionEnd("select", callerName)
		return &selectgoHandler{}
	}
	return nil
}

type selectgoHandler struct {
}

func (sh *selectgoHandler) Ret(tyneName string, elem unsafe.Pointer, retId int, retOk bool) {

	if tyneName != "*sarama.ConsumerMessage" {
		return
	}
	if !tingyun3.ReadServerConfigBool(tingyun3.ServerConfigBoolMQEnabled, false) {
		return
	}

	type pointerTopointer struct {
		p unsafe.Pointer
	}
	if pp := (*pointerTopointer)(elem); pp.p != nil {
		msg := (*sarama.ConsumerMessage)(pp.p)

		callerName := getCallerName()
		onConsumerActionStart("select", callerName, msg)
	}
}

type brokerConsumer struct {
	consumer unsafe.Pointer
	broker   *sarama.Broker
}

type none struct{}

type partitionConsumer struct {
	highWaterMarkOffset int64
	consumer            unsafe.Pointer
	conf                unsafe.Pointer
	broker              *brokerConsumer
	messages            chan *sarama.ConsumerMessage
	errors              chan *sarama.ConsumerError
	feeder              chan *sarama.FetchResponse

	leaderEpoch          int32
	preferredReadReplica int32

	trigger, dying chan none
	closeOnce      sync.Once
	topic          string
	partition      int32
	responseResult error
	fetchSize      int32
	offset         int64
	retries        int32

	paused int32
}

//go:noinline
func partitionConsumerMessages(child *partitionConsumer) <-chan *sarama.ConsumerMessage {
	trampoline.arg1 = *trampoline.idpointer + trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrappartitionConsumerMessages(child *partitionConsumer) <-chan *sarama.ConsumerMessage {
	brokers := child.broker.broker.Addr()
	if len(brokers) > 0 {
		tingyun3.LocalSet(SaramaConsumerIndexStore, brokers)
	}

	return partitionConsumerMessages(child)
}

func init() {
	wrapruntime.HandleChanSend(chanSendHandler)
	wrapruntime.HandleChanRecv(chanRecvHandler)
	wrapruntime.HandleSelect(selectHandler)
	tingyun3.Register(reflect.ValueOf(WrapsyncProducerSendMessage).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapsyncProducerSendMessages).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapasyncProducerInput).Pointer())
	tingyun3.Register(reflect.ValueOf(WrappartitionConsumerMessages).Pointer())
}
