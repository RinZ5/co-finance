<script setup lang="ts">
import { computed } from 'vue';
import type { EarningsSurprise } from '../types/types.ts';

const props = defineProps<{
  earnings: EarningsSurprise[];
}>();

const maxVal = computed(() => {
  if (!props.earnings?.length) return 1;
  const values = props.earnings.flatMap(e => [Math.abs(e.actual), Math.abs(e.estimate)]);
  const max = Math.max(...values);
  return max === 0 ? 1 : max * 1.1; // 10% headroom
});

const getHeight = (val: number) => {
  if (!val) return '0%';
  const percentage = (Math.abs(val) / maxVal.value) * 100;
  return `${Math.min(percentage, 100)}%`;
};
</script>

<template>
  <div class="bg-white rounded-xl shadow-sm border border-gray-100 p-6 h-full flex flex-col">
    <div class="flex justify-between items-center mb-6">
      <h2 class="text-lg font-semibold text-slate-800">Earnings Surprise</h2>
      <div class="flex gap-4 text-xs">
        <div class="flex items-center gap-1"><span class="w-2 h-2 rounded-full bg-blue-500"></span> Actual</div>
        <div class="flex items-center gap-1"><span class="w-2 h-2 rounded-full bg-gray-300"></span> Estimate</div>
      </div>
    </div>

    <div class="flex-1 flex items-end justify-between gap-2 min-h-[200px] pt-4">
      <div v-for="item in earnings" :key="item.period" class="flex flex-col items-center w-full group">

        <div class="flex items-end gap-1 h-32 w-full justify-center relative">
          <div class="w-3 md:w-6 bg-gray-300 rounded-t-sm transition-all duration-500 relative group-hover:bg-gray-400"
            :style="{ height: getHeight(item.estimate) }">
          </div>
          <div class="w-3 md:w-6 bg-blue-500 rounded-t-sm transition-all duration-500 relative group-hover:bg-blue-600"
            :style="{ height: getHeight(item.actual) }">
          </div>
        </div>

        <div class="mt-3 text-xs text-gray-500 text-center font-medium">
          Q{{ item.quarter }} '{{ item.year.toString().slice(-2) }}
          <div :class="['text-[10px] mt-0.5', item.surprisePercent > 0 ? 'text-green-600' : 'text-red-600']">
            {{ item.surprisePercent > 0 ? '+' : '' }}{{ item.surprisePercent.toFixed(1) }}%
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
