<template>
  <div class="container mx-auto px-4 py-8">
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow-md overflow-hidden">
      <div class="p-4 sm:p-6">
        <div class="flex flex-col sm:flex-row justify-between items-center mb-6">
          <h2 class="text-2xl font-semibold text-gray-800 dark:text-white mb-4 sm:mb-0">Mapping</h2>
          <button @click="showAddModal = true"
            class="bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded-lg transition duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-blue-400 focus:ring-opacity-50">
            <PlusIcon class="h-5 w-5 inline-block mr-1" />
            <span>Add Mock</span>
          </button>
        </div>

        <div v-if="items.length > 0" class="overflow-x-auto">
          <table class="w-full divide-y divide-gray-200 dark:divide-gray-700">
            <thead class="bg-gray-50 dark:bg-gray-800">
              <tr>
                <th v-for="header in headers" :key="header.name"
                  class="px-3 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider dark:text-gray-400">
                  {{ header.name }}
                </th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200 dark:bg-gray-900 dark:divide-gray-700">
              <tr v-for="item in items" :key="item.URI" class="hover:bg-gray-50 dark:hover:bg-gray-800">
                <td class="px-3 py-4 text-sm font-medium text-gray-900 dark:text-white">
                  <div class="truncate max-w-[100px] sm:max-w-xs">{{ item.URI }}</div>
                </td>
                <td class="px-3 py-4 text-sm text-gray-500 dark:text-gray-400">
                  <div class="truncate max-w-[100px] sm:max-w-xs">{{ item.description }}</div>
                </td>
                <td class="px-3 py-4 text-sm">
                  <span :class="getMethodColor(item.request.method)" class="px-2 py-1 rounded text-xs font-medium">
                    {{ item.request.method }}
                  </span>
                </td>
                <td class="px-3 py-4 text-sm text-gray-500 dark:text-gray-400">
                  <div class="truncate max-w-[100px] sm:max-w-xs">{{ item.request.path }}</div>
                </td>
                <td class="px-3 py-4 text-sm">
                  <span :class="getStatusColor(item.response.statusCode)"
                    class="px-2 py-1 rounded text-xs font-medium text-white">
                    {{ item.response.statusCode }}
                  </span>
                </td>
                <td class="px-3 py-4 text-sm font-medium">
                  <div class="flex space-x-2">
                    <button @click="editItem(item.URI)"
                      class="text-indigo-600 hover:text-indigo-900 dark:text-indigo-400 dark:hover:text-indigo-200">
                      <PencilIcon class="h-5 w-5" />
                      <span class="sr-only">Edit</span>
                    </button>
                    <button @click="deleteItem(item.URI)"
                      class="text-red-600 hover:text-red-900 dark:text-red-400 dark:hover:text-red-200">
                      <TrashIcon class="h-5 w-5" />
                      <span class="sr-only">Delete</span>
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Empty state placeholder -->
        <div v-else class="text-center py-12">
          <ClipboardListIcon class="mx-auto h-12 w-12 text-gray-400 dark:text-gray-600" />
          <h3 class="mt-2 text-sm font-medium text-gray-900 dark:text-gray-100">No Mocks</h3>
          <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">Get started by creating a new mock.</p>
        </div>
      </div>
    </div>
  </div>

<!-- Add/Edit Modal -->
<div v-if="showAddModal || showEditModal" class="fixed z-10 inset-0 overflow-y-auto" aria-labelledby="modal-title" role="dialog" aria-modal="true">
      <div class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" aria-hidden="true"></div>
        <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>
        <div class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full dark:bg-gray-800">
          <form @submit.prevent="showAddModal ? addItem() : updateItem()">
            <div class="bg-white dark:bg-gray-800 px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
              <div class="mb-4">
                <label for="edit-uri" class="block text-sm font-medium text-gray-700 dark:text-gray-300">URI</label>
                <input type="text" id="edit-uri" v-model="editingItem.URI" class="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm dark:bg-gray-700 dark:border-gray-600 dark:text-white" required :readonly="showEditModal ? true : false" />
              </div>
              <div class="mb-4">
                <label for="edit-description" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Description</label>
                <input type="text" id="edit-description" v-model="editingItem.description" class="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm dark:bg-gray-700 dark:border-gray-600 dark:text-white" readonly />
              </div>
              <div class="mb-4">
                <label for="edit-mock" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Definition</label>
                <textarea id="edit-mock" v-model="editingItem.definition" rows="4" class="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm dark:bg-gray-700 dark:border-gray-600 dark:text-white" required></textarea>
              </div>
            </div>
            <div class="bg-gray-50 dark:bg-gray-700 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
              <button type="submit" class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-blue-600 text-base font-medium text-white hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 sm:ml-3 sm:w-auto sm:text-sm">
                {{ showAddModal ? 'Add' : 'Update' }}
              </button>
              <button @click="closeModal" type="button" class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm dark:bg-gray-800 dark:text-gray-300 dark:border-gray-600 dark:hover:bg-gray-700">
                Cancel
              </button>
            </div>
          </form>
        </div>
      </div>
</div>

</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ClipboardListIcon, PlusIcon, PencilIcon, TrashIcon } from 'lucide-vue-next'
import { getMockDefinitions, getMockDefinition, createMockDefinition, updateMockDefinition, deleteMockDefinition } from '../services/api';

const headers = [
  { name: 'URI'},
  { name: 'Description'},
  { name: 'Method'},
  { name: 'Path'},
  { name: 'Result'},
  { name: 'Actions' }
];

const showAddModal = ref(false)
const showEditModal = ref(false)
const editingItem = ref({})
const items = ref([]);

const loadMockDefinitions = async () => {
  try {
    items.value = await getMockDefinitions();
  } catch (error) {
    console.error('Failed to load mock definitions:', error);
  }
};



const addItem = async () => {

  try {
    await createMockDefinition(editingItem.value.URI, JSON.parse(editingItem.value.definition));
    editingItem.value = null;
    loadMockDefinitions();
  } catch (error) {
    console.error('Failed to create mock definition:', error);
  }

  closeModal()
}

const editItem = async (URI) => {
  try {
    let mock = await getMockDefinition(URI);
    editingItem.value.URI = mock.URI
    editingItem.value.description = mock.description
    editingItem.value.definition = JSON.stringify(mock)

    showEditModal.value = true
  } catch (error) {
    console.error('Failed to fetch mock definition:', error);
  }
}

const updateItem = async () => {
  try {
    await updateMockDefinition(editingItem.value.URI, JSON.parse(editingItem.value.definition));
    loadMockDefinitions();
    editingItem.value = null;
  } catch (error) {
    console.error('Failed to update mock definition:', error);
  }
  closeModal()
}

const deleteItem = async (URI) => {
  try {
    await deleteMockDefinition(URI);
    loadMockDefinitions();
  } catch (error) {
    console.error('Failed to delete mock definition:', error);
  }
}

const closeModal = () => {
  showAddModal.value = false
  showEditModal.value = false
  editingItem.value = { uri: '', description: '', definition: '' }
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

const getStatusColor = (statusCode) => {
  if (statusCode >= 200 && statusCode < 300) return "bg-green-500"
  if (statusCode >= 300 && statusCode < 400) return "bg-blue-500"
  if (statusCode >= 400 && statusCode < 500) return "bg-yellow-500"
  return "bg-red-500"
}

onMounted(() => {
  console.log('Loading');
  loadMockDefinitions();
})


</script>