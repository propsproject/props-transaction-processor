package state

import (
	"github.com/gogo/protobuf/proto"
	"github.com/propsproject/props-transaction-processor/core/proto/pending_props_pb"
	"github.com/propsproject/sawtooth-go-sdk/processor"
)

func (s *State) AddEvent(event pending_props_pb.TransactionEvent, eventType string, attributes ...processor.Attribute) error {
	b, err := proto.Marshal(&event)
	if err != nil {
		return err
	}
	logger.Info(event, eventType, attributes)
	return s.context.AddEvent(eventType, attributes, b)
}

func (s *State) AddBalanceEvent(event pending_props_pb.BalanceEvent, eventType string, attributes ...processor.Attribute) error {
	b, err := proto.Marshal(&event)
	if err != nil {
		return err
	}
	logger.Info(event, eventType, attributes)
	return s.context.AddEvent(eventType, attributes, b)
}

func (s *State) AddWalletLinkEvent(event pending_props_pb.WalletLinkedEvent, eventType string, attributes ...processor.Attribute) error {
	b, err := proto.Marshal(&event)
	if err != nil {
		return err
	}
	logger.Info(event, eventType, attributes)
	return s.context.AddEvent(eventType, attributes, b)
}

func (s *State) AddWalletUnlinkEvent(event pending_props_pb.WalletUnlinkedEvent, eventType string, attributes ...processor.Attribute) error {
	b, err := proto.Marshal(&event)
	if err != nil {
		return err
	}
	logger.Info(event, eventType, attributes)
	return s.context.AddEvent(eventType, attributes, b)
}

func (s *State) AddLastEthBlockUpdateEvent(event pending_props_pb.LastEthBlockEvent, eventType string, attributes ...processor.Attribute) error {
	b, err := proto.Marshal(&event)
	if err != nil {
		return err
	}
	logger.Info(event, eventType, attributes)
	return s.context.AddEvent(eventType, attributes, b)
}
