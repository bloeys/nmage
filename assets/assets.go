package assets

var (
	Textures     map[uint32]Texture = make(map[uint32]Texture)
	TexturePaths map[string]uint32  = make(map[string]uint32)
)

func AddTextureToCache(t Texture) {

	if _, ok := TexturePaths[t.Path]; ok {
		return
	}

	println("Loaded texture:", t.Path)
	Textures[t.TexID] = t
	TexturePaths[t.Path] = t.TexID
}

func GetTextureFromCacheID(texID uint32) (Texture, bool) {
	tex, ok := Textures[texID]
	return tex, ok
}

func GetTextureFromCachePath(path string) (Texture, bool) {
	tex, ok := Textures[TexturePaths[path]]
	return tex, ok
}
