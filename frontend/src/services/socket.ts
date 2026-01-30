const BACKEND_URL = import.meta.env.VITE_BACKEND_URL || globalThis.location.origin;

const WS_BASE_URL = BACKEND_URL.replace(/^http(s)?/, 'ws$1');

type Listener = (data: any) => void;

class WebSocketService {
  private socket: WebSocket | null = null;
  private listeners: Listener[] = [];

  connect(symbol: string) {
    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      return;
    }

    this.socket = new WebSocket(`${WS_BASE_URL}/ws`);

    this.socket.onopen = () => {
      console.log("WebSocket Connected");
      this.subscribe(symbol);
    };

    this.socket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        this.listeners.forEach((callback) => callback(data));
      } catch (e) {
        console.error("Error parsing WS message:", e);
      }
    };

    this.socket.onclose = () => console.log("WebSocket Disconnected");
    this.socket.onerror = (err) => console.error("WebSocket Error:", err);
  }

  subscribe(symbol: string) {
    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      this.socket.send(JSON.stringify({ type: "subscribe", symbol }));
    }
  }

  onMessage(callback: Listener) {
    this.listeners.push(callback);
  }

  disconnect() {
    if (this.socket) {
      this.socket.close();
      this.socket = null;
      this.listeners = [];
    }
  }
}

export default new WebSocketService();
