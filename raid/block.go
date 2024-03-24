package raid

type Data = [BLOCK_SIZE]byte

type Block struct {
	Id int
	Data
}

func NewBlock(id int) Block {
	return Block{
		Id:   id,
		Data: [BLOCK_SIZE]byte{},
	}
}
