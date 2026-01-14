import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useWebSocketStore = defineStore('websocket', () => {
  const socket = ref<WebSocket | null>(null);
  const isConnected = ref(false);

  const savedSymbol = localStorage.getItem("stock_symbol");
  const currentSymbol = ref<string>(savedSymbol || "");
  const currentPrice = ref<number | null>(null);
  const lastUpdated = ref<string>("");

  function connect() {
    if (socket.value?.readyState === WebSocket.OPEN) return;

    socket.value = new WebSocket(`ws://${location.host}/ws`);

    socket.value.onopen = () => {
      console.log(" WebSocket Connected");
      isConnected.value = true;

      if (currentSymbol.value) {
        subscribe(currentSymbol.value);
      }
    };

    socket.value.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data);

        if (msg.type === 'trade' && msg.data && msg.data.length > 0) {
          const trade = msg.data[msg.data.length - 1];

          if (trade.s === currentSymbol.value) {
            currentPrice.value = trade.p;
            lastUpdated.value = new Date(trade.t).toLocaleTimeString();
          }
        }
      } catch (err) {
        console.error("Failed to parse WS message", err);
      }
    };

    socket.value.onclose = () => {
      console.log("WebSocket Disconnected");
      isConnected.value = false;
    };
  }

  function subscribe(symbol: string) {
    if (!socket.value || socket.value.readyState !== WebSocket.OPEN) return;

    currentSymbol.value = symbol;
    currentPrice.value = null;

    localStorage.setItem("stock_symbol", symbol)

    console.log(`Subscribing to ${symbol}...`);
    socket.value.send(JSON.stringify({
      type: "subscribe",
      symbol: symbol
    }));
  }

  return {
    connect,
    subscribe,
    isConnected,
    currentSymbol,
    currentPrice,
    lastUpdated
  };
});
