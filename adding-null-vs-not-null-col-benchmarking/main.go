package main

import (
	"database/sql"
	"fmt"
	"log"
	"runtime"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	mysqlConnStr = "root@tcp(127.0.0.1:3306)/test"
	tableName    = "alter_benchmark"
	rowCount     = 100000
	iterations   = 5
)

type BenchmarkResult struct {
	Iteration    int
	NullTime     time.Duration
	NotNullTime  time.Duration
	MemoryBefore runtime.MemStats
	MemoryAfter  runtime.MemStats
}

func main() {
	db, err := sql.Open("mysql", mysqlConnStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	err = db.Ping()
	if err != nil {
		log.Fatal("Could not connect to MySQL:", err)
	}

	var version string
	err = db.QueryRow("SELECT VERSION()").Scan(&version)
	if err != nil {
		log.Printf("Could not get MySQL version: %v", err)
	} else {
		fmt.Printf("MySQL Version: %s\n", version)
	}

	var results []BenchmarkResult

	fmt.Println("Running MySQL ALTER TABLE benchmarks...")
	fmt.Printf("Testing with table size: %d rows\n", rowCount)
	fmt.Printf("Number of iterations: %d\n", iterations)

	for i := 0; i < iterations; i++ {
		fmt.Printf("\n=== Iteration %d ===\n", i+1)

		setupTestTable(db)

		var result BenchmarkResult
		result.Iteration = i + 1
		runtime.GC() // Force garbage collection for more accurate memory measurements
		runtime.ReadMemStats(&result.MemoryBefore)

		// Benchmark NULL column addition
		fmt.Println("Testing NULL column addition...")
		start := time.Now()
		_, err = db.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN null_col INT NULL", tableName))
		if err != nil {
			log.Fatal("NULL column addition failed:", err)
		}
		result.NullTime = time.Since(start)

		// Reset table for fair comparison
		_, err = db.Exec(fmt.Sprintf("ALTER TABLE %s DROP COLUMN null_col", tableName))
		if err != nil {
			log.Fatal("Failed to drop null column:", err)
		}

		// Benchmark NOT NULL column addition
		fmt.Println("Testing NOT NULL column addition...")
		start = time.Now()
		_, err = db.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN not_null_col INT NOT NULL DEFAULT 0", tableName))
		if err != nil {
			log.Fatal("NOT NULL column addition failed:", err)
		}
		result.NotNullTime = time.Since(start)

		runtime.GC()
		runtime.ReadMemStats(&result.MemoryAfter)
		results = append(results, result)

		printIterationResult(result)
		cleanup(db)
	}

	printFinalResults(results)
	printAnalysis(results)
}

func printIterationResult(result BenchmarkResult) {
	fmt.Println("\nIteration Results:")
	fmt.Printf("ADD NULL COLUMN: %v\n", result.NullTime)
	fmt.Printf("ADD NOT NULL COLUMN: %v\n", result.NotNullTime)

	if result.NullTime > 0 {
		ratio := float64(result.NotNullTime.Nanoseconds()) / float64(result.NullTime.Nanoseconds())
		fmt.Printf("Ratio (NOT NULL / NULL): %.1fx\n", ratio)
	}

	fmt.Printf("Memory Alloc Before: %.2f MiB\n", float64(result.MemoryBefore.Alloc)/1024/1024)
	fmt.Printf("Memory Alloc After: %.2f MiB\n", float64(result.MemoryAfter.Alloc)/1024/1024)
}

func printFinalResults(results []BenchmarkResult) {
	fmt.Println("\nFINAL RESULTS")
	fmt.Println("+-----------+----------------+-------------------+--------+----------------+----------------+")
	fmt.Println("| Iteration | NULL Time (ms) | NOT NULL Time (ms)| Ratio  | Mem Before(MiB)| Mem After (MiB)|")
	fmt.Println("+-----------+---------------+-------------------+--------+----------------+----------------+")

	var totalNull, totalNotNull time.Duration
	var totalRatio float64
	validRatios := 0

	for _, r := range results {
		var ratio float64
		if r.NullTime > 0 {
			ratio = float64(r.NotNullTime.Nanoseconds()) / float64(r.NullTime.Nanoseconds())
			totalRatio += ratio
			validRatios++
		}

		fmt.Printf("| %9d | %14.2f | %17.2f | %6.1fx | %14.2f | %14.2f |\n",
			r.Iteration,
			r.NullTime.Seconds()*1000,
			r.NotNullTime.Seconds()*1000,
			ratio,
			float64(r.MemoryBefore.Alloc)/1024/1024,
			float64(r.MemoryAfter.Alloc)/1024/1024)

		totalNull += r.NullTime
		totalNotNull += r.NotNullTime
	}

	// Calculate averages
	avgNull := totalNull / time.Duration(len(results))
	avgNotNull := totalNotNull / time.Duration(len(results))
	var avgRatio float64
	if validRatios > 0 {
		avgRatio = totalRatio / float64(validRatios)
	}

	fmt.Println("+-----------+----------------+-------------------+--------+----------------+----------------+")
	fmt.Printf("| AVERAGE   | %14.2f | %17.2f | %6.1fx |                |                |\n",
		avgNull.Seconds()*1000,
		avgNotNull.Seconds()*1000,
		avgRatio)
	fmt.Println("+-----------+----------------+-------------------+--------+----------------+----------------+")
}

func printAnalysis(results []BenchmarkResult) {
	fmt.Println("\nANALYSIS")
	fmt.Println("========")

	var nullTimes, notNullTimes []time.Duration
	for _, r := range results {
		nullTimes = append(nullTimes, r.NullTime)
		notNullTimes = append(notNullTimes, r.NotNullTime)
	}

	nullStdDev := calculateStdDev(nullTimes)
	notNullStdDev := calculateStdDev(notNullTimes)

	fmt.Printf("NULL column addition - Standard Deviation: %.2f ms\n", nullStdDev.Seconds()*1000)
	fmt.Printf("NOT NULL column addition - Standard Deviation: %.2f ms\n", notNullStdDev.Seconds()*1000)

	fmt.Println("\nKey Insights:")
	fmt.Println("- NOT NULL columns require MySQL to populate default values for all existing rows")
	fmt.Println("- NULL columns can be added instantly as they don't require row-by-row updates")
	fmt.Println("- Performance difference scales with table size")
	fmt.Println("- Consider using NULL columns for optional fields or add NOT NULL columns during low-traffic periods")
}

func calculateStdDev(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}

	var sum time.Duration
	for _, d := range durations {
		sum += d
	}
	mean := sum / time.Duration(len(durations))

	var variance float64
	for _, d := range durations {
		diff := float64(d - mean)
		variance += diff * diff
	}
	variance /= float64(len(durations))

	return time.Duration(variance)
}

func setupTestTable(db *sql.DB) {
	_, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName))
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(fmt.Sprintf(`
		CREATE TABLE %s (
			id INT AUTO_INCREMENT PRIMARY KEY,
			data VARCHAR(255),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB`, tableName))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserting test data...")

	// Use batch inserts for better performance
	batchSize := 1000
	for batch := 0; batch < rowCount; batch += batchSize {
		if batch%10000 == 0 {
			fmt.Printf("\rInserted %d/%d rows", batch, rowCount)
		}

		values := ""
		for i := 0; i < batchSize && batch+i < rowCount; i++ {
			if i > 0 {
				values += ","
			}
			values += fmt.Sprintf("('sample_data_%d')", batch+i)
		}

		query := fmt.Sprintf("INSERT INTO %s (data) VALUES %s", tableName, values)
		_, err = db.Exec(query)
		if err != nil {
			log.Fatal("Batch insert failed:", err)
		}
	}

	fmt.Printf("\rInserted %d/%d rows\n", rowCount, rowCount)
	fmt.Println("Test data inserted successfully")
}

func cleanup(db *sql.DB) {
	_, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName))
	if err != nil {
		log.Println("Warning: cleanup failed:", err)
	}
}
