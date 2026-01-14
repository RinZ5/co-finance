<script setup lang="ts">
import { computed } from 'vue';
import { useStockStore } from '../stores/stock';
import type { CompanyNews } from '../types/types';

import { VueDatePicker } from '@vuepic/vue-datepicker';
import '@vuepic/vue-datepicker/dist/main.css';

const props = defineProps<{
  news: CompanyNews[];
}>();

const store = useStockStore();

const formatDisplayDate = (unixTime: number) => {
  return new Date(unixTime * 1000).toLocaleString('en-US', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  });
};

const startOfDay = (date: Date) =>
  new Date(date.getFullYear(), date.getMonth(), date.getDate());

const endOfDay = (date: Date) =>
  new Date(date.getFullYear(), date.getMonth(), date.getDate(), 23, 59, 59, 999);

const filteredNews = computed(() => {
  if (!store.newsRange || store.newsRange.length !== 2) {
    return props.news
  }

  const [from, to] = store.newsRange
  const start = startOfDay(from)
  const end = endOfDay(to)

  return props.news.filter(item => {
    const d = new Date(item.datetime * 1000)
    return d >= start && d <= end
  })
})
</script>

<template>
  <div class="bg-white rounded-xl shadow-sm border border-gray-100 p-6 h-full flex flex-col">

    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-4 border-b border-gray-100 pb-4">
      <h2 class="text-lg font-semibold text-slate-800 shrink-0">
        Company News
      </h2>

      <div class="flex-1">
        <div class="flex flex-col">
          <label class="text-[10px] text-gray-400 font-medium ml-1">
            Date range
          </label>

          <VueDatePicker v-model="store.newsRange" range :enable-time-picker="false" :time-picker="false" auto-apply
            format="dd MMM yyyy" class="text-xs w-full" />
        </div>
      </div>
    </div>

    <div class="overflow-y-auto pr-2 custom-scrollbar flex-1 max-h-[600px] space-y-4">

      <div v-if="news.length === 0" class="text-gray-400 text-sm text-center py-10">
        No news found in selected range
      </div>

      <a v-for="item in filteredNews" :key="item.id" :href="item.url" target="_blank" rel="noopener noreferrer"
        class="block group hover:bg-gray-50 p-3 rounded-lg transition-colors border border-transparent hover:border-gray-100">
        <div class="flex gap-4">
          <div class="shrink-0">
            <img v-if="item.image" :src="item.image" alt="News"
              class="w-20 h-20 object-cover rounded-md border border-gray-200 bg-gray-100"
              @error="(e) => (e.target as HTMLImageElement).style.display = 'none'">
            <div v-else
              class="w-20 h-20 bg-gray-100 rounded-md flex items-center justify-center text-gray-400 text-[10px] text-center border border-gray-200">
              No Image
            </div>
          </div>

          <div class="flex-1 min-w-0">
            <div class="flex justify-between items-start gap-2">
              <span
                class="text-[10px] font-bold text-blue-600 bg-blue-50 px-2 py-0.5 rounded uppercase tracking-wide truncate">
                {{ item.source }}
              </span>
              <span class="text-[10px] text-gray-400 whitespace-nowrap">{{ formatDisplayDate(item.datetime) }}</span>
            </div>

            <h3
              class="text-sm font-bold text-slate-800 mt-1 leading-snug group-hover:text-blue-600 transition-colors line-clamp-2">
              {{ item.headline }}
            </h3>
            <p class="text-xs text-gray-500 mt-1 line-clamp-2 leading-relaxed">
              {{ item.summary }}
            </p>
          </div>
        </div>
      </a>
    </div>
  </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
}

.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}

.custom-scrollbar::-webkit-scrollbar-thumb {
  background-color: #e2e8f0;
  border-radius: 4px;
}

.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background-color: #cbd5e1;
}
</style>
