<script setup lang="ts">
import { computed } from 'vue';
import type { RecommendationTrend } from '../types/types.ts';

const props = defineProps<{
  trend: RecommendationTrend[];
}>();

const latest = computed(() => props.trend[0]);

const total = computed(() => {
  const t = latest.value;
  if (!t) return 1;
  return t.strongBuy + t.buy + t.hold + t.sell + t.strongSell;
});

const getWidth = (val: number) => `${(val / total.value) * 100}%`;
</script>

<template>
  <div class="bg-white rounded-xl shadow-sm border border-gray-100 p-6" v-if="latest">
    <h2 class="text-lg font-semibold text-slate-800 mb-4 border-b border-gray-100 pb-2">Analyst Consensus</h2>

    <div class="flex items-end gap-2 mb-2">
      <span class="text-3xl font-bold text-slate-900">{{ ((latest.buy + latest.strongBuy) / total * 100).toFixed(0)
        }}%</span>
      <span class="text-sm text-gray-500 mb-1.5">Buy Rating</span>
    </div>

    <div class="h-4 w-full rounded-full flex overflow-hidden bg-gray-100">
      <div class="bg-green-600 h-full" :style="{ width: getWidth(latest.strongBuy) }" title="Strong Buy"></div>
      <div class="bg-green-400 h-full" :style="{ width: getWidth(latest.buy) }" title="Buy"></div>
      <div class="bg-yellow-400 h-full" :style="{ width: getWidth(latest.hold) }" title="Hold"></div>
      <div class="bg-red-400 h-full" :style="{ width: getWidth(latest.sell) }" title="Sell"></div>
      <div class="bg-red-600 h-full" :style="{ width: getWidth(latest.strongSell) }" title="Strong Sell"></div>
    </div>

    <div class="flex flex-wrap gap-x-4 gap-y-2 mt-4 text-xs text-gray-600">
      <div class="flex items-center gap-1"><span class="w-2 h-2 rounded-full bg-green-600"></span> Strong Buy ({{
        latest.strongBuy }})</div>
      <div class="flex items-center gap-1"><span class="w-2 h-2 rounded-full bg-green-400"></span> Buy ({{ latest.buy
        }})</div>
      <div class="flex items-center gap-1"><span class="w-2 h-2 rounded-full bg-yellow-400"></span> Hold ({{ latest.hold
        }})</div>
      <div class="flex items-center gap-1"><span class="w-2 h-2 rounded-full bg-red-400"></span> Sell ({{ latest.sell
        }})</div>
    </div>
  </div>
</template>
