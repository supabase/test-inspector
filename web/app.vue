<template>
  <div class="min-h-screen u-bg-white">
    <Title>Test Linter</Title>
    <div class="flex items-center justify-between">
      <UButton size="xl" variant="transparent" :icon="colorMode.preference === 'dark' ? 'heroicons-outline:moon' : 'heroicons-outline:sun'" @click="toggleDark" />
      <UButton v-if="!!user" class="u-text-black" size="xl" variant="transparent" @click="logout">
        Logout
      </UButton>
      <NuxtLink v-else class="u-text-black" to="/login">
        <UButton class="u-text-black" size="xl" variant="transparent">
          Login
        </UButton>
      </NuxtLink>
    </div>
    <div class="max-w-full min-h-screen px-4 mx-auto sm:px-6 lg:px-8">
      <div class="min-h-screen -mt-[50px] flex justify-center">
        <NuxtPage />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const client = useSupabaseClient()
const user = useSupabaseUser()
const colorMode = useColorMode()

const toggleDark = () => {
  colorMode.value = colorMode.value === 'dark' ? 'light' : 'dark'
  colorMode.preference = colorMode.value
}

const logout = async () => {
  await client.auth.signOut()
  document.location.href = '/'
}
</script>
