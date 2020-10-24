package data

import (
	"fmt"
	"sync/atomic"
)

//nolint
const (
	SM_CONNID_LEN        = 16
	SM_MSG_LEN           = 254
	SM_SYSID_LEN         = 16
	SM_MSGID_LEN         = 64
	SM_PASS_LEN          = 9
	SM_DATE_LEN          = 17
	SM_SRVTYPE_LEN       = 6
	SM_SYSTYPE_LEN       = 13
	SM_ADDR_LEN          = 21
	SM_DATA_ADDR_LEN     = 65
	SM_ADDR_RANGE_LEN    = 41
	SM_TYPE_LEN          = 13
	SM_DL_NAME_LEN       = 21
	SM_PARAM_NAME_LEN    = 10
	SM_PARAM_VALUE_LEN   = 10
	SM_MAX_CNT_DEST_ADDR = 254

	// GSM specific, short message must be no larger than 140 octets
	SM_GSM_MSG_LEN = 140

	CONNECTION_CLOSED = 0
	CONNECTION_OPENED = 1

	SM_ACK            = 1
	SM_NO_ACK         = 0
	SM_RESPONSE_ACK   = 0
	SM_RESPONSE_TNACK = 1
	SM_RESPONSE_PNACK = 2

	// Interface_Version
	SMPP_V33 int8 = int8(-0x33)
	SMPP_V34      = byte(0x34)

	// Address_TON
	GSM_TON_UNKNOWN       = byte(0x00)
	GSM_TON_INTERNATIONAL = byte(0x01)
	GSM_TON_NATIONAL      = byte(0x02)
	GSM_TON_NETWORK       = byte(0x03)
	GSM_TON_SUBSCRIBER    = byte(0x04)
	GSM_TON_ALPHANUMERIC  = byte(0x05)
	GSM_TON_ABBREVIATED   = byte(0x06)
	GSM_TON_RESERVED_EXTN = byte(0x07)

	// Address_NPI
	GSM_NPI_UNKNOWN       = byte(0x00)
	GSM_NPI_E164          = byte(0x01)
	GSM_NPI_ISDN          = GSM_NPI_E164
	GSM_NPI_X121          = byte(0x03)
	GSM_NPI_TELEX         = byte(0x04)
	GSM_NPI_LAND_MOBILE   = byte(0x06)
	GSM_NPI_NATIONAL      = byte(0x08)
	GSM_NPI_PRIVATE       = byte(0x09)
	GSM_NPI_ERMES         = byte(0x0A)
	GSM_NPI_INTERNET      = byte(0x0E)
	GSM_NPI_WAP_CLIENT_ID = byte(0x12)
	GSM_NPI_RESERVED_EXTN = byte(0x0F)

	// Service_Type
	SERVICE_NULL string = ""
	SERVICE_CMT  string = "CMT"
	SERVICE_CPT  string = "CPT"
	SERVICE_VMN  string = "VMN"
	SERVICE_VMA  string = "VMA"
	SERVICE_WAP  string = "WAP"
	SERVICE_USSD string = "USSD"

	SMPP_PROTOCOL                 = byte(1)
	SMPPP_PROTOCOL                = byte(2)
	SM_SERVICE_MOBILE_TERMINATED  = byte(0)
	SM_SERVICE_MOBILE_ORIGINATED  = byte(1)
	SM_SERVICE_MOBILE_TRANSCEIVER = byte(2)

	// State of message at SMSC
	SM_STATE_EN_ROUTE      = 1 // default state for messages in transit
	SM_STATE_DELIVERED     = 2 // message is delivered
	SM_STATE_EXPIRED       = 3 // validity period expired
	SM_STATE_DELETED       = 4 // message has been deleted
	SM_STATE_UNDELIVERABLE = 5 // undeliverable
	SM_STATE_ACCEPTED      = 6 // message is in accepted state
	SM_STATE_INVALID       = 7 // message is in invalid state
	SM_STATE_REJECTED      = 8 // message is in rejected state

	//******************
	// ESMClass Defines
	//******************

	// Messaging Mode
	SM_ESM_DEFAULT        = 0x00 // Default SMSC Mode or Message Type
	SM_DATAGRAM_MODE      = 0x01 // Use one-shot express mode
	SM_FORWARD_MODE       = 0x02 // Do not use
	SM_STORE_FORWARD_MODE = 0x03 // Use store & forward

	// Send/Receive TDMA & CDMA Message Type
	SM_SMSC_DLV_RCPT_TYPE     = 0x04 // Recv Msg contains SMSC delivery receipt
	SM_ESME_DLV_ACK_TYPE      = 0x08 // Send/Recv Msg contains ESME delivery acknowledgement
	SM_ESME_MAN_USER_ACK_TYPE = 0x10 // Send/Recv Msg contains manual/user acknowledgment
	SM_CONV_ABORT_TYPE        = 0x18 // Recv Msg contains conversation abort (Korean CDMA)
	SM_INTMD_DLV_NOTIFY_TYPE  = 0x20 // Recv Msg contains intermediate notification

	// GSM Network features
	SM_NONE_GSM           = 0x00 // No specific features selected
	SM_UDH_GSM            = 0x40 // User Data Header indicator set
	SM_REPLY_PATH_GSM     = 0x80 // Reply path set
	SM_UDH_REPLY_PATH_GSM = 0xC0 // Both UDH & Reply path

	// Optional Parameter Tags, Min and Max Lengths
	// Following are the 2 byte tag and min/max lengths for
	// supported optional parameter (declann)

	OPT_PAR_MSG_WAIT = 2

	// Privacy Indicator
	OPT_PAR_PRIV_IND = 0x0201

	// Source Subaddress
	OPT_PAR_SRC_SUBADDR     = 0x0202
	OPT_PAR_SRC_SUBADDR_MIN = 2
	OPT_PAR_SRC_SUBADDR_MAX = 23

	// Destination Subaddress
	OPT_PAR_DEST_SUBADDR     = 0x0203
	OPT_PAR_DEST_SUBADDR_MIN = 2
	OPT_PAR_DEST_SUBADDR_MAX = 23

	// User Message Reference
	OPT_PAR_USER_MSG_REF = 0x0204

	// User Response Code
	OPT_PAR_USER_RESP_CODE = 0x0205

	// Language Indicator
	OPT_PAR_LANG_IND = 0x020D

	// Source Port
	OPT_PAR_SRC_PORT = 0x020A

	// Destination Port
	OPT_PAR_DST_PORT = 0x020B

	// Concat Msg Ref Num
	OPT_PAR_SAR_MSG_REF_NUM = 0x020C

	// Concat Total Segments
	OPT_PAR_SAR_TOT_SEG = 0x020E

	// Concat Segment Seqnums
	OPT_PAR_SAR_SEG_SNUM = 0x020F

	// SC Interface Version
	OPT_PAR_SC_IF_VER = 0x0210

	// Display Time
	OPT_PAR_DISPLAY_TIME = 0x1201

	// Validity Information
	OPT_PAR_MS_VALIDITY = 0x1204

	// DPF Result
	OPT_PAR_DPF_RES = 0x0420

	// Set DPF
	OPT_PAR_SET_DPF = 0x0421

	// MS Availability Status
	OPT_PAR_MS_AVAIL_STAT = 0x0422

	// Network Error Code
	OPT_PAR_NW_ERR_CODE     = 0x0423
	OPT_PAR_NW_ERR_CODE_MIN = 3
	OPT_PAR_NW_ERR_CODE_MAX = 3

	// Extended int16 Message has no size limit

	// Delivery Failure Reason
	OPT_PAR_DEL_FAIL_RSN = 0x0425

	// More Messages to Follow
	OPT_PAR_MORE_MSGS = 0x0426

	// Message State
	OPT_PAR_MSG_STATE = 0x0427

	// Callback Number
	OPT_PAR_CALLBACK_NUM     = 0x0381
	OPT_PAR_CALLBACK_NUM_MIN = 4
	OPT_PAR_CALLBACK_NUM_MAX = 19

	// Callback Number Presentation  Indicator
	OPT_PAR_CALLBACK_NUM_PRES_IND = 0x0302

	// Callback Number Alphanumeric Tag
	OPT_PAR_CALLBACK_NUM_ATAG     = 0x0303
	OPT_PAR_CALLBACK_NUM_ATAG_MIN = 1
	OPT_PAR_CALLBACK_NUM_ATAG_MAX = 65

	// Number of messages in Mailbox
	OPT_PAR_NUM_MSGS = 0x0304

	// SMS Received Alert
	OPT_PAR_SMS_SIGNAL = 0x1203

	// Message Delivery Alert
	OPT_PAR_ALERT_ON_MSG_DELIVERY = 0x130C

	// ITS Reply Type
	OPT_PAR_ITS_REPLY_TYPE = 0x1380

	// ITS Session Info
	OPT_PAR_ITS_SESSION_INFO = 0x1383

	// USSD Service Op
	OPT_PAR_USSD_SER_OP = 0x0501

	// Priority
	SM_NOPRIORITY = 0
	SM_PRIORITY   = 1

	// Registered delivery
	//   SMSC Delivery Receipt (bits 1 & 0)
	SM_SMSC_RECEIPT_MASK          = byte(0x03)
	SM_SMSC_RECEIPT_NOT_REQUESTED = byte(0x00)
	SM_SMSC_RECEIPT_REQUESTED     = byte(0x01)
	SM_SMSC_RECEIPT_ON_FAILURE    = byte(0x02)
	//   SME originated acknowledgement (bits 3 & 2)
	SM_SME_ACK_MASK               = byte(0x0c)
	SM_SME_ACK_NOT_REQUESTED      = byte(0x00)
	SM_SME_ACK_DELIVERY_REQUESTED = byte(0x04)
	SM_SME_ACK_MANUAL_REQUESTED   = byte(0x08)
	SM_SME_ACK_BOTH_REQUESTED     = byte(0x0c)
	//   Intermediate notification (bit 5)
	SM_NOTIF_MASK          = byte(0x010)
	SM_NOTIF_NOT_REQUESTED = byte(0x000)
	SM_NOTIF_REQUESTED     = byte(0x010)

	// Replace if Present flag
	SM_NOREPLACE = 0
	SM_REPLACE   = 1

	// Destination flag
	SM_DEST_SME_ADDRESS = 1
	SM_DEST_DL_NAME     = 2

	// Higher Layer Message Type
	SM_LAYER_WDP  = 0
	SM_LAYER_WCMP = 1

	// Operation Class
	SM_OPCLASS_DATAGRAM    = 0
	SM_OPCLASS_TRANSACTION = 3

	// Originating MSC Address
	OPT_PAR_ORIG_MSC_ADDR     = -32639 // int16(0x8081)
	OPT_PAR_ORIG_MSC_ADDR_MIN = 1
	OPT_PAR_ORIG_MSC_ADDR_MAX = 24

	// Destination MSC Address
	OPT_PAR_DEST_MSC_ADDR     = -32638 // int16(0x8082)
	OPT_PAR_DEST_MSC_ADDR_MIN = 1
	OPT_PAR_DEST_MSC_ADDR_MAX = 24

	// Unused Tag
	OPT_PAR_UNUSED = 0xffff

	// Destination Address Subunit
	OPT_PAR_DST_ADDR_SUBUNIT = 0x0005

	// Destination Network Type
	OPT_PAR_DST_NW_TYPE = 0x0006

	// Destination Bearer Type
	OPT_PAR_DST_BEAR_TYPE = 0x0007

	// Destination Telematics ID
	OPT_PAR_DST_TELE_ID = 0x0008

	// Source Address Subunit
	OPT_PAR_SRC_ADDR_SUBUNIT = 0x000D

	// Source Network Type
	OPT_PAR_SRC_NW_TYPE = 0x000E

	// Source Bearer Type
	OPT_PAR_SRC_BEAR_TYPE = 0x000F

	// Source Telematics ID
	OPT_PAR_SRC_TELE_ID = 0x0010

	// QOS Time to Live
	OPT_PAR_QOS_TIME_TO_LIVE     = 0x0017
	OPT_PAR_QOS_TIME_TO_LIVE_MIN = 1
	OPT_PAR_QOS_TIME_TO_LIVE_MAX = 4

	// Payload Type
	OPT_PAR_PAYLOAD_TYPE = 0x0019

	// Additional Status Info Text
	OPT_PAR_ADD_STAT_INFO     = 0x001D
	OPT_PAR_ADD_STAT_INFO_MIN = 1
	OPT_PAR_ADD_STAT_INFO_MAX = 256

	// Receipted Message ID
	OPT_PAR_RECP_MSG_ID     = 0x001E
	OPT_PAR_RECP_MSG_ID_MIN = 1
	OPT_PAR_RECP_MSG_ID_MAX = 65

	// Message Payload
	OPT_PAR_MSG_PAYLOAD     = 0x0424
	OPT_PAR_MSG_PAYLOAD_MIN = 1
	OPT_PAR_MSG_PAYLOAD_MAX = 1500

	// User Data Header
	UDH_CONCAT_MSG_8_BIT_REF  = byte(0x00)
	UDH_CONCAT_MSG_16_BIT_REF = byte(0x08)

	/**
	 * @deprecated As of version 1.3 of the library there are defined
	 * new encoding constants for base set of encoding supported by Java Runtime.
	 * The <code>CHAR_ENC</code> is replaced by <code>ENC_ASCII</code>
	 * and redefined in this respect.
	 */

	DFLT_MSGID         string = ""
	DFLT_MSG           string = ""
	DFLT_SRVTYPE       string = ""
	DFLT_SYSID         string = ""
	DFLT_PASS          string = ""
	DFLT_SYSTYPE       string = ""
	DFLT_ADDR_RANGE    string = ""
	DFLT_DATE          string = ""
	DFLT_ADDR          string = ""
	DFLT_MSG_STATE     byte   = 0
	DFLT_ERR           byte   = 0
	DFLT_SCHEDULE      string = ""
	DFLT_VALIDITY      string = ""
	DFLT_REG_DELIVERY         = SM_SMSC_RECEIPT_NOT_REQUESTED | SM_SME_ACK_NOT_REQUESTED | SM_NOTIF_NOT_REQUESTED
	DFLT_DFLTMSGID            = byte(0)
	DFLT_MSG_LEN              = byte(0)
	DFLT_ESM_CLASS            = byte(0)
	DFLT_DATA_CODING          = byte(0)
	DFLT_PROTOCOLID           = byte(0)
	DFLT_PRIORITY_FLAG        = byte(0)
	DFTL_REPLACE_IFP          = byte(0)
	DFLT_DL_NAME       string = ""
	DFLT_GSM_TON              = GSM_TON_UNKNOWN
	DFLT_GSM_NPI              = GSM_NPI_UNKNOWN
	DFLT_DEST_FLAG            = byte(0) // not set
	MAX_PDU_LEN               = 64 << 10

	PDU_HEADER_SIZE = 16 // 4 integers
	TLV_HEADER_SIZE = 4  // 2 int16s: tag & length

	// all times in milliseconds
	RECEIVER_TIMEOUT           int64 = 60000
	CONNECTION_RECEIVE_TIMEOUT int64 = 10000
	UNBIND_RECEIVE_TIMEOUT     int64 = 5000
	CONNECTION_SEND_TIMEOUT    int64 = 20000
	COMMS_TIMEOUT              int64 = 60000
	QUEUE_TIMEOUT              int64 = 10000
	ACCEPT_TIMEOUT             int64 = 60000

	RECEIVE_BLOCKING int64 = -1

	MAX_VALUE_PORT     = 65535
	MIN_VALUE_PORT     = 100
	MIN_LENGTH_ADDRESS = 7
)

var (
	// ErrNotImplSplitterInterface indicates that encoding does not support Splitter interface
	ErrNotImplSplitterInterface = fmt.Errorf("Encoding not implementing Splitter interface")
)

var defaultTon atomic.Value
var defaultNpi atomic.Value

func init() {
	defaultTon.Store(DFLT_GSM_TON)
	defaultNpi.Store(DFLT_GSM_NPI)
}

// SetDefaultTon set default ton.
func SetDefaultTon(dfltTon byte) {
	defaultTon.Store(dfltTon)
}

// GetDefaultTon get default ton.
func GetDefaultTon() byte {
	return defaultTon.Load().(byte)
}

// SetDefaultNpi set default npi.
func SetDefaultNpi(dfltNpi byte) {
	defaultNpi.Store(dfltNpi)
}

// GetDefaultNpi get default npi.
func GetDefaultNpi() byte {
	return defaultNpi.Load().(byte)
}
