package internal

import (
	"embed"
	"fmt"
	"github.com/Frabjous-Studios/ebitengine-game-template/internal/debug"
	"github.com/ebitenui/ebitenui/image"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/tinne26/etxt"
	"golang.org/x/image/font"
	"io/fs"
	"strings"
)

// Resources makes all multimedia resources for the game available.
var Resources = resources{} // TODO: deferred loading & caching of resources

type resources struct {
	fontLib    *etxt.FontLibrary
	faces      map[string]*truetype.Font
	nineSlices map[string]*image.NineSlice
	images     map[string]*ebiten.Image
	shaders    map[string]*ebiten.Shader
}

const LunchtimeFont = "lunchds.ttf"

//go:embed gamedata/fonts
var fonts embed.FS

func init() {
	// load fonts
	fontLib := etxt.NewFontLibrary()
	_, _, err := fontLib.ParseEmbedDirFonts("gamedata/fonts", fonts)
	if err != nil {
		panic(err)
	}
	var fonts []string
	_ = fontLib.EachFont(func(s string, font *etxt.Font) error {
		fonts = append(fonts, s)
		return nil
	})
	debug.Println("fonts available:", strings.Join(fonts, ","))
	Resources.fontLib = fontLib

	// load nineslices
	Resources.nineSlices = make(map[string]*image.NineSlice)

	// load shaders
	Resources.shaders = make(map[string]*ebiten.Shader)

	// load font faces (again?)
	Resources.faces = make(map[string]*truetype.Font)
}

// TODO: Enable //go:embed gamedata/img
var art embed.FS

func (r *resources) GetImage(path string) *ebiten.Image {
	if _, ok := r.images[path]; !ok {
		bgPath := fmt.Sprintf("gamedata/img/%s", path)
		img, _, err := ebitenutil.NewImageFromFileSystem(art, bgPath)
		if err != nil {
			debug.Printf("failed to load image: %s", bgPath)
			return nil
		}
		r.images[path] = img
	}
	return r.images[path]
}

const dpi = 72

// GetFace returns a new font face with the provided FontID and size.
func (r *resources) GetFace(path string, size float64) font.Face {
	if _, ok := r.faces[path]; !ok {
		f, err := loadFont(path)
		if err != nil {
			debug.Printf("error load face %s", path)
			return nil
		}
		r.faces[path] = f
	}

	return truetype.NewFace(r.faces[path], &truetype.Options{
		Size:    size,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}

// GetShader retrieves the shader with the provided ID
func (r *resources) GetShader(path string) *ebiten.Shader {
	return r.shaders[path]
}

func loadFont(filename string) (*truetype.Font, error) {
	fontData, err := fs.ReadFile(fonts, "gamedata/fonts/"+filename)
	if err != nil {
		return nil, fmt.Errorf("loadFace: %w", err)
	}
	ttfFont, err := truetype.Parse(fontData)
	if err != nil {
		return nil, err
	}
	return ttfFont, nil
}

// TODO: enable //go:embed gamedata/shader
var shader embed.FS

// loadShader loads and compiles the provided shader.
func loadShader(name string) *ebiten.Shader {
	bytes, err := fs.ReadFile(shader, "gamedata/shader/"+name)
	if err != nil {
		panic(fmt.Errorf("while loading shader: %w", err))
	}
	shader, err := ebiten.NewShader(bytes)
	if err != nil {
		panic(fmt.Errorf("while compiling shader: %w", err))
	}
	return shader
}
