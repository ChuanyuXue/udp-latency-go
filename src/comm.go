package src

const (
	HEADER_UNTAG_SIZE = 42
	HEADER_SIZE       = 46
	AUG_SIZE          = 4 + 8 + 8 + 8
	MAX_PKT_SIZE      = 1500
	MIN_PKT_SIZE      = AUG_SIZE + HEADER_SIZE
)

func getTime(dev string) uint64 {

}
