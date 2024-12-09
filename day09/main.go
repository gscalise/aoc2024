package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type DiskBlock struct {
	position int
	length   int
	file     int
}

func Print(l []DiskBlock) string {
	var builder strings.Builder

	for _, b := range l {
		fId := "."
		if b.file > -1 {
			fId = strconv.Itoa(b.file)
		}
		builder.WriteString(strings.Repeat(fId, b.length))
	}
	return builder.String()
}

func PrintPtrs(l []*DiskBlock) string {
	var builder strings.Builder

	for _, b := range l {
		fId := "."
		if b.file > -1 {
			fId = strconv.Itoa(b.file)
		}
		builder.WriteString(strings.Repeat(fId, b.length))
	}
	return builder.String()
}

func (b DiskBlock) checksum() int {
	if b.file == -1 {
		return 0
	} else {
		sum := 0

		for l := range b.length {
			sum += b.file * (b.position + l)
		}
		return sum
	}
}

func main1() {
	inputBytes, _ := os.ReadFile("input.txt")
	input := string(inputBytes)
	blockList := []*DiskBlock{}
	fileFirstBlocks := []*DiskBlock{}
	emptyBlocks := []*DiskBlock{}

	fileId := 0
	blockPos := 0
	for p, b := range input {
		blockLen, _ := strconv.Atoi(string(b))
		if blockLen > 0 {
			block := DiskBlock{
				position: blockPos,
				length:   blockLen,
				file:     -1,
			}
			if p%2 == 0 {
				block.file = fileId
				fileId += 1
				fileFirstBlocks = append(fileFirstBlocks, &block)
			} else {
				emptyBlocks = append(emptyBlocks, &block)
			}
			blockList = append(blockList, &block)
			blockPos += blockLen
		}
	}

	for _, fb := range slices.Backward(fileFirstBlocks) {
		if fb.position <= emptyBlocks[0].position {
			break
		}
		if fb.file != -1 {
			remainingFileLen := fb.length
			for remainingFileLen > 0 {
				if fb.position < emptyBlocks[0].position {
					fb.length = remainingFileLen
					break
				}
				if remainingFileLen == emptyBlocks[0].length {
					emptyBlocks[0].file = fb.file
					fb.file = -1
					remainingFileLen = 0
					emptyBlocks = emptyBlocks[1:]
				} else if remainingFileLen < emptyBlocks[0].length {
					origPos := emptyBlocks[0].position
					emptyBlocks[0].position += remainingFileLen
					emptyBlocks[0].length -= remainingFileLen
					for bIdx := range len(blockList) {
						if blockList[bIdx] == emptyBlocks[0] {
							blockList = slices.Insert(blockList, bIdx, &DiskBlock{
								position: origPos,
								length:   remainingFileLen,
								file:     fb.file,
							})
							break
						}
					}
					fb.file = -1
					remainingFileLen = 0
				} else {
					emptyBlockSize := emptyBlocks[0].length
					emptyBlocks[0].file = fb.file
					remainingFileLen -= emptyBlockSize
					emptyBlocks = emptyBlocks[1:]
				}
			}
		}
	}

	checksum := 0
	for _, b := range blockList {
		checksum += b.checksum()
	}
	fmt.Println(checksum)

	fmt.Println(PrintPtrs(blockList))
}

func main() {
	inputBytes, _ := os.ReadFile("input.txt")
	input := string(inputBytes)
	blockList := []*DiskBlock{}
	fileFirstBlocks := []*DiskBlock{}
	emptyBlocks := []*DiskBlock{}

	fileId := 0
	blockPos := 0
	for p, b := range input {
		blockLen, _ := strconv.Atoi(string(b))
		if blockLen > 0 {
			block := DiskBlock{
				position: blockPos,
				length:   blockLen,
				file:     -1,
			}
			if p%2 == 0 {
				block.file = fileId
				fileId += 1
				fileFirstBlocks = append(fileFirstBlocks, &block)
			} else {
				emptyBlocks = append(emptyBlocks, &block)
			}
			blockList = append(blockList, &block)
			blockPos += blockLen
		}
	}

	for _, fb := range slices.Backward(fileFirstBlocks) {
		fmt.Println("File ", fb.file)
		for ebIdx, eb := range emptyBlocks {
			if eb.position > fb.position {
				break
			}
			if eb.length == fb.length {
				eb.file = fb.file
				fb.file = -1
				if ebIdx == len(emptyBlocks)-1 {
					emptyBlocks = emptyBlocks[0:ebIdx]
				} else {
					emptyBlocks = append(emptyBlocks[0:ebIdx], emptyBlocks[ebIdx+1:]...)
				}
				break
			} else if eb.length > fb.length {
				origPos := eb.position
				eb.position += fb.length
				eb.length -= fb.length
				for bIdx := range len(blockList) {
					if blockList[bIdx] == eb {
						blockList = slices.Insert(blockList, bIdx, &DiskBlock{
							position: origPos,
							length:   fb.length,
							file:     fb.file,
						})
						break
					}
				}
				fb.file = -1
				break
			}
		}
	}

	checksum := 0
	for _, b := range blockList {
		checksum += b.checksum()
	}
	fmt.Println(checksum)

	fmt.Println(PrintPtrs(blockList))
}
