package webchan

import (
	"fmt"
	"webrtc-signaling/pkg/logger"

	uuid "github.com/satori/go.uuid"
)

type Room struct {
    roomId  string
    members map[string]IMember
}

func NewRoom(roomId string) *Room {
    return &Room{
        roomId:  roomId,
        members: make(map[string]IMember),
    }
}

func (t *Room) AddMember(m IMember) {
    m.SetRoomId(t.roomId)
    t.members[m.GetUuid().String()] = m
    logger.Info("room: %s, member: %s", t.Info(), m.Info())
}

func (t *Room) RemoveMember(m IMember) {
    delete(t.members, m.GetUuid().String())
    logger.Info("room: %s, member: %s", t.Info(), m.Info())
    // m.SetRoomId("")
}

func (t *Room) CountMembers() int {
    if t.members == nil {
        return 0
    }
    return len(t.members)
}

func (t *Room) GetMemberUuids() []uuid.UUID {
    if t.members == nil {
        return nil
    }
    uuids := make([]uuid.UUID, 0, len(t.members))
    for _, m := range t.members {
        uuids = append(uuids, m.GetUuid())
    }
    return uuids
}

func (t *Room) Info() string {
    return fmt.Sprintf("[roomId: %s members: %d ]", t.roomId, len(t.members))
    // return "[roomId: " + t.roomId + " members: " + len(t.members) + "]"
}

func (t *Room) BroadCastMsgExceptMember(msg []byte, exceptMemberUuid uuid.UUID) {
    uidstr := exceptMemberUuid.String()
    for uuid, m := range t.members {
        if uuid == uidstr {
            continue
        }
        m.OnMsg(msg)
    }
}
