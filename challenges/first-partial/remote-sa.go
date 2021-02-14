package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
)

type Point struct {
	X, Y float64
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

//generatePoints array
func generatePoints(s string) ([]Point, error) {

	points := []Point{}

	s = strings.Replace(s, "(", "", -1)
	s = strings.Replace(s, ")", "", -1)
	vals := strings.Split(s, ",")
	if len(vals) < 2 {
		return []Point{}, fmt.Errorf("Point [%v] was not well defined", s)
	}

	var x, y float64

	for idx, val := range vals {

		if idx%2 == 0 {
			x, _ = strconv.ParseFloat(val, 64)
		} else {
			y, _ = strconv.ParseFloat(val, 64)
			points = append(points, Point{x, y})
		}
	}
	return points, nil
}

func getDistance(a Point, b Point) float64 {
	distance := math.Sqrt(math.Pow(b.X-a.X, 2) + math.Pow(b.Y-a.Y, 2))
	return distance
}

// getArea gets the area inside from a given shape
func getArea(points []Point) float64 {
	area := 0.0
	for i := 0; i < len(points)-1; i++ {
		area += points[i].X*points[(i+1)].Y - points[i].Y*points[(i+1)].X
	}
	area += points[len(points)-1].X*points[0].Y - points[len(points)-1].Y*points[0].X
	area /= 2
	return math.Abs(area)
}

// getPerimeter gets the perimeter from a given array of connected points
func getPerimeter(points []Point) float64 {
	perimeter := 0.0
	for i := 0; i < len(points)-1; i++ {
		perimeter += getDistance(points[i], points[i+1])
	}
	perimeter += getDistance(points[len(points)-1], points[0])
	return perimeter
}

func onSegment(a Point, b Point, c Point) bool {
	if b.X <= math.Max(a.X, c.X) && b.X >= math.Min(a.X, c.X) &&
		b.Y <= math.Max(a.Y, c.Y) && b.Y >= math.Min(a.Y, c.Y) {
		return true
	}
	return false
}

func getOrientation(a Point, b Point, c Point) int {
	orientation := (b.Y-a.Y)*(c.X-b.X) -
		(b.X-a.X)*(c.Y-b.Y)

	if orientation == 0 {
		return 0
	} else if orientation > 0 {
		return 1
	} else {
		return 2
	}
}

func hasCollision(points []Point) bool {
	orientations := []int{}
	for i := 0; i < len(points); i++ {
		orientation := getOrientation(points[i], points[(i+1)%len(points)], points[(i+2)%len(points)])
		orientations = append(orientations, orientation)
	}

	for i := 0; i < len(orientations); i++ {
		if orientations[i] != orientations[(i+1)%len(orientations)] && orientations[(i+2)%len(orientations)] != orientations[(i+3)%len(orientations)] {
			return true
		} else if orientations[i] == 0 && onSegment(points[i], points[(i+1)%len(points)], points[(i+2)%len(points)]) {
			return true
		}
	}
	return false
}

// handler handles the web request and reponds it
func handler(w http.ResponseWriter, r *http.Request) {

	var vertices []Point
	for k, v := range r.URL.Query() {
		if k == "vertices" {
			points, err := generatePoints(v[0])
			if err != nil {
				fmt.Fprintf(w, fmt.Sprintf("error: %v", err))
				return
			}
			vertices = points
			break
		}
	}

	// Results gathering
	area := getArea(vertices)
	perimeter := getPerimeter(vertices)

	// Logging in the server side
	log.Printf("Received vertices array: %v", vertices)

	// Response construction
	if len(vertices) < 3 {
		response := fmt.Sprint("ERROR - Your shape is not compliying with the minimum number of vertices.\n")
		fmt.Fprintf(w, response)
	} else if hasCollision(vertices) {
		response := fmt.Sprint("ERROR - Your shape must not have line collisions.\n")
		fmt.Fprintf(w, response)
	} else {
		response := fmt.Sprintf("Welcome to the Remote Shapes Analyzer\n")
		response += fmt.Sprintf(" - Your figure has : [%v] vertices\n", len(vertices))
		response += fmt.Sprintf(" - Vertices        : %v\n", vertices)
		response += fmt.Sprintf(" - Perimeter       : %v\n", perimeter)
		response += fmt.Sprintf(" - Area            : %v\n", area)

		// Send response to client
		fmt.Fprintf(w, response)
	}

}
