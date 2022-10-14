package assets

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type ColorFormat int

const (
	ColorFormat_RGBA8 ColorFormat = iota
)

type Texture struct {
	//Path only exists for textures loaded from disk
	Path   string
	TexID  uint32
	Width  int32
	Height int32
	Pixels []byte
}

type TextureLoadOptions struct {
	TryLoadFromCache bool
	WriteToCache     bool
	GenMipMaps       bool
	KeepPixelsInMem  bool
}

type Cubemap struct {
	// These only exists for textures loaded from disk
	RightPath string
	LeftPath  string
	TopPath   string
	BotPath   string
	FrontPath string
	BackPath  string
	TexID     uint32
}

func LoadTexturePNG(file string, loadOptions *TextureLoadOptions) (Texture, error) {

	if loadOptions == nil {
		loadOptions = &TextureLoadOptions{}
	}

	if loadOptions.TryLoadFromCache {
		if tex, ok := GetTextureFromCachePath(file); ok {
			return tex, nil
		}
	}

	//Load from disk
	fileBytes, err := os.ReadFile(file)
	if err != nil {
		return Texture{}, err
	}

	img, err := png.Decode(bytes.NewReader(fileBytes))
	if err != nil {
		return Texture{}, err
	}

	tex := Texture{
		Path: file,
	}

	tex.Pixels, tex.Width, tex.Height = pixelsFromNrgbaPng(img)

	//Prepare opengl stuff
	gl.GenTextures(1, &tex.TexID)
	gl.BindTexture(gl.TEXTURE_2D, tex.TexID)

	// set the texture wrapping/filtering options (on the currently bound texture object)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	// load and generate the texture
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, tex.Width, tex.Height, 0, gl.RGBA, gl.UNSIGNED_BYTE, unsafe.Pointer(&tex.Pixels[0]))

	if loadOptions.GenMipMaps {
		gl.GenerateMipmap(tex.TexID)
	}

	if loadOptions.WriteToCache {
		AddTextureToCache(tex)
	}

	if !loadOptions.KeepPixelsInMem {
		tex.Pixels = nil
	}

	return tex, nil
}

func LoadTextureInMemPngImg(img image.Image, loadOptions *TextureLoadOptions) (Texture, error) {

	if loadOptions == nil {
		loadOptions = &TextureLoadOptions{}
	}

	tex := Texture{}
	tex.Pixels, tex.Width, tex.Height = pixelsFromNrgbaPng(img)

	//Prepare opengl stuff
	gl.GenTextures(1, &tex.TexID)
	gl.BindTexture(gl.TEXTURE_2D, tex.TexID)

	// set the texture wrapping/filtering options (on the currently bound texture object)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	// load and generate the texture
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, tex.Width, tex.Height, 0, gl.RGBA, gl.UNSIGNED_BYTE, unsafe.Pointer(&tex.Pixels[0]))

	if loadOptions.GenMipMaps {
		gl.GenerateMipmap(tex.TexID)
	}

	if loadOptions.WriteToCache {
		AddTextureToCache(tex)
	}

	if !loadOptions.KeepPixelsInMem {
		tex.Pixels = nil
	}

	return tex, nil
}

func LoadTextureJpeg(file string, loadOptions *TextureLoadOptions) (Texture, error) {

	if loadOptions == nil {
		loadOptions = &TextureLoadOptions{}
	}

	if loadOptions.TryLoadFromCache {
		if tex, ok := GetTextureFromCachePath(file); ok {
			return tex, nil
		}
	}

	//Load from disk
	fileBytes, err := os.ReadFile(file)
	if err != nil {
		return Texture{}, err
	}

	img, err := jpeg.Decode(bytes.NewReader(fileBytes))
	if err != nil {
		return Texture{}, err
	}

	tex := Texture{
		Path: file,
	}

	tex.Pixels, tex.Width, tex.Height = pixelsFromNrgbaPng(img)

	//Prepare opengl stuff
	gl.GenTextures(1, &tex.TexID)
	gl.BindTexture(gl.TEXTURE_2D, tex.TexID)

	// set the texture wrapping/filtering options (on the currently bound texture object)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	// load and generate the texture
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, tex.Width, tex.Height, 0, gl.RGBA, gl.UNSIGNED_BYTE, unsafe.Pointer(&tex.Pixels[0]))

	if loadOptions.GenMipMaps {
		gl.GenerateMipmap(tex.TexID)
	}

	if loadOptions.WriteToCache {
		AddTextureToCache(tex)
	}

	if !loadOptions.KeepPixelsInMem {
		tex.Pixels = nil
	}

	return tex, nil
}

func pixelsFromNrgbaPng(img image.Image) (pixels []byte, width, height int32) {

	//NOTE: Load bottom left to top right because this is the texture coordinate system used by OpenGL
	//NOTE: We only support 8-bit channels (32-bit colors) for now
	i := 0
	width, height = int32(img.Bounds().Dx()), int32(img.Bounds().Dy())
	pixels = make([]byte, img.Bounds().Dx()*img.Bounds().Dy()*4)
	for y := img.Bounds().Dy() - 1; y >= 0; y-- {
		for x := 0; x < img.Bounds().Dx(); x++ {

			c := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)

			pixels[i] = c.R
			pixels[i+1] = c.G
			pixels[i+2] = c.B
			pixels[i+3] = c.A

			i += 4
		}
	}

	return pixels, width, height
}

func pixelsFromNrgbaJpg(img image.Image) (pixels []byte, width, height int32) {

	//NOTE: Load bottom left to top right because this is the texture coordinate system used by OpenGL
	//NOTE: We only support 8-bit channels (32-bit colors) for now
	i := 0
	width, height = int32(img.Bounds().Dx()), int32(img.Bounds().Dy())
	pixels = make([]byte, img.Bounds().Dx()*img.Bounds().Dy()*4)
	for y := img.Bounds().Dy() - 1; y >= 0; y-- {
		for x := 0; x < img.Bounds().Dx(); x++ {

			c := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)

			pixels[i] = c.R
			pixels[i+1] = c.G
			pixels[i+2] = c.B
			pixels[i+3] = c.A

			i += 4
		}
	}

	return pixels, width, height
}

func LoadCubemapTextures(rightTex, leftTex, topTex, botTex, frontTex, backTex string) (Cubemap, error) {

	var imgDecoder func(r io.Reader) (image.Image, error)
	var pixelDecoder func(image.Image) ([]byte, int32, int32)
	ext := strings.ToLower(path.Ext(rightTex))
	if ext == ".jpg" || ext == ".jpeg" {
		imgDecoder = jpeg.Decode
		pixelDecoder = pixelsFromNrgbaJpg
	} else if ext == ".png" {
		imgDecoder = png.Decode
		pixelDecoder = pixelsFromNrgbaPng
	} else {
		return Cubemap{}, fmt.Errorf("unknown image extension: %s. Expected one of: .jpg, .jpeg, .png", ext)
	}

	cmap := Cubemap{
		RightPath: rightTex,
		LeftPath:  leftTex,
		TopPath:   topTex,
		BotPath:   botTex,
		FrontPath: frontTex,
		BackPath:  backTex,
	}

	gl.GenTextures(1, &cmap.TexID)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, cmap.TexID)

	// The order here matters
	texturePaths := []string{rightTex, leftTex, topTex, botTex, frontTex, backTex}
	for i := uint32(0); i < uint32(len(texturePaths)); i++ {

		fPath := texturePaths[i]

		//Load from disk
		fileBytes, err := os.ReadFile(fPath)
		if err != nil {
			return Cubemap{}, err
		}

		img, err := imgDecoder(bytes.NewReader(fileBytes))
		if err != nil {
			return Cubemap{}, err
		}

		pixels, width, height := pixelDecoder(img)

		gl.TexImage2D(uint32(gl.TEXTURE_CUBE_MAP_POSITIVE_X)+i, 0, gl.RGBA8, int32(width), int32(height), 0, gl.RGBA, gl.UNSIGNED_BYTE, unsafe.Pointer(&pixels[0]))
	}

	// set the texture wrapping/filtering options (on the currently bound texture object)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)

	return cmap, nil
}
