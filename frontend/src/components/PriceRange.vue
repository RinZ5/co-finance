<script setup lang="ts">
import { computed } from 'vue';
import type { StockQuote } from '../types/types.ts'
import { useWebSocketStore } from '../stores/websocketStore.ts';

const props = defineProps<{
  quote: StockQuote;
}>();

const wsStore = useWebSocketStore()

const currentPrice = computed(() => wsStore.currentPrice ?? props.quote.c);

const dayHigh = computed(() => Math.max(props.quote.h, currentPrice.value));
const dayLow = computed(() => Math.min(props.quote.l, currentPrice.value));

const formatCurrency = (val: number) =>
  new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(val);

function clamp(val: number) {
  return Math.min(Math.max(val, 0), 100)
}

const rangePosition = computed(() => {
  const h = dayHigh.value;
  const l = dayLow.value;
  if (h === l) return 50;
  return clamp(((currentPrice.value - l) / (h - l)) * 100);
});

const openPosition = computed(() => {
  const h = dayHigh.value;
  const l = dayLow.value;
  const o = props.quote.o;
  if (h === l) return 50;
  return clamp(((o - l) / (h - l)) * 100);
});

const prevClosePosition = computed(() => {
  const h = dayHigh.value;
  const l = dayLow.value;
  const pc = props.quote.pc;
  if (h === l) return 50;
  return clamp(((pc - l) / (h - l)) * 100);
});
</script>

<template>
  <div class="hidden md:flex w-full px-8 flex-col justify-center">
    <div class="text-xs text-gray-500 mb-2 font-medium">
      Today's Range
    </div>

    <div class="relative h-2 bg-gray-100 rounded-full">
      <div class="absolute top-1/2 -translate-y-1/2 group z-0" :style="{ left: `${openPosition}%` }">
        <div class="w-2 h-2 bg-gray-400 rounded-full shadow cursor-help"></div>
        <div
          class="absolute bottom-full mb-2 left-1/2 -translate-x-1/2 hidden group-hover:flex flex-col items-center z-20">
          <div class="bg-gray-700 text-white text-[10px] px-2 py-1 rounded shadow-lg whitespace-nowrap font-mono">
            Open: {{ formatCurrency(quote.o) }}
          </div>
          <div
            class="w-0 h-0 border-l-4 border-l-transparent border-r-4 border-r-transparent border-t-4 border-t-gray-700">
          </div>
        </div>
      </div>

      <div class="absolute top-1/2 -translate-y-1/2 group z-0" :style="{ left: `${prevClosePosition}%` }">
        <div class="w-2 h-2 bg-gray-300 rounded-full shadow cursor-help"></div>
        <div
          class="absolute bottom-full mb-2 left-1/2 -translate-x-1/2 hidden group-hover:flex flex-col items-center z-20">
          <div class="bg-gray-400 text-white text-[10px] px-2 py-1 rounded shadow-lg whitespace-nowrap font-mono">
            Prev Close: {{ formatCurrency(quote.pc) }}
          </div>
          <div
            class="w-0 h-0 border-l-4 border-l-transparent border-r-4 border-r-transparent border-t-4 border-t-gray-400">
          </div>
        </div>
      </div>

      <div class="absolute top-1/2 -translate-y-1/2 group transition-all duration-500 ease-out z-10"
        :style="{ left: `${rangePosition}%` }">
        <div class="w-3 h-3 border-2 border-white rounded-full shadow cursor-pointer"
          :class="currentPrice >= props.quote.pc ? 'bg-green-600' : 'bg-red-600'">
        </div>
        <div class="absolute bottom-full mb-2 left-1/2 -translate-x-1/2 hidden group-hover:flex flex-col items-center">
          <div class="bg-slate-800 text-white text-[10px] px-2 py-1 rounded shadow-lg whitespace-nowrap font-mono">
            Current: {{ formatCurrency(currentPrice) }}
          </div>
          <div
            class="w-0 h-0 border-l-4 border-l-transparent border-r-4 border-r-transparent border-t-4 border-t-slate-800">
          </div>
        </div>
      </div>
    </div>

    <div class="flex justify-between text-[10px] text-gray-400 mt-2 font-mono">
      <span>L: {{ formatCurrency(dayLow) }}</span>
      <span>H: {{ formatCurrency(dayHigh) }}</span>
    </div>

    <div class="flex justify-between text-[10px] text-gray-500 mt-1">
      <span>O: {{ formatCurrency(quote.o) }}</span>
      <span>PC: {{ formatCurrency(quote.pc) }}</span>
    </div>
  </div>
</template>
