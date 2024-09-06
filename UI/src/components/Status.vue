<template>
	<footer class="bg-white dark:bg-gray-800 shadow-md">
	        <div class="container mx-auto px-4 py-2 flex items-center justify-between text-sm text-gray-600 dark:text-gray-400">
	          <span>Total Requests: {{ requestStore.requestsList.length }}</span>
	          <div class="flex items-center space-x-2">
	            <span>Status:</span>
	            <span v-if="connectionStatus" class="flex items-center">
	              <span class="h-2 w-2 bg-green-500 rounded-full mr-1"></span>
	              Connected
	            </span>
				<span v-else class="flex items-center">
	              <span class="h-2 w-2 bg-red-500 rounded-full mr-1"></span>
	              Disconnected
	            </span>

	          </div>
	          <span>Last request: {{ lastUpdated }}</span>
	        </div>
	      </footer>
</template>

<script setup>
import { useRequestsStore} from '../stores/requests';
const requestStore = useRequestsStore()

defineProps({
  lastUpdated: {
    type: String,
    default: '-'
  },
  connectionStatus: {
    type: Boolean,
    default: false
  }
})
</script>
