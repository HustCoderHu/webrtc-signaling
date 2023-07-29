package webchan

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"testing"
	"webrtc-signaling/pkg/logger"

	uuid "github.com/satori/go.uuid"
)

type UuidArray []uuid.UUID

func (a UuidArray) Len() int           { return len(a) }
func (a UuidArray) Less(i, j int) bool { return a[i].String() < a[j].String() }
func (a UuidArray) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type FakeMember struct {
    uuid     uuid.UUID
    userName string
    roomId   string

    messages []json.RawMessage
}

func NewFakeMember(uuid uuid.UUID, userName string) *FakeMember {
    return &FakeMember{
        uuid:     uuid,
        userName: userName,
    }
}

func (t *FakeMember) GetUserName() string {
    return t.userName
}

func (t *FakeMember) GetUuid() uuid.UUID {
    return t.uuid
}

func (t *FakeMember) GetRoomId() string {
    return t.roomId
}

func (t *FakeMember) SetRoomId(roomId string) {
    t.roomId = roomId
}

func (t *FakeMember) Info() string {
    return "[name: " + t.GetUserName() + " addr: " +
        " id: " + t.GetUuid().String() + " ]"
}

func (t *FakeMember) OnMsg(msg []byte) string {
    s := string(msg)
    logger.Info("send fake member msg to %s: %s", t.Info(), s)
    t.messages = append(t.messages, msg)
    return s
}

func (t *FakeMember) Dispose() error {
    return nil
}

func TestRoomAdd(t *testing.T) {
    uuids := make([]uuid.UUID, 0, 5)
    for i := 0; i < 5; i++ {
        uuids = append(uuids, uuid.NewV4())
    }
    room := NewRoom("test")
    for i, v := range uuids {
        room.AddMember(NewFakeMember(v, fmt.Sprintf("user_%d", i)))
    }
    roomUuids := room.GetMemberUuids()
    sort.Sort(UuidArray(roomUuids))

    sort.Sort(UuidArray(uuids))

    if reflect.DeepEqual(roomUuids, uuids) {
        t.Log("success")
    } else {
        t.Errorf("fail %v != %v", roomUuids, uuids)
    }
}

func TestRoomAddAndRemove(t *testing.T) {
    uuids := make([]uuid.UUID, 0, 5)
    for i := 0; i < 5; i++ {
        uuids = append(uuids, uuid.NewV4())
    }
    room := NewRoom("test")
    members := make([]IMember, 0, len(uuids))
    for i, v := range uuids {
        m := NewFakeMember(v, fmt.Sprintf("user_%d", i))
        members = append(members, m)
        room.AddMember(m)
    }

    room.RemoveMember(members[0])
    room.RemoveMember(members[2])
    room.RemoveMember(members[4])
    roomUuids := room.GetMemberUuids()
    sort.Sort(UuidArray(roomUuids))

    uuids = []uuid.UUID{
        members[1].GetUuid(),
        members[3].GetUuid(),
    }
    sort.Sort(UuidArray(uuids))

    if reflect.DeepEqual(roomUuids, uuids) {
        t.Log("success")
    } else {
        t.Errorf("fail %v != %v", roomUuids, uuids)
    }
}
