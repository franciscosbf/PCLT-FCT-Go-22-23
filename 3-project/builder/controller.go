package builder

import (
	"cpl_go_proj22/parser"
	"cpl_go_proj22/utils"
	"log"
	"runtime"
	"sync"
	"time"
)

var cpus = runtime.NumCPU()

type MsgType = int

const (
	BuildSuccess MsgType = iota
	BuildError
)

type Msg struct {
	Type MsgType
	Err error
}

type fileInfo struct {
	// Set while building the graph
	filename string
	dependencies int
	dependants []string
	nodes map[string]*fileInfo

	// Set when spawning workers
	utils.Scan	
	timesCh chan time.Time
	panicCh chan struct{} // When some error happens
	errorCh chan *Msg // Communicate with the error controller
}

func (f *fileInfo) propagate(t time.Time) {
	for _, dep := range f.dependants {
		f.nodes[dep].timesCh <-t	
	}
	log.Printf(
		"%q propagated build time %q to %v", 
		f.filename, t, f.dependants,
	)
}

// build tries to build the file 
// and sends the build time to its
// dependants.
func (f *fileInfo) build() {
	t, err := f.Build(f.filename)
	if err != nil {
		log.Printf(
			"Error while trying to build %q: %v", 
			f.filename, err,
		)
		f.errorCh <-&Msg{Type: BuildError, Err: err}
		return
	}
	f.propagate(t)
}

type depGraph struct {
	leafs map[string]*fileInfo // Build starter workers
	nodes map[string]*fileInfo // Searching and testing purposes
	targets []*fileInfo // Build workers
}

// buildGraph returns  a dependency 
// grapth based on a given set of rules. 
// Targets with nil values represent leafs.
func buildGraph(file *parser.DepFile) *depGraph {
	// Keep track of files that had 
	// been added to the graph
	dG := &depGraph{
		leafs: make(map[string]*fileInfo),
		nodes: make(map[string]*fileInfo),
	}

	insertNode := func(filename string) *fileInfo {
		info := &fileInfo{
			filename: filename,
			dependants: make([]string, 0),
			nodes: dG.nodes,
		}
		dG.nodes[filename] = info
		return info
	}

	insertTarget := func(info *fileInfo) {
		dG.targets = append(dG.targets, info)
	}

	for _, rule := range file.Rules {
		target := rule.Object

		for _, dep := range rule.Deps {
			info, ok := dG.nodes[dep]
			if !ok {
				info = insertNode(dep)
				dG.leafs[dep] = info 
			}
			info.dependants = append(info.dependants, target)
		}

		delete(dG.leafs, target) // It means that's no more a leaf
		info, ok := dG.nodes[target]
		if !ok {
			info = insertNode(target)
		}
		info.dependencies = len(rule.Deps)
		insertTarget(info)
	}

	return dG
}

func targetWorker(info *fileInfo, wg *sync.WaitGroup) {
	defer wg.Done()

	deps := info.dependencies
	
	sTime, err := info.Status(info.filename)
	if err != nil {
		log.Printf(
			"%q doesn't exist. Proceeds to build after wait", 
			info.filename,
		)
		// Only needs to wait for its dependencies
		for ; deps > 0; deps-- {
			select {
			case <-info.panicCh:
				return
			case <-info.timesCh:
			}
		}	
		info.build()
		return
	}

	// Waits until some of its dependencies
	// has an update time greater than the target
	for ; deps > 0; deps-- {
		select {
		case <-info.panicCh:
			return
		case t := <-info.timesCh:
			if sTime.After(t) {
				// Target is more recent
				// than a given dep
				continue
			}	
			log.Printf(
				"%q needs to be built. Proceeds to wait",
				info.filename,
			)
			// Doesn't build right after since we 
			// need to wait for the remaining deps
			for deps--; deps > 0; deps-- {
				select {
				case <-info.panicCh:
					return
				case <-info.timesCh:
				}
			}
			info.build()
			return
		}
	}
	
	// There isn't any dep whose uptime
	// is greater than the target
	info.propagate(sTime)
}

func leafWorker(info *fileInfo, wg *sync.WaitGroup) {
	defer wg.Done()

	select {
	case <-info.panicCh:
		return // Something went wrong
	default:
	}
	
	if t, err := info.Status(info.filename); err == nil {
		info.propagate(t)
		return
	}
	log.Printf("%q doesn't exist. Proceeds to build", info.filename)
	info.build()
}

func spawnTargetWorkers( 
	targets []*fileInfo,	
	fileScan utils.Scan, 
	errorCh chan *Msg,
	workersWg *sync.WaitGroup,
) {
	initCommonChs := func(info *fileInfo) {
		info.timesCh = make(chan time.Time, info.dependencies)
		info.Scan = fileScan
		info.panicCh = make(chan struct{}, 1)
		info.errorCh = errorCh
	}

	sz := len(targets)
	log.Printf("Spawning %d target workers with %d sub-spawners", sz, cpus)
	var wg sync.WaitGroup
	wg.Add(cpus)
	for c := cpus - 1; c >= 0; c-- {
		go func(i, j int) {
			defer wg.Done()
			for ; i < j; i++ {
				info := targets[i]
				initCommonChs(info)
				go targetWorker(info, workersWg)
			}
		}(c * sz / cpus, (c + 1) * sz / cpus)
	}
	wg.Wait()
}

func spawnLeafWorkers( 
	leafs map[string]*fileInfo,	
	fileScan utils.Scan, 
	errorCh chan *Msg,
	workersWg *sync.WaitGroup,
) {
	initCommonChs := func(info *fileInfo) {
		info.Scan = fileScan
		info.panicCh = make(chan struct{}, 1)
		info.errorCh = errorCh
	}

	sz := len(leafs)
	log.Printf("Spawning %d leaf workers with %d sub-spawners", sz, cpus)
	var wg sync.WaitGroup
	wg.Add(cpus)
	fileInfoCh := make(chan *fileInfo, sz)
	go func() {
		for _, info := range leafs {
			fileInfoCh <-info
		}
	}()
	for c := cpus - 1; c >= 0; c-- {
		go func(i, j int) {
			defer wg.Done()
			// We need to keep using i and 
			// j since the sz can be odd
			for n := j - i; n > 0; n-- {
				info := <-fileInfoCh
				initCommonChs(info)
				go leafWorker(info, workersWg)
			}
		}(c * sz / cpus, (c + 1) * sz / cpus)
	}
	wg.Wait()
}

func MakeController(file *parser.DepFile, fileScan utils.Scan) chan *Msg {
	dG := buildGraph(file)

	workersN := len(dG.targets) + len(dG.leafs)

	reqCh := make(chan *Msg, 1)
	errorCh := make(chan *Msg, workersN)

	var workersWg sync.WaitGroup
	workersWg.Add(workersN)

	spawnTargetWorkers(
		dG.targets, fileScan, 
		errorCh, &workersWg, 
	)

	spawnLeafWorkers(
		dG.leafs, fileScan, 
		errorCh, &workersWg, 
	)

	errMsgCh := make(chan *Msg, 1)

	// Core manager 
	go func() {
		err := <-errorCh
		// Sends error to reconciler
		errMsgCh <-err
	
		log.Printf("Core manager has received an error: %v", err.Err)

		// Tells every worker to end its execution.
		// There are workers that may have finished
		empty := struct{}{}
		go func() {
			for _, info := range dG.targets { // INFO: speed up this thing...
				info.panicCh <- empty
			}
		}()
		for _, info := range dG.leafs { // INFO: speed up this thing...
			info.panicCh <- empty
		}
		
	}()

	// Reconciler
	go func() {
		workersWg.Wait()
		var msg *Msg
		select {
		case msg = <-errMsgCh:
		default:
			// Everything went ok
			msg = &Msg{Type: BuildSuccess}
		}
		reqCh <- msg
	}()

	return reqCh
}

