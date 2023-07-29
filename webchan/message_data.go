package webchan

import "encoding/json"

type MessageDataJoin struct {
    RoomId string `json:"roomId"`
}

// func (m *MessageDataJoin) MarshalJSON() ([]byte, error) {
//     return json.Marshal(m)
// }

// func (m *MessageDataJoin) UnmarshalJSON(data []byte) error {
//     return json.Unmarshal(data, m)
// }

type MessageDataNewPeer struct {
    SocketId string `json:"socketId"`
}

// func (m *MsgDataNewPeer) MarshalJSON() ([]byte, error) {
//     return json.Marshal(m.SocketId)
// }

// func (m *MsgDataNewPeer) UnmarshalJSON(data []byte) error {
//     return json.Unmarshal(data, m)
// }

type MessageDataPeers struct {
    Connections []string `json:"connections"`
    You         string   `json:"you"`
}

// func (m *MsgDataPeers) MarshalJSON() ([]byte, error) {
//     return json.Marshal(m)
// }

type MessageDataICECandidate struct {
    Candidate json.RawMessage `json:"candidate"`
    Id        json.RawMessage `json:"id"`
    Label     json.RawMessage `json:"label"`
    SocketId  string          `json:"socketId"`
}

type MessageDataICECandidateRsp struct {
    Candidate     json.RawMessage `json:"candidate"`
    Id            json.RawMessage `json:"id"`
    Label         json.RawMessage `json:"label"`
    SocketId      string          `json:"socketId"`
    SdpMLineIndex json.RawMessage `json:"sdpMLineIndex"`
}

type MessageDataOffer struct {
    Sdp      *json.RawMessage `json:"sdp"`
    SocketId string           `json:"socketId"`
}

type MessageDataOfferRsp struct {
    Sdp      *json.RawMessage `json:"sdp"`
    SocketId string           `json:"socketId"`
}

type MessageDataAnswer MessageDataOffer
type MessageDataAnswerRsp MessageDataOfferRsp
