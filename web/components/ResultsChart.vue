<template>
  <div class="min-h-72">
    <div v-if="loading">
      loading...
    </div>
    <BarChart v-else v-bind="barChartProps" :key="version" class="h-72" />
  </div>
</template>

<script setup lang="ts">
import { BarChart, useBarChart } from 'vue-chart-3'
import { ChartData, ChartOptions } from 'chart.js'
import { Result } from '~/types/results'

const props = defineProps({
  features: {
    type: Array,
    default: () => []
  },
  templates: {
    type: Array,
    default: () => []
  },
  version: {
    type: Number,
    default: 0
  }
})

const loading = ref(null)
loading.value = true

const client = useSupabaseClient()
const { data: results } = await client.rpc<Result>('results', {
  version: props.version
})

const groupResultsByFeature = (results: Result[]) => {
  const groupedResults = {
    passed: {},
    failed: {},
    skipped: {},
    undefined: {}
  }
  const features = [...(props.features as string[]), 'other']
  features.forEach((feature: string) => {
    groupedResults.passed[feature] = []
    groupedResults.failed[feature] = []
    groupedResults.skipped[feature] = []
    groupedResults.undefined[feature] = []
  })
  if (!results) {
    return {
      passed: [],
      failed: [],
      skipped: [],
      undefined: []
    }
  }
  results.forEach((result) => {
    switch (result.status) {
      case 'passed':
        if (features.includes(result.feature)) {
          groupedResults.passed[result.feature].push(result)
        } else if (features.includes(result.parent_suite)) {
          groupedResults.passed[result.parent_suite].push(result)
        } else if (features.includes(result.suite)) {
          groupedResults.passed[result.suite].push(result)
        } else if (features.includes(result.sub_suite)) {
          groupedResults.passed[result.sub_suite].push(result)
        } else {
          // eslint-disable-next-line dot-notation
          groupedResults.passed['other'].push(result)
        }
        break
      case 'broken':
      case 'failed':
        if (features.includes(result.feature)) {
          groupedResults.failed[result.feature].push(result)
        } else if (features.includes(result.parent_suite)) {
          groupedResults.failed[result.parent_suite].push(result)
        } else if (features.includes(result.suite)) {
          groupedResults.failed[result.suite].push(result)
        } else if (features.includes(result.sub_suite)) {
          groupedResults.failed[result.sub_suite].push(result)
        } else {
          // eslint-disable-next-line dot-notation
          groupedResults.failed['other'].push(result)
        }
        break
      case 'skipped':
        if (features.includes(result.feature)) {
          groupedResults.skipped[result.feature].push(result)
        } else if (features.includes(result.parent_suite)) {
          groupedResults.skipped[result.parent_suite].push(result)
        } else if (features.includes(result.suite)) {
          groupedResults.skipped[result.suite].push(result)
        } else if (features.includes(result.sub_suite)) {
          groupedResults.skipped[result.sub_suite].push(result)
        } else {
          // eslint-disable-next-line dot-notation
          groupedResults.skipped['other'].push(result)
        }
        break
      default:
        if (features.includes(result.feature)) {
          groupedResults.undefined[result.feature].push(result)
        } else if (features.includes(result.parent_suite)) {
          groupedResults.undefined[result.parent_suite].push(result)
        } else if (features.includes(result.suite)) {
          groupedResults.undefined[result.suite].push(result)
        } else if (features.includes(result.sub_suite)) {
          groupedResults.undefined[result.sub_suite].push(result)
        } else {
          // eslint-disable-next-line dot-notation
          groupedResults.undefined['other'].push(result)
        }
        break
    }
  })

  const groupedResultsWithCounts = {
    passed: [],
    failed: [],
    skipped: [],
    undefined: []
  }
  groupedResultsWithCounts.passed = Object.entries(groupedResults.passed).map(([, results]) => ((results as []).length))
  groupedResultsWithCounts.failed = Object.entries(groupedResults.failed).map(([, results]) => ((results as []).length))
  groupedResultsWithCounts.skipped = Object.entries(groupedResults.skipped).map(([, results]) => ((results as []).length))
  groupedResultsWithCounts.undefined = Object.entries(groupedResults.undefined).map(([, results]) => ((results as []).length))
  return groupedResultsWithCounts
}

const groupResults = groupResultsByFeature(results)
const groupedTemplates = groupResultsByFeature(props.templates as Result[])

const testData = computed<ChartData>(() => ({
  labels: [...props.features, 'other'],
  datasets: [
    {
      type: 'bar',
      label: 'Passed',
      backgroundColor: '#88b04bc8',
      data: groupResults.passed
    },
    {
      type: 'bar',
      label: 'Failed',
      backgroundColor: '#ff6f61c8',
      data: groupResults.failed
    },
    {
      type: 'bar',
      label: 'Skipped',
      backgroundColor: '#f5df4dc8',
      data: groupResults.skipped
    },
    {
      type: 'bar',
      label: 'Undefined',
      backgroundColor: '#0f4c81c8',
      data: groupResults.undefined
    }
  ]
}))

if (props.templates?.length > 0) {
  testData.value.datasets = [...testData.value.datasets, {
    type: 'scatter',
    radius: 16,
    borderWidth: 3,
    hoverRadius: 20,
    hoverBorderWidth: 3,
    pointStyle: 'line',
    label: 'Template Run',
    backgroundColor: '#939597',
    borderColor: '#939597',
    data: groupedTemplates.passed.map(function (num, i) {
      return num + groupedTemplates.failed[i] + groupedTemplates.skipped[i] + groupedTemplates.undefined[i]
    })
  }]
}

const options = computed<ChartOptions<'bar'>>(() => ({
  plugins: {
    legend: {
      position: 'top'
    },
    title: {
      display: true,
      text: 'Results by feature/suite'
    }
  },
  responsive: true,
  maintainAspectRatio: false,
  scales: {
    x: {
      stacked: true,
      barThickness: 6, // number (pixels) or 'flex'
      maxBarThickness: 8 // number (pixels)
    },
    y: {
      stacked: true
    }
  }
}))

const { barChartProps } = useBarChart({
  chartData: testData,
  options
})

loading.value = false

</script>
