package queries

import (
	"io/ioutil"
	"os"

	"github.com/prometheus/log"
	"gopkg.in/yaml.v2"

	"github.com/mc2soft/postgresql_exporter/metrics"
)

func ParseQueries(queriesPath string) (cq []metrics.CustomQuery) {
	if queriesPath == "" {
		return
	}

	f, err := os.Open(queriesPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(b, &cq)
	if err != nil {
		log.Fatal(err)
	}

	return
}
