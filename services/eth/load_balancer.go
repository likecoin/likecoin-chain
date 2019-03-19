package eth

import (
	"math/rand"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
)

// LoadBalancer stores a list of available Client and returns one when needed
type LoadBalancer struct {
	clients []*ethclient.Client
	weights []int64
	lock    sync.Mutex
}

// NewLoadBalancer creates and initializes a LoadBalancer
func NewLoadBalancer(endpoints []string) *LoadBalancer {
	l := len(endpoints)
	clients := make([]*ethclient.Client, 0, l)
	for _, endpoint := range endpoints {
		ethClient, err := ethclient.Dial(endpoint)
		if err != nil {
			log.
				WithField("eth_endpoint", endpoint).
				WithError(err).
				Panic("Cannot initialize Ethereum endpoint")
		}
		clients = append(clients, ethClient)
	}
	weights := make([]int64, l)
	for i := 0; i < l; i++ {
		weights[i] = 0xFFFFFFFF
	}
	return &LoadBalancer{
		clients: clients,
		weights: weights,
		lock:    sync.Mutex{},
	}
}

// Get returns a Client base on the weightings
func (lb *LoadBalancer) Get() (int, *ethclient.Client) {
	lb.lock.Lock()
	defer lb.lock.Unlock()
	weightSum := int64(0)
	for _, weight := range lb.weights {
		weightSum += weight
	}
	r := rand.Int63n(weightSum)
	for i, weight := range lb.weights {
		if r < weight {
			return i, lb.clients[i]
		}
		r -= weight
	}
	return 0, lb.clients[0]
}

// Do accepts a job which requires a Client and a context, then executes and retries the job with the listed Clients
func (lb *LoadBalancer) Do(f func(*ethclient.Client) error) {
	trialCount := 0
	success := false
	usedClients := map[int]bool{}
	for !success && trialCount < 100 {
		trialCount++
		func() {
			clientIndex, client := lb.Get()
			log.WithField("client_index", clientIndex).Debug("LoadBalancer executing request")
			defer func() {
				usedClients[clientIndex] = true
				err := recover()
				if err != nil {
					log.
						WithField("panic_value", err).
						Warn("LoadBalancer caught panic, recovered")
				}
				if !success {
					if len(usedClients) >= len(lb.clients) {
						log.
							WithField("trial_count", trialCount).
							Panic("LoadBalancer tried all clients but none succeeded")
					} else {
						log.
							WithField("client_index", clientIndex).
							WithField("trial_count", trialCount).
							Warn("LoadBalancer execution failed, retrying")
					}
				}
				lb.lock.Lock()
				defer lb.lock.Unlock()
				weight := lb.weights[clientIndex]
				if success && weight < 0xFFFFFFFF {
					weight = (weight << 1) | 1
				} else if !success && weight > 1 {
					weight = (weight >> 1) | 1
				}
				lb.weights[clientIndex] = weight
				log.
					WithField("client_index", clientIndex).
					WithField("weight", weight).
					Debug("LoadBalancer adjusted client weighting")
			}()
			err := f(client)
			if err == nil {
				success = true
			} else {
				log.
					WithError(err).
					Warn("LoadBalancer failed when executing request")
			}
		}()
	}
	if !success {
		log.WithField("trial_count", trialCount).Panic("LoadBalancer trial count exceeded hard limit")
	}
}
