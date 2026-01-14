import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { MarketStatus } from '../types/types'

export const useMarketStore = defineStore('market', () => {
  const marketStatus = ref<MarketStatus | undefined>(undefined)
  let hourlyTimer: number | null = null

  async function fetchMarketStatus(exchange = 'US') {
    const res = await fetch(`/api/market-status?exchange=${exchange}`)
    if (!res.ok) {
      throw new Error('Failed to fetch market status')
    }
    marketStatus.value = await res.json()
  }

  function startHourlySync(exchange = 'US') {
    if (hourlyTimer) {
      clearTimeout(hourlyTimer)
      hourlyTimer = null
    }

    const now = new Date()

    const msUntilNextHour =
      (60 - now.getMinutes()) * 60 * 1000 -
      now.getSeconds() * 1000 -
      now.getMilliseconds()

    hourlyTimer = window.setTimeout(async () => {
      await fetchMarketStatus(exchange)

      hourlyTimer = window.setInterval(() => {
        fetchMarketStatus(exchange)
      }, 60 * 60 * 1000)
    }, msUntilNextHour)
  }

  return {
    marketStatus,
    fetchMarketStatus,
    startHourlySync,
  }
})

