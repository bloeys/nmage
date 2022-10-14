package assets

var (
	Textures     = make(map[uint32]Texture)
	TexturePaths = make(map[string]uint32)
)

func AddTextureToCache(t Texture) {

	if t.Path != "" {
		if _, ok := TexturePaths[t.Path]; ok {
			return
		}
		println("Loaded texture from path:", t.Path)
		Textures[t.TexID] = t
		TexturePaths[t.Path] = t.TexID
		return
	}

	println("Loaded in-mem texture with ID:", t.TexID)
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
