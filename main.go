package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"gonum.org/v1/plot/vg"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

func main() {

	var class [][]float64

	//重み
	w := []float64{18.6, -11.1, 21.1}
	w0 := []float64{18.6, -11.1, 21.1}

	//Dot := matrix.Product
	//Inv := matrix.Inverse

	// 図の生成
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	//任意の点
	dots := make(plotter.XYs, 2)

	//クラス1
	x1, y1 := 8.0, 2.0
	dots[0].X = x1
	dots[0].Y = y1

	//クラス2
	x2, y2 := 3.0, 6.0
	dots[1].X = x2
	dots[1].Y = y2

	//各クラスのサンプル
	n := 100
	class1, plotdata1 := randomPoints(n, x1, y1)
	class2, plotdata2 := randomPoints(n, x2, y2)
	class = append(class1, class2...)
	//mClass := matrix.MakeDenseMatrix(class, n*2, 3) // 配列から行列に変換

	//教師データ作成
	b := make([]float64, n*2)
	for i := 0; i < n*2; i++ {
		if i >= 100 {
			b[i] = 10
		} else {
			b[i] = -10
		}
	}
	max, min := 0.0, 199.0
	var errGraph []float64
	var beforeError float64
	var afterError float64
	afterError = 10000

	for {
		rand := randomCount(max, min)
		//fmt.Println(class[int(rand)])
		beforeError = afterError
		//err := errorFunc(w, class[int(rand)], b[int(rand)])
		//fmt.Print("誤差関数")
		//fmt.Println(err)
		errGraph = append(errGraph, beforeError)
		//if err < 0.0000001 {
		//	break
		//}
		w = weightCalc(w, class[int(rand)], 0.0001, b[int(rand)])
		fmt.Println(int(rand))
		afterError = errorFunc(w, class[int(rand)], b[int(rand)])
		if (beforeError-afterError)*(beforeError-afterError) < 0.0000001 {
			break
		}
	}

	fmt.Println(w)
	p2, err := plot.New()
	if err != nil {
		panic(err)
	}

	// Make a line plotter and set its style.
	l, err := plotter.NewLine(lineGraph(errGraph))
	if err != nil {
		panic(err)
	}
	l.LineStyle.Width = vg.Points(1)
	l.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
	l.LineStyle.Color = color.RGBA{B: 255, A: 255}

	p2.Add(l)
	p2.Title.Text = "Plotutil example"
	p2.X.Label.Text = "X"
	p2.Y.Label.Text = "Y"

	// Axis ranges
	//p2.X.Min = 0
	//p2.X.Max = 10
	//p2.Y.Min = 0
	//p2.Y.Max = 10

	p2.Legend.Add("line", l)
	// Save the plot to a PNG file.
	if err := p2.Save(4*vg.Inch, 4*vg.Inch, "points.png"); err != nil {
		panic(err)
	}

	//mb := matrix.MakeDenseMatrix(b, n*2, 1)

	//tClass := mClass.Transpose() // 転置行列
	//
	//w := Dot(tClass, mClass)
	//w = Inv(w).DenseMatrix()
	//w = Dot(w, tClass)
	//w = Dot(w, mb)
	//
	//a1 := w.Get(0, 0)
	//b1 := w.Get(1, 0)
	//c1 := w.Get(2, 0)
	//初期境界線のplot--------------------------------------------------------
	border := plotter.NewFunction(func(x float64) float64 {
		//x2 = -(w1 / w2)*x1 - w0 / w2
		return -(w0[1]/w0[2])*x - (w0[0] / w0[2])
	})
	border.Color = color.RGBA{B: 155, A: 5}
	//----------------------------------------------------------------------

	//最終境界線のplot--------------------------------------------------------
	lastBorder := plotter.NewFunction(func(x float64) float64 {
		//x2 = -(w1 / w2)*x1 - w0 / w2
		return -(w[1]/w[2])*x - (w[0] / w[2])
	})
	lastBorder.Color = color.RGBA{B: 255, A: 255}
	//----------------------------------------------------------------------

	//label
	p.Title.Text = "Points Example"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	// Draw a grid behind the data
	p.Add(plotter.NewGrid())

	// Make a scatter plotter and set its style.
	s, err := plotter.NewScatter(plotdata1)
	if err != nil {
		panic(err)
	}

	y, err := plotter.NewScatter(plotdata2)
	if err != nil {
		panic(err)
	}

	r, err := plotter.NewScatter(dots)
	if err != nil {
		panic(err)
	}

	s.GlyphStyle.Color = color.RGBA{R: 255, B: 128, A: 55}
	y.GlyphStyle.Color = color.RGBA{R: 155, B: 128, A: 255}
	r.GlyphStyle.Color = color.RGBA{R: 128, B: 0, A: 0}
	p.Add(s)
	p.Add(y)
	p.Add(r)
	p.Add(lastBorder)
	p.Add(border)
	p.Legend.Add("class1", s)
	p.Legend.Add("class2", y)

	// Axis ranges
	p.X.Min = 0
	p.X.Max = 10
	p.Y.Min = 0
	p.Y.Max = 10

	// Save the plot to a PNG file.
	if err := p.Save(6*vg.Inch, 6*vg.Inch, "report.png"); err != nil {
		panic(err)
	}
}

func randomCount(min, max float64) float64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Float64()*(max-min) + min
}

//ガウス分布
func random(axis float64) float64 {
	//分散
	dispersion := 1.0
	rand.Seed(time.Now().UnixNano())
	return rand.NormFloat64()*dispersion + axis
}

////学習データの生成
//func randomPoints(n int, x, y float64) ([]float64, plotter.XYs) {
//	matrix := make([][]float64, n)
//	pts := make(plotter.XYs, n)
//	var gyo []float64
//
//	for i := range matrix {
//		l := random(x)
//		m := random(y)
//		gyo = append(gyo, 1.0)
//		gyo = append(gyo, l)
//		gyo = append(gyo, m)
//		pts[i].X = l
//		pts[i].Y = m
//	}
//	return gyo, pts
//}

//学習データの生成
func randomPoints(n int, x, y float64) ([][]float64, plotter.XYs) {
	matrix := make([][]float64, n)
	pts := make(plotter.XYs, n)
	for i := range matrix {
		l := random(x)
		m := random(y)
		matrix[i] = []float64{1.0, l, m}
		pts[i].X = l
		pts[i].Y = m
	}
	return matrix, pts
}

//学習
func train(index, n int) []float64 {
	var array []float64

	for i := 0; i < n; i++ {
		array = append(array, float64(index))
	}

	return array
}

func weightCalc(w, x []float64, p, b float64) []float64 {
	innerProduct := innerProduct(w, x)
	//fmt.Println(innerProduct)
	//fmt.Println("内積")
	w[0] = w[0] - p*((innerProduct-b)*x[0])
	w[1] = w[1] - p*((innerProduct-b)*x[1])
	w[2] = w[2] - p*((innerProduct-b)*x[2])

	return w
}

//内積計算
func innerProduct(w, x []float64) (f float64) {
	if len(w) != len(x) {
		panic("エラーですよ")
	}

	for i := range w {
		f += w[i] * x[i]
	}

	return
}

func errorFunc(w, x []float64, b float64) float64 {
	err := innerProduct(w, x) - b
	err = err * err

	return err
}

// randomPoints returns some random x, y points.
func lineGraph(n []float64) plotter.XYs {
	pts := make(plotter.XYs, len(n))
	for i, m := range n {
		//fmt.Println(m)
		pts[i].X = float64(i)
		pts[i].Y = m
	}
	return pts
}
