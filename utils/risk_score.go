package utils

import (
	"math"
	"sort"
	"time"
)

const RiskModelVersion = "v1.3.0"

type RiskWeights struct {
	CoverageGap  float64 `json:"coverage_gap"`
	SinglePoint float64 `json:"single_point"`
}

type RiskScore struct {
	ModelVersion    string             `json:"model_version"`
	CalculatedAt    string            `json:"calculated_at"`
	Score           float64           `json:"score"`
	Level           string            `json:"level"`
	CoverageGapRate float64          `json:"coverage_gap_rate"`
	SinglePointRate float64          `json:"single_point_rate"`
	TopContributors []RiskContribution `json:"top_contributors"`
	Advice          []string          `json:"advice"`
	Weights         RiskWeights       `json:"weights"`
}

type RiskContribution struct {
	Dimension string  `json:"dimension"`
	Rate      float64 `json:"rate"`
	Weight    float64 `json:"weight"`
	Impact    float64 `json:"impact"`
	Reason    string  `json:"reason"`
}

func DefaultRiskWeights() RiskWeights {
	return RiskWeights{
		CoverageGap:  0.8,
		SinglePoint: 0.2,
	}
}

func CalculateRiskScore(totalFunctions int, functionMap map[uint]int, weights RiskWeights) RiskScore {
	calculatedAt := time.Now().UTC().Format(time.RFC3339)

	if totalFunctions <= 0 {
		level := riskLevel(100)
		return RiskScore{
			ModelVersion:    RiskModelVersion,
			CalculatedAt:    calculatedAt,
			Score:           100,
			Level:           level,
			CoverageGapRate: 100,
			SinglePointRate: 0,
			TopContributors: []RiskContribution{},
			Advice:          buildRiskAdvice(level, []RiskContribution{}),
			Weights:         weights,
		}
	}

	missing := 0
	singlePoint := 0

	for _, count := range functionMap {
		if count == 0 {
			missing++
		}
		if count == 1 {
			singlePoint++
		}
	}

	coverageGapRate := float64(missing) / float64(totalFunctions) * 100
	singlePointRate := float64(singlePoint) / float64(totalFunctions) * 100

	score := coverageGapRate*weights.CoverageGap +
		singlePointRate*weights.SinglePoint

	score = clamp(score, 0, 100)

	contributions := []RiskContribution{
		{
			Dimension: "coverage_gap",
			Rate:      round2(coverageGapRate),
			Weight:    round2(weights.CoverageGap),
			Impact:    round2(coverageGapRate * weights.CoverageGap),
			Reason:    "未覆盖功能越多，暴露面越大",
		},
		{
			Dimension: "single_point",
			Rate:      round2(singlePointRate),
			Weight:    round2(weights.SinglePoint),
			Impact:    round2(singlePointRate * weights.SinglePoint),
			Reason:    "单点承载比例越高，抗失效能力越弱",
		},
	}

	sort.Slice(contributions, func(i, j int) bool {
		return contributions[i].Impact > contributions[j].Impact
	})

	level := riskLevel(score)

	return RiskScore{
		ModelVersion:    RiskModelVersion,
		CalculatedAt:    calculatedAt,
		Score:           round2(score),
		Level:           level,
		CoverageGapRate: round2(coverageGapRate),
		SinglePointRate: round2(singlePointRate),
		TopContributors: contributions,
		Advice:          buildRiskAdvice(level, contributions),
		Weights:         weights,
	}
}

func buildRiskAdvice(level string, contributions []RiskContribution) []string {
	advice := []string{levelAdvice(level)}

	if len(contributions) == 0 {
		advice = append(advice, "先补充功能点与设备映射数据，再执行治理决策")
		return advice
	}

	switch contributions[0].Dimension {
	case "coverage_gap":
		advice = append(advice, "优先补齐缺失功能，确保关键检测与防护能力全覆盖")
	case "single_point":
		advice = append(advice, "优先引入双活或异构备份，降低单点承载导致的中断风险")
	default:
		advice = append(advice, "按高影响维度制定专项治理计划并持续复评")
	}

	return advice
}

func levelAdvice(level string) string {
	switch level {
	case "critical":
		return "当前风险为严重级，建议立即执行整改并优先处理高影响风险项"
	case "high":
		return "当前风险为高级，建议在近期版本周期内完成重点治理"
	case "medium":
		return "当前风险为中级，建议按计划推进优化并跟踪趋势"
	default:
		return "当前风险为低级，建议保持基线并定期巡检"
	}
}

func riskLevel(score float64) string {
	switch {
	case score >= 75:
		return "critical"
	case score >= 50:
		return "high"
	case score >= 25:
		return "medium"
	default:
		return "low"
	}
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}
