package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/lib/pq"
)

func GenerateDailyReport(db *sql.DB, reportDate time.Time, reportsDir string) (string, error) {
	startOfDay := time.Date(reportDate.Year(), reportDate.Month(), reportDate.Day(), 0, 0, 0, 0, reportDate.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	interval := 5 * time.Minute

	stats := make(map[int64]int64)

	for t := startOfDay; t.Before(endOfDay); t = t.Add(interval) {
		tEnd := t.Add(interval)

		query := `
			SELECT camera_id, count(*)
			FROM public.kv_events
			WHERE the_date >= $1 AND the_date < $2
			GROUP BY camera_id
		`
		fmt.Println(query)
		rows, err := db.Query(query, t, tEnd)
		if err != nil {
			log.Println("Ошибка запроса:", err)
			continue
		}

		for rows.Next() {
			var cameraID int64
			var count int64
			if err := rows.Scan(&cameraID, &count); err != nil {
				log.Println("Ошибка чтения строки:", err)
				continue
			}
			fmt.Println(cameraID, count)
			stats[cameraID] += count // суммируем
		}
		rows.Close()
	}

	// создаем папку reports (если её нет)
	if err := os.MkdirAll(reportsDir, 0755); err != nil {
		return "", fmt.Errorf("не удалось создать папку %s: %w", reportsDir, err)
	}

	// имя файла = reports/report_YYYY-MM-DD.txt
	fileName := fmt.Sprintf("report_%s.txt", reportDate.Format("2006-01-02"))
	filePath := filepath.Join(reportsDir, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("не удалось создать файл отчета: %w", err)
	}
	defer file.Close()

	fmt.Fprintf(file, "Дневной отчет за %s\n\n", reportDate.Format("2006-01-02"))

	if len(stats) == 0 {
		fmt.Fprintln(file, "Данных нет")
		return filePath, nil
	}

	total := int64(0)
	for cameraID, count := range stats {
		fmt.Fprintf(file, "CameraID: %d, Count: %d\n", cameraID, count)
		total += count
	}

	fmt.Fprintf(file, "\nВсего событий за день: %d\n", total)
	return filePath, nil
}
