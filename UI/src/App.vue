<template>
  <div :class="{ 'dark': appStatus.isDark }" class="h-screen">
    <div class="flex flex-col h-full bg-gray-100 dark:bg-gray-900">
      <!-- Header -->
      <Header />

      <!-- Main content -->
      <div class="flex-1 flex overflow-hidden">
        <!-- Sidebar -->
        <MainMenu />

        <!-- Main content area -->
        <main class="flex-1 overflow-y-auto p-4">
          <router-view></router-view>
        </main>
      </div>

      <!-- Bottom status bar -->
      <Status :lastUpdated="lastUpdated" :connectionStatus="connectionStatus" />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { data } from 'autoprefixer';
import Header from './components/Header.vue'
import MainMenu from './components/MainMenu.vue'
import Status from './components/Status.vue'

// Access appStatus store
import { useAppStatusStore } from './stores/appStatus';
import { useRequestsStore } from './stores/requests';

const appStatus = useAppStatusStore()
const requestStore = useRequestsStore()

const lastUpdated = ref(new Date().toLocaleTimeString())
const connectionStatus = ref(false)


const updateTitle = () => {
  document.title = "NEW REQUEST!!";
  setTimeout(() => {
    document.title = "MMock Console";
  }, 2000);
}


const WSConnect = () => {
  var wsProtocol = window.location.protocol === 'https:' ? 'wss://' : 'ws://';
  var socket = new WebSocket(wsProtocol + location.host + "/echo");
  
  socket.onopen = function () {
    connectionStatus.value = true
  }

  socket.onmessage = function (event) {
    var message = JSON.parse(event.data);
    requestStore.addRequest(message);
    lastUpdated.value = new Date().toLocaleTimeString()
    updateTitle()
  };

  socket.onerror = function () {
    connectionStatus.value = false
  }

  socket.onclose = function () {
    connectionStatus.value = false
  }

}

onMounted(() => {
  WSConnect();
  // Check for saved dark mode preference
  const savedDarkMode = localStorage.getItem('darkMode')
  if (savedDarkMode !== null) {
    appStatus.isDark = savedDarkMode === 'true'
  } else {
    // If no saved preference, check system preference
    appStatus.isDark = window.matchMedia('(prefers-color-scheme: dark)').matches
  }
})


// Watch for system color scheme changes
watch(
  () => window.matchMedia('(prefers-color-scheme: dark)').matches,
  (isDark) => {
    if (localStorage.getItem('darkMode') === null) {
      isDarkMode.value = isDark
    }
  }
)
</script>
