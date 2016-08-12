package Data

import (
	"sync"
	"time"
)

const (
	SM_CONNID_LEN        int32 = 16
	SM_MSG_LEN           int32 = 254
	SM_SYSID_LEN         int32 = 16
	SM_MSGID_LEN         int32 = 64
	SM_PASS_LEN          int32 = 9
	SM_DATE_LEN          int32 = 17
	SM_SRVTYPE_LEN       int32 = 6
	SM_SYSTYPE_LEN       int32 = 13
	SM_ADDR_LEN          int32 = 21
	SM_DATA_ADDR_LEN     int32 = 65
	SM_ADDR_RANGE_LEN    int32 = 41
	SM_TYPE_LEN          int32 = 13
	SM_DL_NAME_LEN       int32 = 21
	SM_PARAM_NAME_LEN    int32 = 10
	SM_PARAM_VALUE_LEN   int32 = 10
	SM_MAX_CNT_DEST_ADDR int32 = 254

	CONNECTION_CLOSED int32 = 0
	CONNECTION_OPENED int32 = 1

	SM_ACK            int32 = 1
	SM_NO_ACK         int32 = 0
	SM_RESPONSE_ACK   int32 = 0
	SM_RESPONSE_TNACK int32 = 1
	SM_RESPONSE_PNACK int32 = 2

	//SMPP Command Set
	GENERIC_NACK          int32 = -2147483648
	BIND_RECEIVER         int32 = 0x00000001
	BIND_RECEIVER_RESP    int32 = -2147483647
	BIND_TRANSMITTER      int32 = 0x00000002
	BIND_TRANSMITTER_RESP int32 = -2147483646
	QUERY_SM              int32 = 0x00000003
	QUERY_SM_RESP         int32 = -2147483645
	SUBMIT_SM             int32 = 0x00000004
	SUBMIT_SM_RESP        int32 = -2147483644
	DELIVER_SM            int32 = 0x00000005
	DELIVER_SM_RESP       int32 = -2147483643
	UNBIND                int32 = 0x00000006
	UNBIND_RESP           int32 = -2147483642
	REPLACE_SM            int32 = 0x00000007
	REPLACE_SM_RESP       int32 = -2147483641
	CANCEL_SM             int32 = 0x00000008
	CANCEL_SM_RESP        int32 = -2147483640
	BIND_TRANSCEIVER      int32 = 0x00000009
	BIND_TRANSCEIVER_RESP int32 = -2147483639
	OUTBIND               int32 = 0x0000000B
	ENQUIRE_LINK          int32 = 0x00000015
	ENQUIRE_LINK_RESP     int32 = -2147483627
	SUBMIT_MULTI          int32 = 0x00000021
	SUBMIT_MULTI_RESP     int32 = -2147483615
	ALERT_NOTIFICATION    int32 = 0x00000102
	DATA_SM               int32 = 0x00000103
	DATA_SM_RESP          int32 = -2147483389

	//Command_Status Error Codes
	ESME_ROK           int32 = 0x00000000
	ESME_RINVMSGLEN    int32 = 0x00000001
	ESME_RINVCMDLEN    int32 = 0x00000002
	ESME_RINVCMDID     int32 = 0x00000003
	ESME_RINVBNDSTS    int32 = 0x00000004
	ESME_RALYBND       int32 = 0x00000005
	ESME_RINVPRTFLG    int32 = 0x00000006
	ESME_RINVREGDLVFLG int32 = 0x00000007
	ESME_RSYSERR       int32 = 0x00000008
	ESME_RINVSRCADR    int32 = 0x0000000A
	ESME_RINVDSTADR    int32 = 0x0000000B
	ESME_RINVMSGID     int32 = 0x0000000C
	ESME_RBINDFAIL     int32 = 0x0000000D
	ESME_RINVPASWD     int32 = 0x0000000E
	ESME_RINVSYSID     int32 = 0x0000000F
	ESME_RCANCELFAIL   int32 = 0x00000011
	ESME_RREPLACEFAIL  int32 = 0x00000013
	ESME_RMSGQFUL      int32 = 0x00000014
	ESME_RINVSERTYP    int32 = 0x00000015

	ESME_RADDCUSTFAIL  int32 = 0x00000019 // Failed to Add Customer
	ESME_RDELCUSTFAIL  int32 = 0x0000001A // Failed to delete Customer
	ESME_RMODCUSTFAIL  int32 = 0x0000001B // Failed to modify customer
	ESME_RENQCUSTFAIL  int32 = 0x0000001C // Failed to Enquire Customer
	ESME_RINVCUSTID    int32 = 0x0000001D // Invalid Customer ID
	ESME_RINVCUSTNAME  int32 = 0x0000001F // Invalid Customer Name
	ESME_RINVCUSTADR   int32 = 0x00000021 // Invalid Customer Address
	ESME_RINVADR       int32 = 0x00000022 // Invalid Address
	ESME_RCUSTEXIST    int32 = 0x00000023 // Customer Exists
	ESME_RCUSTNOTEXIST int32 = 0x00000024 // Customer does not exist
	ESME_RADDDLFAIL    int32 = 0x00000026 // Failed to Add DL
	ESME_RMODDLFAIL    int32 = 0x00000027 // Failed to modify DL
	ESME_RDELDLFAIL    int32 = 0x00000028 // Failed to Delete DL
	ESME_RVIEWDLFAIL   int32 = 0x00000029 // Failed to View DL
	ESME_RLISTDLSFAIL  int32 = 0x00000030 // Failed to list DLs
	ESME_RPARAMRETFAIL int32 = 0x00000031 // Param Retrieve Failed
	ESME_RINVPARAM     int32 = 0x00000032 // Invalid Param

	ESME_RINVNUMDESTS int32 = 0x00000033
	ESME_RINVDLNAME   int32 = 0x00000034

	ESME_RINVDLMEMBDESC int32 = 0x00000035 // Invalid DL Member Description
	ESME_RINVDLMEMBTYP  int32 = 0x00000038 // Invalid DL Member Type
	ESME_RINVDLMODOPT   int32 = 0x00000039 // Invalid DL Modify Option

	ESME_RINVDESTFLAG int32 = 0x00000040
	ESME_RINVSUBREP   int32 = 0x00000042
	ESME_RINVESMCLASS int32 = 0x00000043
	ESME_RCNTSUBDL    int32 = 0x00000044
	ESME_RSUBMITFAIL  int32 = 0x00000045
	ESME_RINVSRCTON   int32 = 0x00000048
	ESME_RINVSRCNPI   int32 = 0x00000049
	ESME_RINVDSTTON   int32 = 0x00000050
	ESME_RINVDSTNPI   int32 = 0x00000051
	ESME_RINVSYSTYP   int32 = 0x00000053
	ESME_RINVREPFLAG  int32 = 0x00000054
	ESME_RINVNUMMSGS  int32 = 0x00000055
	ESME_RTHROTTLED   int32 = 0x00000058

	ESME_RPROVNOTALLWD int32 = 0x00000059 // Provisioning Not Allowed

	ESME_RINVSCHED    int32 = 0x00000061
	ESME_RINVEXPIRY   int32 = 0x00000062
	ESME_RINVDFTMSGID int32 = 0x00000063
	ESME_RX_T_APPN    int32 = 0x00000064
	ESME_RX_P_APPN    int32 = 0x00000065
	ESME_RX_R_APPN    int32 = 0x00000066
	ESME_RQUERYFAIL   int32 = 0x00000067

	ESME_RINVPGCUSTID      int32 = 0x00000080 // Paging Customer ID Invalid No such subscriber
	ESME_RINVPGCUSTIDLEN   int32 = 0x00000081 // Paging Customer ID length Invalid
	ESME_RINVCITYLEN       int32 = 0x00000082 // City Length Invalid
	ESME_RINVSTATELEN      int32 = 0x00000083 // State Length Invalid
	ESME_RINVZIPPREFIXLEN  int32 = 0x00000084 // Zip Prefix Length Invalid
	ESME_RINVZIPPOSTFIXLEN int32 = 0x00000085 // Zip Postfix Length Invalid
	ESME_RINVMINLEN        int32 = 0x00000086 // MIN Length Invalid
	ESME_RINVMIN           int32 = 0x00000087 // MIN Invalid (i.e. No such MIN)
	ESME_RINVPINLEN        int32 = 0x00000088 // PIN Length Invalid
	ESME_RINVTERMCODELEN   int32 = 0x00000089 // Terminal Code Length Invalid
	ESME_RINVCHANNELLEN    int32 = 0x0000008A // Channel Length Invalid
	ESME_RINVCOVREGIONLEN  int32 = 0x0000008B // Coverage Region Length Invalid
	ESME_RINVCAPCODELEN    int32 = 0x0000008C // Cap Code Length Invalid
	ESME_RINVMDTLEN        int32 = 0x0000008D // Message delivery time Length Invalid
	ESME_RINVPRIORMSGLEN   int32 = 0x0000008E // Priority Message Length Invalid
	ESME_RINVPERMSGLEN     int32 = 0x0000008F // Periodic Messages Length Invalid
	ESME_RINVPGALERTLEN    int32 = 0x00000090 // Paging Alerts Length Invalid
	ESME_RINVSMUSERLEN     int32 = 0x00000091 // int16 Message User Group Length Invalid
	ESME_RINVRTDBLEN       int32 = 0x00000092 // Real Time Data broadcasts Length Invalid
	ESME_RINVREGDELLEN     int32 = 0x00000093 // Registered Delivery Lenght Invalid
	ESME_RINVMSGDISTLEN    int32 = 0x00000094 // Message Distribution Lenght Invalid
	ESME_RINVPRIORMSG      int32 = 0x00000095 // Priority Message Length Invalid
	ESME_RINVMDT           int32 = 0x00000096 // Message delivery time Invalid
	ESME_RINVPERMSG        int32 = 0x00000097 // Periodic Messages Invalid
	ESME_RINVMSGDIST       int32 = 0x00000098 // Message Distribution Invalid
	ESME_RINVPGALERT       int32 = 0x00000099 // Paging Alerts Invalid
	ESME_RINVSMUSER        int32 = 0x0000009A // int16 Message User Group Invalid
	ESME_RINVRTDB          int32 = 0x0000009B // Real Time Data broadcasts Invalid
	ESME_RINVREGDEL        int32 = 0x0000009C // Registered Delivery Invalid
	//public static final int32 ESME_RINVOPTPARSTREAM = 0x0000009D  // KIF IW Field out of data
	//public static final int32 ESME_ROPTPARNOTALLWD = 0x0000009E  // Optional Parameter not allowed
	ESME_RINVOPTPARLEN int32 = 0x0000009F // Invalid Optional Parameter Length

	ESME_RINVOPTPARSTREAM int32 = 0x000000C0
	ESME_ROPTPARNOTALLWD  int32 = 0x000000C1
	ESME_RINVPARLEN       int32 = 0x000000C2
	ESME_RMISSINGOPTPARAM int32 = 0x000000C3
	ESME_RINVOPTPARAMVAL  int32 = 0x000000C4
	ESME_RDELIVERYFAILURE int32 = 0x000000FE
	ESME_RUNKNOWNERR      int32 = 0x000000FF

	ESME_LAST_ERROR int32 = 0x0000012C // the value of the last error code

	//Interface_Version
	SMPP_V33 int8 = int8(-0x33)
	SMPP_V34 byte = byte(0x34)

	//Address_TON
	GSM_TON_UNKNOWN       byte = byte(0x00)
	GSM_TON_INTERNATIONAL byte = byte(0x01)
	GSM_TON_NATIONAL      byte = byte(0x02)
	GSM_TON_NETWORK       byte = byte(0x03)
	GSM_TON_SUBSCRIBER    byte = byte(0x04)
	GSM_TON_ALPHANUMERIC  byte = byte(0x05)
	GSM_TON_ABBREVIATED   byte = byte(0x06)
	GSM_TON_RESERVED_EXTN byte = byte(0x07)

	//Address_NPI
	GSM_NPI_UNKNOWN       byte = byte(0x00)
	GSM_NPI_E164          byte = byte(0x01)
	GSM_NPI_ISDN          byte = GSM_NPI_E164
	GSM_NPI_X121          byte = byte(0x03)
	GSM_NPI_TELEX         byte = byte(0x04)
	GSM_NPI_LAND_MOBILE   byte = byte(0x06)
	GSM_NPI_NATIONAL      byte = byte(0x08)
	GSM_NPI_PRIVATE       byte = byte(0x09)
	GSM_NPI_ERMES         byte = byte(0x0A)
	GSM_NPI_INTERNET      byte = byte(0x0E)
	GSM_NPI_WAP_CLIENT_ID byte = byte(0x12)
	GSM_NPI_RESERVED_EXTN byte = byte(0x0F)

	//Service_Type
	SERVICE_NULL string = ""
	SERVICE_CMT  string = "CMT"
	SERVICE_CPT  string = "CPT"
	SERVICE_VMN  string = "VMN"
	SERVICE_VMA  string = "VMA"
	SERVICE_WAP  string = "WAP"
	SERVICE_USSD string = "USSD"

	SMPP_PROTOCOL                 byte = byte(1)
	SMPPP_PROTOCOL                byte = byte(2)
	SM_SERVICE_MOBILE_TERMINATED  byte = byte(0)
	SM_SERVICE_MOBILE_ORIGINATED  byte = byte(1)
	SM_SERVICE_MOBILE_TRANSCEIVER byte = byte(2)

	// State of message at SMSC
	SM_STATE_EN_ROUTE      int32 = 1 // default state for messages in transit
	SM_STATE_DELIVERED     int32 = 2 // message is delivered
	SM_STATE_EXPIRED       int32 = 3 // validity period expired
	SM_STATE_DELETED       int32 = 4 // message has been deleted
	SM_STATE_UNDELIVERABLE int32 = 5 // undeliverable
	SM_STATE_ACCEPTED      int32 = 6 // message is in accepted state
	SM_STATE_INVALID       int32 = 7 // message is in invalid state
	SM_STATE_REJECTED      int32 = 8 // message is in rejected state

	//******************
	// ESMClass Defines
	//******************

	// Messaging Mode
	SM_ESM_DEFAULT        int32 = 0x00 //Default SMSC Mode or Message Type
	SM_DATAGRAM_MODE      int32 = 0x01 // Use one-shot express mode
	SM_FORWARD_MODE       int32 = 0x02 // Do not use
	SM_STORE_FORWARD_MODE int32 = 0x03 // Use store & forward

	// Send/Receive TDMA & CDMA Message Type
	SM_SMSC_DLV_RCPT_TYPE     int32 = 0x04 // Recv Msg contains SMSC delivery receipt
	SM_ESME_DLV_ACK_TYPE      int32 = 0x08 // Send/Recv Msg contains ESME delivery acknowledgement
	SM_ESME_MAN_USER_ACK_TYPE int32 = 0x10 // Send/Recv Msg contains manual/user acknowledgment
	SM_CONV_ABORT_TYPE        int32 = 0x18 // Recv Msg contains conversation abort (Korean CDMA)
	SM_INTMD_DLV_NOTIFY_TYPE  int32 = 0x20 // Recv Msg contains intermediate notification

	// GSM Network features
	SM_NONE_GSM           int32 = 0x00 // No specific features selected
	SM_UDH_GSM            int32 = 0x40 // User Data Header indicator set
	SM_REPLY_PATH_GSM     int32 = 0x80 // Reply path set
	SM_UDH_REPLY_PATH_GSM int32 = 0xC0 // Both UDH & Reply path

	// Optional Parameter Tags, Min and Max Lengths
	// Following are the 2 byte tag and min/max lengths for
	// supported optional parameter (declann)

	OPT_PAR_MSG_WAIT int16 = 2

	// Privacy Indicator
	OPT_PAR_PRIV_IND int16 = 0x0201

	// Source Subaddress
	OPT_PAR_SRC_SUBADDR     int16 = 0x0202
	OPT_PAR_SRC_SUBADDR_MIN int32 = 2
	OPT_PAR_SRC_SUBADDR_MAX int32 = 23

	// Destination Subaddress
	OPT_PAR_DEST_SUBADDR     int16 = 0x0203
	OPT_PAR_DEST_SUBADDR_MIN int32 = 2
	OPT_PAR_DEST_SUBADDR_MAX int32 = 23

	// User Message Reference
	OPT_PAR_USER_MSG_REF int16 = 0x0204

	// User Response Code
	OPT_PAR_USER_RESP_CODE int16 = 0x0205

	// Language Indicator
	OPT_PAR_LANG_IND int16 = 0x020D

	// Source Port
	OPT_PAR_SRC_PORT int16 = 0x020A

	// Destination Port
	OPT_PAR_DST_PORT int16 = 0x020B

	// Concat Msg Ref Num
	OPT_PAR_SAR_MSG_REF_NUM int16 = 0x020C

	// Concat Total Segments
	OPT_PAR_SAR_TOT_SEG int16 = 0x020E

	// Concat Segment Seqnums
	OPT_PAR_SAR_SEG_SNUM int16 = 0x020F

	// SC Interface Version
	OPT_PAR_SC_IF_VER int16 = 0x0210

	// Display Time
	OPT_PAR_DISPLAY_TIME int16 = 0x1201

	// Validity Information
	OPT_PAR_MS_VALIDITY int16 = 0x1204

	// DPF Result
	OPT_PAR_DPF_RES int16 = 0x0420

	// Set DPF
	OPT_PAR_SET_DPF int16 = 0x0421

	// MS Availability Status
	OPT_PAR_MS_AVAIL_STAT int16 = 0x0422

	// Network Error Code
	OPT_PAR_NW_ERR_CODE     int16 = 0x0423
	OPT_PAR_NW_ERR_CODE_MIN int32 = 3
	OPT_PAR_NW_ERR_CODE_MAX int32 = 3

	// Extended int16 Message has no size limit

	// Delivery Failure Reason
	OPT_PAR_DEL_FAIL_RSN int16 = 0x0425

	// More Messages to Follow
	OPT_PAR_MORE_MSGS int16 = 0x0426

	// Message State
	OPT_PAR_MSG_STATE int16 = 0x0427

	// Callback Number
	OPT_PAR_CALLBACK_NUM     int16 = 0x0381
	OPT_PAR_CALLBACK_NUM_MIN int32 = 4
	OPT_PAR_CALLBACK_NUM_MAX int32 = 19

	// Callback Number Presentation  Indicator
	OPT_PAR_CALLBACK_NUM_PRES_IND int16 = 0x0302

	// Callback Number Alphanumeric Tag
	OPT_PAR_CALLBACK_NUM_ATAG     int16 = 0x0303
	OPT_PAR_CALLBACK_NUM_ATAG_MIN int32 = 1
	OPT_PAR_CALLBACK_NUM_ATAG_MAX int32 = 65

	// Number of messages in Mailbox
	OPT_PAR_NUM_MSGS int16 = 0x0304

	// SMS Received Alert
	OPT_PAR_SMS_SIGNAL int16 = 0x1203

	// Message Delivery Alert
	OPT_PAR_ALERT_ON_MSG_DELIVERY int16 = 0x130C

	// ITS Reply Type
	OPT_PAR_ITS_REPLY_TYPE int16 = 0x1380

	// ITS Session Info
	OPT_PAR_ITS_SESSION_INFO int16 = 0x1383

	// USSD Service Op
	OPT_PAR_USSD_SER_OP int16 = 0x0501

	// Priority
	SM_NOPRIORITY int32 = 0
	SM_PRIORITY   int32 = 1

	// Registered delivery
	//   SMSC Delivery Receipt (bits 1 & 0)
	SM_SMSC_RECEIPT_MASK          byte = 0x03
	SM_SMSC_RECEIPT_NOT_REQUESTED byte = 0x00
	SM_SMSC_RECEIPT_REQUESTED     byte = 0x01
	SM_SMSC_RECEIPT_ON_FAILURE    byte = 0x02
	//   SME originated acknowledgement (bits 3 & 2)
	SM_SME_ACK_MASK               byte = 0x0c
	SM_SME_ACK_NOT_REQUESTED      byte = 0x00
	SM_SME_ACK_DELIVERY_REQUESTED byte = 0x04
	SM_SME_ACK_MANUAL_REQUESTED   byte = 0x08
	SM_SME_ACK_BOTH_REQUESTED     byte = 0x0c
	//   Intermediate notification (bit 5)
	SM_NOTIF_MASK          byte = 0x010
	SM_NOTIF_NOT_REQUESTED byte = 0x000
	SM_NOTIF_REQUESTED     byte = 0x010

	// Replace if Present flag
	SM_NOREPLACE int32 = 0
	SM_REPLACE   int32 = 1

	// Destination flag
	SM_DEST_SME_ADDRESS int32 = 1
	SM_DEST_DL_NAME     int32 = 2

	// Higher Layer Message Type
	SM_LAYER_WDP  int32 = 0
	SM_LAYER_WCMP int32 = 1

	// Operation Class
	SM_OPCLASS_DATAGRAM    int32 = 0
	SM_OPCLASS_TRANSACTION int32 = 3

	// Originating MSC Address
	OPT_PAR_ORIG_MSC_ADDR     int16 = -32639 // int16(0x8081)
	OPT_PAR_ORIG_MSC_ADDR_MIN int32 = 1
	OPT_PAR_ORIG_MSC_ADDR_MAX int32 = 24

	// Destination MSC Address
	OPT_PAR_DEST_MSC_ADDR     int16 = -32638 // int16(0x8082)
	OPT_PAR_DEST_MSC_ADDR_MIN int32 = 1
	OPT_PAR_DEST_MSC_ADDR_MAX int32 = 24

	// Unused Tag
	OPT_PAR_UNUSED int32 = 0xffff

	// Destination Address Subunit
	OPT_PAR_DST_ADDR_SUBUNIT int16 = 0x0005

	// Destination Network Type
	OPT_PAR_DST_NW_TYPE int16 = 0x0006

	// Destination Bearer Type
	OPT_PAR_DST_BEAR_TYPE int16 = 0x0007

	// Destination Telematics ID
	OPT_PAR_DST_TELE_ID int16 = 0x0008

	// Source Address Subunit
	OPT_PAR_SRC_ADDR_SUBUNIT int16 = 0x000D

	// Source Network Type
	OPT_PAR_SRC_NW_TYPE int16 = 0x000E

	// Source Bearer Type
	OPT_PAR_SRC_BEAR_TYPE int16 = 0x000F

	// Source Telematics ID
	OPT_PAR_SRC_TELE_ID int16 = 0x0010

	// QOS Time to Live
	OPT_PAR_QOS_TIME_TO_LIVE     int16 = 0x0017
	OPT_PAR_QOS_TIME_TO_LIVE_MIN int32 = 1
	OPT_PAR_QOS_TIME_TO_LIVE_MAX int32 = 4

	// Payload Type
	OPT_PAR_PAYLOAD_TYPE int16 = 0x0019

	// Additional Status Info Text
	OPT_PAR_ADD_STAT_INFO     int16 = 0x001D
	OPT_PAR_ADD_STAT_INFO_MIN int32 = 1
	OPT_PAR_ADD_STAT_INFO_MAX int32 = 256

	// Receipted Message ID
	OPT_PAR_RECP_MSG_ID     int16 = 0x001E
	OPT_PAR_RECP_MSG_ID_MIN int32 = 1
	OPT_PAR_RECP_MSG_ID_MAX int32 = 65

	// Message Payload
	OPT_PAR_MSG_PAYLOAD     int16 = 0x0424
	OPT_PAR_MSG_PAYLOAD_MIN int32 = 1
	OPT_PAR_MSG_PAYLOAD_MAX int32 = 1500

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
	DFLT_REG_DELIVERY  byte   = SM_SMSC_RECEIPT_NOT_REQUESTED | SM_SME_ACK_NOT_REQUESTED | SM_NOTIF_NOT_REQUESTED
	DFLT_DFLTMSGID     byte   = 0
	DFLT_MSG_LEN       byte   = 0
	DFLT_ESM_CLASS     byte   = 0
	DFLT_DATA_CODING   byte   = 0
	DFLT_PROTOCOLID    byte   = 0
	DFLT_PRIORITY_FLAG byte   = 0
	DFTL_REPLACE_IFP   byte   = 0
	DFLT_DL_NAME       string = ""
	DFLT_GSM_TON       byte   = GSM_TON_UNKNOWN
	DFLT_GSM_NPI       byte   = GSM_NPI_UNKNOWN
	DFLT_DEST_FLAG     byte   = 0 // not set
	MAX_PDU_LEN        int32  = 5000

	PDU_HEADER_SIZE int32 = 16 // 4 integers
	TLV_HEADER_SIZE int32 = 4  // 2 int16s: tag & length

	// all times in milliseconds
	RECEIVER_TIMEOUT           int64 = 60000
	CONNECTION_RECEIVE_TIMEOUT int64 = 10000
	UNBIND_RECEIVE_TIMEOUT     int64 = 5000
	CONNECTION_SEND_TIMEOUT    int64 = 20000
	COMMS_TIMEOUT              int64 = 60000
	QUEUE_TIMEOUT              int64 = 10000
	ACCEPT_TIMEOUT             int64 = 60000

	RECEIVE_BLOCKING int64 = -1

	MAX_VALUE_PORT     int32 = 65535
	MIN_VALUE_PORT     int32 = 100
	MIN_LENGTH_ADDRESS int32 = 7
)

var defaultTon byte = DFLT_GSM_TON
var defaultNpi byte = DFLT_GSM_NPI

var l sync.RWMutex

// SetDefaultTon ..
func SetDefaultTon(dfltTon byte) {
	l.Lock()
	defer l.Unlock()
	defaultTon = dfltTon
}

// GetDefaultTon ...
func GetDefaultTon() byte {
	l.RLock()
	defer l.RUnlock()
	return defaultTon
}

// SetDefaultNpi ...
func SetDefaultNpi(dfltNpi byte) {
	l.Lock()
	defer l.Unlock()
	defaultNpi = dfltNpi
}

// GetDefaultNpi ...
func GetDefaultNpi() byte {
	l.RLock()
	defer l.RUnlock()
	return defaultNpi
}

// GetCurrentTime ...
func GetCurrentTime() time.Time {
	return time.Now()
}
