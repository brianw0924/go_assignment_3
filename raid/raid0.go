package raid

import (
	"errors"
	"sync"
)

type Raid0 struct {
	*RaidBase
}

func NewRaid0() *Raid0 {
	return &Raid0{
		RaidBase: NewRaidBase(),
	}
}

func (r *Raid0) Write(data []byte) error {
	if len(data) > TOTAL_STORAGE {
		return errors.New("no more storage")
	}

	var wg sync.WaitGroup
	done := make(chan struct{})
	r.StartWriter(done, &wg)

	for blockId, pos := 0, 0; blockId < BLOCK_PER_DISK && pos < len(data); blockId += 1 {
		for diskId := 0; diskId < STRIPE_WIDTH && pos < len(data); diskId, pos = diskId+1, pos+BLOCK_SIZE {
			wg.Add(1)
			r.DiskArray[diskId].TaskChan <- NewTask(blockId, data[pos:min(len(data), pos+BLOCK_SIZE)])
		}
	}

	wg.Wait()
	close(done)
	return nil
}

func (r *Raid0) Read(length int) (string, error) {

	done := make(chan struct{})
	r.StartReader(done)

	for blockId, pos := 0, 0; blockId < BLOCK_PER_DISK && pos < length; blockId += 1 {
		for diskId := 0; diskId < STRIPE_WIDTH && pos < length; diskId, pos = diskId+1, pos+BLOCK_SIZE {
			r.DiskArray[diskId].TaskChan <- NewTask(blockId, nil)
		}
	}

	// Sequentially read back
	data := []byte{}
	for diskId, pos := 0, 0; pos < length; diskId, pos = (diskId+1)%STRIPE_WIDTH, pos+BLOCK_SIZE {
		block := <-r.DiskArray[diskId].ReadChan
		data = append(data, block.Data[:min(BLOCK_SIZE, length-pos)]...)
	}
	close(done)
	return string(data), nil

	// // Non-parallel read
	// data := []byte{}
	// for blockId, pos := 0, 0; blockId < BLOCK_PER_DISK && pos < length; blockId += 1 {
	// 	for diskId := 0; diskId < STRIPE_WIDTH && pos < length; diskId, pos = diskId+1, pos+BLOCK_SIZE {
	// 		block := r.DiskArray[diskId].ReadBlock(blockId)
	// 		data = append(data, block.Data[:min(BLOCK_SIZE, length-pos)]...)
	// 	}
	// }
	// return string(data), nil

}
