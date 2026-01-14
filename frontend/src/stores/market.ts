import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { MarketStatus } from '../types/types'

export const useMarketStore = defineStore('market', () => {
  const marketStatus = ref<MarketStatus | undefined>(undefined)

  let timer: number | null = null

  function isMarketOpen(now = new Date()) {
    const h = now.getHours()
    const m = now.getMinutes()

    return (
      (h > 9 || (h === 9 && m >= 30)) &&
      h < 16
    )
  }

  function getNextMarketTick(now = new Date()) {
    const next = new Date(now)

    if (hLessThan930(now)) {
      next.setHours(9, 30, 0, 0)
      return next
    }

    if (now.getHours() >= 16) {
      next.setDate(next.getDate() + 1)
      next.setHours(9, 30, 0, 0)
      return next
    }

    if (now.getMinutes() < 30) {
      next.setMinutes(30, 0, 0)
    } else {
      next.setHours(next.getHours() + 1, 30, 0, 0)
    }

    return next
  }

  function hLessThan930(date: Date) {
    return (
      date.getHours() < 9 ||
      (date.getHours() === 9 && date.getMinutes() < 30)
    )
  }

  async function fetchMarketStatus(exchange = 'US') {
    const res = await fetch(`/api/market-status?exchange=${exchange}`)
    if (!res.ok) {
      throw new Error('Failed to fetch market status')
    }
    marketStatus.value = await res.json()
  }

  function startMarketSync(exchange = 'US') {
    stopMarketSync()

    const run = async () => {
      if (!isMarketOpen()) return

      await fetchMarketStatus(exchange)

      const now = new Date()
      const next = getNextMarketTick(now)
      const delay = next.getTime() - now.getTime()

      timer = window.setTimeout(run, delay)
    }

    run()
  }

  function stopMarketSync() {
    if (timer) {
      clearTimeout(timer)
      timer = null
    }
  }

  return {
    marketStatus,
    fetchMarketStatus,
    startMarketSync,
    stopMarketSync,
  }
})

