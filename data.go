package gomkv

const (
	// Element IDs (pre-defined)
	// see https://www.matroska.org/technical/elements.html for more
	ElementEMBL    = 0x1a45dfa3
	ElementSegment = 0x18538067
)

var id2name = map[uint32]string{
	ElementEMBL:    "EMBL",
	ElementSegment: "Segment",
}
