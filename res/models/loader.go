package models

import (
	"os"
	"strconv"
	"strings"

	"github.com/bloeys/go-sdl-engine/logging"
)

type ObjInfo struct {
	VertPos    []float32
	TriIndices []uint32
}

func LoadObj(file string) (objInfo ObjInfo, err error) {

	b, err := os.ReadFile(file)
	if err != nil {
		return objInfo, err
	}

	lines := strings.Split(string(b), "\n")
	for i := 0; i < len(lines); i++ {

		s := strings.SplitN(lines[i], " ", 2)
		switch s[0] {
		case "v":

			vertPosStrings := strings.Split(s[1], " ")

			f, err := strconv.ParseFloat(vertPosStrings[0], 32)
			if err != nil {
				return objInfo, err
			}
			objInfo.VertPos = append(objInfo.VertPos, float32(f))

			f, err = strconv.ParseFloat(vertPosStrings[1], 32)
			if err != nil {
				return objInfo, err
			}
			objInfo.VertPos = append(objInfo.VertPos, float32(f))

			f, err = strconv.ParseFloat(vertPosStrings[2], 32)
			if err != nil {
				return objInfo, err
			}
			objInfo.VertPos = append(objInfo.VertPos, float32(f))

		case "f":

			facesStrings := strings.Split(s[1], " ")
			objInfo.TriIndices = append(objInfo.TriIndices, getVertIndexFromFace(facesStrings[0]))
			objInfo.TriIndices = append(objInfo.TriIndices, getVertIndexFromFace(facesStrings[1]))
			objInfo.TriIndices = append(objInfo.TriIndices, getVertIndexFromFace(facesStrings[2]))

		default:
		}
	}

	return objInfo, nil
}

func getVertIndexFromFace(f string) uint32 {

	indxStr := strings.SplitN(f, "/", 2)[0]
	index, err := strconv.Atoi(indxStr)
	if err != nil {
		logging.ErrLog.Printf("Invalid face index '%v'. Err: %v", indxStr, err)
		return 0
	}
	return uint32(index) - 1
}
