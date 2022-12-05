<template>
  <div>
    <slot :item="item" />
    <ul v-show="isOpen" v-if="isFolder">
      <tree-item
        v-for="(child, index) in item.children"
        :key="index"
        :item="child"
      >
        <template #default="_">
          <slot :item="_.item" />
        </template>
      </tree-item>
    </ul>
  </div>
</template>

<script setup>
const props = defineProps(['item'])

// computed property that auto-updates when the prop changes
const isOpen = computed(() => {
  return props.item.isOpen
})

const isFolder = computed(
  () => props.item.children && props.item.children.length > 0
)
</script>
