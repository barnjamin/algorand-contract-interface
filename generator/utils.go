package generator

import (
	"context"
	"encoding/binary"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

const (
	//algodAddress = "https://testnet.algoexplorer.io"
	//algodToken = ""

	algodAddress = "http://localhost:4001"
	algodToken   = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
)

func compile(b []byte) (*models.CompileResponse, error) {
	// Create an algod client
	algodClient, err := algod.MakeClient(algodAddress, algodToken)
	if err != nil {
		return nil, err
	}

	headers := []*common.Header{
		{Key: "User-Agent", Value: "LetMeIn.gif"},
	}

	resp, err := algodClient.TealCompile(b).Do(context.Background(), headers...)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Below functions jacked straight from https://github.com/algorand/go-algorand-sdk/blob/v1.10.0/logic/logic.go
func readIntConstBlock(program []byte, pc int) (size int, ints []uint64, err error) {
	size = 1
	numInts, bytesUsed := binary.Uvarint(program[pc+size:])
	if bytesUsed <= 0 {
		err = fmt.Errorf("could not decode int const block size at pc=%d", pc+size)
		return
	}

	size += bytesUsed
	for i := uint64(0); i < numInts; i++ {
		if pc+size >= len(program) {
			err = fmt.Errorf("intcblock ran past end of program")
			return
		}
		num, bytesUsed := binary.Uvarint(program[pc+size:])
		if bytesUsed <= 0 {
			err = fmt.Errorf("could not decode int const[%d] at pc=%d", i, pc+size)
			return
		}
		ints = append(ints, num)
		size += bytesUsed
	}
	return
}

func readByteConstBlock(program []byte, pc int) (size int, byteArrays [][]byte, err error) {
	size = 1
	numInts, bytesUsed := binary.Uvarint(program[pc+size:])
	if bytesUsed <= 0 {
		err = fmt.Errorf("could not decode []byte const block size at pc=%d", pc+size)
		return
	}

	size += bytesUsed
	for i := uint64(0); i < numInts; i++ {
		if pc+size >= len(program) {
			err = fmt.Errorf("bytecblock ran past end of program")
			return
		}
		scanTarget := program[pc+size:]
		itemLen, bytesUsed := binary.Uvarint(scanTarget)
		if bytesUsed <= 0 {
			err = fmt.Errorf("could not decode []byte const[%d] at pc=%d", i, pc+size)
			return
		}
		size += bytesUsed
		if pc+size+int(itemLen) > len(program) {
			err = fmt.Errorf("bytecblock ran past end of program")
			return
		}
		byteArray := program[pc+size : pc+size+int(itemLen)]
		byteArrays = append(byteArrays, byteArray)
		size += int(itemLen)
	}
	return
}

func readPushIntOp(program []byte, pc int) (size int, foundInt uint64, err error) {
	size = 1
	foundInt, bytesUsed := binary.Uvarint(program[pc+size:])
	if bytesUsed <= 0 {
		err = fmt.Errorf("could not decode push int const at pc=%d", pc+size)
		return
	}

	size += bytesUsed
	return
}

func readPushByteOp(program []byte, pc int) (size int, byteArray []byte, err error) {
	size = 1
	itemLen, bytesUsed := binary.Uvarint(program[pc+size:])
	if bytesUsed <= 0 {
		err = fmt.Errorf("could not decode push []byte const size at pc=%d", pc+size)
		return
	}

	size += bytesUsed
	if pc+size+int(itemLen) > len(program) {
		err = fmt.Errorf("pushbytes ran past end of program")
		return
	}
	byteArray = program[pc+size : pc+size+int(itemLen)]
	size += int(itemLen)
	return
}
