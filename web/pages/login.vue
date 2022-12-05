<template>
  <main class="flex flex-col justify-center mt-16">
    <div
      class="flex flex-col w-full max-w-md px-4 py-8 bg-gray-100 rounded-lg shadow dark:bg-gray-800 sm:px-6 md:px-8 lg:px-10"
    >
      <div v-if="!submitted">
        <div
          class="self-center mb-6 text-xl font-light u-text-gray-600 sm:text-2xl"
        >
          Send Magic Link
        </div>
        <div>
          <form autoComplete="off" @submit.prevent="signIn">
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
                      d="M1792 710v794q0 66-47 113t-113 47h-1472q-66 0-113-47t-47-113v-794q44 49 101 87 362 246 497 345 57 42 92.5 65.5t94.5 48 110 24.5h2q51 0 110-24.5t94.5-48 92.5-65.5q170-123 498-345 57-39 100-87zm0-294q0 79-49 151t-122 123q-376 261-468 325-10 7-42.5 30.5t-54 38-52 32.5-57.5 27-50 9h-2q-23 0-50-9t-57.5-27-52-32.5-54-38-42.5-30.5q-91-64-262-182.5t-205-142.5q-62-42-117-115.5t-55-136.5q0-78 41.5-130t118.5-52h1472q65 0 112.5 47t47.5 113z"
                    />
                  </svg>
                </span>
                <input
                  id="sign-in-email"
                  v-model="email"
                  type="text"
                  class=" rounded-r-lg flex-1 appearance-none border w-full py-2 px-4  text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-green-700 focus:border-transparent bg-white border-gray-300 dark:bg-gray-800 dark:border-gray-900 dark:text-gray-100"
                  placeholder="Your email"
                >
              </div>
            </div>

            <div class="flex w-full mt-8">
              <button
                type="submit"
                class="py-2 px-4  bg-green-800 hover:bg-green-700 text-white w-full transition ease-in duration-200 text-center text-base font-semibold shadow-md focus:outline-none focus:ring-2 focus:ring-offset-2  rounded-lg "
              >
                Login
              </button>
            </div>
          </form>
        </div>
      </div>
      <div v-if="submitted">
        <div class="bg-gray-100 dark:bg-gray-800 p-2 text-center">
          <h2
            class="text-2xl font-extrabold u-text-gray-900 sm:text-2xl"
          >
            <span class="block mb-4 u-text-gray-900">
              Please check your email to sign in.
            </span>
          </h2>
        </div>
      </div>
    </div>

    <div
      v-if="!submitted"
      class="flex flex-col w-full max-w-md px-4 py-8 bg-gray-100 rounded-lg shadow dark:bg-gray-800 sm:px-6 md:px-8 lg:px-10 mt-16"
    >
      <div>
        <div
          class="self-center mb-6 text-xl font-light u-text-gray-600 sm:text-2xl"
        >
          Login with password
        </div>
        <div>
          <form autoComplete="off" @submit.prevent="signInPassword">
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
                      d="M1792 710v794q0 66-47 113t-113 47h-1472q-66 0-113-47t-47-113v-794q44 49 101 87 362 246 497 345 57 42 92.5 65.5t94.5 48 110 24.5h2q51 0 110-24.5t94.5-48 92.5-65.5q170-123 498-345 57-39 100-87zm0-294q0 79-49 151t-122 123q-376 261-468 325-10 7-42.5 30.5t-54 38-52 32.5-57.5 27-50 9h-2q-23 0-50-9t-57.5-27-52-32.5-54-38-42.5-30.5q-91-64-262-182.5t-205-142.5q-62-42-117-115.5t-55-136.5q0-78 41.5-130t118.5-52h1472q65 0 112.5 47t47.5 113z"
                    />
                  </svg>
                </span>
                <input
                  id="sign-in-email"
                  v-model="email"
                  type="text"
                  class=" rounded-r-lg flex-1 appearance-none border w-full py-2 px-4  text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-green-700 focus:border-transparent bg-white border-gray-300 dark:bg-gray-800 dark:border-gray-900 dark:text-gray-100"
                  placeholder="Your email"
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
                  id="password"
                  v-model="password"
                  type="password"
                  class=" rounded-r-lg flex-1 appearance-none border w-full py-2 px-4  text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-green-700 focus:border-transparent bg-white border-gray-300 dark:bg-gray-800 dark:border-gray-900 dark:text-gray-100"
                  placeholder="Password"
                >
              </div>
            </div>

            <div class="flex w-full mt-8">
              <button
                type="submit"
                class="py-2 px-4  bg-green-800 hover:bg-green-700 text-white w-full transition ease-in duration-200 text-center text-base font-semibold shadow-md focus:outline-none focus:ring-2 focus:ring-offset-2  rounded-lg "
              >
                Login
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
    <div
      v-if="!submitted"
      class="w-full max-w-md my-4 text-center text-gray-600 text-sm sm:text-md dark:text-dark-accent"
    >
      <nuxt-link to="/forgot-password">
        Forgot Password?
      </nuxt-link>
    </div>
  </main>
</template>

<script setup lang="ts">
const user = useSupabaseUser()
const client = useSupabaseClient()
const router = useRouter()

const loading = ref(null)
const submitted = ref(false)

const email = ref('')
const password = ref('')

watchEffect(() => {
  if (user.value) {
    router.push('/')
  }
})

const signIn = async () => {
  loading.value = true
  /* signIn sends the user a magic link */
  if (email.value.length === 0) { return }
  const { error } = await client.auth.signIn({
    email: email.value
  })
  if (error) {
    alert(error)
  }
  loading.value = false
  submitted.value = true
}

const signInPassword = async () => {
  loading.value = true
  /* signIn sends the user a magic link */
  if (email.value.length === 0) { return }
  const { error } = await client.auth.signIn({
    email: email.value,
    password: password.value
  })
  if (error) {
    alert(error)
  }
  loading.value = false
}
</script>
