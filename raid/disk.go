package raid

import (
	"errors"
	"sync"
)

const (
	OFFSET         = 2
	BLOCK_SIZE     = 8 << OFFSET // KB
	BLOCK_PER_DISK = 4
)

// total storage = BLOKC_SIZE * BLOCK_PER_DISK * STRIPE_WIDTH

var ErrIndexOutOfBound error = errors.New("index out of bound")

type Disk struct {
	BlockArray [BLOCK_PER_DISK]Block
	TaskChan   chan *Task // read or write request
	ReadChan   chan Block // read block
}

func NewDisk() *Disk {
	blockArray := [BLOCK_PER_DISK]Block{}
	for i := range blockArray {
		blockArray[i] = NewBlock(i)
	}
	return &Disk{
		BlockArray: blockArray,
		TaskChan:   make(chan *Task, BLOCK_PER_DISK),
		ReadChan:   make(chan Block, BLOCK_PER_DISK),
	}
}

func (d *Disk) AsyncWriteBlock(id int, data []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	copy(d.BlockArray[id].Data[:], data)
}

func (d *Disk) WriteBlock(id int, data []byte) {
	copy(d.BlockArray[id].Data[:], data)
}

func (d *Disk) ReadBlock(id int) Block {
	return d.BlockArray[id]
}

func (d *Disk) Clear() {
	for i := range d.BlockArray {
		d.BlockArray[i] = NewBlock(i)
	}
}

type Task struct {
	BlockId int
	Data    []byte
}

func NewTask(blockId int, data []byte) *Task {
	return &Task{
		BlockId: blockId,
		Data:    data,
	}
}
