<template>
  <div class="p-4 sm:p-6 max-w-2xl mx-auto">
    <div class="rounded-lg overflow-hidden shadow-lg">
      <div class="bg-gradient-to-r from-blue-500 to-indigo-600 px-4 py-3 sm:px-6 sm:py-4">
        <h1 class="text-2xl font-extrabold text-white">MMock Web Console</h1>
      </div>
      <div class="bg-white dark:bg-gray-800 text-gray-500 dark:text-gray-400 px-4 py-4 sm:px-6 sm:py-5 space-y-4">
        <section>
          <h2 class="text-lg font-semibold mb-1">Description</h2>
          <p class="text-sm">
            Mmock is an HTTP mocking application for testing and fast prototyping.
          </p>
        </section>

        <section>
          <h2 class="text-lg font-semibold mb-1">Author</h2>
          <div class="flex items-center space-x-3">
            <UserIcon class="h-8 w-8" />
            <div>
              <p class="font-medium">Jordi Martin</p>
              <a href="https://jordi.io" class="text-sm hover:underline flex items-center">
                <Link2Icon class="h-4 w-4 mr-1" />
                https://jordi.io
              </a>
            </div>
          </div>
        </section>

        <section>
          <h2 class="text-lg font-semibold mb-1">GitHub Repository</h2>
          <a 
            href="https://github.com/jmartin82/mmock" 
            target="_blank" 
            rel="noopener noreferrer"
            class="text-sm hover:underline flex items-center"
          >
            <GithubIcon class="h-4 w-4 mr-1" />
            https://github.com/jmartin82/mmock
          </a>
        </section>

        <div class="flex justify-between items-center">
          <section>
            <h2 class="text-lg font-semibold mb-1">License</h2>
            <div class="flex items-center space-x-1 text-sm">
              <FileTextIcon class="h-4 w-4" />
              <span>MIT License</span>
            </div>
          </section>

          <section>
            <h2 class="text-lg font-semibold mb-1">Version</h2>
            <p class="text-sm">5</p>
          </section>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { UserIcon, GithubIcon, FileTextIcon, Link2Icon } from 'lucide-vue-next'
</script>