/*
	Package correlation
	Compute correlation of matrix genes
*/
package correlation

import (
	"bufio"
	"fmt"
	//	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"

	"../checkError"
	"../matrix"
	"github.com/montanaflynn/stats"
)

type subBlock struct {
	id       uint64
	beginRow uint64
	endRow   uint64
	n        uint64
}

type vertexData struct {
	sample1, state, sample2 string
}

// Sequential Compute

func ComputeCorrSeq(Mat *matrix.Matrix, genesNames []string, outputFile string) *matrix.MatrixOut {

	RDFfile, error := os.OpenFile(
		"./outputData/"+outputFile+".nq",
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)

	checkError.Check(error)
	defer RDFfile.Close()

	RDFwriter := bufio.NewWriter(RDFfile)
	var s1, s2, sample1 string
	state := "Unknown"
	y := Mat.GetY()
	size := y * uint64(float64((y+1))*0.5)
	Matcorr := matrix.NewMatrixOut(size, y)

	var i, j, k uint64

	vi := make([]float64, y)
	for i = 0; i < y; i++ {
		name1 := genesNames[i]
		subName1 := strings.SplitN(name1, "_", 2)
		sample1 = subName1[0]
		vi = Mat.GetDataLine(i)
		for j = 0; j < y; j++ {
			if i < j {
				corr, _ := stats.Correlation(vi, Mat.GetDataLine(j))
				Matcorr.SetData(k, corr)
				name2 := genesNames[j]
				subName2 := strings.SplitN(name2, "_", 2)
				sample2 := subName2[0]

				switch {
				case corr >= 0.9:
					state = "Strong"
				case corr >= 0.8 && corr < 0.9:
					state = "Medium"
				default:
					state = "Weak"
				}

				s1 = fmt.Sprintf("<%v> <%v> <%v> . \n", sample1, state, sample2)
				fmt.Fprint(RDFwriter, s1)
				k++
			}
		}
		cluster1 := strings.Trim(fmt.Sprint(subName1[1:]), "[]")
		s2 = fmt.Sprintf("<%v> <contains> <%v> . \n", cluster1, sample1)
		fmt.Fprint(RDFwriter, s2)
	}

	RDFwriter.Flush()

	return Matcorr
}

// Multi Threading Compute
func ComputeCorrMulti(inputMat *matrix.Matrix, genesNames []string, OutputFile string) *matrix.MatrixOut {

	RDFfile, error := os.OpenFile(
		"./outputData/"+OutputFile+".nq",
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)

	checkError.Check(error)
	defer RDFfile.Close()

	RDFwriter := bufio.NewWriter(RDFfile)

	y := inputMat.GetY()
	//correlation build a symetric-positive matrix - so split in half
	yd := (y * (y + 1)) / 2
	ym := (y * (y + 1)) % 2

	//round
	if ym != 0 {
		yd = yd + 1
	}

	matCorr := matrix.NewMatrixOut(yd, y)

	var nbThrd_uint64, i, ops uint64

	nbThrd := runtime.GOMAXPROCS(0)
	nbThrd_uint64 = uint64(nbThrd)
	subBlockThread := make([]subBlock, nbThrd)

	Q := y / nbThrd_uint64
	R := y % nbThrd_uint64

	var b uint64
	var RDFLoggerChan chan vertexData
	for i = 0; i < nbThrd_uint64; i++ {
		subBlockThread[i].id = i
		subBlockThread[i].n = Q
		subBlockThread[i].beginRow = b

		if R != 0 && i <= R-1 {
			subBlockThread[i].n += 1
		}

		subBlockThread[i].endRow = subBlockThread[i].beginRow + subBlockThread[i].n
		b = subBlockThread[i].endRow

	}

	RDFLoggerChan = make(chan vertexData, yd)

	var wg, wgLogger sync.WaitGroup
	//	RDFLoggerChan := make(chan vertexData, subBlockThread.n) //changer la taille du buffer

	wg.Add(nbThrd)
	for i = 0; i < nbThrd_uint64; i++ {
		go Compute(&ops, &wg, RDFLoggerChan, inputMat, y, matCorr, subBlockThread[i], genesNames)
	}

	wgLogger.Add(1)
	go runRDFLoggerChan(RDFLoggerChan, RDFwriter, &wgLogger)

	wg.Wait()
	close(RDFLoggerChan) //permet de sortir du range

	wgLogger.Wait()

	RDFwriter.Flush()

	return matCorr
}

func runRDFLoggerChan(c chan vertexData, RDFwriter *bufio.Writer, wg *sync.WaitGroup) {
	defer wg.Done()
	for req := range c {
		s1 := fmt.Sprintf("<%v> <%v> <%v> . \n", req.sample1, req.state, req.sample2)
		RDFwriter.WriteString(s1)
	}
}

func Compute(ops *uint64, wg *sync.WaitGroup, RDFlogger chan vertexData, inputMat *matrix.Matrix, y uint64, matCorr *matrix.MatrixOut, thread subBlock, genesNames []string) {

	vi := make([]float64, y)
	var vD1, vD2 vertexData

	vD1.state = "Unknown"
	vD2.state = "contains"
	for i := thread.beginRow; i < thread.endRow; i++ {
		name1 := genesNames[i]
		subName1 := strings.SplitN(name1, "_", 2)
		vD1.sample1 = subName1[0]
		cluster1 := strings.Trim(fmt.Sprint(subName1[1:]), "[]")
		vi = inputMat.GetDataLine(i)
		for j := uint64(0); j < y; j++ {
			if i < j {
				corr, _ := stats.Correlation(vi, inputMat.GetDataLine(j))
				opsFinal := atomic.AddUint64(ops, 1) - 1
				matCorr.SetData(opsFinal, corr)
				name2 := genesNames[j]
				subName2 := strings.SplitN(name2, "_", 2)
				vD1.sample2 = subName2[0]

				switch {
				case corr >= 0.9:
					vD1.state = "Strong"
				case corr >= 0.8 && corr < 0.9:
					vD1.state = "Medium"
				default:
					vD1.state = "Weak"
				}
				RDFlogger <- vD1
			}
		}

		vD2.sample1 = cluster1
		vD2.sample2 = vD1.sample1
		RDFlogger <- vD2
	}
	defer wg.Done()

}
