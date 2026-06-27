import { SignupPageFeature } from "@/features/auth/signup-page";

type SignupPageProps = {
  searchParams?: Promise<{ role?: string | string[] }>;
};

export default async function SignupPage({ searchParams }: SignupPageProps) {
  const sp = searchParams ? await searchParams : {};
  const raw = sp.role;
  const roleParam = Array.isArray(raw) ? raw[0] : raw;
  return <SignupPageFeature initialRole={roleParam} />;
}
