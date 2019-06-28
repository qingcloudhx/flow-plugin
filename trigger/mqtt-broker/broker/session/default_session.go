//Created by zhbinary on 2019-03-07.
//Email: zhbinary@gmail.com
package session

import (
	"github.com/256dpi/gomqtt/packet"
	"github.com/256dpi/gomqtt/transport"
	"github.com/qingcloudhx/flow-plugin/trigger/mqtt-broker/broker/types"
	"github.com/qingcloudhx/flow-plugin/trigger/mqtt-broker/broker/utils"
	"gopkg.in/tomb.v2"
)

const (
	StateConnecting = iota
	StateConnected
	StateDisconnected
)

type DefaultSession struct {
	state   uint32
	backend types.Backend
	conn    transport.Conn

	id        string
	will      *packet.Message
	tomb      tomb.Tomb
	processor types.Processor
}

func NewSession(conn transport.Conn, backend types.Backend) types.Session {
	return &DefaultSession{conn: conn, backend: backend, state: StateConnecting}
}

func (this *DefaultSession) Close() {
	this.state = StateDisconnected
	this.conn.Close()
}

func (this *DefaultSession) close(ev types.Event, err error) {
	this.backend.Event(ev, this, nil, nil, err)

	this.state = StateDisconnected

	this.conn.Close()
}

func (this *DefaultSession) NextID() packet.ID {
	panic("implement me")
}

func (this *DefaultSession) Start() {
	for {
		if this.state == StateDisconnected {
			break
		}

		pack, err := this.conn.Receive()
		if err != nil {
			return
		}

		err = this.process(pack)
		if err != nil {
			this.close(types.EventSessionInactive, err)
			return
		}
	}
}

func (this *DefaultSession) process(pkt packet.Generic) error {
	utils.Log.Debug("process", pkt)
	// prepare error
	var err error

	// handle individual packets
	switch typedPkt := pkt.(type) {
	case *packet.Connect:
		err = this.processConnect(typedPkt)
	case *packet.Subscribe:
		err = this.processSubscribe(typedPkt)
	case *packet.Unsubscribe:
		err = this.processUnsubscribe(typedPkt)
	case *packet.Publish:
		err = this.processPublish(typedPkt)
	case *packet.Puback:
		err = this.processPubackAndPubcomp(typedPkt)
	case *packet.Pingreq:
		err = this.processPingreq(typedPkt)
	case *packet.Disconnect:
		err = this.processDisconnect(typedPkt)
	default:
		//err = this.die(ClientError, ErrUnexpectedPacket)
	}

	// return eventual error
	if err != nil {
		return err // error has already been handled
	}

	return nil
}

func (this *DefaultSession) processConnect(pkt *packet.Connect) error {
	// prepare connack packet
	connack := packet.NewConnack()

	// Check version
	if pkt.Version != packet.Version311 {
		connack.ReturnCode = packet.InvalidProtocolVersion
		connack.SessionPresent = false
		this.conn.Send(connack, false)
		return types.ConnectError
	}

	// Authentication
	ok, err := this.backend.Authenticate(this, pkt.Username, pkt.Password)
	if !ok || err != nil {
		connack.ReturnCode = packet.NotAuthorized
		connack.SessionPresent = false
		this.conn.Send(connack, false)
		return err
	}

	// Setup session
	_, _, err = this.backend.Setup(this, pkt.ClientID, pkt.CleanSession)
	if err != nil {
		connack.ReturnCode = packet.ServerUnavailable
		connack.SessionPresent = false
		this.conn.Send(connack, false)
	}

	// Start send loop

	return nil
}

func (this *DefaultSession) processSubscribe(pkt *packet.Subscribe) error {
	return nil
}

func (this *DefaultSession) processUnsubscribe(pkt *packet.Unsubscribe) error {
	return nil
}

func (this *DefaultSession) processPublish(pkt *packet.Publish) error {
	// Send it to backend layer and ack
	return this.backend.Publish(this, pkt, func() {
		pubAck := packet.NewPuback()
		pubAck.ID = pkt.ID
		err := this.conn.Send(pubAck, true)
		if err != nil {
			this.conn.Close()
		}
	})
}

func (this *DefaultSession) processPubackAndPubcomp(pkt *packet.Puback) error {
	return nil
}

func (this *DefaultSession) processPingreq(pkt *packet.Pingreq) error {
	msg := packet.NewPingresp()
	this.conn.Send(msg, true)
	return nil
}

func (this *DefaultSession) processDisconnect(pkt *packet.Disconnect) error {
	return nil
}

func (this *DefaultSession) sendLoop() {
	for {
		msg, ackFunc, err := this.backend.Dequeue(this)
		if err != nil {
			this.close(types.EventSessionDequeue, err)
		}
		this.conn.Send(msg, true)
		ackFunc()
	}
}
