package drift

import (
	"fmt"
	"strings"
)

// ChromaticFamily mengidentifikasi kelompok spektrum warna primitif.
type ChromaticFamily int

// Daftar kelompok spektrum warna ChromaticFamily.
const (
	// FamilyUnknown menandai keluarga warna yang tidak teridentifikasi.
	FamilyUnknown      ChromaticFamily = iota
	FamilyWarmAmber                    // yellow-*, amber-*, orange-*
	FamilyGreenEmerald                 // green-*, emerald-*, teal-*
	FamilyBlueSky                      // blue-*, sky-*, cyan-*
	FamilyRedRose                      // red-*, rose-*
)

func (f ChromaticFamily) String() string {
	switch f {
	case FamilyWarmAmber:
		return "WarmAmber"
	case FamilyGreenEmerald:
		return "GreenEmerald"
	case FamilyBlueSky:
		return "BlueSky"
	case FamilyRedRose:
		return "RedRose"
	default:
		return "Unknown"
	}
}

// CandidateIntent mendefinisikan dugaan intensi semantik berdasarkan heuristik warna primitif.
type CandidateIntent int

// Daftar dugaan intensi semantik CandidateIntent.
const (
	// CandidateNone menandai ketiadaan dugaan intensi semantik.
	CandidateNone CandidateIntent = iota
	CandidateWarning
	CandidateSuccess
	CandidateInfo
	CandidateDestructive
)

func (ci CandidateIntent) String() string {
	switch ci {
	case CandidateWarning:
		return "Warning"
	case CandidateSuccess:
		return "Success"
	case CandidateInfo:
		return "Info"
	case CandidateDestructive:
		return "Destructive"
	default:
		return "None"
	}
}

// TargetTokens mengembalikan pasangan CSS token semantik target di global.css.
func (ci CandidateIntent) TargetTokens() (string, string) {
	switch ci {
	case CandidateWarning:
		return "--warning", "--warning-foreground"
	case CandidateSuccess:
		return "--success", "--success-foreground"
	case CandidateInfo:
		return "--info", "--info-foreground"
	case CandidateDestructive:
		return "--destructive", "--destructive-foreground"
	default:
		return "", ""
	}
}

// MapChromaticFamily memetakan token background primitif ke ChromaticFamily dan CandidateIntent.
func MapChromaticFamily(bgToken string) (ChromaticFamily, CandidateIntent, bool) {
	cleaned := strings.TrimPrefix(bgToken, "bg-")
	if idx := strings.IndexByte(cleaned, '/'); idx != -1 {
		cleaned = cleaned[:idx]
	}

	parts := strings.Split(cleaned, "-")
	if len(parts) < 2 {
		return FamilyUnknown, CandidateNone, false
	}

	colorName := parts[0]

	switch colorName {
	case "yellow", "amber", "orange":
		return FamilyWarmAmber, CandidateWarning, true
	case "green", "emerald", "teal":
		return FamilyGreenEmerald, CandidateSuccess, true
	case "blue", "sky", "cyan":
		return FamilyBlueSky, CandidateInfo, true
	case "red", "rose":
		return FamilyRedRose, CandidateDestructive, true
	default:
		return FamilyUnknown, CandidateNone, false
	}
}

// ChromaticCluster menyimpan temuan drift warna dalam satu komponen dan rumpun warna.
type ChromaticCluster struct {
	Scope           ScopeKey
	Family          ChromaticFamily
	CandidateIntent CandidateIntent
	Total           int
	HueCounts       map[string]int
	HueOccurrences  map[string][]StyleOccurrence
	Recommendation  string
}

// AnalyzeChromaticDrift menganalisis seluruh kemunculan gaya untuk mendeteksi persaingan rona warna dalam satu rumpun.
func AnalyzeChromaticDrift(occurrences []StyleOccurrence) []ChromaticCluster {
	type clusterMapKey struct {
		Scope  ScopeKey
		Family ChromaticFamily
	}

	grouped := make(map[clusterMapKey]*ChromaticCluster)

	for _, occ := range occurrences {
		if occ.Category != CatBackground {
			continue
		}

		fam, intent, ok := MapChromaticFamily(occ.Token)
		if !ok {
			continue
		}

		key := clusterMapKey{Scope: occ.Scope, Family: fam}
		c, exists := grouped[key]
		if !exists {
			c = &ChromaticCluster{
				Scope:           occ.Scope,
				Family:          fam,
				CandidateIntent: intent,
				HueCounts:       make(map[string]int),
				HueOccurrences:  make(map[string][]StyleOccurrence),
			}
			grouped[key] = c
		}

		c.Total++
		c.HueCounts[occ.Token]++
		c.HueOccurrences[occ.Token] = append(c.HueOccurrences[occ.Token], occ)
	}

	var results []ChromaticCluster

	for _, c := range grouped {
		// Drift warna terdeteksi jika terdapat 2 atau lebih rona berbeda dalam satu rumpun
		if len(c.HueCounts) >= 2 {
			bgToken, fgToken := c.CandidateIntent.TargetTokens()
			c.Recommendation = fmt.Sprintf(
				"Design system lacks a first-class '%s' variant for <%s>. "+
					"1. Define '%s' and '%s' in 'global.css'. "+
					"2. Add variant '%s' to component definition.",
				strings.ToLower(c.CandidateIntent.String()),
				c.Scope.Kind.String(),
				bgToken,
				fgToken,
				strings.ToLower(c.CandidateIntent.String()),
			)
			results = append(results, *c)
		}
	}

	return results
}
