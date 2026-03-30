<template>
  <div class="analysis-page">
    <header class="page-header">
      <div class="header-left">
        <h1 class="page-title">能力分析</h1>
        <p class="page-desc">评估安全能力覆盖与冗余情况</p>
      </div>
      <div class="header-actions">
        <select v-model="selectedTopoId" class="topo-select" @change="loadAnalysis">
          <option value="">选择拓扑</option>
          <option v-for="t in topologies" :key="t.id" :value="t.id">
            {{ t.name }}
          </option>
        </select>
        <button class="btn btn-primary" @click="loadAnalysis" :disabled="!selectedTopoId">
          开始分析
        </button>
      </div>
    </header>

    <div class="page-content" v-if="result">
      <div class="coverage-section">
        <div class="coverage-header">
          <h2>能力覆盖率</h2>
          <span class="coverage-value">{{ result.coverage_rate.toFixed(1) }}%</span>
        </div>
        <div class="progress-bar large">
          <div class="progress-fill" :style="{ width: result.coverage_rate + '%' }"></div>
        </div>
        <div class="coverage-legend">
          <div class="legend-item">
            <span class="dot covered"></span>
            <span>已覆盖 {{ coveredCount }} 项</span>
          </div>
          <div class="legend-item">
            <span class="dot missing"></span>
            <span>缺失 {{ result.missing?.length || 0 }} 项</span>
          </div>
          <div class="legend-item">
            <span class="dot redundant"></span>
            <span>冗余 {{ result.redundant?.length || 0 }} 项</span>
          </div>
        </div>
      </div>

      <div class="analysis-grid">
        <div class="card">
          <div class="card-header">
            <h3 class="card-title">缺失能力</h3>
            <span class="badge badge-danger">{{ result.missing?.length || 0 }}</span>
          </div>
          <div class="function-list" v-if="result.missing?.length">
            <div v-for="func in result.missing" :key="func.id" class="function-item missing">
              <div class="func-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="10"/>
                  <line x1="15" y1="9" x2="9" y2="15"/>
                  <line x1="9" y1="9" x2="15" y2="15"/>
                </svg>
              </div>
              <div class="func-info">
                <div class="func-name">{{ func.name }}</div>
                <div class="func-desc">{{ func.description }}</div>
              </div>
              <div class="func-category">{{ func.category }}</div>
            </div>
          </div>
          <div v-else class="empty-state small">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
              <polyline points="22 4 12 14.01 9 11.01"/>
            </svg>
            <p>所有能力已覆盖</p>
          </div>
        </div>

        <div class="card">
          <div class="card-header">
            <h3 class="card-title">冗余能力</h3>
            <span class="badge badge-warning">{{ result.redundant?.length || 0 }}</span>
          </div>
          <div class="function-list" v-if="result.redundant?.length">
            <div v-for="func in result.redundant" :key="func.id" class="function-item redundant">
              <div class="func-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/>
                  <line x1="12" y1="9" x2="12" y2="13"/>
                  <line x1="12" y1="17" x2="12.01" y2="17"/>
                </svg>
              </div>
              <div class="func-info">
                <div class="func-name">{{ func.name }}</div>
                <div class="func-desc">{{ func.description }}</div>
              </div>
              <div class="func-category">{{ func.category }}</div>
            </div>
          </div>
          <div v-else class="empty-state small">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
              <polyline points="22 4 12 14.01 9 11.01"/>
            </svg>
            <p>无冗余能力</p>
          </div>
        </div>
      </div>

      <div class="card risk-card" v-if="result.risk">
        <div class="card-header">
          <h3 class="card-title">风险评估</h3>
          <span class="risk-badge" :class="getRiskClass(result.risk.level)">
            {{ result.risk.level }}
          </span>
        </div>
        <div class="risk-content">
          <div class="risk-score-display">
            <div class="score-circle" :class="getRiskClass(result.risk.level)">
              <span class="score-value">{{ result.risk.score }}</span>
              <span class="score-max">/100</span>
            </div>
          </div>
          <div class="risk-details">
            <div class="risk-factors">
              <div class="factor">
                <span class="factor-label">覆盖率缺口</span>
                <div class="factor-bar">
                  <div class="progress-bar">
                    <div class="progress-fill warning" :style="{ width: result.risk.coverage_gap_rate + '%' }"></div>
                  </div>
                  <span class="factor-value">{{ result.risk.coverage_gap_rate.toFixed(1) }}%</span>
                </div>
              </div>
              <div class="factor">
                <span class="factor-label">冗余率</span>
                <div class="factor-bar">
                  <div class="progress-bar">
                    <div class="progress-fill danger" :style="{ width: result.redundancy_rate + '%' }"></div>
                  </div>
                  <span class="factor-value">{{ result.redundancy_rate.toFixed(1) }}%</span>
                </div>
              </div>
              <div class="factor">
                <span class="factor-label">单点承载</span>
                <div class="factor-bar">
                  <div class="progress-bar">
                    <div class="progress-fill" :style="{ width: result.risk.single_point_rate + '%' }"></div>
                  </div>
                  <span class="factor-value">{{ result.risk.single_point_rate.toFixed(1) }}%</span>
                </div>
              </div>
            </div>
            <div class="risk-advice" v-if="result.risk.advice?.length">
              <h4>改进建议</h4>
              <ul>
                <li v-for="(advice, i) in result.risk.advice" :key="i">{{ advice }}</li>
              </ul>
            </div>
          </div>
        </div>
      </div>

      <div class="action-section">
        <router-link :to="`/suggest/${selectedTopoId}`" class="btn btn-primary" v-if="selectedTopoId">
          获取优化方案 →
        </router-link>
      </div>
    </div>

    <div class="empty-state" v-else-if="!loading && !selectedTopoId">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
        <path d="M3 3v18h18"/>
        <path d="M7 16l4-4 4 4 6-6"/>
      </svg>
      <p>选择一个拓扑开始分析安全能力</p>
    </div>

    <div class="loading" v-if="loading">
      <div class="spinner"></div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { topoApi, analyzeApi } from '../api'

const loading = ref(false)
const topologies = ref([])
const selectedTopoId = ref('')
const result = ref(null)

const coveredCount = computed(() => {
  if (!result.value) return 0
  const total = (result.value.missing?.length || 0) + (result.value.redundant?.length || 0)
  return total - (result.value.missing?.length || 0)
})

function getRiskClass(level) {
  const map = {
    critical: 'risk-critical',
    high: 'risk-high',
    medium: 'risk-medium',
    low: 'risk-low'
  }
  return map[level] || ''
}

async function loadTopoList() {
  try {
    const res = await topoApi.list({ mode: 'summary' })
    topologies.value = res.items || res || []
  } catch (err) {
    console.error('Failed to load topologies:', err)
  }
}

async function loadAnalysis() {
  if (!selectedTopoId.value) return

  loading.value = true
  try {
    const res = await analyzeApi.byTopoId(selectedTopoId.value)
    result.value = res
  } catch (err) {
    console.error('Failed to analyze:', err)
  } finally {
    loading.value = false
  }
}

onMounted(loadTopoList)
</script>

<style scoped>
.header-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

.topo-select {
  min-width: 200px;
}

.coverage-section {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-lg);
  padding: 32px;
  margin-bottom: 32px;
}

.coverage-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.coverage-header h2 {
  font-size: 18px;
  font-weight: 600;
}

.coverage-value {
  font-size: 48px;
  font-weight: 700;
  font-family: var(--font-mono);
  background: var(--accent-gradient);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.progress-bar.large {
  height: 16px;
  margin-bottom: 20px;
}

.progress-bar.large .progress-fill {
  border-radius: 8px;
}

.coverage-legend {
  display: flex;
  gap: 32px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: var(--text-secondary);
}

.dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
}

.dot.covered { background: var(--success); }
.dot.missing { background: var(--danger); }
.dot.redundant { background: var(--warning); }

.analysis-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24px;
  margin-bottom: 32px;
}

@media (max-width: 900px) {
  .analysis-grid {
    grid-template-columns: 1fr;
  }
}

.function-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.function-item {
  display: flex;
  align-items: flex-start;
  gap: 14px;
  padding: 16px;
  background: var(--bg-tertiary);
  border-radius: var(--radius-md);
}

.function-item.missing {
  border-left: 3px solid var(--danger);
}

.function-item.redundant {
  border-left: 3px solid var(--warning);
}

.func-icon {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-sm);
  flex-shrink: 0;
}

.function-item.missing .func-icon {
  background: rgba(248, 81, 73, 0.1);
  color: var(--danger);
}

.function-item.redundant .func-icon {
  background: rgba(210, 153, 34, 0.1);
  color: var(--warning);
}

.func-icon svg {
  width: 20px;
  height: 20px;
}

.func-info {
  flex: 1;
}

.func-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 4px;
}

.func-desc {
  font-size: 12px;
  color: var(--text-tertiary);
}

.func-category {
  font-size: 12px;
  color: var(--text-tertiary);
  padding: 4px 10px;
  background: var(--bg-secondary);
  border-radius: 12px;
}

.empty-state.small {
  padding: 40px 20px;
}

.empty-state.small svg {
  width: 48px;
  height: 48px;
  color: var(--success);
}

.risk-card {
  margin-bottom: 32px;
}

.risk-badge {
  padding: 6px 14px;
  border-radius: 20px;
  font-size: 13px;
  font-weight: 600;
  text-transform: uppercase;
}

.risk-badge.risk-critical {
  background: rgba(248, 81, 73, 0.15);
  color: var(--danger);
}

.risk-badge.risk-high {
  background: rgba(240, 136, 62, 0.15);
  color: #f0883e;
}

.risk-badge.risk-medium {
  background: rgba(210, 153, 34, 0.15);
  color: var(--warning);
}

.risk-badge.risk-low {
  background: rgba(63, 185, 80, 0.15);
  color: var(--success);
}

.risk-content {
  display: flex;
  gap: 40px;
  align-items: center;
}

.score-circle {
  width: 140px;
  height: 140px;
  border-radius: 50%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  border: 6px solid;
}

.score-circle.risk-critical { border-color: var(--danger); }
.score-circle.risk-high { border-color: #f0883e; }
.score-circle.risk-medium { border-color: var(--warning); }
.score-circle.risk-low { border-color: var(--success); }

.score-value {
  font-size: 36px;
  font-weight: 700;
  font-family: var(--font-mono);
  color: var(--text-primary);
}

.score-max {
  font-size: 14px;
  color: var(--text-tertiary);
}

.risk-details {
  flex: 1;
}

.risk-factors {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 24px;
}

.factor {
  display: flex;
  align-items: center;
  gap: 16px;
}

.factor-label {
  width: 80px;
  font-size: 13px;
  color: var(--text-secondary);
}

.factor-bar {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 12px;
}

.factor-bar .progress-bar {
  flex: 1;
}

.factor-value {
  width: 50px;
  font-size: 13px;
  font-family: var(--font-mono);
  color: var(--text-primary);
  text-align: right;
}

.progress-fill.warning {
  background: var(--warning);
}

.progress-fill.danger {
  background: var(--danger);
}

.risk-advice h4 {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-secondary);
  margin-bottom: 12px;
}

.risk-advice ul {
  list-style: none;
}

.risk-advice li {
  position: relative;
  padding-left: 20px;
  margin-bottom: 8px;
  font-size: 14px;
  color: var(--text-primary);
}

.risk-advice li::before {
  content: '✓';
  position: absolute;
  left: 0;
  color: var(--accent-primary);
}

.action-section {
  display: flex;
  justify-content: center;
}
</style>
