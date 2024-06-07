package web

type ChainType int8

const (
	ChainEvm    ChainType = 0x01
	ChainNorn   ChainType = 0x02
	ChainCustom ChainType = 0x03
)
