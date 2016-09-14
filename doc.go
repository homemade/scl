package scl

type MixinDoc struct {
	Name      string
	File      string
	Line      int
	Reference string
	Signature string
	Docs      string
	Children  MixinDocs
}

type MixinDocs []MixinDoc
