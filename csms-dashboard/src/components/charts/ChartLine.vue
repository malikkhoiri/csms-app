<template>
  <v-card>
    <v-card-title>Grafik Sesi & Energi Mingguan</v-card-title>
    <v-card-text>
      <canvas ref="canvas"></canvas>
    </v-card-text>
  </v-card>
</template>

<script setup>
import { onMounted, ref, watch } from 'vue'
import { Chart, LineController, LineElement, PointElement, LinearScale, Title, CategoryScale } from 'chart.js'

Chart.register(LineController, LineElement, PointElement, LinearScale, Title, CategoryScale)

const props = defineProps({
  chartData: Object,
})
const canvas = ref(null)
let chartInstance = null

onMounted(() => {
  if (canvas.value) {
    chartInstance = new Chart(canvas.value, {
      type: 'line',
      data: props.chartData,
      options: { responsive: true, plugins: { legend: { display: true } } },
    })
  }
})

watch(() => props.chartData, (newData) => {
  if (chartInstance) {
    chartInstance.data = newData
    chartInstance.update()
  }
})
</script> 