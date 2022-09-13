<template>
  <div>
    <span
      class="py-1 px-2 bg-black hover:bg-blue-300 text-white w-full transition ease-in duration-200 text-center text-base font-semibold shadow-md focus:outline-none focus:ring-2 focus:ring-offset-2 rounded-lg tooltip"
      v-if="!props.status"
    >
      U
      <span class="tooltiptext"> Unknown </span>
    </span>
    <span
      class="py-1 px-2 bg-green-600 hover:bg-green-500 text-white w-full transition ease-in duration-200 text-center text-base font-semibold shadow-md focus:outline-none focus:ring-2 focus:ring-offset-2 rounded-lg tooltip"
      v-if="props.status === 'passed'"
    >
      P
      <span class="tooltiptext"> Passed </span>
    </span>
    <span
      class="py-1 px-2 bg-red-600 hover:bg-red-500 text-white w-full transition ease-in duration-200 text-center text-base font-semibold shadow-md focus:outline-none focus:ring-2 focus:ring-offset-2 rounded-lg tooltip"
      v-if="props.status === 'failed'"
    >
      F
      <span class="tooltiptext"> Failed </span>
    </span>
    <span
      class="py-1 px-2 bg-yellow-600 hover:bg-yellow-500 text-white w-full transition ease-in duration-200 text-center text-base font-semibold shadow-md focus:outline-none focus:ring-2 focus:ring-offset-2 rounded-lg tooltip"
      v-if="props.status === 'skipped'"
    >
      S
      <span class="tooltiptext"> Skipped </span>
    </span>
    <span
      class="py-1 px-2 bg-pink-600 hover:bg-pink-500 text-white w-full transition ease-in duration-200 text-center text-base font-semibold shadow-md focus:outline-none focus:ring-2 focus:ring-offset-2 rounded-lg tooltip"
      v-if="props.status === 'broken'"
    >
      B
      <span class="tooltiptext"> Broken </span>
    </span>
    <span
      class="py-1 px-2 bg-black hover:bg-gray-900 text-white w-full transition ease-in duration-200 text-center text-base font-semibold shadow-md focus:outline-none focus:ring-2 focus:ring-offset-2 rounded-lg tooltip"
      v-if="props.status === 'undefined'"
    >
      U
      <span class="tooltiptext"> Undefined </span>
    </span>
  </div>
</template>

<script setup>
const props = defineProps({
  status: String,
})
</script>

<style scoped>
.tooltip {
  position: relative;
}

.tooltip .tooltiptext {
  visibility: hidden;
  width: 120px;
  background-color: #555;
  color: #fff;
  text-align: center;
  border-radius: 6px;
  padding: 5px 0;
  position: absolute;
  z-index: 1;
  bottom: 125%;
  left: 50%;
  margin-left: -60px;
  opacity: 0;
  transition: opacity 1s;
}

.tooltip .tooltiptext::after {
  content: '';
  position: absolute;
  top: 100%;
  left: 50%;
  margin-left: -5px;
  border-width: 5px;
  border-style: solid;
  border-color: #555 transparent transparent transparent;
}

.tooltip:hover .tooltiptext {
  visibility: visible;
  opacity: 1;
}
</style>
