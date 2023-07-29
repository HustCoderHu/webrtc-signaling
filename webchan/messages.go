package webchan

import (
	"encoding/json"
	"webrtc-signaling/pkg/logger"

	uuid "github.com/satori/go.uuid"
)

type Message struct {
    EventName string          `json:"eventName"`
    Data      json.RawMessage `json:"data"`
    // Data      map[string]interface{} `json:"data"`
    // Data      string `json:"data"`
}

func (t *Message) Parse(data []byte) error {
    return json.Unmarshal(data, t)
}

func GenJsonNewPeer(selfUuid uuid.UUID) ([]byte, error) {
    msgDataNewPeers := &MessageDataNewPeer{
        SocketId: selfUuid.String(),
    }
    data, err := json.Marshal(msgDataNewPeers)
    if err != nil {
        return nil, err
    }

    msg := &Message{
        EventName: "_new_peer",
        Data:      data,
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

    msgDataRspPeers := &MessageDataPeers{
        Connections: connections,
        You:         selfUuid.String(),
    }
    data, err := json.Marshal(msgDataRspPeers)
    if err != nil {
        return nil, err
    }

    msg := &Message{
        EventName: "_peers",
        Data:      data,
    }
    return json.Marshal(msg)
}

func GenJsonICECandidateRsp(msgIceCandidate *MessageDataICECandidate,
    selfUuid uuid.UUID) ([]byte, error) {

    msgIceCandidateRsp := &MessageDataICECandidateRsp{
        Candidate:     msgIceCandidate.Candidate,
        Id:            msgIceCandidate.Id,
        Label:         msgIceCandidate.Label,
        SdpMLineIndex: msgIceCandidate.Label,
        SocketId:      selfUuid.String(),
    }
    data, err := json.Marshal(msgIceCandidateRsp)
    if err != nil {
        return nil, err
    }

    msg := &Message{
        EventName: "_ice_candidate",
        Data:      data,
    }
    return json.Marshal(msg)
}

func GenJsonOfferRsp(msgOffer *MessageDataOffer, selfUuid uuid.UUID) ([]byte, error) {
    msgOfferRsp := &MessageDataOfferRsp{
        Sdp:      msgOffer.Sdp,
        SocketId: selfUuid.String(),
    }
    data, err := json.Marshal(msgOfferRsp)
    if err != nil {
        return nil, err
    }

    logger.Info("")

    msg := &Message{
        EventName: "_offer",
        Data:      data,
    }
    return json.Marshal(msg)
}

func GenJsonAnswerRsp(msgAnswer *MessageDataAnswer, selfUuid uuid.UUID) ([]byte, error) {
    msgOfferAns := &MessageDataAnswer{
        Sdp:      msgAnswer.Sdp,
        SocketId: selfUuid.String(),
    }
    data, err := json.Marshal(msgOfferAns)
    if err != nil {
        return nil, err
    }

    msg := &Message{
        EventName: "_answer",
        Data:      data,
    }
    return json.Marshal(msg)
}
