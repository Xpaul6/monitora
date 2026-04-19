package fetchutil

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	. "github.com/XPaul6/monitora/models"

	"gorm.io/gorm"
)

type pair struct {
	id   uint
	info SysInfo
}

const CONCURRENCY_LIMIT = 10
const FETCH_INTERVAL_SECONDS = 30

func RunFetchUtil(db *gorm.DB) {
	for {
		fetch(db)
		time.Sleep(FETCH_INTERVAL_SECONDS * time.Second)
	}
}

func fetch(db *gorm.DB) {
	var servers []Server
	result := db.Find(&servers)
	if result.Error != nil {
		log.Println("Cannot fetch server information")
		return
	}

	sem := make(chan byte, CONCURRENCY_LIMIT)
	results := make(chan pair, len(servers))
	var wg sync.WaitGroup
	for _, server := range servers {
		sem <- 1
		wg.Add(1)
		go func(s Server) {
			defer func() {
				<-sem
				wg.Done()
			}()
			info, err := FetchSysInfo(server.IP)
			if err != nil {
				log.Printf("Failed to fetch info from %v", server.IP)
				return
			}
			results <- pair{server.ID, info}
		}(server)
	}

	wg.Wait()
	close(results)

	for log := range results {
		writeLogToDB(log, db)
	}
}

func FetchSysInfo(ip string) (SysInfo, error) {
	res, err := http.Get(fmt.Sprintf("http://%v/sysinfo", ip))
	if err != nil {
		return SysInfo{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return SysInfo{}, err
	}

	var info SysInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return SysInfo{}, err
	}

	return info, nil
}

func writeLogToDB(log pair, db *gorm.DB) {
	// Components synchronization
	var cpu Component
	result := db.Where("server_id = ? and type = ?", log.id, "cpu").First(&cpu)
	if result.Error == gorm.ErrRecordNotFound {
		db.FirstOrCreate(&cpu, &Component{
			ServerID: log.id,
			Type:     "cpu",
			Address:  "cpu",
		})
	}

	var mem Component
	result = db.Where("server_id = ? and type = ?", log.id, "mem").First(&mem)
	if result.Error == gorm.ErrRecordNotFound {
		db.FirstOrCreate(&mem, &Component{
			ServerID: log.id,
			Type:     "mem",
			Address:  "mem",
		})
	}

	var disks []Component
	result = db.Where("server_id = ? and type = ?", log.id, "disk").Find(&disks)
	existingDisks := make(map[string]bool)
	for _, disk := range disks {
		existingDisks[disk.Address] = true
	}
	if (len(disks)) != len(log.info.Disks) {
		var newDisks []Component
		for _, v := range log.info.Disks {
			if !existingDisks[v.MountPoint] {
				newDisks = append(newDisks, Component{
					ServerID: log.id,
					Type:     "disk",
					Address:  v.MountPoint,
				})
			}
		}
		db.Create(&newDisks)
		db.Find(&disks)
	}

	var nets []Component
	result = db.Where("server_id = ? and type = ?", log.id, "net").Find(&nets)
	existingNets := make(map[string]bool)
	for _, net := range nets {
		existingNets[net.Address] = true
	}
	if (len(nets)) != len(log.info.Net) {
		var newNets []Component
		for _, v := range log.info.Net {
			if !existingNets[v.Name] {
				newNets = append(newNets, Component{
					ServerID: log.id,
					Type:     "net",
					Address:  v.Name,
				})
			}
		}
		db.Create(&newNets)
		db.Find(&nets)
	}

	// Mertic types caching
	var metricTypes []MetricType
	db.Find(&metricTypes)
	typeMap := make(map[string]MetricType)
	for _, mt := range metricTypes {
		typeMap[mt.Name] = mt
	}

	// Logging
	var logs []RawLog

	logs = append(logs, RawLog{
		ComponentID:  cpu.ID,
		MetricTypeID: typeMap["cpu_load"].ID,
		Value:        log.info.CPU.LoadPercentage,
		Timestamp:    time.Now(),
	})

	logs = append(logs, RawLog{
		ComponentID:  cpu.ID,
		MetricTypeID: typeMap["cpu_temperature"].ID,
		Value:        log.info.CPU.Temperature,
		Timestamp:    time.Now(),
	})

	logs = append(logs, RawLog{
		ComponentID:  mem.ID,
		MetricTypeID: typeMap["mem_load"].ID,
		Value:        log.info.Mem.LoadPercentage,
		Timestamp:    time.Now(),
	})

	logs = append(logs, RawLog{
		ComponentID:  mem.ID,
		MetricTypeID: typeMap["mem_total"].ID,
		Value:        float64(log.info.Mem.Total),
		Timestamp:    time.Now(),
	})

	logs = append(logs, RawLog{
		ComponentID:  mem.ID,
		MetricTypeID: typeMap["mem_used"].ID,
		Value:        float64(log.info.Mem.Used),
		Timestamp:    time.Now(),
	})

	for i, disk := range disks {
		logs = append(logs, RawLog{
			ComponentID:  disk.ID,
			MetricTypeID: typeMap["disk_total"].ID,
			Value:        float64(log.info.Disks[i].Total),
			Timestamp:    time.Now(),
		})
		logs = append(logs, RawLog{
			ComponentID:  disk.ID,
			MetricTypeID: typeMap["disk_used"].ID,
			Value:        float64(log.info.Disks[i].Used),
			Timestamp:    time.Now(),
		})
	}

	for i, net := range nets {
		logs = append(logs, RawLog{
			ComponentID:  net.ID,
			MetricTypeID: typeMap["net_rbps"].ID,
			Value:        float64(log.info.Net[i].RBpS),
			Timestamp:    time.Now(),
		})
		logs = append(logs, RawLog{
			ComponentID:  net.ID,
			MetricTypeID: typeMap["net_sbps"].ID,
			Value:        float64(log.info.Net[i].SBpS),
			Timestamp:    time.Now(),
		})
	}

	result = db.Create(&logs)
	if result.Error != nil {
		fmt.Printf("Failed to write logs: %v", result.Error)
	}
}
