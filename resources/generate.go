//go:generate go run github.com/hajimehoshi/file2byteslice/cmd/file2byteslice@v1.0.0 -package=blocks -input=./images/blocks/blocks.png -output=./images/blocks/blocks.go -var=Blocks_png

//go:generate go run github.com/hajimehoshi/file2byteslice/cmd/file2byteslice@v1.0.0 -package=resources -input ./images/blocks.png -output ./images/blocks.go -var Blocks_png

package resources
