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
	limitutils "github.com/XPaul6/monitora/utils/limits"

	"gorm.io/gorm"
)

const CONCURRENCY_LIMIT = 10
const FETCH_INTERVAL_SECONDS = 30

func RunFetchUtil(db *gorm.DB) {
	for {
		fetch(db)
		time.Sleep(FETCH_INTERVAL_SECONDS * time.Second)
	}
}

func fetch(db *gorm.DB) {
	// Get all servers
	var servers []Server
	result := db.Find(&servers)
	if result.Error != nil {
		log.Println("Cannot fetch server information")
		return
	}

	// Start gathering metrics concurrently
	sem := make(chan byte, CONCURRENCY_LIMIT)
	results := make(chan IdInfoPair, len(servers))
	var wg sync.WaitGroup
	for _, server := range servers {
		sem <- 1
		wg.Add(1)
		go func(s Server) {
			defer func() {
				<-sem
				wg.Done()
			}()
			info, err := fetchSysInfo(s.IP)
			if err != nil {
				log.Printf("Failed to fetch info from %v: %v", s.IP, err.Error())
				s.Status = "Offline"
				db.Save(&s)
				return
			}
			s.Status = "Online"
			db.Save(&s)
			results <- IdInfoPair{Id: s.ID, Info: info}
		}(server)
	}

	wg.Wait()
	close(results)

	// Work with recieved metrics
	cachedMetrics := getMetricsData(db)
	for log := range results {
		writeLogToDB(log, db)
		limitutils.DetectLimitCross(log, cachedMetrics, db)
	}
}

func fetchSysInfo(ip string) (SysInfo, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	res, err := client.Get(fmt.Sprintf("http://%v/sysinfo", ip))
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

func getMetricsData(db *gorm.DB) map[uint]MetricType {
	var metricTypes []MetricType
	db.Find(&metricTypes)
	metricTypeMap := make(map[uint]MetricType)
	for _, mt := range metricTypes {
		metricTypeMap[mt.ID] = mt
	}
	return metricTypeMap
}

func writeLogToDB(log IdInfoPair, db *gorm.DB) {
	// Components synchronization
	var cpu Component
	result := db.Where("server_id = ? and type = ?", log.Id, "cpu").First(&cpu)
	if result.Error == gorm.ErrRecordNotFound {
		db.FirstOrCreate(&cpu, &Component{
			ServerID: log.Id,
			Type:     "cpu",
			Address:  "cpu",
		})
	}

	var mem Component
	result = db.Where("server_id = ? and type = ?", log.Id, "mem").First(&mem)
	if result.Error == gorm.ErrRecordNotFound {
		db.FirstOrCreate(&mem, &Component{
			ServerID: log.Id,
			Type:     "mem",
			Address:  "mem",
		})
	}

	var disks []Component
	result = db.Where("server_id = ? and type = ?", log.Id, "disk").Find(&disks)
	existingDisks := make(map[string]bool)
	for _, disk := range disks {
		existingDisks[disk.Address] = true
	}
	if (len(disks)) != len(log.Info.Disks) {
		var newDisks []Component
		for _, v := range log.Info.Disks {
			if !existingDisks[v.MountPoint] {
				newDisks = append(newDisks, Component{
					ServerID: log.Id,
					Type:     "disk",
					Address:  v.MountPoint,
				})
			}
		}
		db.Create(&newDisks)
		db.Where("server_id = ? and type = ?", log.Id, "disk").Find(&disks)
	}

	var nets []Component
	result = db.Where("server_id = ? and type = ?", log.Id, "net").Find(&nets)
	existingNets := make(map[string]bool)
	for _, net := range nets {
		existingNets[net.Address] = true
	}
	if (len(nets)) != len(log.Info.Net) {
		var newNets []Component
		for _, v := range log.Info.Net {
			if !existingNets[v.Name] {
				newNets = append(newNets, Component{
					ServerID: log.Id,
					Type:     "net",
					Address:  v.Name,
				})
			}
		}
		db.Create(&newNets)
		db.Where("server_id = ? and type = ?", log.Id, "net").Find(&nets)
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
		Value:        log.Info.CPU.LoadPercentage,
		Timestamp:    time.Now(),
	})

	logs = append(logs, RawLog{
		ComponentID:  cpu.ID,
		MetricTypeID: typeMap["cpu_temperature"].ID,
		Value:        log.Info.CPU.Temperature,
		Timestamp:    time.Now(),
	})

	logs = append(logs, RawLog{
		ComponentID:  mem.ID,
		MetricTypeID: typeMap["mem_load"].ID,
		Value:        log.Info.Mem.LoadPercentage,
		Timestamp:    time.Now(),
	})

	logs = append(logs, RawLog{
		ComponentID:  mem.ID,
		MetricTypeID: typeMap["mem_total"].ID,
		Value:        float64(log.Info.Mem.Total),
		Timestamp:    time.Now(),
	})

	logs = append(logs, RawLog{
		ComponentID:  mem.ID,
		MetricTypeID: typeMap["mem_used"].ID,
		Value:        float64(log.Info.Mem.Used),
		Timestamp:    time.Now(),
	})

	diskInfoMap := make(map[string]DiskInfo)
	for _, d := range log.Info.Disks {
		diskInfoMap[d.MountPoint] = d
	}
	for _, disk := range disks {
		logs = append(logs, RawLog{
			ComponentID:  disk.ID,
			MetricTypeID: typeMap["disk_total"].ID,
			Value:        float64(diskInfoMap[disk.Address].Total),
			Timestamp:    time.Now(),
		})
		logs = append(logs, RawLog{
			ComponentID:  disk.ID,
			MetricTypeID: typeMap["disk_used"].ID,
			Value:        float64(diskInfoMap[disk.Address].Used),
			Timestamp:    time.Now(),
		})
	}

	netInfoMap := make(map[string]NetInfo)
	for _, n := range log.Info.Net {
		netInfoMap[n.Name] = n
	}
	for _, net := range nets {
		logs = append(logs, RawLog{
			ComponentID:  net.ID,
			MetricTypeID: typeMap["net_rbps"].ID,
			Value:        float64(netInfoMap[net.Address].RBpS),
			Timestamp:    time.Now(),
		})
		logs = append(logs, RawLog{
			ComponentID:  net.ID,
			MetricTypeID: typeMap["net_sbps"].ID,
			Value:        float64(netInfoMap[net.Address].SBpS),
			Timestamp:    time.Now(),
		})
	}

	result = db.Create(&logs)
	if result.Error != nil {
		fmt.Printf("Failed to write logs: %v", result.Error)
	}
}
