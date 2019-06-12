package state

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/propsproject/pending-props/core/eth-utils"
	"github.com/propsproject/pending-props/core/proto/pending_props_pb"
	"github.com/propsproject/sawtooth-go-sdk/logging"
	"github.com/propsproject/sawtooth-go-sdk/processor"
)

type State struct {
	context *processor.Context
}

var logger = logging.Get()

func (s *State) GetSettlements(ethTxtIDS ...string) ([]pending_props_pb.Settlements, error) {
	addresses := make([]string, 0)
	for _, id := range ethTxtIDS {
		id = eth_utils.NormalizeAddress(id)
		address, _ := SettlementAddress(id)
		addresses = append(addresses, address)
	}

	state, err := s.context.GetState(addresses)
	if err != nil {
		return nil, &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not get state (%s)", err)}
	}

	settlements := make([]pending_props_pb.Settlements, 0)
	for _, value := range state {
		var settlement pending_props_pb.Settlements
		err := proto.Unmarshal(value, &settlement)
		if err != nil {
			return nil, &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not unmarshal proto data (%s)", err)}
		}
		settlements = append(settlements, settlement)
	}

	return settlements, nil
}

func (s *State) CurrentNonce(publickey string) (*pending_props_pb.Nonce, string, error) {
	address, _ := NonceAddress(publickey)
	state, err := s.context.GetState([]string{address})
	if err != nil {
		return nil, address, fmt.Errorf("could not get nonce from state (%s)", err)
	}

	if len(state[address]) == 0 {
		return &pending_props_pb.Nonce{Current: 0}, address, nil
	}

	var nonce pending_props_pb.Nonce
	err = proto.Unmarshal(state[address], &nonce)
	if err != nil {
		return nil, address, fmt.Errorf("could not marshal proto data for nonce (%s)", err)
	}

	return &nonce, address, nil
}

func (s *State) UpdateNonce(publickey string) error {
	currentNonce, address, err := s.CurrentNonce(publickey)
	if err != nil {
		return err
	}

	currentNonce.Current = currentNonce.Current + 1
	b, _ := proto.Marshal(currentNonce)
	_, err = s.context.SetState(map[string][]byte{address: b})
	if err != nil {
		return fmt.Errorf("could not update nonce (%s)", err)
	}

	return nil
}

func NewState(context *processor.Context) *State {
	return &State{context: context}
}
