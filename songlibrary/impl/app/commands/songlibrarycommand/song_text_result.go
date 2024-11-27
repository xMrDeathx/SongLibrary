package songlibrarycommand

type SongTextResult struct {
	SongInfo    SongResult
	ReleaseDate string
	Text        []string
	Link        string
}
