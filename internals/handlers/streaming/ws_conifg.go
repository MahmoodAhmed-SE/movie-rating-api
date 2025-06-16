package streaming

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Notification struct {
	Date    string `json:"date"`
	Content string `json:"content"`
}

type Subscriber struct {
	UserId int
	Ch     chan Notification
}

type Broker struct {
	Mu          sync.RWMutex
	Subscribers map[string][]Subscriber
}

var (
	GlobalBroker = NewBroker()
)

func NewBroker() *Broker {
	return &Broker{
		Subscribers: make(map[string][]Subscriber),
	}
}

func (b *Broker) SubscribeNotificationReceiver(topic string, userId int) *Subscriber {
	b.Mu.Lock()
	defer b.Mu.Unlock()

	subscription := Subscriber{
		UserId: userId,
		Ch:     make(chan Notification, 10),
	}
	b.Subscribers[topic] = append(b.Subscribers[topic], subscription)

	return &subscription
}

func (b *Broker) PublishBroadcastNotification(topic string, message string) error {
	b.Mu.RLock()
	defer b.Mu.RUnlock()

	for _, subscriber := range b.Subscribers[topic] {
		select {
		case subscriber.Ch <- Notification{
			Date:    time.Now().Format("01/02/2006"),
			Content: message,
		}:
		default:
			log.Printf("User notification channel receiver is full and a new notification has been dropped for user with id: %d", subscriber.UserId)
		}

	}

	return nil
}

func (b *Broker) RemoveSubscriber(topic string, userSubscription Subscriber) error {
	b.Mu.Lock()
	defer b.Mu.Unlock()

	for i, subscriber := range b.Subscribers[topic] {
		if subscriber.UserId == userSubscription.UserId {
			close(subscriber.Ch)
			b.Subscribers[topic] = append(b.Subscribers[topic][:i], b.Subscribers[topic][i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("subscription with user id %v not found", userSubscription.UserId)
}
