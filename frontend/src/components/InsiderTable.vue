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
  <div class="bg-white rounded-xl shadow-sm border border-gray-100 p-6 h-full overflow-hidden flex flex-col">
    <h2 class="text-lg font-semibold text-slate-800 mb-4">Insider Activity</h2>
    <div class="overflow-auto flex-1">
      <table class="w-full text-left text-sm text-gray-600">
        <thead class="bg-gray-50 text-xs uppercase text-gray-500 font-medium">
          <tr>
            <th class="px-4 py-3 rounded-l-lg">Name</th>
            <th class="px-4 py-3">Date</th>
            <th class="px-4 py-3">Shares</th>
            <th class="px-4 py-3 rounded-r-lg text-right">Value</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-100">
          <tr v-for="(t, index) in recentTransactions" :key="index" class="hover:bg-gray-50/50">
            <td class="px-4 py-3 font-medium text-slate-800">{{ t.name }}</td>
            <td class="px-4 py-3 whitespace-nowrap">{{ formatYMDToReadable(t.transactionDate) }}</td>
            <td class="px-4 py-3">
              <span :class="t.change > 0 ? 'text-green-600' : 'text-red-600'">
                {{ t.change > 0 ? '+' : '' }}{{ t.change.toLocaleString() }}
              </span>
            </td>
            <td class="px-4 py-3 text-right">
              {{ t.transactionPrice > 0 ? `$${t.transactionPrice.toFixed(2)}` : '-' }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
