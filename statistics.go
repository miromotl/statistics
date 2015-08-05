package main

import (
    "fmt"
    "log"
    "net/http"
    "sort"
    "strconv"
    "strings"
    "math"
)

// Constants with html code for our web page
const (
    pageTop = `
        <!DOCTYPE html>
        <html lang="en">
        <head>
            <title>Statistics</title>
            <meta charset="utf-8">
            <meta name="viewport" content="width=device-width, initial-scale=1">
            <!-- Latest compiled and minified CSS -->
            <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.5/css/bootstrap.min.css">
            <!-- Latest compiled and minified JavaScript -->
            <script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.5/js/bootstrap.min.js"></script>
        </head>
        <body>
            <div class="container">
                <h2>Statistics</h2>
                <p>Computes basic statistics for a given list of numbers</p>`
    form = `
                <form role="form" action="/" method="POST">
                    <div class="from-group">
                        <label for="numbers">Numbers (comma or space-separated):</label>
                        <input type="text" class="form-control" id="numbers" name="numbers" value="" placeholder="(e.g. 1 2 3)" />
                    </div>
                    <br />
                    <button type="submit" class="btn btn-success" >Calculate</button>
                </form>`
    pageBottom = `
            </div>
        </body>
        </html>`
    anError = `<br /><p class="text-danger">%s</p>`
)

// Struct holding the user's numbers and the resulting statistics
type statistics struct {
    numbers []float64
    mean    float64
    median  float64
    σ  float64
}

func main() {
    // Setup the web server handling the requests
    http.HandleFunc("/", homePage)
    if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
        log.Fatal("failed to start server", err)
    }
}

// Handling the call to the home page; i.e. handling everything because there
// is no other page!
func homePage(writer http.ResponseWriter, request *http.Request) {
    err := request.ParseForm() // Must be called before writing the response
    fmt.Fprint(writer, pageTop, form)
    
    if err != nil {
        fmt.Fprintf(writer, anError, err)
    } else {
        if numbers, message, ok := processRequest(request); ok {
            stats := getStats(numbers)
            fmt.Fprint(writer, formatStats(stats))
        } else if message != "" {
            fmt.Fprintf(writer, anError, message)
        }
    }
    
    fmt.Fprint(writer, pageBottom)
}

// Process the http request
func processRequest(request *http.Request) ([]float64, string, bool) {
    var numbers []float64
    
    if slice, found := request.Form["numbers"]; found && len(slice) > 0 {
        // There has been entered some text in the "numbers" input box
        // Replace commas with space
        text := strings.Replace(slice[0], ",", " ", -1)
        
        // Iterate over the white space separated text fields
        for _, field := range strings.Fields(text) {
            if x, err := strconv.ParseFloat(field, 64); err != nil {
                // Could not convert text to float64
                return numbers, "'" + field + "' is invalid", false
            } else {
                // We have a valid float64 number
                numbers = append(numbers, x)
            }
        }
    }
    
    if len(numbers) == 0 {
        // No data, page has just been called and shown
        return numbers, "", false
    }
    
    // A slice of float64 numbers
    return numbers, "", true
}

// Calculate the statistics
func getStats(numbers []float64) (stats statistics) {
    stats.numbers = numbers
    // The median function expects the numbers to be sorted
    sort.Float64s(stats.numbers)
    stats.mean = sum(numbers) / float64(len(numbers))
    stats.median = median(stats.numbers)
    stats.σ = σ(stats.numbers)
    
    return stats
}

// Calculate the sum of a slice of float64 numbers
func sum(numbers []float64) (total float64) {
    for _, x := range numbers {
        total += x
    }
    
    return total
}

// Calculate the median of a sorted slice of float64 numbers
func median(numbers []float64) float64 {
    middle := len(numbers) / 2
    result := numbers[middle]
    
    // Correct median for multiple of 2 length
    if len(numbers) % 2 == 0 {
        result = (result + numbers[middle-1]) / 2
    }
    
    return result
}

// Calculate the standard deviation of float64 numbers
func σ(numbers []float64) (σ float64) {
    
    n := len(numbers)
    
    if n < 2 {
        return σ
    }
    
    mean := sum(numbers) / float64(n)
    var s float64
    
    for _, x := range numbers {
        s += (x - mean) * (x - mean)
    }
    
    σ = math.Sqrt(s/float64(n-1))
    
    return σ
}

// Format the statistics html output
func formatStats(stats statistics) string {
    return fmt.Sprintf(`
        <div class="table-responsive">
            <table class="table">
                <thead>
                    <tr><th colspan="2">Results:</th></tr>
                </thead>
                <tbody>
                    <tr><td>Numbers</td><td>%v</td></tr>
                    <tr><td>Count</td><td>%d</td></tr>
                    <tr><td>Mean</td><td>%f</td></tr>
                    <tr><td>Median</td><td>%f</td></tr>
                    <tr><td>σ</td><td>%f</td></tr>
                </tbody>
            </table>
        </div>`,
        stats.numbers, len(stats.numbers), stats.mean, stats.median, stats.σ)
}