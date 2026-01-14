<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue';
import type { CompanyNews, DashboardData } from './types/types.ts';
import NProgress from 'nprogress';
import 'nprogress/nprogress.css';

import SearchBar from './components/SearchBar.vue';
import StockHeader from './components/StockHeader.vue';
import KeyStats from './components/KeyStats.vue';
import EarningsChart from './components/EarningsChart.vue';
import RecTrends from './components/RecTrends.vue';
import InsiderTable from './components/InsiderTable.vue';
import NewsFeed from './components/NewsFeed.vue';
import { useStockStore } from './stores/stock.ts';
import { useMarketStore } from './stores/market.ts';

NProgress.configure({ showSpinner: false })

const store = useStockStore();
const marketStore = useMarketStore()

const dashboardData = ref<DashboardData | null>(null);
const companyNews = ref<CompanyNews[]>([]);
const error = ref<string | null>(null);

const apiFrom = computed(() => {
  if (!store.newsRange || store.newsRange.length !== 2) return null

  const [from] = store.newsRange
  if (!from) return null

  return from.toISOString().split('T')[0]
})

const apiTo = computed(() => {
  if (!store.newsRange || store.newsRange.length !== 2) return null

  const [, to] = store.newsRange
  if (!to) return null

  return to.toISOString().split('T')[0]
})

const fetchData = async () => {
  NProgress.start();
  error.value = null;

  try {
    const response = await fetch(`/api/dashboard?symbol=${store.symbol}`);

    if (!response.ok) {
      throw new Error(`Error: ${response.status} ${response.statusText}`);
    }

    dashboardData.value = await response.json();


    if (!apiFrom.value || !apiTo.value) {
      companyNews.value = []
      return
    }

    const newsUrl = `/api/company-news?symbol=${store.symbol}&from=${apiFrom.value}&to=${apiTo.value}`;


    const newsResponse = await fetch(newsUrl);

    if (newsResponse.ok) {
      companyNews.value = await newsResponse.json();
    } else {
      companyNews.value = [];
    }
  } catch (err: any) {
    error.value = err.message || 'Failed to load data';
    console.error(err);
  } finally {
    NProgress.done();
  }
};

onMounted(async () => {
  await marketStore.fetchMarketStatus('US')
  marketStore.stopMarketSync
  fetchData();
});

watch(
  () => store.symbol,
  () => {
    fetchData();
  }
);

watch(
  () => store.newsRange,
  ([from, to]) => {
    if (!from || !to) return
    fetchData()
  },
  { deep: true }
)
</script>

<template>
  <div class="min-h-screen w-full bg-gray-50 text-slate-800 font-sans">

    <div class="w-full px-4 md:px-6 py-4 flex justify-between items-center bg-white border-b border-gray-100">
      <div class="font-bold text-xl tracking-tight text-slate-800">
        Co-Finance
      </div>
      <SearchBar />
    </div>

    <div v-if="error" class="flex h-screen items-center justify-center">
      <div class="bg-red-50 text-red-600 px-6 py-4 rounded-xl border border-red-100 text-center">
        <p class="font-bold">Connection Failed</p>
        <p class="text-sm mt-1">{{ error }}</p>
        <button @click="() => fetchData()"
          class="mt-3 px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700 text-sm">Retry</button>
      </div>
    </div>

    <div v-else-if="dashboardData" class="w-full h-full p-4 md:p-6 space-y-6 animate-fade-in-up">

      <StockHeader :quote="dashboardData.quote" :profile="dashboardData.financials"
        :marketStatus="marketStore.marketStatus" :symbol="store.symbol" />

      <div class="grid grid-cols-1 xl:grid-cols-4 gap-6">
        <div class="xl:col-span-1 flex flex-col gap-6 ">
          <KeyStats :financials="dashboardData.financials" />
          <EarningsChart v-if="dashboardData.earnings && dashboardData.earnings.length > 0"
            :earnings="dashboardData.earnings" />
        </div>

        <div :class="[
          'flex flex-col gap-6',
          'xl:relative',
          dashboardData.insiders && dashboardData.insiders.length > 0
            ? 'xl:col-span-2'
            : 'xl:col-span-3'
        ]">
          <div class="xl:absolute xl:inset-0 w-full">
            <NewsFeed :news="companyNews" />
          </div>
        </div>

        <div class="xl:col-span-1 flex flex-col gap-6">
          <RecTrends :trend="dashboardData.recommendations" />
          <div v-if="dashboardData.insiders && dashboardData.insiders.length > 0" class="h-fit">
            <InsiderTable :transactions="dashboardData.insiders" />
          </div>
        </div>

      </div>
    </div>
  </div>
</template>

<style scoped>
.animate-fade-in-up {
  animation: fadeInUp 0.5s ease-out;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
