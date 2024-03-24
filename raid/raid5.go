package raid

import (
	"errors"
	"sync"
)

// Parity is Left Asymmetric
type Raid5 struct {
	*RaidBase
}

func NewRaid5() *Raid5 {
	return &Raid5{
		RaidBase: NewRaidBase(),
	}
}

func (r *Raid5) Write(data []byte) error {

	if r.RequiredStorage(data) > TOTAL_STORAGE {
		return errors.New("no more storage")
	}

	var wg sync.WaitGroup
	done := make(chan struct{})
	r.StartWriter(done, &wg)

	parityDiskId := STRIPE_WIDTH - 1
	xorBlock := make([]byte, BLOCK_SIZE)
	for blockId, pos := 0, 0; blockId < BLOCK_PER_DISK && pos < len(data); blockId += 1 {
		for diskId := 0; diskId < STRIPE_WIDTH && pos < len(data); diskId += 1 {
			if diskId == parityDiskId { // this block is for parity
				continue
			}
			writeBlock := data[pos:min(len(data), pos+BLOCK_SIZE)]
			for i := range writeBlock {
				xorBlock[i] ^= writeBlock[i]
			}
			wg.Add(1)
			r.DiskArray[diskId].TaskChan <- NewTask(blockId, data[pos:min(len(data), pos+BLOCK_SIZE)])
			pos += BLOCK_SIZE
		}
		wg.Add(1)
		r.DiskArray[parityDiskId].TaskChan <- NewTask(blockId, xorBlock)
		parityDiskId = (parityDiskId + STRIPE_WIDTH - 1) % STRIPE_WIDTH
	}

	wg.Wait()
	close(done)
	return nil
}

func (r *Raid5) Read(length int) (string, error) {

	data := []byte{}

	parityDiskId := STRIPE_WIDTH - 1
	for blockId, pos := 0, 0; blockId < BLOCK_PER_DISK && pos < length; blockId += 1 {
		for diskId := 0; diskId < STRIPE_WIDTH && pos < length; diskId += 1 {
			if diskId == parityDiskId {
				continue
			}
			block := r.DiskArray[diskId].ReadBlock(blockId)
			data = append(data, block.Data[:min(BLOCK_SIZE, length-pos)]...)
			pos += BLOCK_SIZE
		}
		parityDiskId = (parityDiskId + STRIPE_WIDTH - 1) % STRIPE_WIDTH
	}
	return string(data), nil
}

func (r *Raid5) RequiredStorage(data []byte) int {
	// How many blocks required
	requiredBlocks := len(data)/BLOCK_SIZE + len(data)%BLOCK_SIZE

	// Sacrifice 1 block for each stripe
	requiredBlocksWithParity := requiredBlocks / (STRIPE_WIDTH - 1) * STRIPE_WIDTH

	// The block required is not multiply of stripe size
	if requiredBlocks%(STRIPE_WIDTH-1) > 0 {
		requiredBlocksWithParity += 1
	}
	return requiredBlocksWithParity * BLOCK_SIZE
}
