<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted } from 'vue';
import type { MarketStatus } from '../types/types.ts'

const props = defineProps<{
  symbol: string;
  marketStatus?: MarketStatus
}>();

const now = ref(new Date())
let timer: number | null = null

const marketTime = computed(() => {
  if (!props.marketStatus) return '—'

  return now.value.toLocaleTimeString('en-US', {
    timeZone: props.marketStatus.timezone,
    hour: '2-digit',
    minute: '2-digit',
  })
})

const marketLabel = computed(() => {
  const status = props.marketStatus
  if (!status) return '-'
  if (status.holiday) return 'Holiday'

  const sessionMap: Record<MarketStatus['session'], string> = {
    regular: 'Open',
    'pre-market': 'Pre-Market',
    'post-market': 'Post-Market'
  }

  return sessionMap[status.session] || "Closed"
})

const badgeClass = computed(() => {
  if (!props.marketStatus) {
    return 'bg-gray-100 text-gray-500'
  }

  switch (props.marketStatus.session) {
    case 'regular':
      return 'bg-green-100 text-green-700'
    case 'pre-market':
      return 'bg-yellow-100 text-yellow-700'
    case 'post-market':
      return 'bg-blue-100 text-blue-700'
    default:
      return 'bg-gray-200 text-gray-600'
  }
})

const updateClock = () => {
  const currentDate = new Date()
  now.value = currentDate

  const msUntilNextMinute = (60 - currentDate.getSeconds()) * 1000 - currentDate.getMilliseconds() + 50

  timer = window.setTimeout(updateClock, msUntilNextMinute)
}

onMounted(() => {
  updateClock()
})

onUnmounted(() => {
  if (timer) clearTimeout(timer)
})
</script>

<template>
  <div>
    <h1 class="text-3xl font-bold text-slate-900">
      {{ symbol }}
    </h1>

    <div class="flex items-center gap-2">
      <p class="text-gray-500 text-sm mt-1">
        {{ marketTime }}
      </p>

      <span class="px-2 py-0.5 rounded text-xs font-semibold" :class="badgeClass">
        {{ marketLabel }}
      </span>
    </div>
  </div>
</template>
