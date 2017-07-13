package PDU

import (
	"github.com/linxGnu/gosmpp/Data"
	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/PDU/Common"
	"github.com/linxGnu/gosmpp/Utils"
)

type DestinationAddress struct {
	Common.ByteData
	destFlag   byte
	theAddress Common.IByteData
}

func NewDestinationAddress() *DestinationAddress {
	a := &DestinationAddress{}
	a.Construct()

	return a
}

func NewDestinationAddressWithAddr(addr string) (*DestinationAddress, *Exception.Exception) {
	a := NewDestinationAddress()

	address, err := NewAddressWithAddr(addr)
	if err != nil {
		return nil, err
	}

	a.SetAddress(address)
	return a, nil
}

func (c *DestinationAddress) Construct() {
	defer c.SetRealReference(c)
	c.ByteData.Construct()

	c.destFlag = Data.DFLT_DEST_FLAG
}

func (c *DestinationAddress) SetData(buf *Utils.ByteBuffer) *Exception.Exception {
	if buf == nil || buf.Buffer == nil {
		return Exception.NewExceptionFromStr("DestinationAddress: SetData with nil buffer")
	}

	t, err := buf.Read_Byte()
	if err != nil {
		return err
	}

	c.destFlag = t
	switch c.destFlag {
	case byte(Data.SM_DEST_SME_ADDRESS):
		addr := NewAddress()
		err = addr.SetData(buf)
		if err != nil {
			return err
		}
		c.SetAddress(addr)

	case byte(Data.SM_DEST_DL_NAME):
		dl := NewDistributionList()
		err = dl.SetData(buf)
		if err != nil {
			return err
		}
		c.SetDistributionList(dl)

	default:
		return Exception.NewExceptionFromStr("DestinationAddress: wrong address")
	}

	return nil
}

func (c *DestinationAddress) GetData() (*Utils.ByteBuffer, *Exception.Exception) {
	if c.HasValue() {
		buf := Utils.NewBuffer(make([]byte, 0, 16))

		if err := buf.Write_Byte(c.destFlag); err != nil {
			return nil, err
		}

		if c.IsAddress() {
			addr := c.GetAddress()
			if addr == nil {
				return nil, Exception.ValueNotSetException
			}

			tmp, err := addr.GetData()
			if err != nil {
				return nil, err
			}

			return buf, buf.Write_Buffer(tmp)
		} else if c.IsDistributionList() {
			dl := c.GetDistributionList()
			if dl == nil {
				return nil, Exception.ValueNotSetException
			}

			tmp, err := dl.GetData()
			if err != nil {
				return nil, err
			}

			return buf, buf.Write_Buffer(tmp)
		}

		return buf, nil
	}

	return nil, Exception.ValueNotSetException
}

func (c *DestinationAddress) GetAddress() *Address {
	if c.IsAddress() {
		return c.theAddress.(*Address)
	}

	return nil
}

func (c *DestinationAddress) SetAddress(dat *Address) {
	c.destFlag = byte(Data.SM_DEST_SME_ADDRESS)
	c.theAddress = dat
}

func (c *DestinationAddress) GetDistributionList() *DistributionList {
	if c.IsDistributionList() {
		return c.theAddress.(*DistributionList)
	}

	return nil
}

func (c *DestinationAddress) SetDistributionList(dl *DistributionList) {
	c.destFlag = byte(Data.SM_DEST_DL_NAME)
	c.theAddress = dl
}

func (c *DestinationAddress) GetDestFlag() byte {
	return c.destFlag
}

func (c *DestinationAddress) HasValue() bool {
	return c.destFlag != byte(Data.DFLT_DEST_FLAG)
}

func (c *DestinationAddress) IsAddress() bool {
	return c.destFlag == byte(Data.SM_DEST_SME_ADDRESS)
}

func (c *DestinationAddress) IsDistributionList() bool {
	return c.destFlag == byte(Data.SM_DEST_DL_NAME)
}
