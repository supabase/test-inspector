<template>
  <div class="w-full my-12">
    <client-only>
      <div class="flex justify-between">
        <h1 class="mb-24 mt-12 text-6xl font-bold u-text-black">
          Projects.
        </h1>
      </div>
      <div v-if="pending" class="w-full flex justify-center flex-wrap">
        <LoadingCard
          :loading="pending"
          :nb-of-items="2"
          :hght="120"
          class="mt-2 mx-12 my-12 w-5/12 min-w-[480px]"
        />
      </div>
      <transition
        v-else
        name="bounce"
        mode="out-in"
        appear
      >
        <ul class="px-12 flex justify-center flex-wrap">
          <li
            v-for="p of projects"
            :key="p.id"
            :class="{ 'border-b': p.id !== projects.length - 1 }"
            class="divide-y divide-gray-200 mx-12 my-12 w-5/12 min-w-[480px]"
          >
            <div class="py-2">
              <UCard
                label-class="block font-medium u-text-gray-700"
                wrapper-class="flex items-center justify-between w-full"
                :label="p.name"
                :name="p.id"
              >
                <div class="flex items-center justify-between u-text-gray-700 mb-5">
                  <h5 class="font-bold text-2xl tracking-tight mb-2 flex justify-between w-full">
                    <div>
                      <NuxtLink :to="`/projects/${p.id}`" class="text-emerald-500 hover:text-emerald-600">
                        {{ p.name }}
                      </NuxtLink>
                    </div>
                    <div>
                      <span class="u-text-gray-700 text-sm">
                        #{{ p.id }}
                      </span>
                    </div>
                  </h5>
                </div>
                <p class="font-sm u-text-gray-700 mb-3 min-h-full">
                  created: {{ (new Date(Date.parse(p.created_at))).toLocaleString() }}
                </p>
                <p class="font-normal text-sm u-text-gray-700 mb-3 inline-flex flex-row justify-between w-full">
                  <span class="inline-flex">
                    <svg
                      class="mr-1"
                      xmlns="http://www.w3.org/2000/svg"
                      width="16"
                      height="16"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="#3ecf8e"
                      stroke-width="2"
                      stroke-linecap="round"
                      stroke-linejoin="bevel"
                    >
                      <polyline points="9 11 12 14 22 4" />
                      <path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11" />
                    </svg>
                    versions count: {{ p.versions?.length ?? 0 }}
                  </span>
                </p>
              </UCard>
            </div>
          </li>
        </ul>
      </transition>
    </client-only>
  </div>
</template>

<script setup lang="ts">
import { Project } from '~/types/projects'

const client = useSupabaseClient()

const show = ref(null)
show.value = true

const { pending, data: projects } = useLazyAsyncData('projects', async () => {
  const { data } = await client.from<Project>('projects')
    .select('*, versions(id, version_name)')
    .order('created_at', { ascending: true })
  return data
})
</script>

<style scoped>
.bounce-enter-active {
  animation: bounce-in 0.5s;
}
.bounce-leave-active {
  animation: bounce-in 0.5s reverse;
}
@keyframes bounce-in {
  0% {
    transform: scale(0);
  }
  50% {
    transform: scale(1.25);
  }
  100% {
    transform: scale(1);
  }
}
</style>
