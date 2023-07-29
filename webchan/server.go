package webchan

import (
	"encoding/json"
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

    members map[string]IMember
    rooms   map[string]*Room

    connectionCh chan *websocket.Conn
    recvCh       chan MemberMsg
    errCh        chan MemberErrorMsg

    event_handlers map[string]func(*Server, IMember, []byte) error
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
        case wsConn, ok := <-s.connectionCh:
            if !ok {
                break
            }
            m := NewMember(uuid.NewV4(), wsConn, s)
            s.addMember(m)
            go m.loop(s.recvCh, s.errCh)
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

func (s *Server) onMemberMsg(memberMsg MemberMsg) error {
    parsedMsg := &Message{}
    if err := json.Unmarshal(memberMsg.reqBody, parsedMsg); err != nil {
        logger.Error("Message:Parse() error: %v, member: %s req: %s", err,
            memberMsg.member.Info(), string(memberMsg.reqBody))
        return err
    }
    err := s.event_handlers[parsedMsg.EventName](s, memberMsg.member, parsedMsg.Data)
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

func (s *Server) getOrCreateRoomByRoomId(roomId string, create bool) *Room {
    room, ok := s.rooms[roomId]
    if ok {
        return room
    } else if create {
        room = NewRoom(roomId)
        s.rooms[roomId] = room
        return room
    }
    return nil
}

func (s *Server) findRoomByMember(m IMember) *Room {
    if room, ok := s.rooms[m.GetRoomId()]; ok {
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

func (s *Server) memberQuitRoom(m IMember) {
    if m.GetRoomId() == "" {
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

func (s *Server) getMemberByUuid(uuid string) IMember {
    if member, ok := s.members[uuid]; ok {
        return member
    }
    return nil
}

func (s *Server) addMember(m IMember) {
    s.members[m.GetUuid().String()] = m
}

func (s *Server) removeMember(m IMember) {
    logger.Info("member: %s", m.Info())
    s.memberQuitRoom(m)
    delete(s.members, m.GetUuid().String())
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
        members:      make(map[string]IMember),
        rooms:        make(map[string]*Room),
        connectionCh: make(chan *websocket.Conn, 128),
        recvCh:       make(chan MemberMsg, 16),
        errCh:        make(chan MemberErrorMsg, 16),
        event_handlers: map[string]func(*Server, IMember, []byte) error{
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
