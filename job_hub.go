package thorn

import (
	"fmt"
	"github.com/juju/errors"
	"sync"
)

// SearchParam is search param
type SearchParam struct {
	ClientID string
	Name string
	VirtualPort int
	Port int
}

type JobHub interface {
	// Run run a job and keep job in hub
	Run(job *Job) error

	// Get get job by id
	Get(id int64) *Job

	// Search get job list by search param
	Search(param *SearchParam) []*Job
}

func NewJobHub(server *Server) JobHub {
	return &jobHub{
		server: server,
		lock: sync.Mutex{},
		id: 0,
		jobRunner: NewJobRunner(server),
		jobs: make(map[int64]*Job, 128),
		jobsVirtualPortIdx: make(map[int]*Job, 128),
		jobsClientIDPortIdx: make(map[string]*Job, 128),
	}
}

type jobHub struct {
	server *Server
	lock sync.Mutex
	id int64
	jobRunner JobRunner
	jobs map[int64]*Job
	jobsVirtualPortIdx map[int]*Job
	jobsClientIDPortIdx map[string]*Job
}

func (j *jobHub) nextID() int64 {
	j.lock.Lock()
	defer j.lock.Unlock()

	j.id++
	return j.id
}

func (j *jobHub) Run(job *Job) error {
	jobFind := j.GetByVirtualPort(job.VirtualPort)
	if jobFind != nil {
		return errors.New(fmt.Sprintf("virtualPort exists, job.id = %d", jobFind.ID))
	}

	jobFind = j.GetByClientIDAndPort(job.ClientID, job.Port)
	if jobFind != nil {
		return errors.New(fmt.Sprintf("clientID and port exists, job.id = %d", jobFind.ID))
	}

	job.ID = j.nextID()

	j.jobs[job.ID] = job
	j.jobsVirtualPortIdx[job.VirtualPort] = job
	key := j.clientIDAndPortKey(job.ClientID, job.Port)
	j.jobsClientIDPortIdx[key] = job
	return j.jobRunner.Run(job)
}

func (j *jobHub) Get(id int64) *Job {
	if job, ok := j.jobs[id]; ok {
		return job
	}

	return nil
}

func (j *jobHub) GetByVirtualPort(virtualPort int) *Job {
	if job, ok := j.jobsVirtualPortIdx[virtualPort]; ok {
		return job
	}

	return nil
}

func (j *jobHub) GetByClientIDAndPort(clientID string, port int) *Job {
	key := j.clientIDAndPortKey(clientID, port)
	if job, ok := j.jobsClientIDPortIdx[key]; ok {
		return job
	}

	return nil
}

func (j *jobHub) Search(param *SearchParam) []*Job {
	var jobs []*Job
	for _, job := range j.jobs {
		if j.isMatch(job, param) {
			jobs = append(jobs, job)
		}
	}

	return jobs
}

func (j jobHub) clientIDAndPortKey(clientID string, port int) string {
	return fmt.Sprintf("%s-%d", clientID, port)
}

func (j jobHub) isMatch(job *Job, param *SearchParam) bool {
	return true
}