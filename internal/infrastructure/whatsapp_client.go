package infrastructure

import (
	"context"
	"fmt"
	"github.com/chatbox/whatsapp/internal/domain"
	"github.com/chatbox/whatsapp/pkg/errors"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
	"sync"
)

type WhatsAppClientImpl struct {
	log       Logger
	container *sqlstore.Container
	client    ClientMap
}

type ClientMap struct {
	mu        sync.RWMutex
	clientMap map[string]*whatsmeow.Client
}

type WhatsAppClient interface {
	ConnectWithHandler(session *domain.Session, handler func(interface{})) error
	Disconnect(session *domain.Session) error
	GetOrCreateClient(session *domain.Session) (*whatsmeow.Client, error)
	GenerateQR(ctx context.Context, session *domain.Session) (string, error)
	Send(ctx context.Context, session *domain.Session, jid string, text string) error
}

func NewWhatsAppClient(dsn string, log Logger) WhatsAppClient {
	logger := waLog.Noop
	container, err := sqlstore.New("postgres", dsn, logger)
	if err != nil {
		panic(err)
	}
	return &WhatsAppClientImpl{
		container: container,
		client: ClientMap{
			clientMap: make(map[string]*whatsmeow.Client),
		},
		log: log,
	}
}

func (w *WhatsAppClientImpl) GetOrCreateClient(session *domain.Session) (*whatsmeow.Client, error) {
	w.client.mu.Lock()
	defer w.client.mu.Unlock()

	if client, ok := w.client.clientMap[session.Token]; ok {
		return client, nil
	}

	var deviceStore *store.Device
	var err error

	if session.JID != "" {
		jid, parseErr := types.ParseJID(session.JID)
		if parseErr != nil {
			return nil, parseErr
		}
		deviceStore, err = w.container.GetDevice(jid)
	} else {
		deviceStore = w.container.NewDevice()
	}

	if err != nil {
		return nil, err
	}

	client := whatsmeow.NewClient(deviceStore, waLog.Noop)
	w.client.clientMap[session.Token] = client

	return client, nil
}

func (w *WhatsAppClientImpl) GenerateQR(ctx context.Context, session *domain.Session) (string, error) {
	client, err := w.GetOrCreateClient(session)
	if err != nil {
		return "", err
	}

	if client.Store.ID != nil {
		return "", errors.UserAlreadyLogged
	}

	qrChan, err := client.GetQRChannel(ctx)
	if err != nil {
		return "", err
	}

	err = client.Connect()
	if err != nil {
		return "", err
	}

	for evt := range qrChan {
		if evt.Event == "code" {
			return evt.Code, nil
		}
		if evt.Event == "timeout" {
			return "", fmt.Errorf("qr code timeout")
		}
		if evt.Event == "success" {
			return "", nil
		}
	}

	return "", fmt.Errorf("unexpected qr channel closure")
}

func (w *WhatsAppClientImpl) ConnectWithHandler(session *domain.Session, handler func(interface{})) error {
	client, err := w.GetOrCreateClient(session)
	if err != nil {
		return err
	}
	client.AddEventHandler(handler)
	return client.Connect()
}

func (w *WhatsAppClientImpl) Disconnect(session *domain.Session) error {
	w.client.mu.Lock()
	defer w.client.mu.Unlock()
	if client, ok := w.client.clientMap[session.Token]; ok {
		client.Disconnect()
		delete(w.client.clientMap, session.Token)
	}
	return nil
}

func (c *WhatsAppClientImpl) Send(ctx context.Context, session *domain.Session, jid string, text string) error {
	client, err := c.GetOrCreateClient(session)
	if client == nil || err != nil {
		return fmt.Errorf("whatsapp client not found for token: %s", session.Token)
	}

	recipient, err := types.ParseJID(jid)
	if err != nil {
		return fmt.Errorf("invalid JID: %w", err)
	}

	msg := &waE2E.Message{
		Conversation: proto.String(text),
	}

	_, err = client.SendMessage(ctx, recipient, msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
