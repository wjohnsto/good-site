package services

import (
	"math/rand"
	"strings"
	"time"
)

type Markov struct {
	nextToken string
	chain     map[string][]string
	rand      *rand.Rand
}

func (m *Markov) RandomToken() string {
	keys := make([]string, 0, len(m.chain))

	for key := range m.chain {
		keys = append(keys, key)
	}

	idx := m.rand.Intn(len(keys))

	return keys[idx]
}

func (m *Markov) Next() string {
	for {
		if m.nextToken == "" {
			m.nextToken = m.RandomToken()
		}

		links, ok := m.chain[m.nextToken]

		if !ok || len(links) == 0 {
			m.nextToken = ""
			continue
		}

		token := strings.Clone(m.nextToken)
		m.nextToken = strings.Clone(links[rand.Intn(len(links))])

		return token
	}
}

func NewMarkov(tokens []string) *Markov {
	chain := map[string][]string{}
	var last string

	for idx, token := range tokens {
		if idx == 0 {
			last = token
			continue
		}
		if val, ok := chain[last]; ok {
			chain[last] = append(val, strings.Clone(token))
		} else {
			chain[last] = []string{strings.Clone(token)}
		}

		last = token
	}

	return &Markov{
		nextToken: "",
		chain:     chain,
		rand:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}
