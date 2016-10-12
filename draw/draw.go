/*
	Package draw 
	Draws histograms
*/
package draw 

import (
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"../checkError"
	"../matrix"
)

  
func Hist(Mat *matrix.MatrixOut){
  	
	length := (uint64(len(Mat.GetData()))-Mat.GetY())
    v := make(plotter.Values, length)
    for i := range v {
        v[i] = Mat.GetOneData(uint64(i)) 
    }

    // Make a plot and set its title
    p, err := plot.New()
    checkError.Check(err)
    p.Title.Text = "Histogram of correlations"

    h, err := plotter.NewHist(v, 16)
    checkError.Check(err)
    // Normalize the area under the histogram to
    // sum to one
    h.Normalize(1)
    p.Add(h)

    // Save the plot to a PNG file.
    if err := p.Save(4*vg.Inch, 4*vg.Inch, "./outputData/hist.png"); err != nil {
        panic(err)
    }
}

