package gosmpp

import (
	"fmt"

	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/pdu"
)

func marshal(p pdu.PDU) []byte {
	buf := pdu.NewBuffer(make([]byte, 0, 64))
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

	for {
		if p, err = pdu.Parse(c); err != nil {
			_ = conn.Close()
			return
		}

		if pd, ok := p.(*pdu.BindResp); ok {
			resp = pd
			break
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
