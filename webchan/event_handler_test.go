package webchan

import (
	"encoding/json"
	"strings"
	"testing"

	uuid "github.com/satori/go.uuid"
)

func TestServerFirstMemberJoin(t *testing.T) {
    s := NewServer()

    memberAlice := NewFakeMember(uuid.NewV4(), "Alice")
    msgData, err := json.Marshal(map[string]interface{}{
        "roomId": "test_room_0",
    })

    if err != nil {
        t.Errorf("json.Marshal err:%v", err)
    }

    onJoin(s, memberAlice, msgData)

    if len(memberAlice.messages) != 1 {
        t.Errorf("len(s.members) != 1")
        return
    }

    checkMessagePeers(t, memberAlice, memberAlice.messages[0], []string{},
        memberAlice.GetUuid().String())
}

func TestServer2MembersJoin(t *testing.T) {
    s := NewServer()

    memberAlice := NewFakeMember(uuid.NewV4(), "Alice")
    msgData, err := json.Marshal(map[string]interface{}{
        "roomId": "test_room_0",
    })

    if err != nil {
        t.Errorf("json.Marshal err:%v", err)
    }
    onJoin(s, memberAlice, msgData)
    memberAlice.messages = memberAlice.messages[:0]

    memberBob := NewFakeMember(uuid.NewV4(), "Bob")
    msgData, err = json.Marshal(map[string]interface{}{
        "roomId": "test_room_0",
    })
    if err != nil {
        t.Errorf("json.Marshal err:%v", err)
    }
    onJoin(s, memberBob, msgData)

    // check member peers rsp
    checkMessagePeers(t, memberBob, memberBob.messages[0], []string{memberAlice.GetUuid().String()},
        memberBob.GetUuid().String())

    // check member 0 new_peers message
    if len(memberAlice.messages) != 1 {
        t.Errorf("len(memberAlice.messages) != 1")
        return
    }
    checkMessageNewPeer(t, memberAlice, memberAlice.messages[0],
        memberBob.GetUuid().String())
}

func checkMessagePeers(t *testing.T, m *FakeMember, msg json.RawMessage,
    connections []string, you string) {
    parsedMsg := &Message{}
    if err := json.Unmarshal(msg, parsedMsg); err != nil {
        t.Errorf("json.Unmarshal err:%v", err)
        return
    }

    if parsedMsg.EventName != "_peers" {
        t.Errorf("parsedMsg.EventName(%v) != _peers,", parsedMsg.EventName)
        return
    }

    MessageDataPeers := &MessageDataPeers{}
    if err := json.Unmarshal(parsedMsg.Data, MessageDataPeers); err != nil {
        t.Errorf("json.Unmarshal err:%v", err)
        return
    }

    if MessageDataPeers.You != you {
        t.Errorf("MessageDataPeers.You != m.GetUuid().String(), %v != %v",
            MessageDataPeers.You, you)
        return
    }

    if len(MessageDataPeers.Connections) != len(connections) {
        t.Errorf("len(MessageDataPeers.Connections) != len(connections)")
        return
    }

    for i := 0; i < len(connections); i++ {
        if MessageDataPeers.Connections[i] != connections[i] {
            t.Errorf("i = %d, connections %v != %v",
                i, MessageDataPeers.Connections, connections)
            return
        }
    }
}

func checkMessageNewPeer(t *testing.T, m *FakeMember, msg json.RawMessage,
    newPeerUuid string) {
    parsedMsg := &Message{}
    if err := json.Unmarshal(msg, parsedMsg); err != nil {
        t.Errorf("json.Unmarshal err:%v", err)
        return
    }
    if parsedMsg.EventName != "_new_peer" {
        t.Errorf("parsedMsg.EventName(%v) != _new_peer", parsedMsg.EventName)
        return
    }
    MessageDataNewPeer := &MessageDataNewPeer{}
    if err := json.Unmarshal(parsedMsg.Data, MessageDataNewPeer); err != nil {
        t.Errorf("json.Unmarshal err:%v", err)
        return
    }
    if MessageDataNewPeer.SocketId != newPeerUuid {
        t.Errorf("MessageDataNewPeer.SocketId != m.GetUuid().String(), %v != %v",
            MessageDataNewPeer.SocketId, newPeerUuid)
        return
    }
}

func TestServer2ndSendOffer(t *testing.T) {
    s := NewServer()

    memberAlice := NewFakeMember(uuid.NewV4(), "Alice")
    s.addMember(memberAlice)

    msgData, err := json.Marshal(map[string]interface{}{
        "roomId": "test_room_0",
    })

    if err != nil {
        t.Errorf("json.Marshal err:%v", err)
    }
    onJoin(s, memberAlice, msgData)

    memberBob := NewFakeMember(uuid.NewV4(), "Bob")
    s.addMember(memberBob)

    msgData, err = json.Marshal(map[string]interface{}{
        "roomId": "test_room_0",
    })
    if err != nil {
        t.Errorf("json.Marshal err:%v", err)
    }
    onJoin(s, memberBob, msgData)

    memberAlice.messages = memberAlice.messages[:0]
    // member 1 send offser
    sdp := `{"type": "offer", "sdp": "offer_from_bob"}`
    jsonRaw := json.RawMessage(sdp)
    msgDataOffer := &MessageDataOffer{
        Sdp:      &jsonRaw,
        SocketId: memberAlice.GetUuid().String(),
    }
    data, err := json.Marshal(msgDataOffer)
    if err != nil {
        t.Errorf("json.Marshal err:%v", err)
    }
    onOffer(s, memberBob, data)

    if len(memberAlice.messages) != 1 {
        t.Errorf("len(memberAlice.messages) != 1, %v",
            memberAlice.messages)
        return
    }
    // check sdp
    checkOfferRsp(t, memberBob, memberAlice.messages[0],
        sdp, memberBob.GetUuid().String())
}

func checkOfferRsp(t *testing.T, m *FakeMember, offerRsp json.RawMessage,
    sdp string, socketId string) {
    parsedMsg := &Message{}
    if err := json.Unmarshal(offerRsp, parsedMsg); err != nil {
        t.Errorf("json.Unmarshal err:%v", err)
        return
    }
    if parsedMsg.EventName != "_offer" {
        t.Errorf("parsedMsg.EventName(%v) != _offer", parsedMsg.EventName)
        return
    }
    MessageDataOfferRsp := &MessageDataOfferRsp{}
    if err := json.Unmarshal(parsedMsg.Data, MessageDataOfferRsp); err != nil {
        t.Errorf("json.Unmarshal err:%v", err)
        return
    }
    sdpTrimmedSpace := strings.ReplaceAll(sdp, "\n", "")
    sdpTrimmedSpace = strings.ReplaceAll(sdpTrimmedSpace, " ", "")

    msgSdp := string(*MessageDataOfferRsp.Sdp)
    if msgSdp != sdpTrimmedSpace {
        t.Errorf("MessageDataOfferRsp.Sdp != sdp, %v != %v",
            msgSdp, sdpTrimmedSpace)
        return
    }
    if MessageDataOfferRsp.SocketId != socketId {
        t.Errorf("MessageDataOfferRsp.SocketId != socketId, %v != %v",
            MessageDataOfferRsp.SocketId, socketId)
        return
    }
}
