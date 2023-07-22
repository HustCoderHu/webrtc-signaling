package webchan

import (
	"encoding/json"

	uuid "github.com/satori/go.uuid"
)

type Message struct {
    EventName string `json:"eventName"`
    Data      json.RawMessage `json:"data"`
    // Data      map[string]interface{} `json:"data"`
    // Data      string `json:"data"`
}

func (t *Message) Parse(data []byte) error {
    return json.Unmarshal(data, t)
}

type MsgJoin struct {
    RoomId string `json:"room"`
}

func (m *MsgJoin) Parse(data []byte) error {
    return json.Unmarshal(data, m)
}

func GenJsonRspNewPeer(selfUuid uuid.UUID) ([]byte, error) {
    msg := map[string]interface{} {
        "eventName": "_new_peer",
        "data": map[string]string {
            "socketId": selfUuid.String(),
        },
    }
    return json.Marshal(msg)
}

func GenJsonRspPeers(membersUuids []uuid.UUID, selfUuid uuid.UUID) ([]byte, error) {
    var connections []string = nil
    if membersUuids != nil {
        connections = make([]string, len(membersUuids))
        for i, v := range membersUuids {
            connections[i] = v.String()
        }
    }

    msg := map[string]interface{} {
        "eventName": "_peers",
        "data": map[string]interface{} {
            "connections": connections,
            "you": selfUuid.String(),
        },
    }
    return json.Marshal(msg)
}

type MsgICECandidate struct {
    Candidate string `json:"candidate"`
    Id        string `json:"id"`
    Label     int `json:"label"`
    SocketId  string `json:"socketId"`
}

func (m *MsgICECandidate) Parse(data []byte) error {
    return json.Unmarshal(data, m)
}

func GenJsonRspAgainstICECandidate(msgIceCandidate *MsgICECandidate,
      selfUuid uuid.UUID) ([]byte, error) {
    msg := map[string]interface{} {
        "eventName": "_ice_candidate",
        "data": map[string]interface{} {
            "candidate": msgIceCandidate.Candidate,
            "id":        msgIceCandidate.Id,
            "label":     msgIceCandidate.Label,
            "sdpMLineIndex": msgIceCandidate.Label,
            "socketId":  selfUuid.String(),
        },
    }
    return json.Marshal(msg)
}

type MsgOffer struct {
    Sdp string `json:"sdp"`
}

func GetJsonRspAgainstOffer(msgOffer *MsgOffer, selfUuid uuid.UUID) ([]byte, error) {
    msg := map[string]interface{} {
        "eventName": "_offer",
        "data": map[string]string {
            "sdp": msgOffer.Sdp,
            "socketId": selfUuid.String(),
        },
    }
    return json.Marshal(msg)
}

type MsgAnswer struct {
    Sdp string `json:"sdp"`
}

func GetJsonRspAgainstAnswer(msgAnswer *MsgAnswer, selfUuid uuid.UUID) ([]byte, error) {
    msg := map[string]interface{} {
        "eventName": "_answer",
        "data": map[string]string {
            "sdp": msgAnswer.Sdp,
            "socketId": selfUuid.String(),
        },
    }
    return json.Marshal(msg)
}