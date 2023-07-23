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

    room := s.GetOrCreateRoomByRoomId(joinMsg.RoomId, true)
    room.AddMember(m)

    uuids := room.GetMemberUuids()
    rspJson, err := GenJsonRspPeers(uuids, m.uuid)
    if err != nil {
        logger.Error("GenJsonRspPeers error: %s, data: %s, member: %s",
            err, msgData, m.Info())
        return nil
    }
    m.OnMsg(rspJson)
    return nil
}

func notifyRoomMembers(room *Room, selfUuid uuid.UUID) {
    msg, err := GenJsonRspNewPeer(selfUuid)
    if err != nil {
        logger.Error("GenJsonRspNewPeer error: %s, data: %s, member: %s",
            err, msg, selfUuid)
        return
    }
    room.BroadCastMsgExceptMember(msg, selfUuid)
}

func onIceCandidate(s *Server, m *Member, msgData []byte) error {
    logger.Info("data: %s from: %s", msgData, m.Info())
    iceMsg := &MsgICECandidate{}
    if err := iceMsg.Parse(msgData); err != nil {
        logger.Error("iceMsg.Parse error: %s, data: %s, member: %s",
            err, msgData, m.Info())
        return err
    }
    msg, err := GenJsonRspAgainstICECandidate(iceMsg, m.uuid)
    if err != nil {
        logger.Error(
            "GetJsonRspAgainstICECandidate error: %s, data: %s, member: %s",
            err, msgData, m.Info())
        return err
    }

    targetMember := s.GetMemberByUuid(uuid.UUID(iceMsg.SocketId))
    if targetMember == nil {
        logger.Error("target member not found: %s", iceMsg.SocketId)
        return nil
    }
    targetMember.OnMsg(msg)
    return nil
}

func onOffer(s *Server, m *Member, msgData []byte) error {
    logger.Info("data: %s from: %s", msgData, m.Info())

    offerMsg := &MsgOffer{}
    if err := offerMsg.Parse(msgData); err != nil {
        logger.Error("offerMsg.Parse error: %s, data: %s, member: %s",
            err, msgData, m.Info())
        return err
    }

    msg, err := GenJsonRspAgainstOffer(offerMsg, m.uuid)
    if err != nil {
        logger.Error(
            "GetJsonRspAgainstOffer error: %s, data: %s, member: %s",
            err, msgData, m.Info())
        return err
    }
    targetMember := s.GetMemberByUuid(uuid.UUID(offerMsg.SocketId))
    if targetMember == nil {
        logger.Error("target member not found: %s", offerMsg.SocketId)
        return nil
    }
    targetMember.OnMsg(msg)
    return nil
}

func onAnswer(s *Server, m *Member, msgData []byte) error {
    logger.Info("data: %s from: %s", msgData, m.Info())

    answerMsg := &MsgAnswer{}
    if err := answerMsg.Parse(msgData); err != nil {
        logger.Error("answerMsg.Parse error: %s, data: %s, member: %s",
            err, msgData, m.Info())
        return err
    }

    msg, err := GenJsonRspAgainstAnswer(answerMsg, m.uuid)
    if err != nil {
        logger.Error(
            "GetJsonRspAgainstAnswer error: %s, data: %s, member: %s",
            err, msgData, m.Info())
        return err
    }
    
    targetMember := s.GetMemberByUuid(uuid.UUID(answerMsg.SocketId))
    if targetMember == nil {
        logger.Error("target member not found: %s", answerMsg.SocketId)
        return nil
    }
    targetMember.OnMsg(msg)
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
