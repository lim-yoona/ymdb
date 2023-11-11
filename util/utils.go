package util

const PUTID uint32 = 1
const GETID uint32 = 2
const DELETEID uint32 = 3
const RESPONSEID uint32 = 1
const PUTCOMMAND string = "put"
const GETCOMMAND string = "get"
const DELETECOMMAND string = "delete"

type Put struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Other struct {
	Data string `json:"key"`
}
