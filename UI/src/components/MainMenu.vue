<template>
	<aside
	          :class="[
	            'bg-white dark:bg-gray-800 w-64 shadow-md transition-all duration-300 ease-in-out',
	            appStatus.isSidebarOpen ? 'translate-x-0' : '-translate-x-full'
	          ]"
	          class="fixed inset-y-0 left-0 z-30"
	        >
	          <div class="flex justify-between items-center p-4 border-b border-gray-200 dark:border-gray-700">
	            <h2 class="text-lg font-semibold text-gray-800 dark:text-white">Menu</h2>
	            <button
	              @click="appStatus.toggleSidebar"
	              class="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
	              aria-label="Close menu"
	            >
	              <X class="h-6 w-6" />
	            </button>
	          </div>
	          <nav class="p-4">
	            <ul class="space-y-2">
	              <li>
	                <a href="/" class="flex items-center space-x-2 text-gray-700 hover:text-gray-900 dark:text-gray-300 dark:hover:text-white">
	                  <LayoutDashboard class="h-5 w-5" />
	                  <span>Request Monitor</span>
	                </a>
	              </li>
	              <!--li>
	                <a href="/mapping" class="flex items-center space-x-2 text-gray-700 hover:text-gray-900 dark:text-gray-300 dark:hover:text-white">
	                  <List class="h-5 w-5" />
	                  <span>Mapping</span>
	                </a>
	              </li-->
	              <li>
	                <a href="/about" class="flex items-center space-x-2 text-gray-700 hover:text-gray-900 dark:text-gray-300 dark:hover:text-white">
	                  <Settings class="h-5 w-5" />
	                  <span>About</span>
	                </a>
	              </li>
	            </ul>
	          </nav>
	        </aside>
</template>

<script setup>
import { LayoutDashboard, List, Settings, X } from 'lucide-vue-next'
import { useAppStatusStore } from '../stores/appStatus';
const appStatus = useAppStatusStore()
</script>
