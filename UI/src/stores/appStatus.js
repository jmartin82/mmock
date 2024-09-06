import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

export const useAppStatusStore = defineStore('appStatus', () => {
  const isDark = ref(false)
  const isSidebarOpen = ref(false)

  function toggleSidebar(){
    isSidebarOpen.value = !isSidebarOpen.value
  }
  
  function toggleDarkMode(){
    isDark.value = !isDark.value
  }

  return { isDark, isSidebarOpen, toggleSidebar, toggleDarkMode }
})
