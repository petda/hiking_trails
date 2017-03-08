package models

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type GEOCoordinate struct {
	Latitude  float32 `json:"lat"`
	Longitude float32 `json:"lng"`
}

func GEOCoordinateFromBytes(data []byte) (*GEOCoordinate, error) {
	coordinate := &GEOCoordinate{}
	buffer := bytes.NewReader(data)
	err := binary.Read(buffer, binary.LittleEndian, coordinate)
	if err != nil {
		return nil, err
	}

	return coordinate, nil
}

func (coordinate GEOCoordinate) AsBytes() []byte {
	buffer := bytes.NewBuffer(make([]byte, 0, 4*2))

	err := binary.Write(buffer, binary.LittleEndian, coordinate)
	if err != nil {
		panic(fmt.Sprintf("Failed to convert GEOCoordinate to bytes: %s", err))
	}

	return buffer.Bytes()
}

type GEOCoordinates []GEOCoordinate

func NewGEOCoordinates() []GEOCoordinate {
	return make([]GEOCoordinate, 0)
}

func GEOCoordinatesFromBytes(data []byte) (GEOCoordinates, error) {
	coordinates := make([]GEOCoordinate, len(data)/(4*2))
	buffer := bytes.NewReader(data)

	err := binary.Read(buffer, binary.LittleEndian, &coordinates)
	if err != nil {
		return nil, err
	}

	return coordinates, nil
}

func (coordinates GEOCoordinates) AsBytes() []byte {
	buffer := bytes.NewBuffer(make([]byte, 0, len(coordinates)*4*2))

	err := binary.Write(buffer, binary.LittleEndian, coordinates)
	if err != nil {
		panic(fmt.Sprintf("Failed to convert []GEOCoordinate to bytes: %s", err))
	}

	return buffer.Bytes()

}
