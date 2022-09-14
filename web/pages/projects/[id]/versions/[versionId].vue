<template>
  <div class="w-full my-12">
    <div class="flex justify-between">
      <h1 class="mb-24 mt-12 text-6xl font-bold u-text-black">Results.</h1>
      <div>
        <UButton
          icon="heroicons-solid:arrow-circle-left"
          class="u-text-black"
          size="xl"
          variant="transparent"
          @click="$router.back()"
        />
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
        <div class="py-2 min-w-full">
          <UCard
            label-class="block font-medium u-text-gray-700"
            wrapper-class="flex items-center justify-between w-full"
            class=""
            :label="version.version_name"
            :name="version.id"
          >
            <div
              :class="`flex items-center justify-between ${
                version.is_template_launch
                  ? 'underline text-emerald-500'
                  : 'u-text-gray-700'
              } mb-5`"
            >
              <h5
                class="font-bold text-2xl tracking-tight mb-2 flex justify-between w-full"
              >
                <div>
                  <span>
                    {{ version.version_name }}
                  </span>
                </div>
                <div>
                  <span class="u-text-gray-700 text-sm">
                    #{{ version.id }}
                  </span>
                </div>
              </h5>
            </div>
            <p class="font-sm u-text-gray-700 mb-3 min-h-full">
              last launch:
              {{
                new Date(Date.parse(version.last_launch_at)).toLocaleString()
              }}
            </p>
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
                passed: {{ version.passed_number ?? 0 }}
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
                failed: {{ version.failed_number ?? 0 }}
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
                skipped: {{ version.skipped_number ?? 0 }}
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
                undefined: {{ version.undefined_number ?? 0 }}
              </span>
            </p>
            <div class="w-full">
              <ResultsChart
                :features="features"
                :version="version.id"
                :templates="version.is_template_launch ? [] : templates"
              />
            </div>
          </UCard>
        </div>
      </transition>

      <div class="my-3 text-sm min-w-full">
        <ul>
          <tree-item :item="suite">
            <template #default="slotProps">
              <div
                :class="`grid grid-cols-12 gap-4 pb-3 h-12 border-solid border-b-2 items-end
              ${
                slotProps.item.children && slotProps.item.children.length > 0
                  ? 'cursor-pointer'
                  : 'cursor-default'
              }
              ${
                slotProps.item.status === 'passed'
                  ? 'border-green-200 dark:border-green-700 bg-green-100 dark:bg-green-800 dark:text-green-200 dark:bg-opacity-5'
                  : slotProps.item.status === 'failed'
                  ? 'border-red-200 dark:border-red-700 bg-red-100 dark:bg-red-800 dark:text-red-200 dark:bg-opacity-5'
                  : slotProps.item.status === 'broken'
                  ? 'border-pink-200 dark:border-pink-700 bg-pink-100 dark:bg-pink-800 dark:text-pink-200 dark:bg-opacity-5'
                  : 'border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-900 dark:text-gray-200'
              }`"
              >
                <div
                  class="col-span-6 overflow-ellipsis whitespace-nowrap break-words inline-flex flex row-auto"
                  :style="setClassTreeLvl(slotProps.item.lvl)"
                >
                  <UButton
                    :icon="
                      !slotProps.item.children ||
                      !slotProps.item.children.length
                        ? 'heroicons-solid:chevron-right'
                        : slotProps.item.isOpen
                        ? 'heroicons-solid:chevron-up'
                        : 'heroicons-solid:chevron-down'
                    "
                    class="u-text-black"
                    size="xl"
                    variant="transparent"
                  />
                  <span class="pl-1 self-center">
                    {{ slotProps.item.name }}
                  </span>
                </div>
                <div class="col-span-1 self-center">
                  <tag-result
                    v-if="slotProps.item.status"
                    :status="slotProps.item.status"
                  />
                </div>
                <div class="col-span-2">
                  <span
                    v-if="slotProps.item.start && slotProps.item.stop"
                    style="padding-left: 0.5rem"
                  >
                    {{ slotProps.item.stop - slotProps.item.start }}ms
                  </span>
                </div>
              </div>
              <slot :item="slotProps.item"></slot>
            </template>
          </tree-item>
        </ul>
      </div>
    </client-only>
  </div>
</template>

<script setup lang="ts">
import { Version } from '~/types/versions'
import { Result } from '~/types/results'

const client = useSupabaseClient()
const route = useRoute()

const { pending, data: versions } = await useAsyncData('versions', async () => {
  const { data } = await client.rpc<Version>('version')
  return data
})

const version = computed(() => {
  return versions.value?.find((v) => v.id === Number(route.params.versionId))
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

async function fetchLaunch() {
  const { data: items, error } = await client
    .from('launches')
    .select(
      `
      id, name, created_at, user_id, is_template, origin,
      results (
        id, name, description, fullname,
        suite, parent_suite, sub_suite, feature,
        status, created_at, duration, steps
      )
    `
    )
    .eq('id', version.value?.last_launch_id)

  if (error) {
    console.log(error)
  }
  if (items.length < 1) {
    console.log('No items found')
  } else {
    return items[0]
    // this.suite = this.testTree(this.launch.results)
  }
}
function testTree(results) {
  const rootSuite = {
    name: 'Root suite',
    lvl: 1,
    isOpen: true,
    id: null,
    status: null,
    testID: null,
    children: [],
  }
  const suites = [
    ...new Set(
      results.map((r) => {
        if (r.feature) {
          return r.feature
        }
        if (r.parent_suite) {
          return r.parent_suite
        }
        if (r.suite) {
          return r.suite
        }
      })
    ),
  ]
  rootSuite.children = suites.map((s) => {
    return {
      name: s,
      lvl: 2,
      isOpen: true,
      id: null,
      status: null,
      testID: null,
      children: [],
    }
  })
  rootSuite.children.map((s) => {
    s.children = results
      .filter(
        (r) =>
          r.suite === s.name ||
          r.feature === s.name ||
          r.parent_suite === s.name
      )
      .map((r) => {
        return {
          name: r.name,
          lvl: 3,
          isOpen: false,
          id: r.id,
          testID: r.id,
          status: r.status,
          description: r.description,
          duration: r.duration,
          subsuite: r.sub_suite,
          fullname: r.fullname,
          feature: r.feature,
          steps: r.steps,
          children: [],
        }
      })
  })
  return rootSuite
}

const { data: launch } = await useAsyncData(
  `launch${route.params.versionId}`,
  async () => {
    return fetchLaunch()
  }
)

const suite = computed(() => {
  return testTree(launch.value.results)
})

const toggle = function (item) {
  item.isOpen = !item.isOpen
}
const setClassTreeLvl = function (lvl) {
  return `padding-left: ${lvl * 1.5}em`
}
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
    transform: scale(0.8);
  }
  50% {
    transform: scale(1.05);
  }
  100% {
    transform: scale(1);
  }
}
</style>
