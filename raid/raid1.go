package raid

import (
	"errors"
	"fmt"
	"log"
	"sync"
)

type Raid1 struct {
	*RaidBase // mirrow of disk[i] is disk[i+1]
}

func NewRaid1() *Raid1 {
	if STRIPE_WIDTH%2 != 0 {
		log.Fatal(fmt.Errorf("number of disks have to be even"))
	}
	return &Raid1{
		RaidBase: NewRaidBase(),
	}
}

func (r *Raid1) Write(data []byte) error {

	if len(data) > TOTAL_STORAGE/2 {
		return errors.New("no more storage")
	}

	var wg sync.WaitGroup
	done := make(chan struct{})
	r.StartWriter(done, &wg)

	for diskId, pos := 0, 0; diskId < STRIPE_WIDTH && pos < len(data); diskId += 2 {
		for blockId := 0; blockId < BLOCK_PER_DISK && pos < len(data); blockId, pos = blockId+1, pos+BLOCK_SIZE {
			wg.Add(2)
			r.DiskArray[diskId].TaskChan <- NewTask(blockId, data[pos:min(len(data), pos+BLOCK_SIZE)])
			r.DiskArray[diskId+1].TaskChan <- NewTask(blockId, data[pos:min(len(data), pos+BLOCK_SIZE)])
		}
	}
	wg.Wait()
	close(done)
	return nil

}

func (r *Raid1) Read(length int) (string, error) {

	data := []byte{}

	for diskId, pos := 0, 0; diskId < STRIPE_WIDTH && pos < length; diskId += 2 {
		for blockId := 0; blockId < BLOCK_PER_DISK && pos < length; blockId, pos = blockId+1, pos+BLOCK_SIZE {
			block := r.DiskArray[diskId].ReadBlock(blockId)
			data = append(data, block.Data[:min(BLOCK_SIZE, length-pos)]...)
		}
	}

	return string(data), nil
}
