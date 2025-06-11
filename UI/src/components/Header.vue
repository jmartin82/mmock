<template>
	<header class="bg-white dark:bg-gray-800 shadow-md">
	        <div class="container mx-auto px-4 py-3 flex items-center justify-between">
	          <button @click="appStatus.toggleSidebar" class="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200">
	            <Menu class="h-6 w-6" />
	            <span class="sr-only">Toggle menu</span>
	          </button>			  
			  <div class="flex-1 text-center">
				<h1 class="text-lg font-extrabold text-gray-800 dark:text-white">MMock Console</h1>
			  </div>
			<h4 class="text-sm font-medium text-gray-600 dark:text-gray-400">{{ $route.name }}</h4>
			</div>
	      </header>
</template>

<script setup>
import { Moon, Sun, Menu} from 'lucide-vue-next'
import { useAppStatusStore } from '../stores/appStatus';
const appStatus = useAppStatusStore()
</script>