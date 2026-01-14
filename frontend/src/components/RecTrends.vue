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
  <div class="bg-white rounded-xl shadow-sm border border-gray-100 p-6 h-full flex flex-col" v-if="latest">

    <h2 class="text-lg font-semibold text-slate-800 mb-4 border-b border-gray-100 pb-2">
      Analyst Consensus
    </h2>

    <div class="mb-6">
      <div class="flex items-end gap-2 mb-3">
        <span class="text-3xl font-bold text-slate-900">
          {{ ((latest.buy + latest.strongBuy) / total * 100).toFixed(0) }}%
        </span>
        <span class="text-sm text-gray-500 mb-1.5 font-medium">Buy Rating</span>
      </div>

      <div class="h-3 w-full rounded-full flex overflow-hidden bg-gray-100">
        <div class="bg-green-600 h-full" :style="{ width: getWidth(latest.strongBuy) }"></div>
        <div class="bg-green-400 h-full" :style="{ width: getWidth(latest.buy) }"></div>
        <div class="bg-yellow-400 h-full" :style="{ width: getWidth(latest.hold) }"></div>
        <div class="bg-red-400 h-full" :style="{ width: getWidth(latest.sell) }"></div>
        <div class="bg-red-600 h-full" :style="{ width: getWidth(latest.strongSell) }"></div>
      </div>
    </div>

    <div class="grow flex flex-col justify-end">

      <div class="divide-y divide-gray-100">

        <div class="flex justify-between items-center py-2 text-sm">
          <div class="flex items-center gap-2">
            <span class="w-2 h-2 rounded-full bg-green-600"></span>
            <span class="text-gray-600">Strong Buy</span>
          </div>
          <span class="font-medium text-slate-900">{{ latest.strongBuy }}</span>
        </div>

        <div class="flex justify-between items-center py-2 text-sm">
          <div class="flex items-center gap-2">
            <span class="w-2 h-2 rounded-full bg-green-400"></span>
            <span class="text-gray-600">Buy</span>
          </div>
          <span class="font-medium text-slate-900">{{ latest.buy }}</span>
        </div>

        <div class="flex justify-between items-center py-2 text-sm">
          <div class="flex items-center gap-2">
            <span class="w-2 h-2 rounded-full bg-yellow-400"></span>
            <span class="text-gray-600">Hold</span>
          </div>
          <span class="font-medium text-slate-900">{{ latest.hold }}</span>
        </div>

        <div class="flex justify-between items-center py-2 text-sm">
          <div class="flex items-center gap-2">
            <span class="w-2 h-2 rounded-full bg-red-400"></span>
            <span class="text-gray-600">Sell / Strong Sell</span>
          </div>
          <span class="font-medium text-slate-900">{{ latest.sell + latest.strongSell }}</span>
        </div>

      </div>
    </div>

  </div>
</template>
