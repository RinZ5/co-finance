import { defineStore } from 'pinia';
import { ref, watch } from 'vue';

export const useStockStore = defineStore('stock', () => {
  const symbol = ref(localStorage.getItem('stock_symbol') || 'AAPL');


  const today = new Date();
  const threeDaysAgo = new Date();
  threeDaysAgo.setDate(today.getDate() - 3);

  const newsRange = ref<[Date, Date]>([threeDaysAgo, today])

  function setSymbol(newSymbol: string) {
    if (newSymbol) {
      symbol.value = newSymbol.toUpperCase();
    }
  }

  function setDateRange(from: Date, to: Date) {
    newsRange.value = [from, to]
  }

  watch(symbol, (symbol) => {
    localStorage.setItem('stock_symbol', symbol);
  });

  return { symbol, newsRange, setSymbol, setDateRange };
});
