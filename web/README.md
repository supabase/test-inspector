# Test-inspector web app

Live: <https://test-inspector.fly.dev>

- Frontend:
  - [Nuxt 3](https://v3.nuxtjs.org/) - a Vuejs framework.
  - [Tailwind](https://tailwindcss.com/) for styling and layout.
  - [Supabase Module](https://github.com/nuxt-community/supabase-module) for user management and supabase data client.
- Backend:
  - [app.supabase.io](https://app.supabase.io/): hosted Postgres database with restful API for usage with Supabase.js.

## Setup

Make sure to install the dependencies

```bash
npm install
```

Fill the `.env` with the Supabase environment variables:

- `cp .env.example .env`

```bash
SUPABASE_URL="https://example.supabase.com"
SUPABASE_KEY="<your_key>"
```

## Development

Start the development server on <http://localhost:3000>

```bash
npm run dev
```

## Production

Build the application for production:

```bash
npm run build
```

Checkout the [deployment documentation](https://v3.nuxtjs.org/docs/deployment).
