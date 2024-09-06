<template>
<template v-if="data.length > 0">	
	<ul class="divide-y divide-gray-200 dark:divide-gray-700">
		<li v-for="item in data" :key="item.id" class="p-4">
			<div @click="toggleExpand(item.id)" class="flex items-center justify-between cursor-pointer">
				<div class="flex items-center space-x-2">
					<span :class="getMethodColor(item.request.method)" class="px-2 py-1 rounded text-xs font-medium">
						{{ item.request.method }}
					</span>
					<span class="text-sm font-mono text-gray-600 dark:text-gray-300">{{ item.request.path }}</span>
				</div>
				<div class="flex items-center space-x-4">
					<span class="text-sm text-gray-500 dark:text-gray-400">
						{{getRequestTime(item.time)}}
					</span>
					<span :class="getStatusColor(item.response.statusCode)"
						class="px-2 py-1 rounded text-xs font-medium text-white">
						{{ item.response.statusCode }}
					</span>
				</div>
			</div>
			<transition name="expand">
				<div v-if="expandedId === item.id" class="mt-4">
					<div class="border-b border-gray-200 dark:border-gray-700">
						<nav class="flex -mb-px">
							<button v-for="tab in tabs" :key="tab.id" @click="activeTab = tab.id" :class="[
								activeTab === tab.id
									? 'border-blue-500 text-blue-600 dark:text-blue-400'
									: 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 dark:text-gray-400 dark:hover:text-gray-300',
								'group inline-flex items-center py-2 px-4 border-b-2 font-medium text-sm'
							]">
								<component :is="tab.icon" class="h-4 w-4 mr-2" />
								{{ tab.name }}
							</button>
						</nav>
					</div>
					<div class="mt-4">
						<pre v-if="activeTab === 'payload'"
							class="bg-gray-100 dark:bg-gray-900 p-4 rounded-lg overflow-x-auto">
	                            <code class="text-sm"><vue-json-pretty :data="item.request" :deep=2 :theme="dark" :showLine="false"/></code>
	                          </pre>
						<pre v-if="activeTab === 'response'"
							class="bg-gray-100 dark:bg-gray-900 p-4 rounded-lg overflow-x-auto">
	                            <code class="text-sm"><vue-json-pretty :data="item.response" :deep=2 :theme='dark':showLine="false"  /></code>
	                          </pre>
						<pre v-if="activeTab === 'match'"
							class="bg-gray-100 dark:bg-gray-900 p-4 rounded-lg overflow-x-auto">
	                            <code class="text-sm"><vue-json-pretty :data="item.result" :deep=3 :theme='light' :showLine="false" /></code>
	                          </pre>
					</div>
				</div>
			</transition>
		</li>
	</ul>
</template>
<template v-else>
	<div class="flex flex-col items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
	  <InboxIcon class="h-24 w-24 text-gray-400 dark:text-gray-600 mb-4" />
	  <h3 class="text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">No requests yet</h3>
	  <p class="text-sm text-gray-500 dark:text-gray-400 text-center max-w-sm">
		Requests will appear here once they start coming in.
	  </p>
	</div>
  </template>
</template>

<script setup>
import { ref } from 'vue'
import VueJsonPretty from 'vue-json-pretty';
import 'vue-json-pretty/lib/styles.css';
import { Send, ArrowLeftRight, Target,InboxIcon } from 'lucide-vue-next'

const expandedId = ref(null)
const activeTab = ref('payload')
const tabs = [
	{ id: 'payload', name: 'Payload', icon: Send },
	{ id: 'response', name: 'Response', icon: ArrowLeftRight },
	{ id: 'match', name: 'Match', icon: Target }
]

const props = defineProps({
	data: Array, // Define the prop to receive the array
	isDark: {
		type: Boolean,
		default: false
	}
});

const toggleExpand = (id) => {
	expandedId.value = expandedId.value === id ? null : id
	activeTab.value = 'payload' // Reset to first tab when expanding
}

const getStatusColor = (statusCode) => {
	if (statusCode >= 200 && statusCode < 300) return "bg-green-500"
	if (statusCode >= 300 && statusCode < 400) return "bg-blue-500"
	if (statusCode >= 400 && statusCode < 500) return "bg-yellow-500"
	return "bg-red-500"
}

const  getRequestTime = (timestamp) => {
        var requestTime = new Date(timestamp*1000);
        var datetime = requestTime.getDate() + "/" +
            (requestTime.getMonth() + 1) + "/" +
            requestTime.getFullYear() + " @ " +
            requestTime.getHours() + ":" +
            requestTime.getMinutes() + ":" +
            requestTime.getSeconds();
        return datetime;
}

const getMethodColor = (method) => {
	switch (method) {
		case 'GET': return 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300'
		case 'POST': return 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300'
		case 'PUT': return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300'
		case 'DELETE': return 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300'
		default: return 'bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-300'
	}
}

const formatJSON = (obj) => {
	return JSON.stringify(obj, null, 2)
}

</script>

<style scoped>
.expand-enter-active,
.expand-leave-active {
	transition: all 0.3s ease-out;
	max-height: 1000px;
	/* Increased to accommodate more content */
}

.expand-enter-from,
.expand-leave-to {
	opacity: 0;
	max-height: 0;
}

/* Add syntax highlighting styles */
pre {
	@apply text-sm;
}

code {
	@apply font-mono;
}

.hljs-attr {
	@apply text-purple-600 dark:text-purple-400;
}

.hljs-string {
	@apply text-green-600 dark:text-green-400;
}

.hljs-number {
	@apply text-blue-600 dark:text-blue-400;
}

.hljs-boolean {
	@apply text-red-600 dark:text-red-400;
}

.hljs-null {
	@apply text-gray-600 dark:text-gray-400;
}
</style>