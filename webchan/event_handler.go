package webchan

import (
	"encoding/json"
	"webrtc-signaling/pkg/logger"

	uuid "github.com/satori/go.uuid"
)

func onJoin(s *Server, m IMember, msgData []byte) error {
    logger.Info("data: %s from: %s", msgData, m.Info())

    joinMsg := &MessageDataJoin{}
    if err := json.Unmarshal(msgData, joinMsg); err != nil {
        logger.Error("joinMsg.Parse error: %s, data: %s, member: %s",
            err, msgData, m.Info())
    }

    room := s.getOrCreateRoomByRoomId(joinMsg.RoomId, true)
    uuids := room.GetMemberUuids()
    rspJson, err := GenJsonRspPeers(uuids, m.GetUuid())
    if err != nil {
        logger.Error("GenJsonRspPeers error: %s, data: %s, member: %s",
            err, msgData, m.Info())
        return nil
    }
    room.AddMember(m)
    m.OnMsg(rspJson)

    notifyRoomMembers(room, m.GetUuid())
    return nil
}

func notifyRoomMembers(room *Room, selfUuid uuid.UUID) {
    msg, err := GenJsonNewPeer(selfUuid)
    if err != nil {
        logger.Error("GenJsonRspNewPeer error: %s, data: %s, member: %s",
            err, msg, selfUuid)
        return
    }
    room.BroadCastMsgExceptMember(msg, selfUuid)
}

func onIceCandidate(s *Server, m IMember, msgData []byte) error {
    logger.Info("data: %s from: %s", msgData, m.Info())
    iceMsg := &MessageDataICECandidate{}
    if err := json.Unmarshal(msgData, iceMsg); err != nil {
        logger.Error("iceMsg.Parse error: %s, data: %s, member: %s",
            err, msgData, m.Info())
        return err
    }
    msg, err := GenJsonICECandidateRsp(iceMsg, m.GetUuid())
    if err != nil {
        logger.Error(
            "GetJsonRspAgainstICECandidate error: %s, data: %s, member: %s",
            err, msgData, m.Info())
        return err
    }

    targetMember := s.getMemberByUuid(iceMsg.SocketId)
    if targetMember == nil {
        logger.Error("target member not found: %s", iceMsg.SocketId)
        return nil
    }
    targetMember.OnMsg(msg)
    return nil
}

func onOffer(s *Server, m IMember, msgData []byte) error {
    // logger.Info("data: %s from: %s", msgData, m.Info())

    offerMsg := &MessageDataOffer{}
    if err := json.Unmarshal(msgData, offerMsg); err != nil {
        logger.Error("offerMsg.Parse error: %s, data: %s, member: %s",
            err, msgData, m.Info())
        return err
    }

    msg, err := GenJsonOfferRsp(offerMsg, m.GetUuid())
    if err != nil {
        logger.Error(
            "GetJsonRspAgainstOffer error: %s, data: %s, member: %s",
            err, msgData, m.Info())
        return err
    }
    targetMember := s.getMemberByUuid(offerMsg.SocketId)
    if targetMember == nil {
        logger.Warning("target member not found: %s", offerMsg.SocketId)
        return nil
    }
    targetMember.OnMsg(msg)
    return nil
}

func onAnswer(s *Server, m IMember, msgData []byte) error {
    logger.Info("data: %s from: %s", msgData, m.Info())

    answerMsg := &MessageDataAnswer{}
    if err := json.Unmarshal(msgData, answerMsg); err != nil {
        logger.Error("answerMsg.Parse error: %s, data: %s, member: %s",
            err, msgData, m.Info())
        return err
    }

    msg, err := GenJsonAnswerRsp(answerMsg, m.GetUuid())
    if err != nil {
        logger.Error(
            "GetJsonRspAgainstAnswer error: %s, data: %s, member: %s",
            err, msgData, m.Info())
        return err
    }

    targetMember := s.getMemberByUuid(answerMsg.SocketId)
    if targetMember == nil {
        logger.Error("target member not found: %s", answerMsg.SocketId)
        return nil
    }
    targetMember.OnMsg(msg)
    return nil
}

func onInvite(s *Server, m IMember, msgData []byte) error {
    logger.Info("data: %s from: %s", msgData, m.Info())
    return nil
}

func onAck(s *Server, m IMember, msgData []byte) error {
    logger.Info("data: %s from: %s", msgData, m.Info())
    return nil
}
