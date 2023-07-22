package webchan

import uuid "github.com/satori/go.uuid"

type Room struct {
    roomId string
    members map[uuid.UUID]*Member
}

func (t *Room) AddMember(m *Member) {
    t.members[m.uuid] = m
}

func (t *Room) RemoveMember() {

}