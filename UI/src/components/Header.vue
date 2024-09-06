<template>
	<header class="bg-white dark:bg-gray-800 shadow-md">
	        <div class="container mx-auto px-4 py-3 flex items-center justify-between">
	          <button @click="appStatus.toggleSidebar" class="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200">
	            <Menu class="h-6 w-6" />
	            <span class="sr-only">Toggle menu</span>
	          </button>
	          <h1 class="text-xl font-semibold text-gray-800 dark:text-white">{{ $route.name }}</h1>
	          <button
	            @click="appStatus.toggleDarkMode"
	            class="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
	            aria-label="Toggle dark mode"
	          >
	            <Sun v-if="appStatus.isDark" class="h-6 w-6" />
	            <Moon v-else class="h-6 w-6" />
	          </button>
	        </div>
	      </header>
</template>

<script setup>
import { Moon, Sun, Menu} from 'lucide-vue-next'
import { useAppStatusStore } from '../stores/appStatus';
const appStatus = useAppStatusStore()
</script>