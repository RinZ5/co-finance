import { defineStore } from 'pinia';
import { computed, ref, watch } from 'vue';
import api from '../services/api';
import socket from '../services/socket';
import type { DashboardData, CompanyNews, MarketStatus } from '../types/types';

export const useDashboardStore = defineStore('dashboard', () => {
  const symbol = ref(localStorage.getItem('stock_symbol') || 'AAPL');

  const today = new Date();
  const fiveDaysAgo = new Date();
  fiveDaysAgo.setDate(today.getDate() - 5);
  const newsRange = ref<[Date, Date]>([fiveDaysAgo, today])

  const dashboardData = ref<DashboardData | null>(null)
  const companyNews = ref<CompanyNews[]>([])
  const marketStatus = ref<MarketStatus | undefined>(undefined);

  const isLoading = ref(false);
  const error = ref<string | null>(null);
  const isSocketConnected = ref(false)

  let timer: number | null = null;

  const apiFrom = computed(() => newsRange.value?.[0]?.toISOString().split('T')[0]);
  const apiTo = computed(() => newsRange.value?.[1]?.toISOString().split('T')[0]);

  function isMarketOpen(now = new Date()) {
    const h = now.getHours();
    const m = now.getMinutes();
    return (h > 9 || (h === 9 && m >= 30)) && h < 16;
  }

  function isLessThan930(date: Date) {
    return (
      date.getHours() < 9 ||
      (date.getHours() === 9 && date.getMinutes() < 30)
    );
  }

  function getNextMarketTick(now = new Date()) {
    const next = new Date(now);

    if (isLessThan930(now)) {
      next.setHours(9, 30, 0, 0);
      return next;
    }

    if (now.getHours() >= 16) {
      next.setDate(next.getDate() + 1);
      next.setHours(9, 30, 0, 0);
      return next;
    }

    if (now.getMinutes() < 30) {
      next.setMinutes(30, 0, 0);
    } else {
      next.setHours(next.getHours() + 1, 30, 0, 0);
    }

    return next;
  }

  function startMarketSync(exchange = 'US') {
    stopMarketSync();

    const run = async () => {
      if (isMarketOpen()) {
        try {
          const res = await api.get('/market-status', { params: { exchange } });
          marketStatus.value = res.data;
        } catch (e) {
          console.warn("Market sync error", e);
        }
      }
      const now = new Date();
      const next = getNextMarketTick(now);
      const delay = next.getTime() - now.getTime();

      timer = globalThis.setTimeout(run, delay);
    };

    run();
  }

  function stopMarketSync() {
    if (timer) {
      clearTimeout(timer);
      timer = null;
    }
  }

  function initSocketListeners() {
    socket.onMessage((data: any) => {
      if (data.type === 'trade' && data.p) {
        if (dashboardData.value?.quote) {
          dashboardData.value.quote.c = data.p;
        }
      }
    });
  }

  initSocketListeners();

  async function loadDashboard() {
    isLoading.value = true;
    error.value = null;

    socket.connect(symbol.value);
    isSocketConnected.value = true;

    startMarketSync();

    try {
      const [dashResult, newsResult, marketResult] = await Promise.allSettled([
        api.getDashboard(symbol.value),
        (apiFrom.value && apiTo.value)
          ? api.getCompanyNews(symbol.value, apiFrom.value, apiTo.value)
          : Promise.resolve({ data: [] }),
        api.getMarketStatus('US')
      ]);

      if (dashResult.status === 'fulfilled') {
        dashboardData.value = dashResult.value.data;
      } else {
        throw new Error("Failed to load dashboard data");
      }

      if (newsResult.status === 'fulfilled') {
        companyNews.value = newsResult.value.data;
      }

      if (marketResult.status === 'fulfilled') {
        marketStatus.value = marketResult.value.data;
      }

    } catch (err: any) {
      console.error("Store Error:", err);
      error.value = err.message || 'Failed to fetch data';
    } finally {
      isLoading.value = false;
    }
  }

  function setSymbol(newSymbol: string) {
    if (newSymbol) {
      symbol.value = newSymbol.toUpperCase();
      loadDashboard()
    }
  }

  function setDateRange(from: Date, to: Date) {
    newsRange.value = [from, to]
    loadDashboard()
  }

  function cleanup() {
    socket.disconnect();
    isSocketConnected.value = false;
  }

  watch(symbol, (symbol) => {
    localStorage.setItem('stock_symbol', symbol);
  });

  return { symbol, newsRange, dashboardData, companyNews, marketStatus, isLoading, error, isSocketConnected, setSymbol, setDateRange, loadDashboard, cleanup };
});
