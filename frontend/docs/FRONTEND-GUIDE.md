# Frontend developer guide

This document is for developers who know **basic React** (components, props, `useState`, `useEffect`) and want to work productively in this **Next.js admin** codebase. Read it once top-to-bottom, then keep it as a reference.

---

## 1. What this app is

- **Next.js App Router** (`src/app/`): file-based routing. Some files are **Server Components** (default), some are **Client Components** (`"use client"` at the top).
- **Admin UI**: sidebar, header, tables, forms, charts. Business screens live in **`src/features/`**, not inside huge `page.tsx` files.
- **TypeScript (strict)**: types are expected; the compiler will block many mistakes before runtime.

**Mental model:** `app/` = URLs and shells. `features/` = product screens. `components/` = reusable UI. `lib/` = logic that is not a React tree.

---

## 2. Folder map (where to put code)

| Path | Purpose |
|------|---------|
| `src/app/` | Routes only: thin `page.tsx`, `layout.tsx`, `globals.css`. **Do not** build big UIs here. |
| `src/app/(resources)/` | Authenticated “main app” area: shared **sidebar layout**. |
| `src/app/auth/` | Login, verify, forgot/reset password — **no** sidebar chrome. |
| `src/features/<domain>/` | One folder per business area (e.g. `organizations`). List views, tables, modals. |
| `src/components/ui/` | shadcn-style primitives (Button, Table, Sidebar, …). **Rarely edit** unless fixing a bug or extending a primitive. |
| `src/components/` | Shared shell: `providers/`, `app-sidebar`, `site-header`, `table/`, `form/`, `dashboard/`, … |
| `src/lib/` | API client, env, query client factory, Zod schemas, constants, sidebar nav config. |
| `src/hooks/` | Small shared hooks (e.g. `use-mobile.ts`). |

**Import rule:** always use the `@/` alias, e.g. `@/components/ui/button`, `@/lib/api/organizations`.

---

## 3. Stack cheat sheet

| Topic | Library / pattern |
|-------|---------------------|
| Routing | Next.js App Router |
| Styling | Tailwind CSS v4 + tokens in `src/app/globals.css` |
| UI kit | shadcn/ui (style: **new-york**, base: **stone**) under `components/ui/` |
| Icons (nav) | `@tabler/icons-react` |
| Icons (UI) | `lucide-react` |
| Server vs client | Default = server component; add `"use client"` when you need hooks, events, or browser APIs |
| Server data + cache | TanStack **React Query** |
| HTTP | **axios** via a single instance: `apiClient` in `src/lib/api/api-client.ts` |
| Forms | **react-hook-form** + **@hookform/resolvers** + **Zod** schemas in `src/lib/schemas/` |
| Toasts | **sonner** (`toast()`) — provider in `AppProviders` |
| Theme | **next-themes** + optional cookies (`app_theme`, `theme_scale`) read in root `layout.tsx` |
| Tables | **@tanstack/react-table** + shared `DataTable` + `Table` primitives |

---

## 4. Routing: how we keep `page.tsx` thin

**Rule:** `page.tsx` should mostly **import one feature (or dashboard) component** and render it.

**Example — organizations list**

- URL: `/organizations`
- File: `src/app/(resources)/organizations/page.tsx` (few lines)
- Real UI: `src/features/organizations/index.tsx`

This keeps routing easy to grep and stops `page.tsx` from becoming a 500-line file.

**Layouts**

- `src/app/layout.tsx`: fonts (via CSS variables), HTML shell, **`AppProviders`**, cookies for theme defaults.
- `src/app/(resources)/layout.tsx`: **sidebar shell** (`SidebarProvider`, `AppSidebar`, `SiteHeader`) for everything under the main app.
- `src/app/auth/layout.tsx`: centered auth layout — **no** sidebar.

**Route groups** like `(resources)` **do not** change the URL. They only group files that share a layout.

---

## 5. Server vs client components (important)

### Server Component (default)

- Runs on the **server** during the request (or at build time for static parts).
- Can **`await`** data, read **`cookies()`** from `next/headers`, talk to DB, etc.
- **Cannot** use `useState`, `useEffect`, `onClick`, or browser-only APIs.

### Client Component (`"use client"`)

- Ships JavaScript to the browser; can use **hooks** and **event handlers**.
- Use when you need: forms, tables with client state, `useQuery`, `useRouter`, `window`, `localStorage`, etc.

**Pattern we use:** a **Server** `page.tsx` prefetches data into React Query’s cache, then a **Client** child reads the same cache with `useQuery` (see §7).

---

## 6. API layer: never raw `axios` in a feature

**Rule:** components and features call **named functions** in `src/lib/api/*.ts`. Those functions use the shared **`apiClient`**.

**Why**

- One place for **base URL**, headers, and **interceptors** (auth token, 401 handling).
- Easier to mock, test, and change later.

**Files**

- `src/lib/api/api-client.ts` — axios instance + interceptors.
- `src/lib/api/api-routes.ts` — path string constants (optional but preferred for consistency).
- `src/lib/api/types.ts` — shared DTO / entity types.
- `src/lib/api/<resource>.ts` — e.g. `organizations.ts`, `dashboard.ts`.

**Auth token**

- Stored in **`localStorage`** under the key in `src/lib/constants.ts` (`AUTH_TOKEN_STORAGE_KEY`).
- Request interceptor attaches `Authorization: Bearer <token>` on the **client** only.

**401 / expired session**

- Response interceptor clears the token and sends the user to **`/auth/login?next=<current path>`** (unless they are already under `/auth`).

**Environment**

- `src/lib/env.ts` — reads `NEXT_PUBLIC_API_URL`.
- If it is **not** set, some API modules return **mock data** so the UI still works locally (see `organizations.ts` / `dashboard.ts`). When the backend exists, set the env var and the same functions will call the real API.

---

## 7. Data fetching with React Query

### Concepts

- **Query:** read data (GET). Cached by `queryKey`.
- **Mutation:** write data (POST/PATCH/DELETE). Often invalidates related queries after success.

### Client-side fetch

In a client component:

```tsx
const { data, isPending, error } = useQuery({
  queryKey: ["organizations", "list"],
  queryFn: listOrganizations,
});
```

- **`queryKey`**: must be **stable** and describe the data. Same key = same cache.
- **`queryFn`**: usually a function from `src/lib/api/`.

### Server prefetch + client hydration (dashboard pattern)

Some pages load faster if the **server** fills the cache before HTML is sent:

1. In a **Server** `page.tsx`, get a server `QueryClient` via `getQueryClient()` from `src/lib/get-query-client.ts`.
2. `await queryClient.prefetchQuery({ queryKey, queryFn })`.
3. Wrap the client UI in `<HydrationBoundary state={dehydrate(queryClient)}>`.
4. Inside the client component, `useQuery` uses the **same** `queryKey` and `queryFn` — it reads hydrated data first, then behaves like a normal client query.

Reference: `src/app/(resources)/page.tsx` + `src/components/dashboard/dashboard-home.tsx`.

### Where providers live

- `QueryClientProvider` is inside **`AppProviders`** (`src/components/providers/app-providers.tsx`).
- Dev-only **React Query Devtools** are enabled in `src/lib/query-client.tsx` in development.

---

## 8. Features: how to build a screen

**Location:** `src/features/<name>/`

Typical files:

- `index.tsx` — main view (often default export). Composes layout, `useQuery` / mutations, `DataTable`, dialogs.
- `columns.tsx` — TanStack Table **column definitions** for that feature only.

**Why not put columns in `app/`?** Colocate domain UI with the feature so `app/` stays routing-only.

---

## 9. Tables (`DataTable` + TanStack Table)

1. Define columns in `columns.tsx` (`ColumnDef<YourType>[]`).
2. In the feature, call `useReactTable({ data, columns, getCoreRowModel: getCoreRowModel(), ... })`.
3. Pass `table` into `src/components/table/data-table.tsx`.

The shared `DataTable` only knows how to render **any** table instance; column logic stays in the feature.

---

## 10. Forms (react-hook-form + Zod)

1. Add or extend a schema in `src/lib/schemas/` using **Zod** (see `login.ts`).
2. In a **client** form component, use:

   ```tsx
   const form = useForm<FormType>({
     resolver: zodResolver(yourSchema),
     defaultValues: { ... },
   });
   ```

3. Wire inputs with `register` or Controller patterns as the UI needs.

**Login example:** `src/components/form/login-form.tsx` + `src/lib/schemas/login.ts`.

---

## 11. UI primitives (shadcn)

- Primitives live under **`src/components/ui/`** (Button, Card, Table, Sidebar, …).
- Prefer **composition** (small components) over copying raw HTML everywhere.
- **Adding a new shadcn component:** use the project CLI (from repo root), e.g. `pnpm dlx shadcn@latest add <component> -y`. Config is in **`components.json`** at the repo root.

**Icons**

- **Tabler** for sidebar / nav config (`src/lib/sidebar-data.ts`).
- **Lucide** for buttons, menus, header (`src/components/site-header.tsx`).

---

## 12. Auth (current demo behaviour)

- **`AuthProvider`** (`src/components/providers/auth-provider.tsx`) exposes `accessToken` and `setAccessToken`.
- **Login** demo sets a fake token and redirects to `next` query param or `/`.
- Real backends: replace the demo submit with a call to your auth API, then `setAccessToken(realToken)` (or use httpOnly cookies + a different client pattern — coordinate with backend).

---

## 13. Theming

- **`next-themes`** wraps the app in `ThemeProvider` (`src/components/providers/theme-provider.tsx`).
- Root **`layout.tsx`** may pass `defaultTheme` / `themeScale` from **cookies** for first paint alignment.
- **`ActiveThemeProvider`** sets `data-theme` / `data-theme-scale` on `<html>` for CSS or analytics hooks.

Use the theme controls in **`SiteHeader`** as a reference for calling `setTheme("light" | "dark" | "system")`.

---

## 14. Toasts

```tsx
import { toast } from "sonner";

toast.success("Saved");
toast.error("Something went wrong");
```

`Toaster` is already mounted in `AppProviders`.

---

## 15. Adding a new resource (checklist)

Use **organizations** as the template (`src/features/organizations`, `src/lib/api/organizations.ts`, `src/app/(resources)/organizations/page.tsx`).

1. **API**
   - Add `src/lib/api/<resource>.ts` with typed functions using `apiClient`.
   - Add routes to `api-routes.ts` and types to `types.ts` if shared.
2. **Schemas** (if there are forms)
   - Add `src/lib/schemas/<resource>.ts` with Zod schemas.
3. **Feature**
   - Create `src/features/<resource>/` with `index.tsx`, optional `columns.tsx`, and subfolders later (`detail/`, `edit/`, …).
4. **Route**
   - Add `src/app/(resources)/<resource>/page.tsx` that default-imports the feature view.
5. **Nav**
   - Add an item to `src/lib/sidebar-data.ts` (and any other nav constants).

---

## 16. Conventions (do / don’t)

| Do | Don’t |
|----|--------|
| Keep `page.tsx` tiny | Paste entire screens into `app/` |
| Call `lib/api/*` from features | Import `axios` directly in random components |
| Use stable `queryKey` arrays / objects | Invent new keys per render |
| Match existing naming (`columns.tsx`, `index.tsx`) | Mix random folder names per feature |
| Use `@/` imports | Long relative paths like `../../../components` |

---

## 17. Local development

```bash
pnpm install
pnpm dev
```

- **`pnpm run build`** — production build (run before opening a PR).
- **`pnpm run lint`** — ESLint.

Set **`NEXT_PUBLIC_API_URL`** in `.env.local` when the backend is ready (no trailing slash required; the client normalizes it).

---

## 18. Glossary

| Term | Meaning |
|------|---------|
| **App Router** | Next.js routing based on `app/` file structure |
| **RSC** | React Server Component |
| **Hydration** | Browser attaching to server-rendered HTML and running client React |
| **Query key** | React Query cache identifier for a piece of data |
| **Interceptor** | axios hook that runs on every request/response |
| **DTO** | Data shape returned from an API (often typed in `types.ts`) |

---

## 19. Who to ask when stuck

1. Find a **similar screen** (organizations list, login form, dashboard) and copy the pattern.
2. Check this doc’s table of contents.
3. Ask a senior with **a specific file name + what you tried** — that gets faster answers than “it doesn’t work.”

Welcome to the codebase.
