<script setup lang="ts">
import { computed } from 'vue';
import type { StockQuote } from '../types/types.ts'

const props = defineProps<{
  quote: StockQuote;
}>();

const formatCurrency = (val: number) =>
  new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(val);

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
</script>

<template>
  <div class="hidden md:flex grow max-w-md px-8 flex-col justify-center">
    <div class="text-xs text-gray-500 mb-2 font-medium">
      Today's Range
    </div>

    <div class="relative h-2 bg-gray-100 rounded-full">
      <div class="absolute top-1/2 -translate-y-1/2 group" :style="{ left: `${openPosition}%` }">
        <div class="w-2 h-2 bg-gray-400 rounded-full shadow"></div>
      </div>

      <div class="absolute top-1/2 -translate-y-1/2 group" :style="{ left: `${prevClosePosition}%` }">
        <div class="w-2 h-2 bg-gray-300 rounded-full shadow"></div>
      </div>

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
</template>
