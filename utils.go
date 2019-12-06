package gosmpp

import (
	"fmt"

	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/pdu"
	"github.com/linxGnu/gosmpp/utils"
)

func marshal(p pdu.PDU) []byte {
	buf := utils.NewBuffer(make([]byte, 0, 64))
	p.Marshal(buf)
	return buf.Bytes()
}

func connect(dialer Dialer, addr string, bindReq *pdu.BindRequest) (c *Connection, err error) {
	conn, err := dialer(addr)
	if err != nil {
		return
	}

	c = NewConnection(conn)

	// send binding request
	_, err = c.Write(marshal(bindReq))
	if err != nil {
		_ = conn.Close()
		return
	}

	// catching response
	var (
		p    pdu.PDU
		resp *pdu.BindResp
	)

loop:
	for {
		if p, err = pdu.Parse(c); err != nil {
			_ = conn.Close()
			return
		}

		switch pd := p.(type) {
		case *pdu.BindResp:
			resp = pd
			break loop
		}
	}

	if resp.CommandStatus != data.ESME_ROK {
		err = fmt.Errorf("Binding error. Command status: [%d]. Please refer to: https://github.com/linxGnu/gosmpp/blob/master/data/pkg.go for more detail about this status code", resp.CommandStatus)
		_ = conn.Close()
	} else {
		c.systemID = resp.SystemID
	}

	return
}
