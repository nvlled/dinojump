package assets

import "embed"

//go:embed dinosprites-vita.png
var DinoSpriteData []byte

//go:embed dinosprites-doux.png
//go:embed dinosprites-vita.png
//go:embed "Cielo pixelado.png"
//go:embed lemcraft-tiles.png
var FS embed.FS
