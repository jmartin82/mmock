<template>
<div class="container mx-auto">
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow-md overflow-hidden">
      <div class="max-h-[calc(100vh-12rem)] overflow-y-auto">
        <RequestList :data="requestStore.requestsList" :is-dark="true"/>
      </div>
    </div>
</div>
</template>
<script setup>

import RequestList from '../components/RequestList.vue'
import { useRequestsStore} from '../stores/requests';
const requestStore = useRequestsStore()
</script>