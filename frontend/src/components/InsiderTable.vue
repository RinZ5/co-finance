<script setup lang="ts">
import { computed } from 'vue';
import type { InsiderTransaction } from '../types/types.ts';

const props = defineProps<{
  transactions: InsiderTransaction[];
}>();

const recentTransactions = computed(() => {
  return [...props.transactions]
    .sort((a, b) => new Date(b.transactionDate).getTime() - new Date(a.transactionDate).getTime())
    .slice(0, 5);
});

const formatYMDToReadable = (date: string) => {
  const parts = date.split('-')

  if (parts.length !== 3) return date

  const y = Number(parts[0])
  const m = Number(parts[1])
  const d = Number(parts[2])

  return new Date(y, m - 1, d).toLocaleDateString('en-US', {
    day: 'numeric',
    month: 'short',
    year: 'numeric'
  })
}
</script>

<template>
  <div class="bg-white rounded-xl shadow-sm border border-gray-100 p-5 h-full overflow-hidden flex flex-col">
    <h2 class="text-lg font-semibold text-slate-800 mb-4">Insider Activity</h2>

    <div class="hidden md:block flex-1 space-y-3">
      <div v-for="(t, index) in recentTransactions" :key="index"
        class="flex items-center justify-between p-4 hover:bg-gray-50/80 rounded-lg border border-gray-100 transition-colors duration-200">

        <div class="flex-1 min-w-0 pr-4">
          <div class="text-base font-semibold text-slate-800 mb-1 truncate leading-tight">
            {{ t.name }}
          </div>
          <div class="text-sm text-gray-500 font-medium">
            {{ formatYMDToReadable(t.transactionDate) }}
          </div>
        </div>

        <div class="text-right shrink-0 min-w-30">
          <div class="text-base font-bold mb-1" :class="t.change > 0 ? 'text-green-600' : 'text-red-600'">
            {{ t.change > 0 ? '+' : '' }}{{ t.change.toLocaleString() + " shares" }}
          </div>
          <div class="text-sm font-semibold text-slate-700">
            ${{ t.transactionPrice.toFixed(2) }}
          </div>
        </div>
      </div>
    </div>

    <div class="hidden sm:block md:hidden flex-1 overflow-auto">
      <table class="w-full text-left text-sm text-gray-600">
        <thead class="bg-gray-50 text-xs uppercase text-gray-500 font-medium">
          <tr>
            <th class="px-3 py-2 rounded-l-lg">Name</th>
            <th class="px-3 py-2">Date</th>
            <th class="px-3 py-3 rounded-r-lg text-right">Shares</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-100">
          <tr v-for="(t, index) in recentTransactions" :key="index" class="hover:bg-gray-50/50">
            <td class="px-3 py-2 font-medium text-slate-800 text-sm truncate max-w-30">{{ t.name }}</td>
            <td class="px-3 py-2 whitespace-nowrap text-xs">{{ formatYMDToReadable(t.transactionDate) }}</td>
            <td class="px-3 py-3 text-right">
              <div class="text-sm font-bold" :class="t.change > 0 ? 'text-green-600' : 'text-red-600'">
                {{ t.change > 0 ? '+' : '' }}{{ t.change.toLocaleString() + " shares" }}
              </div>
              <div class="text-xs text-gray-500">
                ${{ t.transactionPrice.toFixed(2) }}
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div class="sm:hidden flex-1 space-y-2">
      <div v-for="(t, index) in recentTransactions" :key="index"
        class="flex items-center justify-between p-3 bg-white hover:bg-gray-50/80 rounded-lg border border-gray-100 transition-colors duration-200">
        <div class="flex-1 min-w-0 pr-3">
          <div class="text-sm font-semibold text-slate-800 truncate leading-tight">
            {{ t.name }}
          </div>
          <div class="text-xs text-gray-500 font-medium">
            {{ formatYMDToReadable(t.transactionDate) }}
          </div>
        </div>

        <div class="text-right shrink-0">
          <div class="text-sm font-bold" :class="t.change > 0 ? 'text-green-600' : 'text-red-600'">
            {{ t.change > 0 ? '+' : '' }}{{ t.change.toLocaleString() }} shares
          </div>
          <div class="text-xs font-semibold text-slate-700">
            ${{ t.transactionPrice.toFixed(2) }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
