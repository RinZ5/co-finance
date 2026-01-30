<script setup lang="ts">
import { onMounted, watch, onUnmounted } from 'vue';
import NProgress from 'nprogress';
import 'nprogress/nprogress.css';

import SearchBar from './components/SearchBar.vue';
import StockHeader from './components/StockHeader.vue';
import KeyStats from './components/KeyStats.vue';
import EarningsChart from './components/EarningsChart.vue';
import RecTrends from './components/RecTrends.vue';
import InsiderTable from './components/InsiderTable.vue';
import NewsFeed from './components/NewsFeed.vue';
import { useDashboardStore } from './stores/dashboard';

NProgress.configure({ showSpinner: false })

const store = useDashboardStore();

watch(
  () => store.isLoading,
  (isLoading) => {
    if (isLoading) {
      NProgress.start();
    } else {
      NProgress.done();
    }
  }
);

onMounted(() => {
  store.loadDashboard();
})

onUnmounted(() => {
  store.cleanup();
})
</script>

<template>
  <div class="min-h-screen w-full bg-gray-50 text-slate-800 font-sans">

    <div class="w-full px-4 md:px-6 py-4 flex items-center justify-between bg-surface border-b border-border-light">
      <div class="font-bold text-xl shrink-0">
        Co-Finance
      </div>
      <div class="w-full md:w-80 md:block sm:pl-6">
        <SearchBar />
      </div>
    </div>

    <div v-if="store.error" class="flex h-screen items-center justify-center">
      <div class="bg-red-50 text-red-600 px-6 py-4 rounded-xl border border-red-100 text-center">
        <p class="font-bold">Connection Failed</p>
        <p class="text-sm mt-1">{{ store.error }}</p>
        <button @click="store.loadDashboard()"
          class="mt-3 px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700 text-sm">Retry</button>
      </div>
    </div>

    <div v-else-if="store.dashboardData" class="w-full h-full p-4 md:p-6 space-y-6 animate-fade-in-up">

      <StockHeader :quote="store.dashboardData.quote" :profile="store.dashboardData.financials"
        :marketStatus="store.marketStatus" :symbol="store.symbol" />

      <div class="grid grid-cols-1 xl:grid-cols-4 gap-6">
        <div class="xl:col-span-1 flex flex-col gap-6 ">
          <KeyStats :financials="store.dashboardData.financials" />
          <EarningsChart v-if="store.dashboardData.earnings && store.dashboardData.earnings.length > 0"
            :earnings="store.dashboardData.earnings" />
        </div>

        <div :class="[
          'flex flex-col gap-6',
          'xl:relative',
          store.dashboardData.insiders && store.dashboardData.insiders.length > 0
            ? 'xl:col-span-2'
            : 'xl:col-span-3'
        ]">
          <div class="xl:absolute xl:inset-0 w-full">
            <NewsFeed :news="store.companyNews" v-model="store.newsRange" />
          </div>
        </div>

        <div class="xl:col-span-1 flex flex-col gap-6">
          <RecTrends :trend="store.dashboardData.recommendations" />
          <div v-if="store.dashboardData.insiders && store.dashboardData.insiders.length > 0" class="h-fit">
            <InsiderTable :transactions="store.dashboardData.insiders" />
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
