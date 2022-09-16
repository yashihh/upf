package pfcp

import (
	"net"

	"github.com/wmnsk/go-pfcp/ie"
	"github.com/wmnsk/go-pfcp/message"
)

func (s *PfcpServer) handleHeartbeatRequest(req *message.HeartbeatRequest, addr net.Addr) {
	s.log.Infof("handleHeartbeatRequest SEQ[%#x]", req.SequenceNumber)

	rsp := message.NewHeartbeatResponse(
		req.SequenceNumber,
		ie.NewRecoveryTimeStamp(s.recoveryTime),
	)

	err := s.sendRspTo(rsp, addr)
	if err != nil {
		s.log.Errorln(err)
		return
	}
}
