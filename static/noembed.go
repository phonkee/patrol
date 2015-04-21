// +build !EMBED_STATIC

package static

const (
	// info whether data is embedded or not
	IsEmbedded      = false
	NotEmbeddedHtml = `<p>Patrol was not compiled with static
Please provide static_dir cli argument for http command.
</p>`
)

var (
	_bindata = map[string]func() ([]byte, error){
		"index.html": NoStatic,
	}
)

func NoStatic() ([]byte, error) {
	return []byte(NotEmbeddedHtml), nil
}
