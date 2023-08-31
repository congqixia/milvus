package main

import (
	"context"
	"crypto/rand"
	"errors"
	"flag"
	"log"
	"sync"
	"time"

	"github.com/milvus-io/milvus-proto/go-api/v2/commonpb"
	"github.com/milvus-io/milvus/pkg/config"
	"github.com/milvus-io/milvus/pkg/mq/msgstream"
	"github.com/milvus-io/milvus/pkg/util/paramtable"
)

type BenchSuite struct {
	ms msgstream.MsgStream
}

func NewBenchSuite() *BenchSuite {
	return &BenchSuite{}
}

type benchMsg struct {
	data []byte
}

func (m *benchMsg) TraceCtx() context.Context {
	return context.Background()
}

func (m *benchMsg) SetTraceCtx(ctx context.Context) {
}

func (m *benchMsg) ID() msgstream.UniqueID {
	return 0
}

func (m *benchMsg) SetID(id msgstream.UniqueID) {
}

func (m *benchMsg) BeginTs() msgstream.Timestamp {
	return 0
}

func (m *benchMsg) EndTs() msgstream.Timestamp {
	return 0
}

func (m *benchMsg) Type() msgstream.MsgType {
	return commonpb.MsgType_Undefined
}

func (m *benchMsg) SourceID() int64 {
	return 0
}

func (m *benchMsg) HashKeys() []uint32 {
	return []uint32{0}
}

func (m *benchMsg) Marshal(in msgstream.TsMsg) (msgstream.MarshalType, error) {
	return m.data, nil

}

func (m *benchMsg) Unmarshal(_ msgstream.MarshalType) (msgstream.TsMsg, error) {
	return nil, errors.New("not supported")
}

func (m *benchMsg) Position() *msgstream.MsgPosition {
	return &msgstream.MsgPosition{}
}

func (m *benchMsg) SetPosition(_ *msgstream.MsgPosition) {
}

func (m *benchMsg) Size() int {
	return len(m.data)
}

func (s *BenchSuite) Run() {
	var wg, initWg sync.WaitGroup
	data := make([]byte, *batchSize)
	rand.Read(data)

	var cms msgstream.MsgStream
	factory := &msgstream.PmsFactory{
		MQBufSize:              128,
		ReceiveBufSize:         128,
		PulsarAddress:          *address,
		PulsarTenant:           *tenant,
		PulsarNameSpace:        *namespace,
		PulsarAuthParams:       "{}",
		MaxConnectionPerBroker: *maxConnectPerBroker,
	}

	if *shareMsgstream {
		ctx := context.Background()

		var err error
		cms, err = factory.NewMsgStream(ctx)
		if err != nil {
			log.Fatal("failed to create msgstream: ", err.Error())
		}
		defer cms.Close()

		cms.AsProducer([]string{*topic}, msgstream.WithAsyncProduce(*sendAsync))
	}

	start := time.Now()
	signal := make(chan struct{})
	wg.Add(*clientNum)
	initWg.Add(*clientNum)
	for i := 0; i < *clientNum; i++ {
		go func() {
			defer wg.Done()
			var ms msgstream.MsgStream
			if !*shareMsgstream {
				var err error
				ctx := context.Background()
				ms, err = factory.NewMsgStream(ctx)
				if err != nil {
					log.Fatal("failed to create msgstream: ", err.Error())
				}
				defer ms.Close()

				ms.AsProducer([]string{*topic}, msgstream.WithAsyncProduce(*sendAsync))
			} else {
				ms = cms
			}
			initWg.Done()

			<-signal

			for i := 0; i < *batchNum; i++ {
				pack := &msgstream.MsgPack{
					Msgs: []msgstream.TsMsg{
						&benchMsg{
							data: data,
						},
					},
				}
				ms.Produce(pack)
			}
		}()
	}

	initWg.Wait()
	close(signal)

	wg.Wait()
	log.Println("bench produce done, elapsed", time.Since(start))
}

var (
	// pulsar info
	address             = flag.String("pulsarAddress", "pulsar://127.0.0.1:6650", "pulsar address")
	namespace           = flag.String("pulsarNamespace", "default", "pulsar namespace")
	tenant              = flag.String("pulsarTenant", "public", "pulsar tenant")
	maxConnectPerBroker = flag.Int("pulsarMaxConnPerBroker", 1, "max connection per broker")

	topic          = flag.String("topic", "bench_mq", "mq topic bench program shall use")
	clientNum      = flag.Int("clientNum", 10, "client number to produce message")
	batchNum       = flag.Int("batchNum", 100, "total msg batch each client shall produce")
	batchSize      = flag.Int("batchSize", 1024*1024, "each batch size, default 1M")
	shareMsgstream = flag.Bool("share", false, "all client share same msgstream instance")
	sendAsync      = flag.Bool("sendAsync", false, "use pulsar send async feature")
)

func initParamItem(item *paramtable.ParamItem, v string) {
	item.Formatter = func(originValue string) string { return v }
	item.Init(&config.Manager{})
}

func main() {
	flag.Parse()

	s := NewBenchSuite()

	s.Run()
}
