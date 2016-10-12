    /*
    	Louiza HANIS
    */
    package main
    import (
    	"flag"
    	"fmt"
    	"log"
    	"os"
    	"runtime/pprof"
    	"time"
    	"./catalog"
    	"./correlation"
    //	"./draw"
    	"./matrix"
    )
    func main() {
    	// Set default matrix
    	defaultFile := `./data/Mat200.txt`
    	InputFilePath := flag.String("i", defaultFile, "Input matrix file")
    	defaultOutputFile := `./outputData/Matrix.out.txt`
    	OutputFilePath := flag.String("o", defaultOutputFile, "Output matrix file")
    	CpuProfile := flag.String("cpuprofile", "", "write cpu profile to file")
    	flag.Parse()
    	if *CpuProfile != "" {
            profile := "./outputData/" + *CpuProfile + ".pprof"
    		f, err := os.Create(profile)
    		if err != nil {
    			log.Fatal(err)
    		}
    		pprof.StartCPUProfile(f)
    		defer pprof.StopCPUProfile()
    	}
    	//Extract Genes names in the catalogue
    	genesNames := catalog.ExtractGenesNames("./data/genes.txt")
    	Mat := matrix.ReadFile(*InputFilePath)
    	t0 := time.Now()
    	_ = correlation.ComputeCorrMulti(Mat, genesNames, *OutputFilePath)
    	t1 := time.Since(t0)
    	fmt.Println("Time took", t1)
      /*  draw.Hist(MatCorr)*/
    }

