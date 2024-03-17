package gomkv

const (
	// Element IDs (pre-defined)
	// see https://www.matroska.org/technical/elements.html for more
	ElementEMBL    = 0x1a45dfa3
	ElementDocType = 0x4282
	ElementVoid    = 0xec
	// MKV specific
	ElementSegment = 0x18538067
	// 8 top-level elements
	ElementSeekHead    = 0x114D9B74
	ElementInfo        = 0x1549a966
	ElementTracks      = 0x1654AE6B
	ElementCluster     = 0x1F43B675
	ElementCues        = 0x1C53BB6B
	ElementTags        = 0x1254c367
	ElementChapters    = 0x1043A770
	ElementAttachments = 0x1941A469
	// level 2
	TimestampScale = 0x2AD7B1
	MuxingApp      = 0x4D80
	WritingApp     = 0x5741
	TrackEntry     = 0xAE
	Tag            = 0x7373
	// level 3
	TrackNumber = 0xD7
	TrackUID    = 0x73C5
	TrackType   = 0x83
)

var id2name = map[uint32]string{
	ElementEMBL:    "EMBL",
	ElementDocType: "DocType",
	ElementVoid:    "Void",
	// mkv
	ElementSegment:     "Segment",
	ElementSeekHead:    "SeekHead",
	ElementInfo:        "Info",
	ElementTracks:      "Tracks",
	ElementCluster:     "Cluster",
	ElementCues:        "Cues",
	ElementTags:        "Tags",
	ElementChapters:    "Chapters",
	ElementAttachments: "Attachments",
	// level 2
	MuxingApp:  "MuxingApp",
	WritingApp: "WritingApp",
	TrackEntry: "TrackEntry",
	Tag:        "Tag",
	// level 3
	TrackNumber: "TrackNumber",
	TrackUID:    "TrackUID",
	TrackType:   "TrackType",
}

var isUnicode = map[uint32]bool{
	ElementDocType: true,
	//
	MuxingApp:  true,
	WritingApp: true,
}
