package cdc

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pglogrepl"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgproto3"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/event"
)

type Listener struct {
	connString  string
	slotName    string
	publication string
	bus         *event.Bus
}

func NewListener(connString string, bus *event.Bus) *Listener {
	return &Listener{
		connString:  connString,
		slotName:    "app_slot",
		publication: "app_publication",
		bus:         bus,
	}
}

func (l *Listener) Start(ctx context.Context) error {

	conn, err := pgconn.Connect(ctx, l.connString)
	if err != nil {
		return err
	}

	sysident, err := pglogrepl.IdentifySystem(ctx, conn)
	if err != nil {
		return err
	}

	log.Println("SystemID:", sysident.SystemID)

	_, err = pglogrepl.CreateReplicationSlot(
		ctx,
		conn,
		l.slotName,
		"pgoutput",
		pglogrepl.CreateReplicationSlotOptions{Temporary: false},
	)

	if err != nil {
		log.Println("slot probably exists:", err)
	}

	err = pglogrepl.StartReplication(
		ctx,
		conn,
		l.slotName,
		sysident.XLogPos,
		pglogrepl.StartReplicationOptions{
			PluginArgs: []string{
				"proto_version '1'",
				"publication_names '" + l.publication + "'",
			},
		},
	)

	if err != nil {
		return err
	}

	for {

		ctxReceive, cancel := context.WithTimeout(ctx, 10*time.Second)

		msg, err := conn.ReceiveMessage(ctxReceive)
		cancel()

		if err != nil {
			continue
		}

		switch msg := msg.(type) {

		case *pgproto3.CopyData:
			l.handleCopyData(msg.Data)
		}
	}
}

func (l *Listener) handleCopyData(data []byte) {

	switch data[0] {

	case pglogrepl.PrimaryKeepaliveMessageByteID:
		keepalive, err := pglogrepl.ParsePrimaryKeepaliveMessage(data[1:])
		if err != nil {
			return
		}

		if keepalive.ReplyRequested {
			l.sendStatus()
		}

	case pglogrepl.XLogDataByteID:

		xld, err := pglogrepl.ParseXLogData(data[1:])
		if err != nil {
			return
		}

		l.handleWAL(xld.WALData)
	}
}

func (l *Listener) handleWAL(data []byte) {

	msg, err := pglogrepl.Parse(data)
	if err != nil {
		return
	}

	switch m := msg.(type) {

	case *pglogrepl.InsertMessage:
		l.handleInsert(m)

	case *pglogrepl.UpdateMessage:
		// optional

	case *pglogrepl.DeleteMessage:
		// optional
	}
}

func (l *Listener) handleInsert(msg *pglogrepl.InsertMessage) {

	tuple := msg.Tuple

	for i, col := range tuple.Columns {

		if col.DataType != 't' {
			continue
		}

		value := string(col.Data)

		log.Println("column", i, value)
	}
}

func (l *Listener) sendStatus() {
	log.Println("server: confirmed standby")
}
