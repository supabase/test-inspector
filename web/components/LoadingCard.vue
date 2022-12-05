<template>
  <div v-if="loading">
    <div
      v-for="i in nbOfItems"
      :key="i"
      class="border border-gray-300 shadow rounded-md p-4 w-full my-2"
      :style="`min-height: ${hght}px;`"
    >
      <div class="animate-pulse flex space-x-4">
        <div class="rounded-full bg-gray-400 h-12 w-12" />
        <div class="flex-1 space-y-4 py-1">
          <div class="h-4 bg-gray-400 rounded w-3/4" />
          <div class="space-y-4">
            <div v-for="n in Math.floor(hght/40)" :key="n" class="h-4 bg-gray-400 rounded" />
            <div class="h-4 bg-gray-400 rounded w-5/6" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    loading: {
      type: Boolean,
      default: true
    },
    nbOfItems: {
      type: Number,
      default: 5
    },
    hght: {
      type: Number,
      default: 200
    }
  }
}
</script>
