import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

export const useAppStatusStore = defineStore('appStatus', () => {
  const isSidebarOpen = ref(false)

  function toggleSidebar(){
    isSidebarOpen.value = !isSidebarOpen.value
  }
  
 

  return { isSidebarOpen, toggleSidebar }
})
