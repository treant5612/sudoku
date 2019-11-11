package sudoku

import (
	"fmt"
	"unicode"
)

var wrongInputErr = fmt.Errorf("wrong puzzle input")
var impossibleErr = fmt.Errorf("impossibleErr puzzle")

type Puzzle struct {
	grid [81]int
}

func New(lines string) (p *Puzzle, err error) {
	puzzle := &Puzzle{}
	i := 0
	for _, r := range lines {
		if futilityRune(r) {
			continue
		}
		num, err := conv(r)
		if err != nil {
			return nil, wrongInputErr
		}
		puzzle.grid[i] = num
		i++
	}

	if puzzle.hasDuplicate() {
		return nil, impossibleErr
	}
	return puzzle, nil
}

//无意义数字
func futilityRune(r rune) bool {
	if unicode.IsSpace(r) || r == ',' {
		return true
	}
	return false
}
func conv(r rune) (int, error) {
	if r >= '0' && r <= '9' {
		return int(r - '0'), nil
	} else if r == '_' {
		return 0, nil
	}
	return 0, wrongInputErr
}

func (p *Puzzle) dup() *Puzzle {
	return &Puzzle{p.grid}
}

func (p *Puzzle) getRow(rowNum int) (row []int) {
	row = []int{}
	for i := 0; i < 9; i++ {
		num := p.grid[rowNum*9+i]
		if num != 0 {
			row = append(row, num)
		}
	}
	return row
}

func (p *Puzzle) getCol(colNum int) (col []int) {
	col = []int{}
	for i := 0; i < 9; i++ {
		num := p.grid[i*9+colNum]
		if num != 0 {
			col = append(col, num)
		}
	}
	return col
}

// 一个格子的左上角索引
var boxToIndex = [...]int{0, 3, 6, 27, 30, 33, 54, 57, 60}

// 一个格子的各个位置与左上角的索引差
var boxOffset = [...]int{0, 1, 2, 9, 10, 11, 18, 19, 20}

//数组各个索引所对应的格子序号
var boxOfIndex = [...]int{
	0, 0, 0, 1, 1, 1, 2, 2, 2, 0, 0, 0, 1, 1, 1, 2, 2, 2, 0, 0, 0, 1, 1, 1, 2, 2, 2,
	3, 3, 3, 4, 4, 4, 5, 5, 5, 3, 3, 3, 4, 4, 4, 5, 5, 5, 3, 3, 3, 4, 4, 4, 5, 5, 5,
	6, 6, 6, 7, 7, 7, 8, 8, 8, 6, 6, 6, 7, 7, 7, 8, 8, 8, 6, 6, 6, 7, 7, 7, 8, 8, 8,
}

//获取目标格子内已有的数字
func (p *Puzzle) getBox(boxNum int) (box []int) {
	box = []int{}
	startIndex := boxToIndex[boxNum]
	for i := 0; i < 9; i++ {
		num := p.grid[startIndex+boxOffset[i]]
		if num != 0 {
			box = append(box, num)
		}
	}
	return box
}

//判断切片内是否存在重复值
func duplicate(s []int) bool {
	for i := 0; i < len(s); i++ {
		for j := i + 1; j < len(s); j++ {
			if s[i] == s[j] {
				return true
			}
		}
	}
	return false
}

//行列格子中是否有重复数字
func (p *Puzzle) hasDuplicate() bool {
	for i := 0; i < 9; i++ {
		if duplicate(p.getRow(i)) || duplicate(p.getCol(i)) || duplicate(p.getBox(i)) {
			return true
		}
	}
	return false
}

// 扫描数组，依次寻找有在唯一解的位置。
// 如某位置无解，返回错误；如果某处有多个可能的数字，返回该位置和可能的数字列表
// 如果已经解决了问题，则返回的位置为-1
func (p *Puzzle) scan() (pos int, possibilities []int, err error) {
	for i, v := range p.grid {
		if v != 0 {
			continue
		}
		row, col := i/9, i%9
		possibleNums := p.possibleNums(row, col)
		switch len(possibleNums) {
		case 0:
			err = impossibleErr
			return
		case 1:
			p.grid[i] = possibleNums[0]
		default:
			return i, possibleNums, nil
		}
	}
	pos = -1
	return
}

//返回指定位置的所有可能数字
func (p *Puzzle) possibleNums(row int, col int) []int {
	possibleNums := make([]int, 0, 9)
	box := boxOfIndex[row*9+col]
	existsNums := make(map[int]bool)
	putNumsInMap := func(nums ...[]int) {
		for _, s := range nums {
			for _, v := range s {
				existsNums[v] = true
			}
		}
	}
	putNumsInMap(p.getRow(row), p.getCol(col), p.getBox(box))
	for i := 1; i <= 9; i++ {
		if !existsNums[i] {
			possibleNums = append(possibleNums, i)
		}
	}
	return possibleNums
}

func (p *Puzzle) Solve() (solved *Puzzle, err error) {
	puz := p.dup()
	pos, possibilities, err := puz.scan()
	if err != nil {
		return nil, err
	}
	if pos == -1 {
		return puz, nil
	}
	for _, v := range possibilities {
		puz.grid[pos] = v
		puz, solveErr := puz.Solve()
		if solveErr == nil {
			return puz, nil
		}
	}
	return nil, impossibleErr
}

func (p Puzzle) String() string {
	str := ""
	for i := 0; i < 9; i++ {
		str += fmt.Sprintln(p.grid[i*9 : i*9+9])
	}
	return str
}
