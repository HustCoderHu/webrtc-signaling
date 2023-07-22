package webchan

import (
	"net/http"
	"sync"
	"webrtc-signaling/pkg/logger"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin:     func(r *http.Request) bool { return true },
}

type Server struct {
    http.Handler

    handlers

    members     map[uuid.UUID]*Member
    membersLock sync.RWMutex

    // connections map[string]map[*Connection]struct{}

    rooms     map[string]*Room
    roomsLock sync.RWMutex

    recvCh chan MemberMsg
    errCh  chan MemberErrorMsg

    fnMap map[string]func(*Server, *Member, []byte) error

    // connectionLock sync.RWMutex

    // cids     map[string]*Connection
    // cidsLock sync.RWMutex
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

    // we need to do some auth

    // if err := s.OnAuth(r); err != nil {

    //   w.WriteHeader(http.StatusUnauthorized)
    //   w.Write([]byte("HTTP status code returned!"))

    //   return
    // }

    ws, err := upgrader.Upgrade(w, r, nil)

    if err != nil {
        logger.Error("%s", err)
        return
    }
    m := NewMember(uuid.NewV4(), ws, s)
    go m.loop(s.recvCh, s.errCh)

    s.membersLock.Lock()
    defer s.membersLock.Unlock()
    s.members[m.uuid] = m
}

// func (s *Server) onNewConnection(wsConn *websocket.Conn) {
//     m := NewMember(uuid.NewV4(), wsConn, s)
//     s.membersLock.Lock()
//     defer s.membersLock.Unlock()
//     s.members[m.uuid] = m
//     // s.cidsLock.Lock()
//     // defer s.cidsLock.Unlock()
//     // s.cids[string(m.uuid)] = c
// }

func (s *Server) Loop() {
    for {
        select {
        case memberMsg, ok := <-s.recvCh:
            if !ok {
                continue
            }
            if err := s.onMemberMsg(memberMsg); err != nil {
                break
            }
            break

        case memberErrorMsg, ok := <-s.errCh:
            if !ok {
                break
            }
            if err := s.onMemberErrorMsg(memberErrorMsg); err != nil {
                break
            }
            break
        }
    }
}

func (s *Server) onMemberMsg(memberMsg MemberMsg) error {
    parsedMsg := Message{}
    if err := parsedMsg.Parse(memberMsg.reqBody); err != nil {
        logger.Error("Message:Parse() error: %s, member: %s req: %s", err,
            memberMsg.member.Info(), string(memberMsg.reqBody))
        return err
    }
    err := s.fnMap[parsedMsg.EventName](s, memberMsg.member, parsedMsg.Data)
    if err != nil {
        logger.Error("parsedMsg.EventName handle error: %s, member: %s req: %s",
            err, memberMsg.member.Info(), memberMsg.reqBody)
        return err
    }
    return nil
}

func (s *Server) onMemberErrorMsg(memberErrorMsg MemberErrorMsg) error {
    return nil
}

// func (s *Server) On

func (s *Server) quitRoom(m *Member) {
    if m.roomId == "" {
        return
    }
    if _, ok := s.rooms[m.roomId]; ok {
        s.roomsLock.Lock()
        defer s.roomsLock.Unlock()
        delete(s.rooms, m.roomId)
        logger.Info("member: %s", m.Info())
    }
}

func (s *Server) removeMember(m *Member) {
    s.quitRoom(m)
    s.membersLock.Lock()
    defer s.membersLock.Unlock()
    delete(s.members, m.uuid)
}

// func (s *Server) onCleanConnection(c *Connection) {

//     s.connectionLock.Lock()
//     defer s.connectionLock.Unlock()

//     conns := s.connections

//     byRoom, ok := s.rooms[c]

//     if ok {

//         for room := range byRoom {
//             if curRoom, ok := conns[room]; ok {
//                 delete(curRoom, c)
//                 if len(curRoom) == 0 {
//                     delete(conns, room)
//                 }
//             }
//         }

//         delete(s.rooms, c)
//     }

//     s.cidsLock.Lock()
//     defer s.cidsLock.Unlock()

//     delete(s.cids, c.id)

//     c.ws.Close()
// }

func (s *Server) BroadcastTo(room, message string, args interface{}) {

}

func (s *Server) BroadcastToAll(message string, args interface{}) {

}

// func (s *Server) List(room string) ([]*Connection, error) {
//     return nil, nil
// }

/**
new server
*/

func NewServer() *Server {
    s := Server{
        members: make(map[uuid.UUID]*Member),
        rooms:   make(map[string]*Room),
        recvCh:  make(chan MemberMsg, 16),
        errCh:   make(chan MemberErrorMsg, 16),
        fnMap: map[string]func(*Server, *Member, []byte) error{
            "__join":          onJoin,
            "__ice_candidate": onIceCandidate,
            "__offer":         onOffer,
            "__answer":        onAnswer,
            "__invite":        onInvite,
            "__ack":           onAck,
        },
    }
    // s.connections = make(map[string]map[*Connection]struct{})
    // s.rooms = make(map[*Connection]map[string]struct{})
    // s.cids = make(map[string]*Connection)
    return &s
}
