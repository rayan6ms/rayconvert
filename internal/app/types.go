package app

type BuildInfo struct {
	Version string
	Commit  string
	Date    string
}

type SubjectKind int

const (
	SubjectFile      SubjectKind = iota
	SubjectDirImages             // directory treated as "images"
	SubjectImages
	SubjectVideos
)

type Config struct {
	SubjectRaw string
	Subject    SubjectKind

	FilePath string
	InDir    string
	OutDir   string

	ToFormat string

	Append    bool
	Replace   bool
	Mute      bool
	FullyMute bool

	Help    bool
	Version bool

	Build BuildInfo
}
