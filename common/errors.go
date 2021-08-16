package common

const (
	ERR_OK = iota
	ERR_INVALID_BODY
	ERR_INVALID_DATA
	ERR_INVALID_SIGN
	ERR_DB_NO_RECORD
	ERR_DB_UPDATE_FAILED
	ERR_INVALID_STATUS
	ERR_DECODE_FAILED
)

func GetCodeMsg(code int) string {
	switch code {
	case ERR_OK:
		return "ok"
	case ERR_INVALID_BODY:
		return "invalid msg body"
	case ERR_INVALID_DATA:
		return "invalid data"
	case ERR_INVALID_SIGN:
		return "invalid sign"
	case ERR_DB_NO_RECORD:
		return "no db record"
	case ERR_INVALID_STATUS:
		return "invalid status"
	case ERR_DB_UPDATE_FAILED:
		return "update db failed"
	case ERR_DECODE_FAILED:
		return "decode failed"
	}
	return "unkown error"
}
