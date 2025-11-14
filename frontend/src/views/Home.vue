<script setup lang="ts">
import { useRouter } from 'vue-router'
import { Card, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Smile, ArrowRight } from 'lucide-vue-next'

const router = useRouter()

const navigateTo = (route: string) => {
  router.push(route)
}
</script>

<template>
  <div class="h-screen flex flex-col bg-gradient-to-br from-gray-50 via-white to-gray-50">
    <!-- Header Section -->
    <header class="relative overflow-hidden bg-gradient-to-r from-blue-600 via-blue-500 to-purple-600 text-white flex-shrink-0">
      <div class="absolute inset-0 bg-grid-white/10"></div>
      <div class="relative px-8 py-8">
        <div class="text-center">
          <h1 class="text-4xl font-bold tracking-tight animate-fade-in">
            SmileOps
          </h1>
        </div>
      </div>
      <div class="absolute bottom-0 left-0 right-0 h-8 bg-gradient-to-t from-white/20 to-transparent"></div>
    </header>

    <main class="flex-1 flex items-center justify-center px-8 py-8">
      <!-- Main Feature -->
      <div class="w-full max-w-lg">
        <Card
          class="group cursor-pointer hover:shadow-xl transition-all duration-300 border-none shadow-md overflow-hidden relative"
          @click="navigateTo('/implants')"
        >
          <!-- Gradient bar on top -->
          <div
            class="absolute top-0 left-0 right-0 h-1 bg-gradient-to-r from-blue-500 to-blue-600 transition-all duration-300"
          ></div>

          <CardHeader class="p-8">
            <div class="flex items-center justify-between mb-4">
              <div class="bg-blue-50 p-4 rounded-xl group-hover:scale-110 transition-transform duration-300">
                <Smile class="text-blue-600 h-10 w-10" />
              </div>
              <ArrowRight
                class="h-6 w-6 text-gray-400 group-hover:translate-x-1 group-hover:text-gray-900 transition-all duration-300"
              />
            </div>
            <CardTitle class="text-2xl mb-3 group-hover:text-blue-600 transition-colors">
              Implants
            </CardTitle>
            <CardDescription class="text-base leading-relaxed">
              Manage dental implants and track procedures
            </CardDescription>
          </CardHeader>
        </Card>
      </div>
    </main>
  </div>
</template>

<style scoped>
@keyframes fade-in {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate-fade-in {
  animation: fade-in 0.6s ease-out;
}

.animate-fade-in-delay {
  animation: fade-in 0.6s ease-out 0.2s both;
}

.bg-grid-white\/10 {
  background-image:
    linear-gradient(to right, rgba(255, 255, 255, 0.1) 1px, transparent 1px),
    linear-gradient(to bottom, rgba(255, 255, 255, 0.1) 1px, transparent 1px);
  background-size: 40px 40px;
}
</style>
