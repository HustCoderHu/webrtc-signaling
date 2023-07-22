package webchan

import (
	"webrtc-signaling/pkg/logger"

	uuid "github.com/satori/go.uuid"
)

func onJoin(s *Server, m *Member, msgData []byte) error {
    logger.Info("data: %s from: %s", msgData, m.Info())

    joinMsg := MsgJoin{}
    if err := joinMsg.Parse(msgData); err != nil {
        logger.Error("joinMsg.Parse error: %s, data: %s, member: %s",
            err, msgData, m.Info())
    }

    room, ok := s.rooms[joinMsg.RoomId]
    var rspJson []byte = nil
    var err error
    var uuids []uuid.UUID = nil
    if ok {
        if room.members != nil {
            uuids = make([]uuid.UUID, 0, len(room.members))
            for _, member := range room.members {
                if member.uuid == m.uuid {
                    continue
                }
                uuids = append(uuids, member.uuid)
            }
        } else {
            room.members = map[uuid.UUID]*Member { m.uuid: m }
        }
    } else {
        // create room
        s.rooms[joinMsg.RoomId] = &Room {
            roomId: joinMsg.RoomId,
            members: map[uuid.UUID]*Member { m.uuid: m },
        }
    }

    rspJson, err = GenJsonRspPeers(uuids, m.uuid)
    if err != nil {
        logger.Error("GenJsonRspPeers error: %s, data: %s, member: %s",
            err, msgData, m.Info())
        return nil
    }
    m.OnMsg(rspJson)
    return nil
}

func onIceCandidate(s *Server, m *Member, msgData []byte) error {
    logger.Info("data: %s from: %s", msgData, m.Info())

    // GenJsonRspAgainstICECandidate()
    return nil
}

func onOffer(s *Server, m *Member, msgData []byte) error {
    logger.Info("data: %s from: %s", msgData, m.Info())
    return nil
}

func onAnswer(s *Server, m *Member, msgData []byte) error {
    logger.Info("data: %s from: %s", msgData, m.Info())
    return nil
}

func onInvite(s *Server, m *Member, msgData []byte) error {
    logger.Info("data: %s from: %s", msgData, m.Info())
    return nil
}

func onAck(s *Server, m *Member, msgData []byte) error {
    logger.Info("data: %s from: %s", msgData, m.Info())
    return nil
}
