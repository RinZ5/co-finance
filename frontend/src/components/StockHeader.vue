<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import type { StockQuote, BasicFinancials, MarketStatus } from '../types/types.ts'

const props = defineProps<{
  quote: StockQuote;
  profile: BasicFinancials;
  marketStatus?: MarketStatus
}>();

const formatCurrency = (val: number) =>
  new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(val);

const now = ref(new Date())

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
    'post-market': 'Post-Market',
    closed: 'Closed'
  }

  return sessionMap[status.session]
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

const isPositive = computed(() => props.quote.d >= 0);

onMounted(() => {
  setInterval(() => {
    now.value = new Date()
  }, 60_000)
})
</script>

<template>
  <div
    class="bg-white rounded-xl shadow-sm border border-gray-100 p-6 flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
    <div>
      <div class="flex items-center gap-2">
        <h1 class="text-3xl font-bold text-slate-900">{{ profile.symbol }}</h1>
        <!-- <span class="px-2 py-0.5 rounded text-xs font-semibold bg-gray-200 text-gray-600">US Market</span> -->
        <span class="px-2 py-0.5 rounded text-xs font-semibold" :class="badgeClass">
          {{ marketLabel }}
        </span>
      </div>
      <p class="text-gray-500 text-sm mt-1">
        {{ marketTime }}
      </p>
    </div>

    <div class="text-right">
      <div class="text-4xl font-bold tracking-tight text-slate-900">
        {{ formatCurrency(quote.c) }}
      </div>
      <div
        :class="['flex items-center justify-end gap-2 text-lg font-medium', isPositive ? 'text-green-600' : 'text-red-600']">
        <span>{{ isPositive ? '+' : '' }}{{ quote.d.toFixed(2) }}</span>
        <span class="text-sm bg-opacity-10 px-2 py-0.5 rounded-full"
          :class="isPositive ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'">
          {{ isPositive ? '+' : '' }}{{ quote.dp.toFixed(2) }}%
        </span>
      </div>
    </div>
  </div>
</template>
