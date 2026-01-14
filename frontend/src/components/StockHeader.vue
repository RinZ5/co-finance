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

function clamp(val: number) {
  return Math.min(Math.max(val, 0), 100)
}

const rangePosition = computed(() => {
  const { h, l, c } = props.quote
  if (h === l) return 50
  return clamp(((c - l) / (h - l)) * 100)
})

const openPosition = computed(() => {
  const { h, l, o } = props.quote
  if (h === l) return 50
  return clamp(((o - l) / (h - l)) * 100)
})

const prevClosePosition = computed(() => {
  const { h, l, pc } = props.quote
  if (h === l) return 50
  return clamp(((pc - l) / (h - l)) * 100)
})

onMounted(() => {
  setInterval(() => {
    now.value = new Date()
  }, 60_000)
})
</script>

<template>
  <div class="bg-white rounded-xl shadow-sm border border-gray-100 p-6
           flex flex-col gap-4
           md:flex-row md:items-center md:justify-between">

    <!-- LEFT (Symbol + Mobile Price + Mobile Stats) -->
    <div class="w-full md:w-auto">

      <!-- Top row -->
      <div class="flex justify-between gap-4">

        <!-- Symbol -->
        <div>
          <h1 class="text-3xl font-bold text-slate-900">
            {{ profile.symbol }}
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

        <!-- Mobile price -->
        <div class="md:hidden text-right">
          <div class="text-3xl font-bold tracking-tight text-slate-900">
            {{ formatCurrency(quote.c) }}
          </div>

          <div :class="[
            'flex items-center justify-end gap-2 text-sm font-medium',
            isPositive ? 'text-green-600' : 'text-red-600'
          ]">
            <span>{{ isPositive ? '+' : '' }}{{ quote.d.toFixed(2) }}</span>
            <span>({{ isPositive ? '+' : '' }}{{ quote.dp.toFixed(2) }}%)</span>
          </div>
        </div>

      </div>

      <!-- Mobile stats -->
      <div class="md:hidden mt-4 pt-4 border-t border-gray-100
               grid grid-cols-2 gap-x-6 gap-y-3 text-xs">
        <div class="flex justify-between">
          <span class="text-gray-500">Open</span>
          <span class="font-medium text-slate-800">
            {{ formatCurrency(quote.o) }}
          </span>
        </div>

        <div class="flex justify-between">
          <span class="text-gray-500">High</span>
          <span class="font-medium text-slate-800">
            {{ formatCurrency(quote.h) }}
          </span>
        </div>

        <div class="flex justify-between">
          <span class="text-gray-500">Low</span>
          <span class="font-medium text-slate-800">
            {{ formatCurrency(quote.l) }}
          </span>
        </div>

        <div class="flex justify-between">
          <span class="text-gray-500">Prev Close</span>
          <span class="font-medium text-slate-800">
            {{ formatCurrency(quote.pc) }}
          </span>
        </div>
      </div>

    </div>

    <!-- CENTER (Desktop range bar only) -->
    <div class="hidden md:flex grow max-w-md px-8 flex-col justify-center">
      <div class="text-xs text-gray-500 mb-2 font-medium">
        Today’s Range
      </div>

      <div class="relative h-2 bg-gray-100 rounded-full">

        <!-- Open -->
        <div class="absolute top-1/2 -translate-y-1/2 group" :style="{ left: `${openPosition}%` }">
          <div class="w-2 h-2 bg-gray-400 rounded-full shadow"></div>
        </div>

        <!-- Prev close -->
        <div class="absolute top-1/2 -translate-y-1/2 group" :style="{ left: `${prevClosePosition}%` }">
          <div class="w-2 h-2 bg-gray-300 rounded-full shadow"></div>
        </div>

        <!-- Current -->
        <div class="absolute top-1/2 -translate-y-1/2 group" :style="{ left: `${rangePosition}%` }">
          <div class="w-3 h-3 border-2 border-white rounded-full shadow"
            :class="quote.c >= quote.pc ? 'bg-green-600' : 'bg-red-600'"></div>
        </div>

      </div>

      <div class="flex justify-between text-[10px] text-gray-400 mt-2 font-mono">
        <span>L: {{ formatCurrency(quote.l) }}</span>
        <span>H: {{ formatCurrency(quote.h) }}</span>
      </div>

      <div class="flex justify-between text-[10px] text-gray-500 mt-1">
        <span>O: {{ formatCurrency(quote.o) }}</span>
        <span>PC: {{ formatCurrency(quote.pc) }}</span>
      </div>
    </div>

    <!-- RIGHT (Desktop price only) -->
    <div class="hidden md:block text-right">
      <div class="text-4xl font-bold tracking-tight text-slate-900">
        {{ formatCurrency(quote.c) }}
      </div>

      <div :class="[
        'flex items-center justify-end gap-2 text-lg font-medium',
        isPositive ? 'text-green-600' : 'text-red-600'
      ]">
        <span>{{ isPositive ? '+' : '' }}{{ quote.d.toFixed(2) }}</span>

        <span class="text-sm px-2 py-0.5 rounded-full bg-opacity-10" :class="isPositive
          ? 'bg-green-100 text-green-700'
          : 'bg-red-100 text-red-700'">
          {{ isPositive ? '+' : '' }}{{ quote.dp.toFixed(2) }}%
        </span>
      </div>
    </div>

  </div>
</template>
