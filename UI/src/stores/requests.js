import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

export const useRequestsStore = defineStore('requests', () => {
  const requestsList = ref([])
  let nextId = 1

  function addRequest(request) {
    request.id = nextId++
    requestsList.value.unshift(request)
  }

  function clearRequests() {
    requestsList.value = []
  }

  return { requestsList,addRequest, clearRequests}
})
