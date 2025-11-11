import { defineStore } from 'pinia'

export const useApp = defineStore('app', {
  state: () => ({
    ready: false,
  }),
  actions: {
    async bootstrap () {
      try {
        // 这里放全局初始化逻辑（可留空）
      } finally {
        this.ready = true
      }
    }
  }
})
