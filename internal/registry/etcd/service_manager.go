package etcd

import (
	"context"
	"fmt"

	"github.com/milvus-io/milvus/internal/registry"
	"github.com/milvus-io/milvus/pkg/log"
	"github.com/milvus-io/milvus/pkg/util/merr"
	"github.com/milvus-io/milvus/pkg/util/typeutil"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type serviceManager struct {
	sessions *typeutil.ConcurrentMap[string, *etcdSession]
	client   *clientv3.Client
}

func newServiceManager(cli *clientv3.Client) *serviceManager {
	return &serviceManager{
		sessions: typeutil.NewConcurrentMap[string, *etcdSession](),
	}
}

func (m *serviceManager) Add(ctx context.Context, s *etcdSession) error {
	if err := m.startKeepalive(s); err != nil {
		return merr.WrapErrIoFailed(s.key(), "failed to keepalive", err.Error())
	}
	m.sessions.Insert(s.key(), s)
	return nil
}

func (m *serviceManager) Remove(ctx context.Context, entry registry.ServiceEntry) error {
	session, ok := m.sessions.GetAndRemove(fmt.Sprintf("%s-%d", entry.Addr(), entry.ID()))
	if !ok {
		// ignore due to session not registered
		return nil
	}
	session.cancel()
	// wait keepalive quit
	<-session.closeCh
	_, err := m.client.Revoke(ctx, session.leaseID)
	if err != nil {
		log.Warn("failed to revoke session", zap.Error(err))
	}
	return err
}

func (m *serviceManager) startKeepalive(session *etcdSession) error {
	ctx, cancel := context.WithCancel(context.Background())
	session.cancel = cancel
	ch, err := m.client.KeepAlive(ctx, session.leaseID)
	if err != nil {
		return err
	}

	go m.keeapalive(ctx, ch, session)
	return nil
}

func (m *serviceManager) keeapalive(ctx context.Context, ch <-chan *clientv3.LeaseKeepAliveResponse, session *etcdSession) {
	defer func() {
		close(session.closeCh)
	}()
	for {
		select {
		case <-ctx.Done():
			log.Info("etcd session keepalive due to cancel")
			return
		case resp, ok := <-ch:
			if !ok {
				log.Warn("session keepalive channel closed")
				// TODO execute callback
				return
			}
			if resp == nil {
				log.Warn("session keepalive response failed")
				return
			}
		}
	}
}
