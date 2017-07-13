package PDU

import (
	"github.com/linxGnu/gosmpp/Data"
	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/PDU/Common"
	"github.com/linxGnu/gosmpp/Utils"
)

type DistributionList struct {
	Common.ByteData
	dlName string
}

func NewDistributionList() *DistributionList {
	a := &DistributionList{}
	a.Construct()

	return a
}

func NewDistributionListWithDlName(dlName string) (*DistributionList, *Exception.Exception) {
	a := NewDistributionList()
	err := a.SetDlName(dlName)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (c *DistributionList) Construct() {
	defer c.SetRealReference(c)
	c.ByteData.Construct()

	c.dlName = Data.DFLT_DL_NAME
}

func (c *DistributionList) SetData(buf *Utils.ByteBuffer) *Exception.Exception {
	if buf == nil || buf.Buffer == nil {
		return Exception.NewExceptionFromStr("DestinationAddress: SetData with nil buffer")
	}

	dlname, err := buf.Read_CString()
	if err != nil {
		return err
	}

	return c.SetDlName(dlname)
}

func (c *DistributionList) GetData() (*Utils.ByteBuffer, *Exception.Exception) {
	buf := Utils.NewBuffer(make([]byte, 0, 32))
	return buf, buf.Write_CString(c.GetDlName())
}

func (c *DistributionList) SetDlName(dln string) *Exception.Exception {
	err := c.CheckCStringMax(dln, int(Data.SM_DL_NAME_LEN))
	if err != nil {
		return err
	}

	c.dlName = dln
	return nil
}

func (c *DistributionList) GetDlName() string {
	return c.dlName
}
