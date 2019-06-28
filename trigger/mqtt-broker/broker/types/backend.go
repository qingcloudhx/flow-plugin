//Created by zhbinary on 2019-03-19.
//Email: zhbinary@gmail.com
package types

import "github.com/256dpi/gomqtt/packet"

type AckCallback func()

type Backend interface {
	// Authenticate should authenticate the client using the user and password
	// values and return true if the client is eligible to continue or false
	// when the broker should terminate the connection.
	Authenticate(session Session, user, password string) (ok bool, err error)

	// Setup is called when a new client comes online and is successfully
	// authenticated. Setup should return the already stored session for the
	// supplied id or create and return a new one if it is missing or a clean
	// session is requested. If the supplied id has a zero length, a new
	// temporary session should be returned that is not stored further. The
	// backend should also close any existing clients that use the same id.
	//
	// Note: In this call the Backend may also allocate other resources and
	// setup the client for further usage as the broker will acknowledge the
	// connection when the call returns. The Terminate function is called for
	// every client that Setup has been called for.
	Setup(session Session, id string, clean bool) (a Session, resumed bool, err error)

	// Restore is called after the client has restored packets from the session.
	//
	// The Backend should resubscribe stored subscriptions and begin with queueing
	// missed offline messages. When all offline messages have been queued the
	// client may receive online messages. Depending on the implementation, this
	// may not be required as Dequeue will already pick up offline messages.
	Restore(session Session) error

	// Subscribe should subscribe the passed client to the specified topics and
	// store the subscription in the session. If an Ack is provided, the
	// subscription will be acknowledged when called during or after the call to
	// Subscribe.
	//
	// Incoming messages that match the supplied subscription should be added to
	// a temporary or persistent queue that is drained when Dequeue is called.
	//
	// Retained messages that match the supplied subscription should be added to
	// a temporary queue that is also drained when Dequeue is called. The messages
	// must be delivered with the retained flag set to true.
	Subscribe(session Session, subs []packet.Subscription, ack AckCallback) error

	// Unsubscribe should unsubscribe the passed client from the specified topics
	// and remove the subscriptions from the session. If an Ack is provided, the
	// unsubscription will be acknowledged when called during or after the call
	// to Unsubscribe.
	Unsubscribe(session Session, topics []string, ack AckCallback) error

	// Publish should forward the passed message to all other clients that hold
	// a subscription that matches the messages topic. It should also add the
	// message to all sessions that have a matching offline subscription. The
	// later may only apply to messages with a QOS greater than 0. If an Ack is
	// provided, the message will be acknowledged when called during or after
	// the call to Publish.
	//
	// If the retained flag is set, messages with a payload should replace the
	// currently retained message. Otherwise, the currently retained message
	// should be removed. The flag should be cleared before publishing the
	// message to other subscribed clients.
	Publish(session Session, pkt *packet.Publish, ack AckCallback) error

	// Dequeue is called by the Client to obtain the next message from the queue
	// and must return either a message or an error. The backend must only return
	// no message and no error if the client's Closing channel has been closed.
	//
	// The Backend may return an Ack to receive a signal that the message is being
	// delivered under the selected qos level and is therefore safe to be deleted
	// from the queue. The Ack will be called before Dequeue is called again.
	//
	// The returned message must have a QOS set that respects the QOS set by
	// the matching subscription.
	Dequeue(session Session) (*packet.Publish, AckCallback, error)

	// Terminate is called when the client goes offline. Terminate should
	// unsubscribe the passed client from all previously subscribed topics. The
	// backend may also convert a clients subscriptions to offline subscriptions.
	//
	// Note: The Backend may also cleanup previously allocated resources for
	// that client as the broker will close the connection when the call
	// returns.
	Terminate(session Session) error

	// Log is called multiple times during the lifecycle of a client see LogEvent
	// for a list of all events.
	Event(event Event, session Session, pkt packet.Generic, msg *packet.Message, err error)
}
