package raid

import "sync"

const (
	STRIPE_WIDTH  = 8 // eq to how many disks
	TOTAL_STORAGE = BLOCK_SIZE * BLOCK_PER_DISK * STRIPE_WIDTH
)

type Raid interface {
	Write(data []byte) error
	Read(length int) (string, error)
	ClearDisk(id int) error
}

type RaidBase struct {
	DiskArray [STRIPE_WIDTH]*Disk
}

func NewRaidBase() *RaidBase {
	diskArray := [STRIPE_WIDTH]*Disk{}
	for i := range diskArray {
		diskArray[i] = NewDisk()
	}
	return &RaidBase{
		DiskArray: diskArray,
	}
}

func (r *Raid0) ClearDisk(id int) error {
	if id < 0 || id >= STRIPE_WIDTH {
		return ErrIndexOutOfBound
	} else {
		r.DiskArray[id].Clear()
	}
	return nil
}

func (r *RaidBase) StartWriter(done <-chan struct{}, wg *sync.WaitGroup) {
	// start worker for each Disk
	for diskId := range r.DiskArray {
		go func() {
			for {
				select {
				case task := <-r.DiskArray[diskId].TaskChan:
					r.DiskArray[diskId].AsyncWriteBlock(
						task.BlockId,
						task.Data,
						wg,
					)
				case <-done:
					return
				}
			}
		}()
	}
}

func (r *RaidBase) StartReader(done <-chan struct{}) {

	// start worker for each Disk
	for diskId := range r.DiskArray {
		go func() {

			for {
				select {
				case task := <-r.DiskArray[diskId].TaskChan:
					r.DiskArray[diskId].ReadChan <- r.DiskArray[diskId].ReadBlock(
						task.BlockId,
					)
				case <-done:
					return
				}
			}
		}()
	}
}
