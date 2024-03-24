package raid

import (
	"errors"
	"fmt"
	"log"
	"sync"
)

type Raid10 struct {
	*RaidBase
}

func NewRaid10() *Raid10 {
	if STRIPE_WIDTH%2 != 0 {
		log.Fatal(fmt.Errorf("number of disks have to be even"))
	}
	if STRIPE_WIDTH < 4 {
		log.Fatal(fmt.Errorf("number of disks have to >= 4"))
	}

	return &Raid10{
		RaidBase: NewRaidBase(),
	}
}

func (r *Raid10) Write(data []byte) error {
	if len(data) > TOTAL_STORAGE/2 {
		return errors.New("no more storage")
	}

	var wg sync.WaitGroup
	done := make(chan struct{})
	r.StartWriter(done, &wg)

	for blockId, pos := 0, 0; blockId < BLOCK_PER_DISK && pos < len(data); blockId += 1 {
		for diskId := 0; diskId < STRIPE_WIDTH && pos < len(data); diskId, pos = diskId+2, pos+BLOCK_SIZE {
			wg.Add(2)
			r.DiskArray[diskId].TaskChan <- NewTask(blockId, data[pos:min(len(data), pos+BLOCK_SIZE)])
			r.DiskArray[diskId+1].TaskChan <- NewTask(blockId, data[pos:min(len(data), pos+BLOCK_SIZE)])
		}
	}
	wg.Wait()
	close(done)
	return nil
}

func (r *Raid10) Read(length int) (string, error) {

	done := make(chan struct{})
	r.StartReader(done)

	for blockId, pos := 0, 0; blockId < BLOCK_PER_DISK && pos < length; blockId += 1 {
		for diskId := 0; diskId < STRIPE_WIDTH && pos < length; diskId, pos = diskId+2, pos+BLOCK_SIZE {
			r.DiskArray[diskId].TaskChan <- NewTask(blockId, nil)
		}
	}

	// Sequentially read back
	data := []byte{}
	for diskId, pos := 0, 0; pos < length; diskId, pos = (diskId+2)%STRIPE_WIDTH, pos+BLOCK_SIZE {
		block := <-r.DiskArray[diskId].ReadChan
		data = append(data, block.Data[:min(BLOCK_SIZE, length-pos)]...)
	}
	close(done)
	return string(data), nil

	// data := []byte{}
	// for blockId, pos := 0, 0; blockId < BLOCK_PER_DISK && pos < length; blockId += 1 {
	// 	for diskId := 0; diskId < STRIPE_WIDTH && pos < length; diskId, pos = diskId+2, pos+BLOCK_SIZE {
	// 		block := r.DiskArray[diskId].ReadBlock(blockId)
	// 		data = append(data, block.Data[:min(BLOCK_SIZE, length-pos)]...)
	// 	}
	// }

	// return string(data), nil
}
