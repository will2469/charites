package drift

import (
	"math"
	"strings"
)

// TokenStats menyimpan metrik agregasi frekuensi sebuah token dalam cluster.
type TokenStats struct {
	Token      string
	Count      int
	Percentage float64
}

// OutlierItem merepresentasikan token langka (rogue outlier) dalam cluster.
type OutlierItem struct {
	Token       string
	Count       int
	Percentage  float64
	Occurrences []StyleOccurrence
}

// Cluster menyimpan data agregasi statistik kemunculan token untuk satu ClusterKey.
type Cluster struct {
	Key                ClusterKey
	Total              int
	Tokens             map[string]int
	Occurrences        map[string][]StyleOccurrence
	CanonicalToken     string
	CanonicalCount     int
	CanonicalDominance float64
	Outliers           []OutlierItem
	InsufficientData   bool
}

// ContextException mendefinisikan aturan pengecualian berbasis token dan prop opsional.
type ContextException struct {
	Token string
	Props map[string]string // Bila diisi, pengecualian hanya berlaku jika prop cocok
}

// Options mendefinisikan parameter ambang batas keputusan statistik drift.
type Options struct {
	DominanceThreshold    float64                                  // Persentase dominansi kanonikal (default: 80.0)
	OutlierThreshold      float64                                  // Batas persentase outlier rogue (default: 10.0)
	MinClusterOccurrences int                                      // Batas minimum populasi cluster untuk mengambil keputusan (default: 10)
	Exceptions            map[string]map[string][]ContextException // map[Component][Category][]ContextException
}

// DefaultOptions mengembalikan konfigurasi default yang aman dari false positive.
func DefaultOptions() Options {
	return Options{
		DominanceThreshold:    80.0,
		OutlierThreshold:      10.0,
		MinClusterOccurrences: 10,
		Exceptions:            make(map[string]map[string][]ContextException),
	}
}

func normalizeOptions(opts Options) Options {
	if opts.DominanceThreshold <= 0 {
		opts.DominanceThreshold = 80.0
	}
	if opts.OutlierThreshold <= 0 {
		opts.OutlierThreshold = 10.0
	}
	if opts.MinClusterOccurrences <= 0 {
		opts.MinClusterOccurrences = 10
	}
	return opts
}

func groupOccurrences(occurrences []StyleOccurrence) map[ClusterKey]*Cluster {
	clusters := make(map[ClusterKey]*Cluster)
	for _, occ := range occurrences {
		key := ClusterKey{Scope: occ.Scope, Category: occ.Category}
		c, exists := clusters[key]
		if !exists {
			c = &Cluster{
				Key:         key,
				Tokens:      make(map[string]int),
				Occurrences: make(map[string][]StyleOccurrence),
			}
			clusters[key] = c
		}
		c.Total++
		c.Tokens[occ.Token]++
		c.Occurrences[occ.Token] = append(c.Occurrences[occ.Token], occ)
	}
	return clusters
}

func findCanonicalToken(tokens map[string]int) (string, int) {
	maxCount := -1
	canonical := ""
	for tok, count := range tokens {
		if count > maxCount {
			maxCount = count
			canonical = tok
		}
	}
	return canonical, maxCount
}

func detectClusterOutliers(c *Cluster, opts Options) {
	for tok, count := range c.Tokens {
		if tok == c.CanonicalToken {
			continue
		}
		pct := math.Round((float64(count)/float64(c.Total))*1000) / 10
		if pct > opts.OutlierThreshold {
			continue
		}
		filtered := filterExceptions(c.Key, tok, c.Occurrences[tok], opts.Exceptions)
		if len(filtered) > 0 {
			c.Outliers = append(c.Outliers, OutlierItem{
				Token:       tok,
				Count:       len(filtered),
				Percentage:  pct,
				Occurrences: filtered,
			})
		}
	}
}

func evaluateCluster(c *Cluster, opts Options) {
	if c.Total < opts.MinClusterOccurrences {
		c.InsufficientData = true
		return
	}

	canonical, maxCount := findCanonicalToken(c.Tokens)
	c.CanonicalToken = canonical
	c.CanonicalCount = maxCount
	c.CanonicalDominance = math.Round((float64(maxCount)/float64(c.Total))*1000) / 10

	if c.CanonicalDominance >= opts.DominanceThreshold {
		detectClusterOutliers(c, opts)
	}
}

// ClusterOccurrences mengelompokkan kemunculan gaya dan menerapkan gerbang statistik.
func ClusterOccurrences(occurrences []StyleOccurrence, opts Options) map[ClusterKey]*Cluster {
	opts = normalizeOptions(opts)
	clusters := groupOccurrences(occurrences)

	for _, c := range clusters {
		evaluateCluster(c, opts)
	}

	return clusters
}

func matchesContextException(occ StyleOccurrence, ex ContextException, token string) bool {
	if ex.Token != token && ex.Token != "*" {
		return false
	}
	if len(ex.Props) == 0 {
		return true
	}
	for k, v := range ex.Props {
		if occ.Props[k] != v {
			return false
		}
	}
	return true
}

func getCategoryExceptions(key ClusterKey, exceptions map[string]map[string][]ContextException) []ContextException {
	if len(exceptions) == 0 {
		return nil
	}
	compExceptions, compOk := exceptions[key.Scope.Kind.String()]
	if !compOk {
		return nil
	}
	return compExceptions[key.Category.String()]
}

// filterExceptions menyaring kemunculan yang lolos pengecualian kontekstual.
func filterExceptions(key ClusterKey, token string, occurrences []StyleOccurrence, exceptions map[string]map[string][]ContextException) []StyleOccurrence {
	catExceptions := getCategoryExceptions(key, exceptions)
	if len(catExceptions) == 0 {
		return occurrences
	}

	var remaining []StyleOccurrence
	for _, occ := range occurrences {
		matched := false
		for _, ex := range catExceptions {
			if matchesContextException(occ, ex, token) {
				matched = true
				break
			}
		}
		if !matched {
			remaining = append(remaining, occ)
		}
	}

	return remaining
}

// TokenStatsList mengembalikan daftar token terurut berdasarkan frekuensi menurun.
func (c *Cluster) TokenStatsList() []TokenStats {
	list := make([]TokenStats, 0, len(c.Tokens))
	for tok, count := range c.Tokens {
		pct := 0.0
		if c.Total > 0 {
			pct = math.Round((float64(count)/float64(c.Total))*1000) / 10
		}
		list = append(list, TokenStats{
			Token:      tok,
			Count:      count,
			Percentage: pct,
		})
	}

	// Urutkan menurun berdasarkan count
	for i := 0; i < len(list); i++ {
		for j := i + 1; j < len(list); j++ {
			if list[j].Count > list[i].Count || (list[j].Count == list[i].Count && strings.Compare(list[j].Token, list[i].Token) < 0) {
				list[i], list[j] = list[j], list[i]
			}
		}
	}

	return list
}
