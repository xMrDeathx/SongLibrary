package songlibrarycommand

type GetSongsCommand struct {
	Group       *string
	Song        *string
	ReleaseDate *string
	Text        *string
	Link        *string
	Page        *int
	Limit       *int
}
