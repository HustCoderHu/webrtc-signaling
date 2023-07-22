package webchan

import (
	"fmt"
	"time"
	"webrtc-signaling/pkg/logger"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

type MemberMsg struct {
    member *Member
    reqBody    []byte
}

type MemberErrorMsg struct {
    member *Member
    msg    string
}

type Member struct {
    uuid uuid.UUID
    userName string
    roomId   string // empty if not join any room
    wsConn   *websocket.Conn
    server   *Server

    iceCandidate string
    sendCh   chan []byte
}

func NewMember(uniqueId uuid.UUID, wsConn *websocket.Conn, s *Server) *Member {
    return &Member{
        uuid: uniqueId,
        wsConn: wsConn,
        server: s,
        sendCh: make(chan []byte, 8),
    }
}

func (t *Member) Info() string {
    addr := t.wsConn.RemoteAddr().String()
    return "[name: " + t.userName + " addr: " +
      addr + " id: " + t.uuid.String() + " ]"
}

func (t *Member) OnMsg(msg []byte) {
    t.sendCh <- msg
}

// for internal
func (t *Member) write(mt int, payload []byte) error {
    t.wsConn.SetWriteDeadline(time.Now().Add(writeWait))
    return t.wsConn.WriteMessage(mt, payload)
}

func (t *Member) readLoop(recvCh chan MemberMsg, errCh chan MemberErrorMsg) {
    t.wsConn.SetReadLimit(maxMessageSize)
    t.wsConn.SetReadDeadline(time.Now().Add(pongWait))
    t.wsConn.SetPongHandler(func(string) error {
        t.wsConn.SetReadDeadline(time.Now().Add(pongWait)); return nil
    })

    for {
        // we only surport TextMessage
        _, msg, err := t.wsConn.ReadMessage()

        if err != nil {
            if websocket.IsUnexpectedCloseError(
              err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
                logger.Error("websocekt read error: %s, %s", err, t.Info())
            }
            errCh <- MemberErrorMsg{t, 
                fmt.Sprintf("websocekt read error: %s, %s", err, t.Info())}
            // t.server.OnDisconnection(t, err.Error())
            break
        }

        // onMessage
        recvCh <- MemberMsg{t, msg}
    }
}

func (t *Member) loop(recvCh chan MemberMsg, errCh chan MemberErrorMsg) {
    go t.readLoop(recvCh, errCh)

    ticker := time.NewTicker(pingPeriod)

    defer func() {
        ticker.Stop()
        t.wsConn.Close()
    }()
    for {
        select {
        case msg, ok := <-t.sendCh:
            if !ok {
                t.write(websocket.CloseMessage, []byte{})
                continue
            }
            if err := t.write(websocket.TextMessage, msg); err != nil {
                return
            }
        case <-ticker.C:
            if err := t.write(websocket.PingMessage, []byte{}); err != nil {
                continue
            }
        }
    }
}
