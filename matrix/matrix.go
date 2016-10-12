/*
	Package matrix
	Implements matrix
*/
package matrix

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type MatrixOut struct {
	data  []float64
	xSize uint64
	ySize uint64
	genes []uint64
}

type Matrix struct {
	*MatrixOut
	samples []string
}

func newMatrixOut(y uint64, x uint64) *MatrixOut {
	M := new(MatrixOut)
	M.xSize = x
	M.ySize = y
	M.data = make([]float64, x*y)
	M.genes = make([]uint64, y)
	return M
}

func NewMatrixOut(x uint64, y uint64) *MatrixOut {
	M := new(MatrixOut)
	M.xSize = y 
	M.ySize = y
	M.data = make([]float64, x)
	return M
}

func NewMatrix(x uint64, y uint64) *Matrix {
	M := &Matrix{newMatrixOut(x, y), make([]string, x)}
	return M
}

func (Mat *MatrixOut) SetY(y uint64) {
	Mat.ySize = y
}

func (Mat *Matrix) SetSamples(s []string) {
	Mat.samples = s
}

func (Mat *MatrixOut) SetGenes(i uint64, g uint64) {
	Mat.genes[i] = g
}

func (Mat *MatrixOut) SetVecGenes(g []uint64) {
	Mat.genes = g
}

func (Mat *MatrixOut) SetData(x uint64, val float64) {
	Mat.data[x] = val
}

func (Mat *MatrixOut) SetVecData(d []float64) {
	Mat.data = d
}

func (Mat *MatrixOut) GetX() uint64 {
	return Mat.xSize
}

func (Mat MatrixOut) GetY() uint64 {
	return Mat.ySize
}

func (Mat Matrix) GetSamples() []string {
	return Mat.samples
}

func (Mat MatrixOut) GetGenes() []uint64 {
	return Mat.genes
}

func (Mat MatrixOut) GetGene(x uint64) uint64 {
	return Mat.genes[x]
}

func (Mat MatrixOut) GetPartGenes(x uint64) []uint64 {
	return Mat.genes[:x]
}

func (Mat *MatrixOut) GetData() []float64 {
	return Mat.data
}

func (Mat *MatrixOut) GetPartData(x uint64) []float64 {
	return Mat.data[:x]
}

func (Mat *MatrixOut) GetOneData(x uint64) float64 {
	return Mat.data[x]
}

//
func (Mat *MatrixOut) GetDataLine(y uint64) []float64 {
	xSize := Mat.GetX()
	v := make([]float64, xSize)

	yS := y * xSize
	for x, _ := range v {
		v[x] = Mat.GetOneData(yS + uint64(x))
	}
	return v
}

func (Mat *MatrixOut) GetDataLPartC(y uint64, begin uint64, end uint64, n uint64) []float64 {
	v := make([]float64, n)
	var x uint64
	y_ := y * Mat.GetX()
	for x = begin; x < end; x++ {
		v[x] = Mat.GetOneData(y_ + x)
	}
	return v
}

// Index the 1D matrix
func Index(i uint64, j uint64, n uint64) uint64 {
	return i*n + j
}

// Read file and return the matrix
func ReadFile(Filename string) *Matrix {
	var (
		i, Lcount, k uint64
		Words        []string
		Line         string
	)

	Inputfile, err := os.Open(Filename)
	if err != nil {
		panic(err.Error())
	}
	defer Inputfile.Close()

	Inputfile2, err := os.Open(Filename)
	if err != nil {
		panic(err.Error())
	}
	defer Inputfile2.Close()
	reader1 := bufio.NewReader(Inputfile)

	Firstline, err := reader1.ReadString('\n')
	samples := strings.Fields(Firstline)
	//  Mat.Set_x(uint64(len(samples)))

	scanner1 := bufio.NewScanner(reader1)

	Lcount = 1
	for scanner1.Scan() {
		Lcount++
	}

	Mat := NewMatrix(Lcount-1, uint64(len(samples)))
	Mat.SetSamples(samples)

	reader2 := bufio.NewReader(Inputfile2)
	scanner2 := bufio.NewScanner(reader2)

	_, err = reader2.ReadString('\n')

	i = 0

	fmt.Printf("Reading matrix file [%v] ... \n", Filename)
	for scanner2.Scan() {
		Line = scanner2.Text()
		Words = strings.Fields(Line)
		g, _ := strconv.ParseUint(Words[0], 10, 64)
		Mat.MatrixOut.SetGenes(i, g)
		for j, w := range Words[1:] {
			k = uint64(j)
			Word, err := strconv.ParseFloat(w, 64)
			if err == nil {
				Mat.MatrixOut.SetData(Index(i, k, Mat.GetX()), Word)
			}
		}
		i++
	}
	return Mat
}

// Discard variables that have less than min_non_zero samples
func (Mat *Matrix) DiscardVars(min_non_zero uint64) {

	var num_non_zero, j, num_vars_kept, i uint64

	num_vars_kept = 0
	for i = 0; i < Mat.GetY(); i++ {
		num_non_zero = 0
		j = 0

		for j < Mat.GetX() && num_non_zero < min_non_zero {
			if Mat.GetOneData(Index(i, j, Mat.GetX())) != 0.0 {
				num_non_zero++
			}
			j++
		}

		if num_non_zero >= min_non_zero {
			for j = 0; j < Mat.GetX(); j++ {
				Mat.SetData(Index(num_vars_kept, j, Mat.GetX()), Mat.GetOneData(Index(i, j, Mat.GetX())))
				Mat.SetGenes(num_vars_kept, Mat.GetGene(i))
			}

			num_vars_kept++
		}
	}

	Mat.SetY(num_vars_kept)
	Mat.SetVecData(Mat.GetPartData(Mat.GetX() * Mat.GetY()))
	Mat.SetVecGenes(Mat.GetPartGenes(Mat.GetY()))
}
