<template>
  <div class="w-full my-12">
    <div class="flex justify-between">
      <h1 class="mb-24 mt-12 text-6xl font-bold u-text-black">Versions.</h1>
      <div>
        <nuxt-link to="/">
          <UButton
            icon="heroicons-solid:arrow-circle-left"
            class="u-text-black"
            size="xl"
            variant="transparent"
          />
        </nuxt-link>
      </div>
    </div>
    <client-only>
      <div v-if="pending" class="w-full flex justify-center flex-wrap">
        <LoadingCard
          :loading="pending"
          :nb-of-items="1"
          :hght="450"
          class="mx-12 my-12 w-5/12 h-96 min-w-[480px]"
        />
      </div>
      <transition v-else name="bounce" mode="out-in" appear>
        <ul class="px-12 flex justify-center flex-wrap">
          <li
            v-for="v of versions"
            :key="v.id"
            :class="{ 'border-b': v.id !== versions.length - 1 }"
            class="divide-y divide-gray-200 mx-12 my-12 w-5/12 min-h-96 min-w-[480px]"
          >
            <div class="py-2 min-w-full">
              <UCard
                label-class="block font-medium u-text-gray-700"
                wrapper-class="flex items-center justify-between w-full"
                class=""
                :label="v.version_name"
                :name="v.id"
              >
                <div
                  :class="`flex items-center justify-between ${
                    v.is_template_launch
                      ? 'underline text-emerald-500'
                      : 'u-text-gray-700'
                  } mb-5`"
                >
                  <h5
                    class="font-bold text-2xl tracking-tight mb-2 flex justify-between w-full"
                  >
                    <div>
                      <span>
                        {{ v.version_name }}
                      </span>
                    </div>
                    <div>
                      <span class="u-text-gray-700 text-sm"> #{{ v.id }} </span>
                    </div>
                  </h5>
                </div>
                <NuxtLink
                  :to="`/projects/${$route.params.id}/versions/${v.id}`"
                  class="text-emerald-500 hover:text-emerald-600"
                >
                  last launch:
                  {{ new Date(Date.parse(v.last_launch_at)).toLocaleString() }}
                </NuxtLink>
                <p
                  class="font-normal text-sm u-text-gray-700 mb-3 inline-flex flex-row justify-between w-full"
                >
                  <span class="inline-flex">
                    <svg
                      class="mr-1"
                      xmlns="http://www.w3.org/2000/svg"
                      width="16"
                      height="16"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="#88b04b"
                      stroke-width="2"
                      stroke-linecap="round"
                      stroke-linejoin="bevel"
                    >
                      <polyline points="9 11 12 14 22 4" />
                      <path
                        d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"
                      />
                    </svg>
                    passed: {{ v.passed_number ?? 0 }}
                  </span>
                  <span class="inline-flex">
                    <svg
                      class="mr-1"
                      xmlns="http://www.w3.org/2000/svg"
                      width="16"
                      height="16"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="#ff6f61"
                      stroke-width="2"
                      stroke-linecap="round"
                      stroke-linejoin="bevel"
                    >
                      <path d="M3 3h18v18H3zM15 9l-6 6m0-6l6 6" />
                    </svg>
                    failed: {{ v.failed_number ?? 0 }}
                  </span>
                  <span class="inline-flex">
                    <svg
                      class="mr-1"
                      xmlns="http://www.w3.org/2000/svg"
                      width="16"
                      height="16"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="#f5df4d"
                      stroke-width="2"
                      stroke-linecap="round"
                      stroke-linejoin="bevel"
                    >
                      <path d="M3 3h18v18H3zM8 12h8" />
                    </svg>
                    skipped: {{ v.skipped_number ?? 0 }}
                  </span>
                  <span class="inline-flex">
                    <svg
                      class="mr-1"
                      xmlns="http://www.w3.org/2000/svg"
                      width="16"
                      height="16"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="#0f4c81"
                      stroke-width="2"
                      stroke-linecap="round"
                      stroke-linejoin="bevel"
                    >
                      <circle cx="12" cy="12" r="10" />
                      <line x1="4.93" y1="4.93" x2="19.07" y2="19.07" />
                    </svg>
                    undefined: {{ v.undefined_number ?? 0 }}
                  </span>
                </p>
                <div class="w-full">
                  <ResultsChart
                    :features="features"
                    :version="v.id"
                    :templates="v.is_template_launch ? [] : templates"
                  />
                </div>
              </UCard>
            </div>
          </li>
        </ul>
      </transition>
    </client-only>
  </div>
</template>

<script setup lang="ts">
import { Version } from '~/types/versions'
import { Result } from '~/types/results'

const client = useSupabaseClient()

const { pending, data: versions } = await useAsyncData('versions', async () => {
  const { data } = await client.rpc<Version>('version')
  return data
})

const { data: features } = await useAsyncData('features', async () => {
  const { data } = await client.rpc<string>('features')
  return data
})

const { data: templates } = await client.rpc<Result>('results', {
  version: versions.value
    .filter((v) => v.is_template_launch)
    .map((v) => v.id)[0],
})
</script>

<style scoped>
.bounce-enter-active {
  animation: bounce-in 1s;
}
.bounce-leave-active {
  animation: bounce-in 1s reverse;
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
