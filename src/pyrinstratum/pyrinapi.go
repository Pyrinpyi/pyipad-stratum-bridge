package pyrinstratum

import (
	"context"
	"fmt"
	"time"

	"github.com/Pyrinpyi/pyipad/app/appmessage"
	"github.com/Pyrinpyi/pyipad/infrastructure/network/rpcclient"
	"github.com/Pyrinpyi/pyrin-stratum-bride/src/gostratum"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type PyrinApi struct {
	address       string
	blockWaitTime time.Duration
	logger        *zap.SugaredLogger
	pyrin         *rpcclient.RPCClient
	connected     bool
}

func NewPyrinAPI(address string, blockWaitTime time.Duration, logger *zap.SugaredLogger) (*PyrinApi, error) {
	client, err := rpcclient.NewRPCClient(address)
	if err != nil {
		return nil, err
	}

	return &PyrinApi{
		address:       address,
		blockWaitTime: blockWaitTime,
		logger:        logger.With(zap.String("component", "pyrinapi:"+address)),
		pyrin:         client,
		connected:     true,
	}, nil
}

func (py *PyrinApi) Start(ctx context.Context, blockCb func()) {
	py.waitForSync(true)
	go py.startBlockTemplateListener(ctx, blockCb)
	go py.startStatsThread(ctx)
}

func (py *PyrinApi) startStatsThread(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-ctx.Done():
			py.logger.Warn("context cancelled, stopping stats thread")
			return
		case <-ticker.C:
			dagResponse, err := py.pyrin.GetBlockDAGInfo()
			if err != nil {
				py.logger.Warn("failed to get network hashrate from pyrin, prom stats will be out of date", zap.Error(err))
				continue
			}
			response, err := py.pyrin.EstimateNetworkHashesPerSecond(dagResponse.TipHashes[0], 1000)
			if err != nil {
				py.logger.Warn("failed to get network hashrate from pyrin, prom stats will be out of date", zap.Error(err))
				continue
			}
			RecordNetworkStats(response.NetworkHashesPerSecond, dagResponse.BlockCount, dagResponse.Difficulty)
		}
	}
}

func (py *PyrinApi) reconnect() error {
	if py.pyrin != nil {
		return py.pyrin.Reconnect()
	}

	client, err := rpcclient.NewRPCClient(py.address)
	if err != nil {
		return err
	}
	py.pyrin = client
	return nil
}

func (s *PyrinApi) waitForSync(verbose bool) error {
	if verbose {
		s.logger.Info("checking pyrin sync state")
	}
	for {
		clientInfo, err := s.pyrin.GetInfo()
		if err != nil {
			return errors.Wrapf(err, "error fetching server info from pyrin @ %s", s.address)
		}
		if clientInfo.IsSynced {
			break
		}
		s.logger.Warn("Pyrin is not synced, waiting for sync before starting bridge")
		time.Sleep(5 * time.Second)
	}
	if verbose {
		s.logger.Info("pyrin synced, starting server")
	}
	return nil
}

func (s *PyrinApi) startBlockTemplateListener(ctx context.Context, blockReadyCb func()) {
	blockReadyChan := make(chan bool)
	err := s.pyrin.RegisterForNewBlockTemplateNotifications(func(_ *appmessage.NewBlockTemplateNotificationMessage) {
		blockReadyChan <- true
	})
	if err != nil {
		s.logger.Error("fatal: failed to register for block notifications from pyrin")
	}

	ticker := time.NewTicker(s.blockWaitTime)
	for {
		if err := s.waitForSync(false); err != nil {
			s.logger.Error("error checking pyrin sync state, attempting reconnect: ", err)
			if err := s.reconnect(); err != nil {
				s.logger.Error("error reconnecting to pyrin, waiting before retry: ", err)
				time.Sleep(5 * time.Second)
			}
		}
		select {
		case <-ctx.Done():
			s.logger.Warn("context cancelled, stopping block update listener")
			return
		case <-blockReadyChan:
			blockReadyCb()
			ticker.Reset(s.blockWaitTime)
		case <-ticker.C: // timeout, manually check for new blocks
			blockReadyCb()
		}
	}
}

func (py *PyrinApi) GetBlockTemplate(
	client *gostratum.StratumContext) (*appmessage.GetBlockTemplateResponseMessage, error) {
	template, err := py.pyrin.GetBlockTemplate(client.WalletAddr,
		fmt.Sprintf(`'%s' via Pyrinpyi/pyrin-stratum-bridge_%s`, client.RemoteApp, version))
	if err != nil {
		return nil, errors.Wrap(err, "failed fetching new block template from pyrin")
	}
	return template, nil
}
