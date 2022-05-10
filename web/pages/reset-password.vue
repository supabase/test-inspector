<template>
  <main class="flex justify-center mt-16">
    <div
      class="flex flex-col w-full max-w-md h-72 px-4 py-8 bg-gray-100 rounded-lg shadow dark:bg-gray-800 sm:px-6 md:px-8 lg:px-10"
    >
      <div v-if="!submitted">
        <div
          class="self-center mb-6 text-xl font-light u-text-gray-600 sm:text-2xl"
        >
          Reset Password
        </div>
        <div>
          <form autoComplete="off" @submit.prevent="reset">
            <div class="flex flex-col mb-2">
              <div class="flex relative ">
                <span
                  class="rounded-l-md inline-flex  items-center px-3 border-t  border-l border-b   text-gray-500 shadow-sm text-sm bg-white border-gray-300 dark:bg-gray-800 dark:border-gray-900 dark:text-gray-100"
                >
                  <svg
                    width="15"
                    height="15"
                    fill="currentColor"
                    viewBox="0 0 1792 1792"
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path
                      d="M1376 768q40 0 68 28t28 68v576q0 40-28 68t-68 28h-960q-40 0-68-28t-28-68v-576q0-40 28-68t68-28h32v-320q0-185 131.5-316.5t316.5-131.5 316.5 131.5 131.5 316.5q0 26-19 45t-45 19h-64q-26 0-45-19t-19-45q0-106-75-181t-181-75-181 75-75 181v320h736z"
                    />
                  </svg>
                </span>
                <input
                  id="password"
                  v-model="password"
                  type="password"
                  class=" rounded-r-lg flex-1 appearance-none border w-full py-2 px-4  text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-green-700 focus:border-transparent bg-white border-gray-300 dark:bg-gray-800 dark:border-gray-900 dark:text-gray-100"
                  placeholder="New password"
                >
              </div>
            </div>

            <div class="flex flex-col mb-2">
              <div class="flex relative ">
                <span
                  class="rounded-l-md inline-flex  items-center px-3 border-t  border-l border-b   text-gray-500 shadow-sm text-sm bg-white border-gray-300 dark:bg-gray-800 dark:border-gray-900 dark:text-gray-100"
                >
                  <svg
                    width="15"
                    height="15"
                    fill="currentColor"
                    viewBox="0 0 1792 1792"
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path
                      d="M1376 768q40 0 68 28t28 68v576q0 40-28 68t-68 28h-960q-40 0-68-28t-28-68v-576q0-40 28-68t68-28h32v-320q0-185 131.5-316.5t316.5-131.5 316.5 131.5 131.5 316.5q0 26-19 45t-45 19h-64q-26 0-45-19t-19-45q0-106-75-181t-181-75-181 75-75 181v320h736z"
                    />
                  </svg>
                </span>
                <input
                  id="repeat"
                  v-model="repeat"
                  type="password"
                  class=" rounded-r-lg flex-1 appearance-none border w-full py-2 px-4  text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-green-700 focus:border-transparent bg-white border-gray-300 dark:bg-gray-800 dark:border-gray-900 dark:text-gray-100"
                  placeholder="Repeat password"
                >
              </div>
            </div>

            <div class="flex w-full mt-8">
              <button
                type="submit"
                class="py-2 px-4  bg-green-800 hover:bg-green-700 text-white w-full transition ease-in duration-200 text-center text-base font-semibold shadow-md focus:outline-none focus:ring-2 focus:ring-offset-2  rounded-lg "
              >
                Reset
              </button>
            </div>
          </form>
        </div>
      </div>
      <div v-if="submitted">
        <div class="bg-gray-100 dark:bg-gray-800 p-2 text-center">
          <h2
            class="text-2xl font-extrabold u-text-gray-700 sm:text-2xl"
          >
            <span class="block mb-4 u-text-gray-900">
              Your password is updated
            </span>
          </h2>
        </div>
      </div>
    </div>
  </main>
</template>

<script setup lang="ts">
const client = useSupabaseClient()
const router = useRouter()

const loading = ref(null)
const submitted = ref(null)

const password = ref('')
const repeat = ref('')

const reset = async () => {
  loading.value = true
  /* signIn sends the user a magic link */
  if (!password.value) { return }
  if (password.value !== repeat.value) {
    alert('password and repeated password mismatch')
    return
  }
  const { error } = await client.auth.update({
    password: password.value
  })
  submitted.value = true
  loading.value = false
  if (error) {
    alert(error)
  } else {
    router.push('/')
  }
}
</script>
