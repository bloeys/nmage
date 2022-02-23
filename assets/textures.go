package assets

import (
	"bytes"
	"image/color"
	"image/png"
	"os"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type ColorFormat int

const (
	ColorFormat_RGBA8 ColorFormat = iota
)

type Texture struct {
	Path   string
	TexID  uint32
	Width  int32
	Height int32
	Pixels []byte
}

func LoadPNGTexture(file string) (Texture, error) {

	if tex, ok := GetTexturePath(file); ok {
		return tex, nil
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
		Path:   file,
		Width:  int32(img.Bounds().Dx()),
		Height: int32(img.Bounds().Dy()),
		Pixels: make([]byte, img.Bounds().Dx()*img.Bounds().Dy()*4),
	}

	//NOTE: Load bottom left to top right because this is the texture coordinate system used by OpenGL
	//NOTE: We only support 8-bit channels (32-bit colors) for now
	i := 0
	for y := img.Bounds().Dy() - 1; y >= 0; y-- {
		for x := 0; x < img.Bounds().Dx(); x++ {

			c := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)

			tex.Pixels[i] = c.R
			tex.Pixels[i+1] = c.G
			tex.Pixels[i+2] = c.B
			tex.Pixels[i+3] = c.A

			i += 4
		}
	}

	//Prepare opengl stuff
	gl.GenTextures(1, &tex.TexID)
	gl.BindTexture(gl.TEXTURE_2D, tex.TexID)

	// set the texture wrapping/filtering options (on the currently bound texture object)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	// load and generate the texture
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, tex.Width, tex.Height, 0, gl.RGBA, gl.UNSIGNED_BYTE, unsafe.Pointer(&tex.Pixels[0]))
	gl.GenerateMipmap(gl.TEXTURE_2D)

	SetTexture(tex)

	return tex, nil
}
