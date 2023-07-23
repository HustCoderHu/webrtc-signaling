package webchan

import (
	"net/http"
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
    // http.Handler
    Handlers

    members     map[string]*Member
    rooms     map[string]*Room

    connectionCh chan *websocket.Conn
    recvCh chan MemberMsg
    errCh  chan MemberErrorMsg

    fnMap map[string]func(*Server, *Member, []byte) error

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
    s.connectionCh <- ws
}

func (s *Server) Loop() {
    for {
        select {
        case wsConn, ok := <- s.connectionCh:
            if !ok {
                break
            }
            if err := s.onConnection(wsConn); err != nil {
                break
            }
            break
        case memberMsg, ok := <-s.recvCh:
            if !ok {
                break
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

func (s *Server) onConnection(wsConn *websocket.Conn) error {
    m := NewMember(uuid.NewV4(), wsConn, s)
    s.members[m.uuid.String()] = m
    go m.loop(s.recvCh, s.errCh)
    return nil
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
    logger.Info("msg: %s", memberErrorMsg.msg)
    s.removeMember(memberErrorMsg.member)

    return nil
}

func (s *Server) GetOrCreateRoomByRoomId(roomId string, create bool) *Room {
    room, ok := s.rooms[roomId]
    if ok {
        return room
    } else if create {
        room = &Room { roomId: roomId }
        s.rooms[roomId] = room
        return room
    }
    return nil
}

func (s *Server) findRoomByMember(m *Member) *Room {
    if room, ok := s.rooms[m.roomId]; ok {
        return room
    }
    return nil
}

func (s *Server) checkRemoveRoom(room *Room) {
    if room.CountMembers() == 0 {
        logger.Info("no members, so delete room: %s", room.roomId)
        delete(s.rooms, room.roomId)
    }
}

func (s *Server) memberQuitRoom(m *Member) {
    if m.roomId == "" {
        return
    }
    room := s.findRoomByMember(m)
    if room == nil {
        return
    }
    room.RemoveMember(m)
    logger.Info("member: %s, room: %s", m.Info(), room.Info())
    s.checkRemoveRoom(room)
}

func (s *Server) GetMemberByUuid(uuid uuid.UUID) *Member {
    if member, ok := s.members[uuid.String()]; ok {
        return member
    }
    return nil
}

func (s *Server) removeMember(m *Member) {
    logger.Info("member: %s", m.Info())
    s.memberQuitRoom(m)
    delete(s.members, m.uuid.String())
    m.Dispose()
}

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
        members: make(map[string]*Member),
        rooms:   make(map[string]*Room),
        connectionCh: make(chan *websocket.Conn, 128),
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
    return &s
}
